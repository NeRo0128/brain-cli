package task_test

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/NeRo0128/brain-cli/internal/adapters/database"
	"github.com/NeRo0128/brain-cli/internal/adapters/database/repositories"
	coretask "github.com/NeRo0128/brain-cli/internal/core/task"
	"github.com/NeRo0128/brain-cli/internal/core/tool"
	uc "github.com/NeRo0128/brain-cli/internal/usecases/task"
)

func newManagerEnv(t *testing.T) (*uc.Manager, *repositories.TaskRepository, *repositories.ToolRepository) {
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
	return uc.NewManager(taskRepo, toolRepo), taskRepo, toolRepo
}

func seedTool(t *testing.T, repo *repositories.ToolRepository, name string) *tool.Tool {
	t.Helper()
	tl := &tool.Tool{
		Name:           name,
		ScriptType:     tool.ScriptTypeBash,
		Category:       tool.CategoryUtils,
		ScriptContent:  "echo hi",
		TimeoutSeconds: 30,
		Version:        1,
	}
	if err := repo.Create(context.Background(), tl); err != nil {
		t.Fatalf("creando tool: %v", err)
	}
	return tl
}

func TestManager_Create_OK(t *testing.T) {
	mgr, taskRepo, toolRepo := newManagerEnv(t)
	ctx := context.Background()
	tl := seedTool(t, toolRepo, "test-tool")

	tk, err := mgr.Create(ctx, uc.TaskInput{
		ID:       "nueva-task",
		Name:     "Nueva Task",
		Type:     coretask.TaskTypeScript,
		ToolID:   &tl.ID,
		Priority: coretask.PriorityMedium,
		IsActive: true,
		Tags:     []string{"dev", "test"},
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if tk.ID != "nueva-task" {
		t.Errorf("ID = %q", tk.ID)
	}

	// Verificar persistencia
	got, err := taskRepo.GetByID(ctx, "nueva-task")
	if err != nil {
		t.Fatalf("no persistida: %v", err)
	}
	if got.Name != "Nueva Task" {
		t.Errorf("Name en DB = %q", got.Name)
	}
}

func TestManager_Create_RequiresID(t *testing.T) {
	mgr, _, _ := newManagerEnv(t)
	_, err := mgr.Create(context.Background(), uc.TaskInput{
		Name: "sin id",
		Type: coretask.TaskTypeAI,
	})
	if !errors.Is(err, uc.ErrTaskIDRequired) {
		t.Fatalf("esperaba ErrTaskIDRequired, dio: %v", err)
	}
}

func TestManager_Create_RequiresName(t *testing.T) {
	mgr, _, _ := newManagerEnv(t)
	_, err := mgr.Create(context.Background(), uc.TaskInput{
		ID:   "sin-nombre",
		Type: coretask.TaskTypeAI,
	})
	if !errors.Is(err, uc.ErrTaskNameRequired) {
		t.Fatalf("esperaba ErrTaskNameRequired, dio: %v", err)
	}
}

func TestManager_Create_RejectsInvalidToolID(t *testing.T) {
	mgr, _, _ := newManagerEnv(t)
	missing := 9999
	_, err := mgr.Create(context.Background(), uc.TaskInput{
		ID:       "tool-malo",
		Name:     "Tool malo",
		Type:     coretask.TaskTypeScript,
		ToolID:   &missing,
		Priority: coretask.PriorityMedium,
	})
	if !errors.Is(err, tool.ErrNotFound) {
		t.Fatalf("esperaba tool.ErrNotFound, dio: %v", err)
	}
}
func TestManager_Create_RequiresToolForScriptType(t *testing.T) {
	mgr, _, _ := newManagerEnv(t)
	_, err := mgr.Create(context.Background(), uc.TaskInput{
		ID:       "sin-tool",
		Name:     "Sin Tool",
		Type:     coretask.TaskTypeScript,
		Priority: coretask.PriorityMedium,
	})
	if err == nil {
		t.Fatal("esperaba error por tool_id faltante")
	}
	// El error viene de la entidad, no del use case.
	// Solo verificamos que mencione tool_id.
	if !strings.Contains(err.Error(), "tool_id") {
		t.Errorf("error debería mencionar tool_id: %v", err)
	}
}

func TestManager_Create_AITask_NoTool(t *testing.T) {
	mgr, _, _ := newManagerEnv(t)
	tk, err := mgr.Create(context.Background(), uc.TaskInput{
		ID:         "chat-nueva",
		Name:       "Chat Nueva",
		Type:       coretask.TaskTypeAI,
		RequiresAI: true,
		AIPrompt:   "{{.Input}}",
		Priority:   coretask.PriorityLow,
		IsActive:   true,
	})
	if err != nil {
		t.Fatalf("Create AI: %v", err)
	}
	if tk.ToolID != nil {
		t.Errorf("ToolID debería ser nil")
	}
}

func TestManager_Update_OK(t *testing.T) {
	mgr, taskRepo, toolRepo := newManagerEnv(t)
	ctx := context.Background()
	tl := seedTool(t, toolRepo, "test-tool")

	// Crear
	_, err := mgr.Create(ctx, uc.TaskInput{
		ID:       "para-editar",
		Name:     "Original",
		Type:     coretask.TaskTypeScript,
		ToolID:   &tl.ID,
		Priority: coretask.PriorityLow,
		IsActive: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Editar
	_, err = mgr.Update(ctx, uc.TaskInput{
		ID:       "para-editar",
		Name:     "Editada",
		Type:     coretask.TaskTypeScript,
		ToolID:   &tl.ID,
		Priority: coretask.PriorityHigh,
		IsActive: true,
		Tags:     []string{"editada"},
	})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}

	got, _ := taskRepo.GetByID(ctx, "para-editar")
	if got.Name != "Editada" {
		t.Errorf("Name = %q", got.Name)
	}
	if got.Priority != coretask.PriorityHigh {
		t.Errorf("Priority = %q", got.Priority)
	}
}

func TestManager_Delete_SoftDeletes(t *testing.T) {
	mgr, taskRepo, toolRepo := newManagerEnv(t)
	ctx := context.Background()
	tl := seedTool(t, toolRepo, "test-tool")

	_, err := mgr.Create(ctx, uc.TaskInput{
		ID:       "para-borrar",
		Name:     "Para Borrar",
		Type:     coretask.TaskTypeScript,
		ToolID:   &tl.ID,
		Priority: coretask.PriorityMedium,
		IsActive: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := mgr.Delete(ctx, "para-borrar"); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	_, err = taskRepo.GetByID(ctx, "para-borrar")
	if !errors.Is(err, coretask.ErrNotFound) {
		t.Fatalf("esperaba ErrNotFound, dio: %v", err)
	}
}

func TestManager_Delete_NotFound(t *testing.T) {
	mgr, _, _ := newManagerEnv(t)
	err := mgr.Delete(context.Background(), "no-existe")
	if !errors.Is(err, coretask.ErrNotFound) {
		t.Fatalf("esperaba ErrNotFound, dio: %v", err)
	}
}
