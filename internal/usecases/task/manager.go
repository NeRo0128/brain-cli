package task

import (
	"context"
	"errors"
	"fmt"

	"github.com/NeRo0128/brain-cli/internal/core/task"
	"github.com/NeRo0128/brain-cli/internal/core/tool"
)

var (
	ErrTaskIDRequired   = errors.New("task.id es obligatorio")
	ErrTaskNameRequired = errors.New("task.name es obligatorio")
)

// Manager gestiona el ciclo CRUD de Tasks.
type Manager struct {
	tasks task.Repository
	tools tool.Repository
}

// NewManager construye el gestor CRUD.
func NewManager(tasks task.Repository, tools tool.Repository) *Manager {
	return &Manager{tasks: tasks, tools: tools}
}

// TaskInput son los campos editables de una Task.
type TaskInput struct {
	ID          string
	Name        string
	Description string
	Type        task.TaskType
	ToolID      *int
	Params      map[string]string
	RequiresAI  bool
	AIPrompt    string
	Tags        []string
	Priority    task.Priority
	IsActive    bool
	IsFavorite  bool
}

// Create valida y persiste una nueva Task.
func (m *Manager) Create(ctx context.Context, in TaskInput) (*task.Task, error) {
	tk, err := m.validate(ctx, in)
	if err != nil {
		return nil, err
	}
	if err := m.tasks.Create(ctx, tk); err != nil {
		return nil, fmt.Errorf("creando task: %w", err)
	}
	return tk, nil
}

// Update valida y persiste los cambios de una Task existente.
func (m *Manager) Update(ctx context.Context, in TaskInput) (*task.Task, error) {
	tk, err := m.validate(ctx, in)
	if err != nil {
		return nil, err
	}
	if err := m.tasks.Update(ctx, tk); err != nil {
		return nil, fmt.Errorf("actualizando task: %w", err)
	}
	return tk, nil
}

// Delete hace soft-delete de una Task por ID.
func (m *Manager) Delete(ctx context.Context, id string) error {
	if id == "" {
		return ErrTaskIDRequired
	}
	if err := m.tasks.Delete(ctx, id); err != nil {
		return fmt.Errorf("borrando task %q: %w", id, err)
	}
	return nil
}

// validate construye la Task, ejecuta Validate() y verifica la FK al Tool.
func (m *Manager) validate(ctx context.Context, in TaskInput) (*task.Task, error) {
	if in.ID == "" {
		return nil, ErrTaskIDRequired
	}
	if in.Name == "" {
		return nil, ErrTaskNameRequired
	}

	tk := &task.Task{
		ID:          in.ID,
		Name:        in.Name,
		Description: in.Description,
		Type:        in.Type,
		ToolID:      in.ToolID,
		Params:      in.Params,
		RequiresAI:  in.RequiresAI,
		AIPrompt:    in.AIPrompt,
		Tags:        in.Tags,
		Priority:    in.Priority,
		IsActive:    in.IsActive,
		IsFavorite:  in.IsFavorite,
	}

	if err := tk.Validate(); err != nil {
		return nil, fmt.Errorf("task inválida: %w", err)
	}

	if tk.Type.RequiresTool() {

		if _, err := m.tools.GetByID(ctx, *tk.ToolID); err != nil {
			if errors.Is(err, tool.ErrNotFound) {
				return nil, fmt.Errorf("%w: tool_id=%d", tool.ErrNotFound, *tk.ToolID)
			}
			return nil, fmt.Errorf("verificando tool: %w", err)
		}
	}

	return tk, nil
}
