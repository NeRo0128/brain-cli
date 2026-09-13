package tool_test

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/NeRo0128/brain-cli/internal/adapters/database"
	"github.com/NeRo0128/brain-cli/internal/adapters/database/repositories"
	coretool "github.com/NeRo0128/brain-cli/internal/core/tool"
	uc "github.com/NeRo0128/brain-cli/internal/usecases/tool"
)

func newManagerEnv(t *testing.T) (*uc.Manager, *repositories.ToolRepository) {
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

	// Limpiar el seed de tools para tests aislados.
	for _, tbl := range []string{"executions", "task_tags", "tasks", "tools", "tags"} {
		if _, err := db.DB().Exec("DELETE FROM " + tbl); err != nil {
			t.Fatalf("limpiando %s: %v", tbl, err)
		}
	}

	repo := repositories.NewToolRepository(db.DB())
	return uc.NewManager(repo), repo
}

// --- Create ---

func TestManager_Create_BashOK(t *testing.T) {
	mgr, repo := newManagerEnv(t)
	ctx := context.Background()

	tl, err := mgr.Create(ctx, uc.ToolInput{
		Name:           "test-bash",
		Description:    "un tool de prueba",
		ScriptType:     coretool.ScriptTypeBash,
		Category:       coretool.CategoryUtils,
		ScriptContent:  "#!/bin/bash\necho hi",
		TimeoutSeconds: 30,
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if tl.ID <= 0 {
		t.Error("ID no fue rellenado")
	}

	fromDB, err := repo.GetByID(ctx, tl.ID)
	if err != nil {
		t.Fatalf("no persistido: %v", err)
	}
	if fromDB.Name != "test-bash" {
		t.Errorf("Name = %q", fromDB.Name)
	}
}

func TestManager_Create_NativeOK(t *testing.T) {
	mgr, _ := newManagerEnv(t)
	ctx := context.Background()

	tl, err := mgr.Create(ctx, uc.ToolInput{
		Name:       "git-status",
		ScriptType: coretool.ScriptTypeNative,
		Command:    "git status --short",
	})
	if err != nil {
		t.Fatalf("Create native: %v", err)
	}
	if tl.Command != "git status --short" {
		t.Errorf("Command = %q", tl.Command)
	}
}

func TestManager_Create_RequiresName(t *testing.T) {
	mgr, _ := newManagerEnv(t)
	_, err := mgr.Create(context.Background(), uc.ToolInput{
		ScriptType: coretool.ScriptTypeNative,
		Command:    "echo hi",
		Name:       "   ", // solo espacios
	})
	if !errors.Is(err, uc.ErrNameRequired) {
		t.Fatalf("esperaba ErrNameRequired, dio: %v", err)
	}
}

func TestManager_Create_DefaultsApplied(t *testing.T) {
	mgr, _ := newManagerEnv(t)
	ctx := context.Background()

	tl, err := mgr.Create(ctx, uc.ToolInput{
		Name:       "no-defaults",
		ScriptType: coretool.ScriptTypeNative,
		Command:    "echo hi",
		// sin Category, sin Timeout, sin Version
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if tl.Category != coretool.CategoryCustom {
		t.Errorf("Category default = %q", tl.Category)
	}
	if tl.TimeoutSeconds != 300 {
		t.Errorf("Timeout default = %d", tl.TimeoutSeconds)
	}
	if tl.Version != 1 {
		t.Errorf("Version default = %d", tl.Version)
	}
}

func TestManager_Create_RejectsInvalidCombo(t *testing.T) {
	mgr, _ := newManagerEnv(t)
	// native + script_content es inválido por Validate()
	_, err := mgr.Create(context.Background(), uc.ToolInput{
		Name:          "invalid",
		ScriptType:    coretool.ScriptTypeNative,
		ScriptContent: "#!/bin/bash",
		Command:       "echo hi",
	})
	if err == nil {
		t.Fatal("esperaba error por combo inválido")
	}
	if !strings.Contains(err.Error(), "no admite") {
		t.Errorf("error no menciona conflicto: %v", err)
	}
}

// --- Update ---

func TestManager_Update_OK(t *testing.T) {
	mgr, repo := newManagerEnv(t)
	ctx := context.Background()

	created, err := mgr.Create(ctx, uc.ToolInput{
		Name:       "before",
		ScriptType: coretool.ScriptTypeNative,
		Command:    "echo 1",
	})
	if err != nil {
		t.Fatal(err)
	}

	updated, err := mgr.Update(ctx, uc.ToolInput{
		ID:         created.ID,
		Name:       "after",
		ScriptType: coretool.ScriptTypeNative,
		Command:    "echo 2",
	})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if updated.Name != "after" {
		t.Errorf("Name = %q", updated.Name)
	}

	fromDB, _ := repo.GetByID(ctx, created.ID)
	if fromDB.Name != "after" {
		t.Errorf("Name en DB = %q", fromDB.Name)
	}
}

func TestManager_Update_RequiresID(t *testing.T) {
	mgr, _ := newManagerEnv(t)
	_, err := mgr.Update(context.Background(), uc.ToolInput{
		ID:         0,
		Name:       "x",
		ScriptType: coretool.ScriptTypeNative,
		Command:    "echo",
	})
	if !errors.Is(err, uc.ErrIDRequired) {
		t.Fatalf("esperaba ErrIDRequired, dio: %v", err)
	}
}

// --- Delete ---

func TestManager_Delete_SoftDelete(t *testing.T) {
	mgr, repo := newManagerEnv(t)
	ctx := context.Background()

	created, _ := mgr.Create(ctx, uc.ToolInput{
		Name:       "to-delete",
		ScriptType: coretool.ScriptTypeNative,
		Command:    "echo",
	})

	if err := mgr.Delete(ctx, created.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	_, err := repo.GetByID(ctx, created.ID)
	if !errors.Is(err, coretool.ErrNotFound) {
		t.Fatalf("esperaba ErrNotFound tras Delete, dio: %v", err)
	}
}

func TestManager_Delete_RequiresID(t *testing.T) {
	mgr, _ := newManagerEnv(t)
	if err := mgr.Delete(context.Background(), 0); !errors.Is(err, uc.ErrIDRequired) {
		t.Fatalf("esperaba ErrIDRequired, dio: %v", err)
	}
}
