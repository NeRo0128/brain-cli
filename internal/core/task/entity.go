package task

import (
	"errors"
	"fmt"
	"regexp"
	"time"
)

// idPattern: slugs en minúsculas, alfanuméricos y guiones.
// Ej: "wifi-vpn", "start-dev", "backup-db".
var idPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*[a-z0-9]$`)

// Task es la entidad central del dominio: una unidad de trabajo ejecutable.
type Task struct {
	// ID es un slug legible y estable (ej: "wifi-vpn").
	// Debe cumplir idPattern.
	ID          string   `db:"id"          json:"id"          yaml:"id"`
	Name        string   `db:"name"        json:"name"        yaml:"name"`
	Description string   `db:"description" json:"description" yaml:"description,omitempty"`
	Type        TaskType `db:"type"      json:"type"        yaml:"type"`

	// ToolID referencia a tools.id. nil si Type=ai.
	ToolID *int `db:"tool_id" json:"tool_id,omitempty" yaml:"tool_id,omitempty"`

	// Params se inyectan al Tool. JSON en DB.
	Params map[string]string `db:"params" json:"params,omitempty" yaml:"params,omitempty"`

	RequiresAI bool   `db:"requires_ai" json:"requires_ai" yaml:"requires_ai"`
	AIPrompt   string `db:"ai_prompt"   json:"ai_prompt,omitempty" yaml:"ai_prompt,omitempty"`

	Tags     []string `db:"tags"     json:"tags,omitempty" yaml:"tags,omitempty"`
	Priority Priority `db:"priority" json:"priority"       yaml:"priority"`

	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at,omitempty"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

// Validate comprueba la coherencia de la Task.
func (t *Task) Validate() error {
	var errs []error

	// ID: slug válido y no vacío
	if t.ID == "" {
		errs = append(errs, errors.New("task.id no puede estar vacío"))
	} else if !idPattern.MatchString(t.ID) {
		errs = append(errs, fmt.Errorf("task.id inválido: %q (usa minúsculas, números y guiones)", t.ID))
	}

	if t.Name == "" {
		errs = append(errs, errors.New("task.name no puede estar vacío"))
	}
	if !t.Type.IsValid() {
		errs = append(errs, fmt.Errorf("task.type inválido: %q", t.Type))
	}
	if !t.Priority.IsValid() {
		errs = append(errs, fmt.Errorf("task.priority inválido: %q", t.Priority))
	}

	switch t.Type {
	case TaskTypeScript, TaskTypeCommand:
		if t.ToolID == nil || *t.ToolID <= 0 {
			errs = append(errs, fmt.Errorf("task.type=%q requiere tool_id válido", t.Type))
		}
		if t.AIPrompt != "" {
			errs = append(errs, fmt.Errorf("task.type=%q no admite ai_prompt", t.Type))
		}
	case TaskTypeAI:
		if t.ToolID != nil && *t.ToolID > 0 {
			errs = append(errs, errors.New("task.type=ai no admite tool_id"))
		}
		if t.AIPrompt == "" {
			errs = append(errs, errors.New("task.type=ai requiere ai_prompt"))
		}
	}

	if t.RequiresAI != (t.Type == TaskTypeAI) {
    errs = append(errs, errors.New("task.requires_ai debe coincidir con type=ai"))
	}

	seen := make(map[string]struct{}, len(t.Tags))
	for _, tag := range t.Tags {
		if tag == "" {
			errs = append(errs, errors.New("task.tags no admite valores vacíos"))
			continue
		}
		if _, dup := seen[tag]; dup {
			errs = append(errs, fmt.Errorf("task.tags duplicado: %q", tag))
		}
		seen[tag] = struct{}{}
	}

	return errors.Join(errs...)
}
