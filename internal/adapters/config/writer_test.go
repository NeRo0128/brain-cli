package config_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/NeRo0128/brain-cli/internal/adapters/config"
	core "github.com/NeRo0128/brain-cli/internal/core/config"
)

func TestWrite_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	if err := config.NewYAMLWriter().Write(path, core.Default()); err != nil {
		t.Fatalf("Write: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("archivo no creado: %v", err)
	}
}

func TestWrite_BackupOnlyOnce(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	original := []byte("# original\napp:\n  name: original\n")
	if err := os.WriteFile(path, original, 0o644); err != nil {
		t.Fatal(err)
	}

	w := config.NewYAMLWriter()

	// Primera escritura → debe crear backup con el contenido original.
	if err := w.Write(path, core.Default()); err != nil {
		t.Fatal(err)
	}
	bak1, err := os.ReadFile(path + ".bak")
	if err != nil {
		t.Fatalf("backup no creado: %v", err)
	}
	if string(bak1) != string(original) {
		t.Errorf("backup no preserva el original:\n%s", bak1)
	}

	// Modificar el config y volver a escribir.
	cfg := core.Default()
	cfg.App.Name = "modificado"
	if err := w.Write(path, cfg); err != nil {
		t.Fatal(err)
	}

	// El backup NO debe cambiar (sigue siendo el original).
	bak2, _ := os.ReadFile(path + ".bak")
	if string(bak2) != string(original) {
		t.Errorf("backup fue sobreescrito:\n%s", bak2)
	}
}

func TestWrite_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	cfg := core.Default()
	cfg.App.Name = "my-app"
	cfg.UI.Theme = "nord"
	cfg.Logging.Level = "warn"

	if err := config.NewYAMLWriter().Write(path, cfg); err != nil {
		t.Fatalf("Write: %v", err)
	}

	loaded, err := config.NewYAMLLoader().Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if loaded.App.Name != "my-app" {
		t.Errorf("App.Name = %q", loaded.App.Name)
	}
	if loaded.UI.Theme != "nord" {
		t.Errorf("UI.Theme = %q", loaded.UI.Theme)
	}
	if loaded.Logging.Level != "warn" {
		t.Errorf("Logging.Level = %q", loaded.Logging.Level)
	}
}

func TestWrite_InvalidPath_ReturnsError(t *testing.T) {
	// /dev/null es un archivo, no un directorio → MkdirAll falla.
	path := "/dev/null/subdir/config.yaml"

	err := config.NewYAMLWriter().Write(path, core.Default())
	if err == nil {
		t.Fatal("esperaba error por path inválido")
	}
}

func TestWrite_RejectsInvalidConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	cfg := core.Default()
	cfg.App.Name = "" // inválido

	err := config.NewYAMLWriter().Write(path, cfg)
	if err == nil {
		t.Fatal("esperaba error por config inválida")
	}
	if !strings.Contains(err.Error(), "inválida") {
		t.Errorf("error no menciona 'inválida': %v", err)
	}

	// El archivo NO debe haberse creado.
	if _, statErr := os.Stat(path); statErr == nil {
		t.Error("archivo creado a pesar de config inválida")
	}
}

func TestWrite_AtomicOverwrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	w := config.NewYAMLWriter()

	// Primera escritura.
	cfgA := core.Default()
	cfgA.App.Name = "first"
	if err := w.Write(path, cfgA); err != nil {
		t.Fatal(err)
	}

	// Segunda escritura.
	cfgB := core.Default()
	cfgB.App.Name = "second"
	if err := w.Write(path, cfgB); err != nil {
		t.Fatal(err)
	}

	// Load debe devolver los valores de B.
	loaded, err := config.NewYAMLLoader().Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.App.Name != "second" {
		t.Errorf("App.Name = %q, esperaba 'second'", loaded.App.Name)
	}

	// No debe quedar ningún .tmp.
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".tmp") {
			t.Errorf("quedó un temporal: %s", e.Name())
		}
	}
}
