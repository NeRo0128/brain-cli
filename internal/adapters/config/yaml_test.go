package config_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/NeRo0128/brain-cli/internal/adapters/config"
)

func writeConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLoad_MergesDefaults(t *testing.T) {
	path := writeConfig(t, `
app:
  name: test-app
logging:
  level: debug
`)

	cfg, err := config.NewYAMLLoader().Load(path)
	if err != nil {
		t.Fatalf("Load falló: %v", err)
	}

	if cfg.App.Name != "test-app" {
		t.Errorf("App.Name = %q, quiero %q", cfg.App.Name, "test-app")
	}
	if cfg.Logging.Level != "debug" {
		t.Errorf("Logging.Level = %q, quiero %q", cfg.Logging.Level, "debug")
	}
	// El YAML no toca database.path → debe quedar el default
	if cfg.Database.Path != "data/brain.db" {
		t.Errorf("default Database.Path no aplicado: %q", cfg.Database.Path)
	}
}

func TestLoad_ExpandsEnvVar(t *testing.T) {
	t.Setenv("TEST_DB_PATH", "/tmp/custom.db")

	path := writeConfig(t, `
database:
  path: ${TEST_DB_PATH:-fallback.db}
`)

	cfg, err := config.NewYAMLLoader().Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Database.Path != "/tmp/custom.db" {
		t.Errorf("env no expandido: %q", cfg.Database.Path)
	}
}

func TestLoad_UsesDefaultWhenEnvMissing(t *testing.T) {
	os.Unsetenv("DEFINITELY_NOT_SET_VAR")

	path := writeConfig(t, `
database:
  path: ${DEFINITELY_NOT_SET_VAR:-fallback.db}
`)

	cfg, err := config.NewYAMLLoader().Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Database.Path != "fallback.db" {
		t.Errorf("default no aplicado: %q", cfg.Database.Path)
	}
}

func TestLoad_RejectsInvalidLogLevel(t *testing.T) {
	path := writeConfig(t, `
logging:
  level: verbose
`)

	_, err := config.NewYAMLLoader().Load(path)
	if err == nil {
		t.Fatal("esperaba error por level inválido")
	}
	if !strings.Contains(err.Error(), "logging.level") {
		t.Errorf("error no menciona el campo: %v", err)
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := config.NewYAMLLoader().Load("/no/existe/config.yaml")
	if err == nil {
		t.Fatal("esperaba error de archivo no encontrado")
	}
}
