package repositories_test

import (
	"context"
	"errors"
	"path/filepath"
	"sort"
	"testing"

	"github.com/NeRo0128/brain-cli/internal/adapters/database"
	"github.com/NeRo0128/brain-cli/internal/adapters/database/repositories"
	"github.com/NeRo0128/brain-cli/internal/core/task"
	"github.com/NeRo0128/brain-cli/internal/core/tool"
)

func newTestTaskRepo(t *testing.T) (*repositories.TaskRepository, *repositories.ToolRepository) {
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

	// Orden FK: tasks → tools
	if _, err := db.DB().Exec("DELETE FROM task_tags"); err != nil {
		t.Fatalf("limpiando task_tags: %v", err)
	}
	if _, err := db.DB().Exec("DELETE FROM tasks"); err != nil {
		t.Fatalf("limpiando tasks: %v", err)
	}
	if _, err := db.DB().Exec("DELETE FROM tools"); err != nil {
		t.Fatalf("limpiando tools: %v", err)
	}
	if _, err := db.DB().Exec("DELETE FROM tags"); err != nil {
		t.Fatalf("limpiando tags: %v", err)
	}

	return repositories.NewTaskRepository(db.DB()), repositories.NewToolRepository(db.DB())
}

// createTool crea un tool dummy y devuelve su ID.
func createTool(t *testing.T, repo *repositories.ToolRepository, name string) int {
	t.Helper()
	tl := &tool.Tool{
		Name:           name,
		ScriptType:     tool.ScriptTypeBash,
		Category:       tool.CategoryUtils,
		ScriptContent:  "echo test",
		TimeoutSeconds: 30,
		Version:        1,
	}
	if err := repo.Create(context.Background(), tl); err != nil {
		t.Fatalf("creando tool %q: %v", name, err)
	}
	return tl.ID
}

func sampleScriptTask(toolID int) *task.Task {
	return &task.Task{
		ID:       "sample-task",
		Name:     "Sample Task",
		Type:     task.TaskTypeScript,
		ToolID:   &toolID,
		Priority: task.PriorityMedium,
		IsActive: true,
		Params:   map[string]string{"key": "value"},
		Tags:     []string{"dev", "test"},
	}
}

func TestTaskCreate_FillsCreatedAt(t *testing.T) {
	taskRepo, toolRepo := newTestTaskRepo(t)
	toolID := createTool(t, toolRepo, "dummy")

	tk := sampleScriptTask(toolID)
	if err := taskRepo.Create(context.Background(), tk); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if tk.CreatedAt.IsZero() {
		t.Error("CreatedAt no fue rellenado")
	}
	if tk.UpdatedAt != nil {
		t.Error("UpdatedAt debería ser nil tras Create")
	}
}

