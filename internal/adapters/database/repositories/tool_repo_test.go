package repositories_test

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/NeRo0128/brain-cli/internal/adapters/database"
	"github.com/NeRo0128/brain-cli/internal/adapters/database/repositories"
	"github.com/NeRo0128/brain-cli/internal/core/tool"
)

// newTestRepo crea una DB temporal, aplica migraciones,
// limpia el seed y devuelve un repo listo para tests aislados.
func newTestRepo(t *testing.T) *repositories.ToolRepository {
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

	// Quitamos el seed para no arrastrar los 6 tools builtin.
	if _, err := db.DB().Exec("DELETE FROM tasks"); err != nil {
		t.Fatalf("limpiando tasks: %v", err)
	}
	if _, err := db.DB().Exec("DELETE FROM tools"); err != nil {
		t.Fatalf("limpiando tools: %v", err)
	}

	return repositories.NewToolRepository(db.DB())
}

// sampleTool devuelve un Tool bash válido.
func sampleTool() *tool.Tool {
	return &tool.Tool{
		Name:           "test-tool",
		Description:    "tool de prueba",
		ScriptType:     tool.ScriptTypeBash,
		Category:       tool.CategoryUtils,
		ScriptContent:  "#!/bin/bash\necho hi",
		TimeoutSeconds: 30,
		Version:        1,
	}
}

func TestCreate_FillsIDAndCreatedAt(t *testing.T) {
	repo := newTestRepo(t)

	tl := sampleTool()
	if err := repo.Create(context.Background(), tl); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if tl.ID <= 0 {
		t.Errorf("ID no fue rellenado: %d", tl.ID)
	}
	if tl.CreatedAt.IsZero() {
		t.Error("CreatedAt no fue rellenado")
	}
}

func TestCreate_DuplicateName_ReturnsErrDuplicateName(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	if err := repo.Create(ctx, sampleTool()); err != nil {
		t.Fatal(err)
	}
	err := repo.Create(ctx, sampleTool())
	if !errors.Is(err, tool.ErrDuplicateName) {
		t.Fatalf("esperaba ErrDuplicateName, dio: %v", err)
	}
}

func TestGetByID_Roundtrip(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	orig := sampleTool()
	orig.RequiresSudo = true
	orig.TimeoutSeconds = 42
	if err := repo.Create(ctx, orig); err != nil {
		t.Fatal(err)
	}

	got, err := repo.GetByID(ctx, orig.ID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if got.Name != orig.Name {
		t.Errorf("Name = %q", got.Name)
	}
	if !got.RequiresSudo {
		t.Error("RequiresSudo no se persistió")
	}
	if got.TimeoutSeconds != 42 {
		t.Errorf("TimeoutSeconds = %d", got.TimeoutSeconds)
	}
	if got.ScriptType != tool.ScriptTypeBash {
		t.Errorf("ScriptType = %q", got.ScriptType)
	}
	if got.DeletedAt != nil {
		t.Error("DeletedAt debería ser nil")
	}
}

func TestGetByID_NotFound(t *testing.T) {
	repo := newTestRepo(t)
	_, err := repo.GetByID(context.Background(), 9999)
	if !errors.Is(err, tool.ErrNotFound) {
		t.Fatalf("esperaba ErrNotFound, dio: %v", err)
	}
}

func TestGetByName(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	if err := repo.Create(ctx, sampleTool()); err != nil {
		t.Fatal(err)
	}
	got, err := repo.GetByName(ctx, "test-tool")
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "test-tool" {
		t.Errorf("Name = %q", got.Name)
	}
}

func TestList_ExcludesDeleted(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	for _, name := range []string{"a", "b", "c"} {
		tl := sampleTool()
		tl.Name = name
		if err := repo.Create(ctx, tl); err != nil {
			t.Fatal(err)
		}
	}

	// Soft-delete "b"
	b, _ := repo.GetByName(ctx, "b")
	if err := repo.Delete(ctx, b.ID); err != nil {
		t.Fatal(err)
	}

	got, err := repo.List(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Errorf("List = %d tools, quiero 2", len(got))
	}
	for _, tl := range got {
		if tl.Name == "b" {
			t.Error("tool borrado apareció en List")
		}
	}
}

func TestDelete_SoftDelete_HidesFromGetByID(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	tl := sampleTool()
	if err := repo.Create(ctx, tl); err != nil {
		t.Fatal(err)
	}
	if err := repo.Delete(ctx, tl.ID); err != nil {
		t.Fatal(err)
	}

	_, err := repo.GetByID(ctx, tl.ID)
	if !errors.Is(err, tool.ErrNotFound) {
		t.Fatalf("esperaba ErrNotFound tras Delete, dio: %v", err)
	}

	// Segunda Delete sobre el mismo → ErrNotFound
	if err := repo.Delete(ctx, tl.ID); !errors.Is(err, tool.ErrNotFound) {
		t.Fatalf("segunda Delete: esperaba ErrNotFound, dio: %v", err)
	}
}

func TestUpdate_ChangesFieldsAndSetsUpdatedAt(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	tl := sampleTool()
	if err := repo.Create(ctx, tl); err != nil {
		t.Fatal(err)
	}
	if tl.UpdatedAt != nil {
		t.Fatal("UpdatedAt debería ser nil tras Create")
	}

	tl.Description = "actualizado"
	tl.TimeoutSeconds = 99
	if err := repo.Update(ctx, tl); err != nil {
		t.Fatalf("Update: %v", err)
	}
	if tl.UpdatedAt == nil {
		t.Fatal("UpdatedAt no fue rellenado")
	}

	got, _ := repo.GetByID(ctx, tl.ID)
	if got.Description != "actualizado" {
		t.Errorf("Description = %q", got.Description)
	}
	if got.TimeoutSeconds != 99 {
		t.Errorf("TimeoutSeconds = %d", got.TimeoutSeconds)
	}
}

func TestUpdate_NotFound(t *testing.T) {
	repo := newTestRepo(t)
	tl := sampleTool()
	tl.ID = 9999
	err := repo.Update(context.Background(), tl)
	if !errors.Is(err, tool.ErrNotFound) {
		t.Fatalf("esperaba ErrNotFound, dio: %v", err)
	}
}

func TestListByCategory_Filters(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	dev := sampleTool()
	dev.Name, dev.Category = "dev-tool", tool.CategoryDev
	net := sampleTool()
	net.Name, net.Category = "net-tool", tool.CategoryNetwork
	util := sampleTool()
	util.Name, util.Category = "util-tool", tool.CategoryUtils

	for _, tl := range []*tool.Tool{dev, net, util} {
		if err := repo.Create(ctx, tl); err != nil {
			t.Fatal(err)
		}
	}

	got, err := repo.ListByCategory(ctx, tool.CategoryDev)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Name != "dev-tool" {
		t.Errorf("ListByCategory dev = %+v", got)
	}
}

func TestListBuiltin_Filters(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	builtin := sampleTool()
	builtin.Name, builtin.IsBuiltin = "builtin-tool", true
	custom := sampleTool()
	custom.Name, custom.IsBuiltin = "custom-tool", false

	for _, tl := range []*tool.Tool{builtin, custom} {
		if err := repo.Create(ctx, tl); err != nil {
			t.Fatal(err)
		}
	}

	got, err := repo.ListBuiltin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Name != "builtin-tool" {
		t.Errorf("ListBuiltin = %+v", got)
	}
}
