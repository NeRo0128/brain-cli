package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	core "github.com/NeRo0128/brain-cli/internal/core/config"
)

// YAMLWriter persiste la configuración a un archivo YAML.
type YAMLWriter struct{}

func NewYAMLWriter() *YAMLWriter {
	return &YAMLWriter{}
}

var _ core.Writer = (*YAMLWriter)(nil)

// Write valida, hace backup (una sola vez) y escribe atómicamente.
//
// El backup se crea como <path>.bak solo la primera vez que se
// sobrescribe un archivo existente. Backups posteriores no lo pisan,
// así preservamos la config original del usuario.
func (w *YAMLWriter) Write(path string, cfg *core.Config) error {
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("config inválida: %w", err)
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creando directorio %q: %w", dir, err)
	}

	if err := backupOnce(path); err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("serializando config: %w", err)
	}

	return atomicWrite(path, data)
}

// backupOnce copia path a path.bak si aún no existe un backup.
// Silencioso si el archivo original no existe o si ya hay backup.
func backupOnce(path string) error {
	if _, err := os.Stat(path + ".bak"); err == nil {
		return nil
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	src, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("abriendo original para backup: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(path + ".bak")
	if err != nil {
		return fmt.Errorf("creando backup: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("copiando a backup: %w", err)
	}
	return nil
}

// atomicWrite escribe data a path de forma atómica:
// escribe a path.tmp, sincroniza y renombra sobre path.
func atomicWrite(path string, data []byte) error {
	tmp := path + ".tmp"

	f, err := os.Create(tmp)
	if err != nil {
		return fmt.Errorf("creando temporal %q: %w", tmp, err)
	}

	if _, err := f.Write(data); err != nil {
		f.Close()
		os.Remove(tmp)
		return fmt.Errorf("escribiendo temporal: %w", err)
	}
	if err := f.Sync(); err != nil {
		f.Close()
		os.Remove(tmp)
		return fmt.Errorf("sincronizando temporal: %w", err)
	}
	if err := f.Close(); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("cerrando temporal: %w", err)
	}

	if err := os.Rename(tmp, path); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("renombrando temporal: %w", err)
	}
	return nil
}
