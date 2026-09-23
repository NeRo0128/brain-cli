package repositories_test

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/NeRo0128/brain-cli/internal/adapters/database"
	"github.com/NeRo0128/brain-cli/internal/adapters/database/repositories"
	"github.com/NeRo0128/brain-cli/internal/core/settings"
)

func newTestSettingsRepo(t *testing.T) *repositories.SettingsRepository {
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

	return repositories.NewSettingsRepository(db.DB())
}

func TestSet_InsertAndGet(t *testing.T) {
	repo := newTestSettingsRepo(t)
	ctx := context.Background()

	if err := repo.Set(ctx, "ui.theme", "nord"); err != nil {
		t.Fatalf("Set: %v", err)
	}

	got, err := repo.Get(ctx, "ui.theme")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Value != "nord" {
		t.Errorf("Value = %q", got.Value)
	}
	if got.Category != "ui" {
		t.Errorf("Category = %q, quiero 'ui'", got.Category)
	}
}

func TestSet_UpsertOverwrites(t *testing.T) {
	repo := newTestSettingsRepo(t)
	ctx := context.Background()

	if err := repo.Set(ctx, "ui.theme", "nord"); err != nil {
		t.Fatal(err)
	}
	if err := repo.Set(ctx, "ui.theme", "catppuccin"); err != nil {
		t.Fatal(err)
	}

	got, _ := repo.Get(ctx, "ui.theme")
	if got.Value != "catppuccin" {
		t.Errorf("Value = %q", got.Value)
	}
}

func TestGet_NotFound(t *testing.T) {
	repo := newTestSettingsRepo(t)
	_, err := repo.Get(context.Background(), "no.existe")
	if !errors.Is(err, settings.ErrNotFound) {
		t.Fatalf("esperaba ErrNotFound, dio: %v", err)
	}
}

func TestDelete_Idempotent(t *testing.T) {
	repo := newTestSettingsRepo(t)
	ctx := context.Background()

	// Borrar algo que no existe no falla.
	if err := repo.Delete(ctx, "no.existe"); err != nil {
		t.Fatalf("Delete sin existir: %v", err)
	}

	_ = repo.Set(ctx, "ui.theme", "nord")
	if err := repo.Delete(ctx, "ui.theme"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := repo.Get(ctx, "ui.theme"); !errors.Is(err, settings.ErrNotFound) {
		t.Errorf("esperaba ErrNotFound tras delete")
	}
}

func TestSetMany_Atomic(t *testing.T) {
	repo := newTestSettingsRepo(t)
	ctx := context.Background()

	values := map[string]string{
		"ui.theme":       "nord",
		"ui.icons":       "nerd-b",
		"ui.brand_style": "slim",
		"logging.level":  "debug",
	}
	if err := repo.SetMany(ctx, values); err != nil {
		t.Fatalf("SetMany: %v", err)
	}

	all, err := repo.List(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 4 {
		t.Errorf("List = %d, quiero 4", len(all))
	}
}
