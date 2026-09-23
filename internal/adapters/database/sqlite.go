package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

// SQLite envuelve la conexión a la base de datos.
type SQLite struct {
	db *sql.DB
}

// Options controla la apertura de la conexión.
type Options struct {
	Path         string // Path al archivo .db. ":memory:" para tests.
	AutoMigrate  bool   // AutoMigrate ejecuta las migraciones al abrir.
	MaxOpenConns int    // MaxOpenConns: SQLite recomienda 1 para evitar locks.
}

// New abre la conexión, aplica PRAGMAs y, si AutoMigrate=true,
// ejecuta las migraciones pendientes.
func New(ctx context.Context, opts Options) (*SQLite, error) {

	// Crear el directorio padre si no existe (excepto :memory:).
	if opts.Path != ":memory:" && !strings.HasPrefix(opts.Path, "file::memory:") {
		dir := filepath.Dir(opts.Path)
		if dir != "" && dir != "." {
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return nil, fmt.Errorf("creando directorio de DB %q: %w", dir, err)
			}
		}
	}
	// _pragma se pasa como query param al driver modernc.
	// Importante: foreign_keys está OFF por defecto en SQLite.
	dsn := fmt.Sprintf("file:%s?_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)", opts.Path)

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("abriendo SQLite %q: %w", opts.Path, err)
	}

	// SQLite escala mal con writes concurrentes. 1 conexión = sin locks.
	maxConns := opts.MaxOpenConns
	if maxConns <= 0 {
		maxConns = 1
	}
	db.SetMaxOpenConns(maxConns)
	db.SetMaxIdleConns(maxConns)
	db.SetConnMaxLifetime(0) // sin expiración

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping a SQLite: %w", err)
	}

	// WAL: lecturas concurrentes con escrituras.
	// synchronous=NORMAL: balance velocidad/seguridad con WAL.
	pragmas := []string{
		"PRAGMA journal_mode = WAL",
		"PRAGMA synchronous  = NORMAL",
		"PRAGMA foreign_keys = ON",
		"PRAGMA temp_store   = MEMORY",
		"PRAGMA cache_size   = -64000", // 64MB
	}
	for _, p := range pragmas {
		if _, err := db.ExecContext(ctx, p); err != nil {
			db.Close()
			return nil, fmt.Errorf("aplicando %q: %w", p, err)
		}
	}

	s := &SQLite{db: db}

	if opts.AutoMigrate {
		migCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
		if err := newMigrator(db).Run(migCtx); err != nil {
			db.Close()
			return nil, fmt.Errorf("migraciones: %w", err)
		}
	}

	return s, nil
}

// DB expone el *sql.DB subyacente para los repositorios.
func (s *SQLite) DB() *sql.DB { return s.db }

// Ping verifica que la conexión sigue viva.
func (s *SQLite) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

// Close cierra la conexión.
func (s *SQLite) Close() error {
	if s.db == nil {
		return nil
	}
	return s.db.Close()
}
