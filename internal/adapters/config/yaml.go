package config

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"

	core "github.com/NeRo0128/brain-cli/internal/core/config"
)

// YAMLLoader implementa core.Manager leyendo configuración desde YAML.
type YAMLLoader struct{}

// NewYAMLLoader construye un loader de configuración YAML.
func NewYAMLLoader() *YAMLLoader {
	return &YAMLLoader{}
}

// Compile-time check: YAMLLoader cumple core.Manager.
var _ core.Manager = (*YAMLLoader)(nil)

// Load lee path, expande ${VAR} y ${VAR:-default}, parsea YAML
// sobre los defaults y valida el resultado.
func (l *YAMLLoader) Load(path string) (*core.Config, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("leyendo config %q: %w", path, err)
	}

	expanded := expandEnv(string(raw))

	cfg := core.Default()
	if err := yaml.Unmarshal([]byte(expanded), cfg); err != nil {
		return nil, fmt.Errorf("parseando YAML %q: %w", path, err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config inválida: %w", err)
	}

	return cfg, nil
}

// envPattern matchea ${VAR} o ${VAR:-default}.
// Grupo 1 = nombre, grupo 2 = default (puede ser vacío).
var envPattern = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)(?::-([^}]*))?\}`)

// expandEnv reemplaza placeholders por el valor del entorno.
//
//	${VAR}         → valor de VAR ("" si no está definida)
//	${VAR:-value}  → valor de VAR, o "value" si no está o está vacía
//
// Sigue la semántica POSIX de `:-` (unset o empty → default).
func expandEnv(input string) string {
	return envPattern.ReplaceAllStringFunc(input, func(match string) string {
		groups := envPattern.FindStringSubmatch(match)
		name := groups[1]
		def := groups[2]
		hasDefault := strings.Contains(match, ":-")

		v, set := os.LookupEnv(name)
		if set && (v != "" || !hasDefault) {
			return v
		}
		return def
	})
}
