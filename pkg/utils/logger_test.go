package utils_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/NeRo0128/brain-cli/pkg/utils"
)

func TestNewLogger_RejectsInvalidLevel(t *testing.T) {
	_, err := utils.NewLogger(utils.LoggerOptions{Level: "verbose"})
	if err == nil {
		t.Fatal("esperaba error por nivel inválido")
	}
	if !strings.Contains(err.Error(), "verbose") {
		t.Errorf("error no menciona el nivel: %v", err)
	}
}

func TestNewLogger_WritesJSONToFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "brain.log")

	log, err := utils.NewLogger(utils.LoggerOptions{
		Level:      "info",
		Format:     "json",
		OutputPath: path,
	})
	if err != nil {
		t.Fatal(err)
	}

	log.Info().Str("foo", "bar").Msg("hola")
	log.Debug().Msg("no debe aparecer") // nivel info lo descarta

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)

	if !strings.Contains(content, `"foo":"bar"`) {
		t.Errorf("log no escrito: %s", content)
	}
	if !strings.Contains(content, `"message":"hola"`) {
		t.Errorf("mensaje ausente: %s", content)
	}
	if strings.Contains(content, "no debe aparecer") {
		t.Errorf("debug no debería estar en nivel info: %s", content)
	}
}

func TestNewLogger_CreatesParentDirectory(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "deep", "brain.log")

	_, err := utils.NewLogger(utils.LoggerOptions{
		Level:      "info",
		Format:     "json",
		OutputPath: path,
	})
	if err != nil {
		t.Fatalf("esperaba que creara el directorio: %v", err)
	}
	if _, err := os.Stat(filepath.Dir(path)); err != nil {
		t.Errorf("directorio no creado: %v", err)
	}
}
