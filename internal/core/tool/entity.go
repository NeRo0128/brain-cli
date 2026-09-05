package tool

import "time"

// Tool representa una herramienta/script ejecutable
type Tool struct {
	ID             int       `db:"id"`
	Name           string    `db:"name"`
	Description    string    `db:"description"`
	ScriptContent  string    `db:"script_content"`
	ScriptType     ToolType  `db:"script_type"`
	Category       Category  `db:"category"`
	RequiresSudo   bool      `db:"requires_sudo"`
	TimeoutSeconds int       `db:"timeout_seconds"`
	IsBuiltin      bool      `db:"is_builtin"`
	Version        int       `db:"version"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      *time.Time `db:"updated_at"`
}

// ToolType define los tipos de scripts soportados
type ToolType string

const (
	ToolTypeBash   ToolType = "bash"
	ToolTypePython ToolType = "python"
	ToolTypeNative ToolType = "native"
	ToolTypeGo     ToolType = "go"
)

// Category define las categorías de herramientas
type Category string

const (
	CategorySystem      Category = "system"
	CategoryDev         Category = "dev"
	CategoryAI          Category = "ai"
	CategoryUtils       Category = "utils"
	CategoryCustom      Category = "custom"
	CategoryNetwork     Category = "network"
	CategoryMaintenance Category = "maintenance"
)

// TODO: Implementar métodos de validación
// func (t *Tool) Validate() error
// func (t *Tool) GetExecutor() (executor.Executor, error)
// func (t *Tool) Execute(ctx context.Context, params map[string]string) (string, error)
