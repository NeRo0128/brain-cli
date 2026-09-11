package tool

import (
	"errors"
	"fmt"
	"time"
)

// Tool representa una herramienta/script ejecutable
type Tool struct {
	ID          int        `db:"id"          json:"id"`
	Name        string     `db:"name"        json:"name"        yaml:"name"`
	Description string     `db:"description" json:"description" yaml:"description"`
	ScriptType  ScriptType `db:"script_type" json:"script_type" yaml:"script_type"`
	Category    Category   `db:"category"    json:"category"    yaml:"category"`

	// Source: el script/comando. Debe estar presente EXACTAMENTE uno de:
	// ScriptContent (código embebido en DB) o ScriptPath (ruta a archivo externo).
	ScriptContent string `db:"script_content" json:"script_content,omitempty" yaml:"script_content,omitempty"`
	ScriptPath    string `db:"script_path"    json:"script_path,omitempty"    yaml:"script_path,omitempty"`

	// Command: solo aplica a ScriptTypeNative.
	Command string `db:"command" json:"command,omitempty" yaml:"command,omitempty"`

	TimeoutSeconds int        `db:"timeout_seconds" json:"timeout_seconds" yaml:"timeout_seconds"`
	IsBuiltin      bool       `db:"is_builtin"      json:"is_builtin"      yaml:"is_builtin"`
	Version        int        `db:"version"         json:"version"         yaml:"version"`
	CreatedAt      time.Time  `db:"created_at"      json:"created_at"`
	UpdatedAt      *time.Time `db:"updated_at"      json:"updated_at,omitempty"`
	DeletedAt      *time.Time `db:"deleted_at"      json:"deleted_at,omitempty"`
}

// ToolType define los tipos de scripts soportados
type ScriptType string

const (
	ScriptTypeBash   ScriptType = "bash"
	ScriptTypePython ScriptType = "python"
	ScriptTypeNative ScriptType = "native"
	ScriptTypeGo     ScriptType = "go"
)

// Category define las categorías de herramientas
func (s ScriptType) IsValid() bool {
	switch s {
	case ScriptTypeBash, ScriptTypePython, ScriptTypeNative, ScriptTypeGo:
		return true
	}
	return false
}

// Category agrupa herramientas por área funcional.
type Category string

const (
	CategorySystem      Category = "system"
	CategoryDev         Category = "dev"
	CategoryAI          Category = "ai"
	CategoryUtils       Category = "utils"
	CategoryNetwork     Category = "network"
	CategoryMaintenance Category = "maintenance"
	CategoryCustom      Category = "custom"
)

func (c Category) IsValid() bool {
	switch c {
	case CategorySystem, CategoryDev, CategoryAI, CategoryUtils,
		CategoryNetwork, CategoryMaintenance, CategoryCustom:
		return true
	}
	return false
}

// Validate comprueba la coherencia del Tool.
func (t *Tool) Validate() error {
	var errs []error

	if t.Name == "" {
		errs = append(errs, errors.New("tool.name no puede estar vacío"))
	}
	if !t.ScriptType.IsValid() {
		errs = append(errs, fmt.Errorf("tool.script_type inválido: %q", t.ScriptType))
	}
	if t.Category != "" && !t.Category.IsValid() {
		errs = append(errs, fmt.Errorf("tool.category inválido: %q", t.Category))
	}
	if t.TimeoutSeconds < 0 {
		errs = append(errs, errors.New("tool.timeout_seconds no puede ser negativo"))
	}

	// Regla: exactamente una fuente según el tipo.
	switch t.ScriptType {
	case ScriptTypeBash, ScriptTypePython, ScriptTypeGo:
		if t.ScriptContent == "" && t.ScriptPath == "" {
			errs = append(errs, fmt.Errorf("script_type=%q requiere script_content o script_path", t.ScriptType))
		}
		if t.ScriptContent != "" && t.ScriptPath != "" {
			errs = append(errs, fmt.Errorf("script_type=%q no admite script_content y script_path a la vez", t.ScriptType))
		}
		if t.Command != "" {
			errs = append(errs, fmt.Errorf("script_type=%q no admite command", t.ScriptType))
		}
	case ScriptTypeNative:
		if t.Command == "" {
			errs = append(errs, errors.New("script_type=native requiere command"))
		}
		if t.ScriptContent != "" || t.ScriptPath != "" {
			errs = append(errs, errors.New("script_type=native no admite script_content ni script_path"))
		}
	}

	return errors.Join(errs...)
}
