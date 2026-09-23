package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// defaultConfigTemplate es el YAML inicial.
// %s se reemplaza por la ruta por defecto de la DB.
const defaultConfigTemplate = `# Brain CLI — config generado automáticamente.
#
# Edita este archivo a mano si prefieres; la app respeta tus cambios.
# Para resetear a defaults: borra el archivo y reinicia la app.
#
# Variables de entorno:
#   ${VAR}         → valor de VAR, o vacío
#   ${VAR:-value}  → valor de VAR, o "value" si no está seteada

app:
  name: brain-cli

database:
  path: ${BRAIN_DB_PATH:-%s}
  auto_migrate: true

logging:
  level: ${BRAIN_LOG_LEVEL:-info}
  format: ${BRAIN_LOG_FORMAT:-pretty}
  output_path: ""
  max_size_mb: 100

ui:
  theme: ${BRAIN_THEME:-brain}
  icons: ${BRAIN_ICONS:-nerd-b}
  brand_style: ${BRAIN_BRAND:-slim}

github:
  client_id: ${GITHUB_CLIENT_ID:-}
`

// EnsureConfig crea el archivo con defaults si no existe.
// defaultDBPath se inyecta en el template.
func EnsureConfig(path, defaultDBPath string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("stat %q: %w", path, err)
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creando directorio %q: %w", dir, err)
	}

	content := fmt.Sprintf(defaultConfigTemplate, defaultDBPath)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("escribiendo %q: %w", path, err)
	}
	return nil
}
