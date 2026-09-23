package settings

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("setting no encontrado")

// Repository define la persistencia de overrides.
type Repository interface {
	Get(ctx context.Context, key string) (*Setting, error)
	Set(ctx context.Context, key, value string) error
	// SetMany aplica varios cambios en una transacción atómica.
	SetMany(ctx context.Context, values map[string]string) error
	Delete(ctx context.Context, key string) error
	// Reset borra TODOS los overrides (vuelve a defaults + YAML).
	Reset(ctx context.Context) error
	List(ctx context.Context) ([]Setting, error)
	ListByCategory(ctx context.Context, category string) ([]Setting, error)
}
