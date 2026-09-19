package task

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/NeRo0128/brain-cli/internal/core/task"
	"github.com/NeRo0128/brain-cli/internal/core/tool"
)

var (
	ErrTaskIDRequired   = errors.New("task.id es obligatorio")
	ErrTaskNameRequired = errors.New("task.name es obligatorio")
	ErrToolTypeMismatch = errors.New("script_type del tool no coincide con el tipo de task")
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
	Command     string
	ToolID      *int

	ScriptContent string
	ScriptType    tool.ScriptType

	Params     map[string]string
	RequiresAI bool
	AIPrompt   string
	Tags       []string
	Priority   task.Priority
	IsActive   bool
	IsFavorite bool
}

// Create valida y persiste una nueva Task.
func (m *Manager) Create(ctx context.Context, in TaskInput) (*task.Task, error) {
	resolved, err := m.resolveTool(ctx, in) // ← nuevo: auto-crea el tool
	if err != nil {
		return nil, err
	}
	tk, err := m.validate(ctx, resolved)
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
	resolved, err := m.resolveTool(ctx, in) // [NUEVO]
	if err != nil {
		return nil, err
	}
	tk, err := m.validate(ctx, resolved)
	if err != nil {
		return nil, err
	}
	if err := m.tasks.Update(ctx, tk); err != nil {
		return nil, fmt.Errorf("actualizando task: %w", err)
	}
	return tk, nil
}

func (m *Manager) Delete(ctx context.Context, id string) error {
	if id == "" {
		return ErrTaskIDRequired
	}

	// Borrar el tool inline asociado (si existe).
	inlineName := "inline-" + id
	if tl, err := m.tools.GetByName(ctx, inlineName); err == nil {
		_ = m.tools.Delete(ctx, tl.ID) // best-effort
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
		tl, err := m.tools.GetByID(ctx, *tk.ToolID)

		if err != nil {
			if errors.Is(err, tool.ErrNotFound) {
				return nil, fmt.Errorf("%w: tool_id=%d", tool.ErrNotFound, *tk.ToolID)
			}
			return nil, fmt.Errorf("verificando tool: %w", err)
		}
		switch tk.Type {
		case task.TaskTypeScript:
			if tl.ScriptType == tool.ScriptTypeNative {
				return nil, fmt.Errorf("%w: task.type=script no admite script_type=native",
					ErrToolTypeMismatch)
			}
		case task.TaskTypeCommand:
			if tl.ScriptType != tool.ScriptTypeNative {
				return nil, fmt.Errorf("%w: task.type=command requiere script_type=native (tool %q es %q)",
					ErrToolTypeMismatch, tl.Name, tl.ScriptType)
			}
		}
	}

	return tk, nil
}

// resolveTool decide qué tool asociar a la task según su tipo.
//
//   - script con ScriptContent: crea/actualiza un tool inline automático.
//   - command con Command: crea/reutiliza un tool native (ya existente).
//   - ToolID ya seteado: se respeta (modo avanzado, editando task existente).
func (m *Manager) resolveTool(ctx context.Context, in TaskInput) (TaskInput, error) {
	// Modo avanzado: ya tiene ToolID explícito, no auto-crear.
	if in.ToolID != nil && *in.ToolID > 0 {
		return in, nil
	}

	switch in.Type {
	case task.TaskTypeScript:
		return m.resolveInlineScript(ctx, in)

	case task.TaskTypeCommand:
		return m.resolveCommand(ctx, in)
	}

	return in, nil
}

// resolveInlineScript crea o actualiza el tool asociado a un script inline.
// El nombre del tool es determinista: "inline-{task-id}", así al editar
// la misma task se actualiza el mismo tool.
func (m *Manager) resolveInlineScript(ctx context.Context, in TaskInput) (TaskInput, error) {
	if in.ScriptContent == "" {
		// Sin contenido: la validación fallará después con mensaje claro.
		return in, nil
	}

	name := "inline-" + in.ID
	st := in.ScriptType
	if st == "" {
		st = tool.ScriptTypeBash // default
	}

	// ¿Existe ya un tool inline para esta task?
	existing, err := m.tools.GetByName(ctx, name)
	if err == nil {
		// Actualizar el existente.
		existing.ScriptContent = in.ScriptContent
		existing.ScriptType = st
		if err := m.tools.Update(ctx, existing); err != nil {
			return in, fmt.Errorf("actualizando tool inline: %w", err)
		}
		in.ToolID = &existing.ID
		return in, nil
	}
	if !errors.Is(err, tool.ErrNotFound) {
		return in, fmt.Errorf("buscando tool inline: %w", err)
	}

	// No existe: crear.
	tl := &tool.Tool{
		Name:           name,
		Description:    "Script inline de la task " + in.ID,
		ScriptType:     st,
		Category:       tool.CategoryCustom,
		ScriptContent:  in.ScriptContent,
		TimeoutSeconds: 300,
		IsBuiltin:      false,
		Version:        1,
	}
	if err := tl.Validate(); err != nil {
		return in, fmt.Errorf("validando tool inline: %w", err)
	}
	if err := m.tools.Create(ctx, tl); err != nil {
		return in, fmt.Errorf("creando tool inline: %w", err)
	}
	in.ToolID = &tl.ID
	return in, nil
}

// resolveCommand extrae la lógica vieja del command a su propia función.
func (m *Manager) resolveCommand(ctx context.Context, in TaskInput) (TaskInput, error) {
	cmd := strings.TrimSpace(in.Command)
	if cmd == "" {
		return in, nil
	}
	tl, err := m.findOrCreateCommandTool(ctx, cmd)
	if err != nil {
		return in, err
	}
	in.ToolID = &tl.ID
	return in, nil
}
func (m *Manager) findOrCreateCommandTool(ctx context.Context, command string) (*tool.Tool, error) {
	existing, err := m.tools.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("listando tools: %w", err)
	}
	for _, t := range existing {
		if t.ScriptType == tool.ScriptTypeNative && t.Command == command {
			return t, nil
		}
	}

	// Nombre derivado del comando (slug estable)
	name := "auto-" + slugify(command)
	tl := &tool.Tool{
		Name:           name,
		Description:    "Auto-generado para: " + command,
		ScriptType:     tool.ScriptTypeNative,
		Category:       tool.CategoryCustom,
		Command:        command,
		TimeoutSeconds: 300,
		IsBuiltin:      false,
		Version:        1,
	}
	if err := tl.Validate(); err != nil {
		return nil, fmt.Errorf("validando tool auto-generado: %w", err)
	}
	if err := m.tools.Create(ctx, tl); err != nil {
		if errors.Is(err, tool.ErrDuplicateName) {
			// Race: otro proceso lo creó entre el List y el Create.
			// Reintentamos la búsqueda.
			if found, err := m.tools.GetByName(ctx, name); err == nil {
				return found, nil
			}
		}
		return nil, fmt.Errorf("creando tool auto-generado: %w", err)
	}
	return tl, nil
}
func slugify(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(s) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == ' ' || r == '-' || r == '_' || r == '.':
			b.WriteRune('-')
		}
	}
	out := strings.Trim(b.String(), "-")
	if len(out) > 40 {
		out = out[:40]
	}
	if out == "" {
		out = "cmd"
	}
	return out
}
