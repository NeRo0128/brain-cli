package task

import (
	"context"
	"errors"
)

var (
	ErrNotFound    = errors.New("task no encontrada")
	ErrDuplicateID = errors.New("task.id ya existe")
)

type Repository interface {
	Create(ctx context.Context, t *Task) error                       // Create inserta y rellena t.CreatedAt.
	GetByID(ctx context.Context, id string) (*Task, error)           // GetByID devuelve ErrNotFound si no existe o está soft-deleted.
	List(ctx context.Context) ([]*Task, error)                       // List no incluye soft-deleted.
	ListByTag(ctx context.Context, tags []string) ([]*Task, error)   // ListByTag devuelve Tasks que contengan TODAS las tags dadas.
	ListByPriority(ctx context.Context, p Priority) ([]*Task, error) // ListByPriority filtra por prioridad. Sin soft-deleted.
	Update(ctx context.Context, t *Task) error                       // Update modifica la task. Rellena t.UpdatedAt.
	Delete(ctx context.Context, id string) error                     // Delete hace soft-delete.
}
