package database

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	
	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// SQLite representa la conexión a la base de datos
type SQLite struct {
	db *sql.DB
}

// New crea una nueva conexión a SQLite
func New(dbPath string) (*SQLite, error) {
	// TODO: Implementar conexión a SQLite
	// - Abrir conexión con CGO_ENABLED=1
	// - Configurar PRAGMA statements (WAL, foreign_keys, etc.)
	// - Ejecutar migraciones automáticas
	return nil, fmt.Errorf("not implemented")
}

// RunMigrations ejecuta todas las migraciones pendientes
func (s *SQLite) RunMigrations(ctx context.Context) error {
	// TODO: Implementar sistema de migraciones
	// - Leer archivos de migraciones desde embed.FS
	// - Ejecutar en orden (001_, 002_, etc.)
	// - Guardar versión de migración aplicada
	return fmt.Errorf("not implemented")
}

// DB retorna la instancia de *sql.DB
func (s *SQLite) DB() *sql.DB {
	return s.db
}

// Close cierra la conexión a la base de datos
func (s *SQLite) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// TODO: Implementar métodos helper
// - Ping() error
// - BeginTx(ctx context.Context) (*sql.Tx, error)
// - GetVersion() (int, error)
