package task_test

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/NeRo0128/brain-cli/internal/adapters/database"
	"github.com/NeRo0128/brain-cli/internal/adapters/database/repositories"
	"github.com/NeRo0128/brain-cli/internal/adapters/executor"
	coretask "github.com/NeRo0128/brain-cli/internal/core/task"
	"github.com/NeRo0128/brain-cli/internal/core/tool"
	uc "github.com/NeRo0128/brain-cli/internal/usecases/task"
)

// testEnv agrupa todas las dependencias de un test.
type testEnv struct {
	executorUC *uc.Executor
	taskRepo   *repositories.TaskRepository
	toolRepo   *repositories.ToolRepository
	execRepo   *repositories.ExecutionRepository
}

func newTestEnv(t *testing.T) *testEnv {
	t.Helper()
	dir := t.TempDir()
	db, err := database.New(context.Background(), database.Options{
		Path:        filepath.Join(dir, "test.db"),
		AutoMigrate: true,
	})
	if err != nil {
		t.Fatalf("abriendo DB: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	for _, tbl := range []string{"executions", "task_tags", "tasks", "tools", "tags"} {
		if _, err := db.DB().Exec("DELETE FROM " + tbl); err != nil {
			t.Fatalf("limpiando %s: %v", tbl, err)
		}
	}

	taskRepo := repositories.NewTaskRepository(db.DB())
	toolRepo := repositories.NewToolRepository(db.DB())
	execRepo := repositories.NewExecutionRepository(db.DB())

	// Inyectamos la factory real. El use case NO importa el adapter,
	// solo recibe una función con la firma ExecutorFactory.
	return &testEnv{
		executorUC: uc.NewExecutor(taskRepo, toolRepo, execRepo, executor.New),
		taskRepo:   taskRepo,
		toolRepo:   toolRepo,
		execRepo:   execRepo,
	}
}

// --- helpers de creación de fixtures -----------------------------------

func mkTool(t *testing.T, repo *repositories.ToolRepository, name, content string) *tool.Tool {
	t.Helper()
	tl := &tool.Tool{
		Name:           name,
		ScriptType:     tool.ScriptTypeBash,
		Category:       tool.CategoryUtils,
		ScriptContent:  content,
		TimeoutSeconds: 5,
		Version:        1,
	}
	if err := repo.Create(context.Background(), tl); err != nil {
		t.Fatalf("creando tool %q: %v", name, err)
	}
	return tl
}

func mkTask(t *testing.T, repo *repositories.TaskRepository, id string, toolID int) *coretask.Task {
	t.Helper()
	tk := &coretask.Task{
		ID:       id,
		Name:     "Task " + id,
		Type:     coretask.TaskTypeScript,
		ToolID:   &toolID,
		Priority: coretask.PriorityMedium,
		IsActive: true,
	}
	if err := repo.Create(context.Background(), tk); err != nil {
		t.Fatalf("creando task %q: %v", id, err)
	}
	return tk
}

// --- tests -------------------------------------------------------------

func TestExecute_Success(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	tl := mkTool(t, env.toolRepo, "echo-tool", "echo hola-mundo")
	mkTask(t, env.taskRepo, "task-ok", tl.ID)

	exec, err := env.executorUC.Execute(ctx, "task-ok")
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if exec.Status != "completed" {
		t.Errorf("status = %q, quiero completed", exec.Status)
	}
	if !strings.Contains(exec.Output, "hola-mundo") {
		t.Errorf("output: %q", exec.Output)
	}
	if exec.ExitCode == nil || *exec.ExitCode != 0 {
		t.Errorf("exit code = %v", exec.ExitCode)
	}
	if exec.FinishedAt == nil {
		t.Error("FinishedAt no seteado")
	}
}

func TestExecute_FailedExit(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	tl := mkTool(t, env.toolRepo, "fail-tool", "echo oops; exit 7")
	mkTask(t, env.taskRepo, "task-fail", tl.ID)

	exec, err := env.executorUC.Execute(ctx, "task-fail")
	if err != nil {
		t.Fatalf("Execute no debe fallar por exit != 0: %v", err)
	}
	if exec.Status != "failed" {
		t.Errorf("status = %q, quiero failed", exec.Status)
	}
	if exec.ExitCode == nil || *exec.ExitCode != 7 {
		t.Errorf("exit code = %v, quiero 7", exec.ExitCode)
	}
	if !strings.Contains(exec.Error, "7") {
		t.Errorf("error: %q", exec.Error)
	}
}

func TestExecute_TaskNotFound(t *testing.T) {
	env := newTestEnv(t)
	_, err := env.executorUC.Execute(context.Background(), "no-existe")
	if !errors.Is(err, uc.ErrTaskNotFound) {
		t.Fatalf("esperaba ErrTaskNotFound, dio: %v", err)
	}
}

func TestExecute_AITask_Unsupported(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	tk := &coretask.Task{
		ID:         "task-ai",
		Name:       "AI Task",
		Type:       coretask.TaskTypeAI,
		RequiresAI: true,
		AIPrompt:   "{{.Input}}",
		Priority:   coretask.PriorityLow,
		IsActive:   true,
	}
	if err := env.taskRepo.Create(ctx, tk); err != nil {
		t.Fatal(err)
	}

	_, err := env.executorUC.Execute(ctx, "task-ai")
	if !errors.Is(err, uc.ErrUnsupportedTaskType) {
		t.Fatalf("esperaba ErrUnsupportedTaskType, dio: %v", err)
	}
}

func TestExecute_Cancellation(t *testing.T) {
	env := newTestEnv(t)

	tl := mkTool(t, env.toolRepo, "slow-tool", "sleep 5; echo never")
	mkTask(t, env.taskRepo, "task-slow", tl.ID)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	exec, err := env.executorUC.Execute(ctx, "task-slow")
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if exec.Status != "cancelled" {
		t.Errorf("status = %q, quiero cancelled", exec.Status)
	}
	if exec.ExitCode == nil || *exec.ExitCode != -1 {
		t.Errorf("exit code = %v, quiero -1", exec.ExitCode)
	}
}

func TestExecute_PersistsExecution(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	tl := mkTool(t, env.toolRepo, "echo-tool", "echo persisted")
	mkTask(t, env.taskRepo, "task-persist", tl.ID)

	exec, err := env.executorUC.Execute(ctx, "task-persist")
	if err != nil {
		t.Fatal(err)
	}

	// Recuperar de la DB (no del objeto en memoria)
	fromDB, err := env.execRepo.GetByID(ctx, exec.ID)
	if err != nil {
		t.Fatalf("execution no persistida: %v", err)
	}
	if fromDB.Status != "completed" {
		t.Errorf("status en DB = %q", fromDB.Status)
	}
	if !strings.Contains(fromDB.Output, "persisted") {
		t.Errorf("output en DB: %q", fromDB.Output)
	}
}
