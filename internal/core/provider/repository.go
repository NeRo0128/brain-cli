package provider

import (
	"context"
	"errors"
)

var (
	ErrNotFound      = errors.New("provider no encontrado")
	ErrDuplicateName = errors.New("provider.name ya existe")
	ErrNoActive      = errors.New("no hay provider activo")
)

// Repository define las operaciones de persistencia para proveedores
type Repository interface {
	// Create crea un nuevo proveedor
	Create(ctx context.Context, provider *Provider) error

	// GetByID obtiene un proveedor por su ID
	GetByID(ctx context.Context, id int) (*Provider, error)

	// GetByName obtiene un proveedor por su nombre
	GetByName(ctx context.Context, name string) (*Provider, error)

	// GetActive obtiene el proveedor activo actual
	GetActive(ctx context.Context) (*Provider, error)

	// List lista todos los proveedores
	List(ctx context.Context) ([]*Provider, error)

	// Update actualiza un proveedor existente
	Update(ctx context.Context, provider *Provider) error

	// Delete elimina un proveedor
	Delete(ctx context.Context, id int) error

	// SetActive establece un proveedor como activo (desactiva los demás)
	SetActive(ctx context.Context, id int) error
}

// TODO: Implementar SQLite repository en internal/adapters/database/repositories/provider_repo.go
