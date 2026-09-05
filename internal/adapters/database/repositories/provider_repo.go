package repositories

import (
	"context"
	"database/sql"
	
	"github.com/nero/brain-cli/internal/core/provider"
)

// SQLiteProviderRepository implementa provider.Repository usando SQLite
type SQLiteProviderRepository struct {
	db *sql.DB
}

// NewProviderRepository crea un nuevo repositorio de proveedores
func NewProviderRepository(db *sql.DB) *SQLiteProviderRepository {
	return &SQLiteProviderRepository{db: db}
}

// Create crea un nuevo proveedor
func (r *SQLiteProviderRepository) Create(ctx context.Context, p *provider.Provider) error {
	// TODO: Implementar INSERT
	// - Encriptar API key antes de guardar
	// - Retornar ID generado
	panic("not implemented")
}

// GetByID obtiene un proveedor por su ID
func (r *SQLiteProviderRepository) GetByID(ctx context.Context, id int) (*provider.Provider, error) {
	// TODO: Implementar SELECT por ID
	// - Desencriptar API key al recuperar
	panic("not implemented")
}

// GetByName obtiene un proveedor por su nombre
func (r *SQLiteProviderRepository) GetByName(ctx context.Context, name string) (*provider.Provider, error) {
	// TODO: Implementar SELECT por nombre
	panic("not implemented")
}

// GetActive obtiene el proveedor activo actual
func (r *SQLiteProviderRepository) GetActive(ctx context.Context) (*provider.Provider, error) {
	// TODO: Implementar SELECT WHERE is_active = 1
	panic("not implemented")
}

// List lista todos los proveedores
func (r *SQLiteProviderRepository) List(ctx context.Context) ([]*provider.Provider, error) {
	// TODO: Implementar SELECT * ORDER BY name
	panic("not implemented")
}

// Update actualiza un proveedor existente
func (r *SQLiteProviderRepository) Update(ctx context.Context, p *provider.Provider) error {
	// TODO: Implementar UPDATE
	// - Encriptar API key si cambió
	// - Actualizar updated_at automáticamente (trigger)
	panic("not implemented")
}

// Delete elimina un proveedor
func (r *SQLiteProviderRepository) Delete(ctx context.Context, id int) error {
	// TODO: Implementar DELETE
	// - Verificar que no sea el proveedor activo
	panic("not implemented")
}

// SetActive establece un proveedor como activo
func (r *SQLiteProviderRepository) SetActive(ctx context.Context, id int) error {
	// TODO: Implementar transacción
	// - UPDATE providers SET is_active = 0 WHERE is_active = 1
	// - UPDATE providers SET is_active = 1 WHERE id = ?
	// - O usar trigger (ya existe en 001_initial.sql)
	panic("not implemented")
}
