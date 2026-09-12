package repositories_test

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/NeRo0128/brain-cli/internal/adapters/database"
	"github.com/NeRo0128/brain-cli/internal/adapters/database/repositories"
	"github.com/NeRo0128/brain-cli/internal/core/execution"
	"github.com/NeRo0128/brain-cli/internal/core/task"
	"github.com/NeRo0128/brain-cli/internal/core/tool"
)

func newTestExecutionRepo(t *testing.T) (*repositories.ExecutionRepository, *repositories.TaskRepository, *repositories.ToolRepository) {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "test.db")

	db, err := database.New(context.Background(), database.Options{
		Path:        path,
		AutoMigrate: true,
	})
	if err != nil {
		t.Fatalf("abriendo DB: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	// Orden FK: executions → tasks → tools
	for _, tbl := range []string{"executions", "task_tags", "tasks", "tools", "tags"} {
		if _, err := db.DB().Exec("DELETE FROM " + tbl); err != nil {
			t.Fatalf("limpiando %s: %v", tbl, err)
		}
	}

	return repositories.NewExecutionRepository(db.DB()),
		repositories.NewTaskRepository(db.DB()),
		repositories.NewToolRepository(db.DB())
}

// setupTask crea tool + task y devuelve el slug de la task.
func setupTask(t *testing.T, taskRepo *repositories.TaskRepository, toolRepo *repositories.ToolRepository) string {
	t.Helper()
	ctx := context.Background()

	uniq := fmt.Sprintf("%d", time.Now().UnixNano())

	tl := &tool.Tool{
		Name:           "t-" + uniq,
		ScriptType:     tool.ScriptTypeBash,
		Category:       tool.CategoryUtils,
		ScriptContent:  "echo hi",
		TimeoutSeconds: 30,
		Version:        1,
	}
	if err := toolRepo.Create(ctx, tl); err != nil {
		t.Fatalf("creando tool: %v", err)
	}

	tk := &task.Task{
		ID:       "task-" + uniq,
		Name:     "Test Task" + uniq,
		Type:     task.TaskTypeScript,
		ToolID:   &tl.ID,
		Priority: task.PriorityMedium,
		IsActive: true,
	}
	if err := taskRepo.Create(ctx, tk); err != nil {
		t.Fatalf("creando task: %v", err)
	}
	return tk.ID
}

func runningExecution(taskID string) *execution.Execution {
	return &execution.Execution{
		TaskID:      taskID,
		Status:      execution.StatusRunning,
		TriggeredBy: execution.TriggerManual,
		StartedAt:   time.Now().UTC(),
	}
}

func TestExecutionCreate_FillsID(t *testing.T) {
	execRepo, taskRepo, toolRepo := newTestExecutionRepo(t)
	ctx := context.Background()
	taskID := setupTask(t, taskRepo, toolRepo)

	e := runningExecution(taskID)
	if err := execRepo.Create(ctx, e); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if e.ID <= 0 {
		t.Errorf("ID no fue rellenado: %d", e.ID)
	}
}

func TestExecutionCreate_RoundtripWithParams(t *testing.T) {
	execRepo, taskRepo, toolRepo := newTestExecutionRepo(t)
	ctx := context.Background()
	taskID := setupTask(t, taskRepo, toolRepo)

	e := runningExecution(taskID)
	e.Params = map[string]string{"SSID": "MiCasa"}
	if err := execRepo.Create(ctx, e); err != nil {
		t.Fatal(err)
	}

	got, err := execRepo.GetByID(ctx, e.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Params["SSID"] != "MiCasa" {
		t.Errorf("Params mal: %+v", got.Params)
	}
	if got.ExitCode != nil {
		t.Errorf("ExitCode debería ser nil, es %v", *got.ExitCode)
	}
	if got.FinishedAt != nil {
		t.Errorf("FinishedAt debería ser nil")
	}
}

func TestExecutionUpdate_CompletesRun(t *testing.T) {
	execRepo, taskRepo, toolRepo := newTestExecutionRepo(t)
	ctx := context.Background()
	taskID := setupTask(t, taskRepo, toolRepo)

	e := runningExecution(taskID)
	if err := execRepo.Create(ctx, e); err != nil {
		t.Fatal(err)
	}

	// Simulamos que terminó OK
	now := time.Now().UTC()
	exit := 0
	e.Status = execution.StatusCompleted
	e.Output = "todo bien"
	e.ExitCode = &exit
	e.FinishedAt = &now

	if err := execRepo.Update(ctx, e); err != nil {
		t.Fatalf("Update: %v", err)
	}

	got, _ := execRepo.GetByID(ctx, e.ID)
	if got.Status != execution.StatusCompleted {
		t.Errorf("Status = %q", got.Status)
	}
	if got.ExitCode == nil || *got.ExitCode != 0 {
		t.Errorf("ExitCode mal: %v", got.ExitCode)
	}
	if got.FinishedAt == nil {
		t.Error("FinishedAt no se persistió")
	}
	if got.Output != "todo bien" {
		t.Errorf("Output = %q", got.Output)
	}
}

func TestExecutionUpdate_FailedRun(t *testing.T) {
	execRepo, taskRepo, toolRepo := newTestExecutionRepo(t)
	ctx := context.Background()
	taskID := setupTask(t, taskRepo, toolRepo)

	e := runningExecution(taskID)
	if err := execRepo.Create(ctx, e); err != nil {
		t.Fatal(err)
	}

	now := time.Now().UTC()
	exit := 1
	e.Status = execution.StatusFailed
	e.Error = "exit status 1"
	e.ExitCode = &exit
	e.FinishedAt = &now

	if err := execRepo.Update(ctx, e); err != nil {
		t.Fatal(err)
	}

	got, _ := execRepo.GetByID(ctx, e.ID)
	if got.Error != "exit status 1" {
		t.Errorf("Error = %q", got.Error)
	}
	if got.ExitCode == nil || *got.ExitCode != 1 {
		t.Errorf("ExitCode mal: %v", got.ExitCode)
	}
}

func TestExecutionUpdate_NotFound(t *testing.T) {
	execRepo, _, _ := newTestExecutionRepo(t)
	e := runningExecution("x")
	e.ID = 99999
	err := execRepo.Update(context.Background(), e)
	if !errors.Is(err, execution.ErrNotFound) {
		t.Fatalf("esperaba ErrNotFound, dio: %v", err)
	}
}

func TestExecutionGetByID_NotFound(t *testing.T) {
	execRepo, _, _ := newTestExecutionRepo(t)
	_, err := execRepo.GetByID(context.Background(), 99999)
	if !errors.Is(err, execution.ErrNotFound) {
		t.Fatalf("esperaba ErrNotFound, dio: %v", err)
	}
}

func TestExecutionListByTask(t *testing.T) {
	execRepo, taskRepo, toolRepo := newTestExecutionRepo(t)
	ctx := context.Background()
	taskID := setupTask(t, taskRepo, toolRepo)

	for i := 0; i < 3; i++ {
		e := runningExecution(taskID)
		e.StartedAt = time.Now().UTC().Add(time.Duration(i) * time.Second)
		if err := execRepo.Create(ctx, e); err != nil {
			t.Fatal(err)
		}
	}

	got, err := execRepo.ListByTask(ctx, taskID, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 3 {
		t.Errorf("ListByTask = %d, quiero 3", len(got))
	}
	// Orden DESC por started_at
	for i := 1; i < len(got); i++ {
		if got[i].StartedAt.After(got[i-1].StartedAt) {
			t.Error("orden no es DESC")
		}
	}
}

func TestExecutionListByTask_RespectsLimit(t *testing.T) {
	execRepo, taskRepo, toolRepo := newTestExecutionRepo(t)
	ctx := context.Background()
	taskID := setupTask(t, taskRepo, toolRepo)

	for i := 0; i < 5; i++ {
		e := runningExecution(taskID)
		e.StartedAt = time.Now().UTC().Add(time.Duration(i) * time.Second)
		if err := execRepo.Create(ctx, e); err != nil {
			t.Fatal(err)
		}
	}

	got, err := execRepo.ListByTask(ctx, taskID, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Errorf("limit no respetado: %d", len(got))
	}
}

func TestExecutionDeleteBefore(t *testing.T) {
	execRepo, taskRepo, toolRepo := newTestExecutionRepo(t)
	ctx := context.Background()
	taskID := setupTask(t, taskRepo, toolRepo)

	old := runningExecution(taskID)
	old.StartedAt = time.Now().UTC().Add(-48 * time.Hour)
	if err := execRepo.Create(ctx, old); err != nil {
		t.Fatal(err)
	}

	recent := runningExecution(taskID)
	recent.StartedAt = time.Now().UTC()
	if err := execRepo.Create(ctx, recent); err != nil {
		t.Fatal(err)
	}

	cutoff := time.Now().UTC().Add(-24 * time.Hour)
	n, err := execRepo.DeleteBefore(ctx, cutoff)
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Errorf("DeleteBefore borró %d, quiero 1", n)
	}

	remaining, _ := execRepo.ListByTask(ctx, taskID, 0)
	if len(remaining) != 1 {
		t.Errorf("quedan %d executions, quiero 1", len(remaining))
	}
}

func TestExecutionListRecent(t *testing.T) {
	execRepo, taskRepo, toolRepo := newTestExecutionRepo(t)
	ctx := context.Background()

	// Creamos 3 tasks distintas, 1 execution cada una
	for i := 0; i < 3; i++ {
		taskID := setupTask(t, taskRepo, toolRepo)
		e := runningExecution(taskID)
		e.StartedAt = time.Now().UTC().Add(time.Duration(i) * time.Second)
		if err := execRepo.Create(ctx, e); err != nil {
			t.Fatal(err)
		}
	}

	got, err := execRepo.ListRecent(ctx, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Errorf("ListRecent = %d, quiero 2", len(got))
	}
}