func TestTaskCreate_RoundtripWithParamsAndTags(t *testing.T) {
	taskRepo, toolRepo := newTestTaskRepo(t)
	ctx := context.Background()
	toolID := createTool(t, toolRepo, "dummy")

	orig := sampleScriptTask(toolID)
	orig.Params = map[string]string{"SSID": "MiCasa", "VPN": "office"}
	orig.Tags = []string{"network", "automation"}

	if err := taskRepo.Create(ctx, orig); err != nil {
		t.Fatal(err)
	}

	got, err := taskRepo.GetByID(ctx, orig.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Params["SSID"] != "MiCasa" || got.Params["VPN"] != "office" {
		t.Errorf("Params mal: %+v", got.Params)
	}
	sort.Strings(got.Tags)
	if len(got.Tags) != 2 || got.Tags[0] != "automation" || got.Tags[1] != "network" {
		t.Errorf("Tags mal: %+v", got.Tags)
	}
}

func TestTaskCreate_DuplicateID(t *testing.T) {
	taskRepo, toolRepo := newTestTaskRepo(t)
	ctx := context.Background()
	toolID := createTool(t, toolRepo, "dummy")

	tk := sampleScriptTask(toolID)
	if err := taskRepo.Create(ctx, tk); err != nil {
		t.Fatal(err)
	}
	err := taskRepo.Create(ctx, tk)
	if !errors.Is(err, task.ErrDuplicateID) {
		t.Fatalf("esperaba ErrDuplicateID, dio: %v", err)
	}
}

func TestTaskCreate_AITask_NoTool(t *testing.T) {
	taskRepo, _ := newTestTaskRepo(t)
	ctx := context.Background()

	tk := &task.Task{
		ID:         "chat-test",
		Name:       "Chat Test",
		Type:       task.TaskTypeAI,
		RequiresAI: true,
		AIPrompt:   "{{.Input}}",
		Priority:   task.PriorityLow,
		IsActive:   true,
	}
	if err := taskRepo.Create(ctx, tk); err != nil {
		t.Fatalf("Create ai task: %v", err)
	}

	got, err := taskRepo.GetByID(ctx, "chat-test")
	if err != nil {
		t.Fatal(err)
	}
	if got.ToolID != nil {
		t.Errorf("ToolID debería ser nil, es %v", *got.ToolID)
	}
	if !got.RequiresAI {
		t.Error("RequiresAI no se persistió")
	}
}

func TestTaskCreate_SchemaRejectsInconsistentAI(t *testing.T) {
	taskRepo, toolRepo := newTestTaskRepo(t)
	ctx := context.Background()
	toolID := createTool(t, toolRepo, "dummy")

	// CHECK constraint: type=ai no debe tener tool_id
	tk := &task.Task{
		ID:         "bad-ai",
		Name:       "Bad AI",
		Type:       task.TaskTypeAI,
		ToolID:     &toolID, // ← viola CHECK
		AIPrompt:   "x",
		RequiresAI: true,
		Priority:   task.PriorityLow,
		IsActive:   true,
	}
	err := taskRepo.Create(ctx, tk)
	if err == nil {
		t.Fatal("esperaba error de CHECK constraint")
	}
}

func TestTaskUpdate_SyncsTags(t *testing.T) {
	taskRepo, toolRepo := newTestTaskRepo(t)
	ctx := context.Background()
	toolID := createTool(t, toolRepo, "dummy")

	tk := sampleScriptTask(toolID)
	tk.Tags = []string{"dev", "test"}
	if err := taskRepo.Create(ctx, tk); err != nil {
		t.Fatal(err)
	}

	// Cambiar tags
	tk.Tags = []string{"production"}
	if err := taskRepo.Update(ctx, tk); err != nil {
		t.Fatal(err)
	}

	got, _ := taskRepo.GetByID(ctx, tk.ID)
	if len(got.Tags) != 1 || got.Tags[0] != "production" {
		t.Errorf("Tags mal tras Update: %+v", got.Tags)
	}
	if tk.UpdatedAt == nil {
		t.Error("UpdatedAt no fue rellenado")
	}
}

func TestTaskList_ExcludesDeleted(t *testing.T) {
	taskRepo, toolRepo := newTestTaskRepo(t)
	ctx := context.Background()
	toolID := createTool(t, toolRepo, "dummy")

	for _, id := range []string{"a", "b", "c"} {
		tk := sampleScriptTask(toolID)
		tk.ID, tk.Name = id, "Task "+id
		tk.Tags = nil
		if err := taskRepo.Create(ctx, tk); err != nil {
			t.Fatal(err)
		}
	}

	if err := taskRepo.Delete(ctx, "b"); err != nil {
		t.Fatal(err)
	}

	got, err := taskRepo.List(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Errorf("List = %d tasks, quiero 2", len(got))
	}
}

func TestTaskListByTag_RequiresAllTags(t *testing.T) {
	taskRepo, toolRepo := newTestTaskRepo(t)
	ctx := context.Background()
	toolID := createTool(t, toolRepo, "dummy")

	t1 := sampleScriptTask(toolID)
	t1.ID, t1.Name = "t1", "T1"
	t1.Tags = []string{"a", "b"}
	t2 := sampleScriptTask(toolID)
	t2.ID, t2.Name = "t2", "T2"
	t2.Tags = []string{"a"}

	for _, tk := range []*task.Task{t1, t2} {
		if err := taskRepo.Create(ctx, tk); err != nil {
			t.Fatal(err)
		}
	}

	got, err := taskRepo.ListByTag(ctx, []string{"a", "b"})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].ID != "t1" {
		t.Errorf("ListByTag [a,b] = %+v", got)
	}
}

func TestTaskSoftDelete_ThenCreateSameID_Fails(t *testing.T) {
	taskRepo, toolRepo := newTestTaskRepo(t)
	ctx := context.Background()
	toolID := createTool(t, toolRepo, "dummy")

	tk := sampleScriptTask(toolID)
	tk.Tags = nil
	if err := taskRepo.Create(ctx, tk); err != nil {
		t.Fatal(err)
	}
	if err := taskRepo.Delete(ctx, tk.ID); err != nil {
		t.Fatal(err)
	}

	// Recrear con mismo ID: PK es único (no partial index)
	// → Esto es una limitación conocida, falla hasta que hagamos migración 003.
	tk2 := sampleScriptTask(toolID)
	tk2.Tags = nil
	err := taskRepo.Create(ctx, tk2)
	if !errors.Is(err, task.ErrDuplicateID) {
		t.Logf("Nota: recrear tras soft-delete dio: %v", err)
	}
}
