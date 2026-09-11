package tool

import (
	"context"
	"errors"
)

var (
	ErrNotFound      = errors.New("tool no encontrada")
	ErrDuplicateName = errors.New("tool.name ya existe")
)

type Repository interface {
	Create(ctx context.Context, t *Tool) error                              // Create inserta el tool y RELLENA t.ID, t.CreatedAt.
	GetByID(ctx context.Context, id int) (*Tool, error)                     // GetByID devuelve ErrNotFound si no existe o está soft-deleted.
	GetByName(ctx context.Context, name string) (*Tool, error)              // GetByName devuelve ErrNotFound si no existe.
	List(ctx context.Context) ([]*Tool, error)                              // List no incluye soft-deleted.
	ListByCategory(ctx context.Context, category Category) ([]*Tool, error) // ListByCategory filtra por categoría. Sin soft-deleted.
	ListBuiltin(ctx context.Context) ([]*Tool, error)                       // ListBuiltin lista solo las builtin (IsBuiltin=true).
	Update(ctx context.Context, t *Tool) error                              // Update modifica un tool existente. Rellena t.UpdatedAt.
	Delete(ctx context.Context, id int) error                               // Delete hace soft-delete (UPDATE deleted_at = now).
}
