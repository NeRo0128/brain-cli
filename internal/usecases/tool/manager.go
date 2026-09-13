// Package tool provee los casos de uso de gestión de Tools:
// creación, edición y borrado, con validación previa a persistir.
package tool

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/NeRo0128/brain-cli/internal/core/tool"
)

var (
	ErrNameRequired = errors.New("tool.name es obligatorio")
	ErrIDRequired   = errors.New("tool.id inválido")
)

// Manager gestiona el ciclo CRUD de Tools.
//
// Es un thin wrapper sobre tool.Repository: valida, normaliza y delega.
// Mantiene la consistencia con task.Manager y permite añadir reglas
// de negocio sin tocar los adapters.
type Manager struct {
	tools tool.Repository
}

func NewManager(tools tool.Repository) *Manager {
	return &Manager{tools: tools}
}

// ToolInput son los campos editables de un Tool.
type ToolInput struct {
	ID             int
	Name           string
	Description    string
	ScriptType     tool.ScriptType
	Category       tool.Category
	ScriptContent  string
	ScriptPath     string
	Command        string
	RequiresSudo   bool
	TimeoutSeconds int
	IsBuiltin      bool
	Version        int
}

// Create valida y persiste un nuevo Tool.
func (m *Manager) Create(ctx context.Context, in ToolInput) (*tool.Tool, error) {
	tl, err := m.validate(in)
	if err != nil {
		return nil, err
	}
	// El ID es autoincremental: ignorar el del input.
	tl.ID = 0
	if err := m.tools.Create(ctx, tl); err != nil {
		return nil, fmt.Errorf("creando tool: %w", err)
	}
	return tl, nil
}

// Update valida y persiste los cambios de un Tool existente.
func (m *Manager) Update(ctx context.Context, in ToolInput) (*tool.Tool, error) {
	if in.ID <= 0 {
		return nil, ErrIDRequired
	}
	tl, err := m.validate(in)
	if err != nil {
		return nil, err
	}
	tl.ID = in.ID
	if err := m.tools.Update(ctx, tl); err != nil {
		return nil, fmt.Errorf("actualizando tool: %w", err)
	}
	return tl, nil
}

// Delete hace soft-delete de un Tool por ID.
func (m *Manager) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return ErrIDRequired
	}
	if err := m.tools.Delete(ctx, id); err != nil {
		return fmt.Errorf("borrando tool: %w", err)
	}
	return nil
}

// validate construye el Tool, normaliza campos y ejecuta Validate().
func (m *Manager) validate(in ToolInput) (*tool.Tool, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, ErrNameRequired
	}

	tl := &tool.Tool{
		ID:             in.ID,
		Name:           name,
		Description:    strings.TrimSpace(in.Description),
		ScriptType:     in.ScriptType,
		Category:       in.Category,
		ScriptContent:  in.ScriptContent, // no trim: scripts pueden tener espacios significativos
		ScriptPath:     strings.TrimSpace(in.ScriptPath),
		Command:        strings.TrimSpace(in.Command),
		RequiresSudo:   in.RequiresSudo,
		TimeoutSeconds: in.TimeoutSeconds,
		IsBuiltin:      in.IsBuiltin,
		Version:        in.Version,
	}

	// Defaults sensatos.
	if tl.Category == "" {
		tl.Category = tool.CategoryCustom
	}
	if tl.Version <= 0 {
		tl.Version = 1
	}
	if tl.TimeoutSeconds <= 0 {
		tl.TimeoutSeconds = 300
	}

	if err := tl.Validate(); err != nil {
		return nil, fmt.Errorf("tool inválida: %w", err)
	}
	return tl, nil
}
