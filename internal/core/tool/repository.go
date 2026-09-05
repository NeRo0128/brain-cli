package tool

import "context"

// Repository define las operaciones de persistencia para herramientas
type Repository interface {
	// Create crea una nueva herramienta
	Create(ctx context.Context, tool *Tool) error
	
	// GetByID obtiene una herramienta por su ID
	GetByID(ctx context.Context, id int) (*Tool, error)
	
	// GetByName obtiene una herramienta por su nombre
	GetByName(ctx context.Context, name string) (*Tool, error)
	
	// List lista todas las herramientas
	List(ctx context.Context) ([]*Tool, error)
	
	// ListByCategory lista herramientas por categoría
	ListByCategory(ctx context.Context, category Category) ([]*Tool, error)
	
	// ListBuiltin lista solo herramientas builtin
	ListBuiltin(ctx context.Context) ([]*Tool, error)
	
	// Update actualiza una herramienta existente
	Update(ctx context.Context, tool *Tool) error
	
	// Delete elimina una herramienta
	Delete(ctx context.Context, id int) error
}

// TODO: Implementar SQLite repository en internal/adapters/database/repositories/tool_repo.go
