package database

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strconv"
	"strings"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// migrator aplica migraciones embebidas de forma secuencial.
type migrator struct {
	db *sql.DB
}

func newMigrator(db *sql.DB) *migrator {
	return &migrator{db: db}
}

// Run aplica todas las migraciones pendientes.
// Cada migración va en su propia transacción (atómica).
func (m *migrator) Run(ctx context.Context) error {
	if err := m.ensureTable(ctx); err != nil {
		return fmt.Errorf("creando tabla schema_migrations: %w", err)
	}

	applied, err := m.appliedVersions(ctx)
	if err != nil {
		return fmt.Errorf("leyendo migraciones aplicadas: %w", err)
	}

	entries, err := fs.Glob(migrationsFS, "migrations/*.sql")
	if err != nil {
		return fmt.Errorf("listando migraciones: %w", err)
	}
	sort.Strings(entries)

	for _, path := range entries {
		version, err := parseVersion(path)
		if err != nil {
			return fmt.Errorf("nombre de migración inválido %q: %w", path, err)
		}
		if applied[version] {
			continue
		}

		content, err := migrationsFS.ReadFile(path)
		if err != nil {
			return fmt.Errorf("leyendo %q: %w", path, err)
		}

		if err := m.applyOne(ctx, version, string(content)); err != nil {
			return fmt.Errorf("aplicando migración %d (%s): %w", version, path, err)
		}
	}
	return nil
}

func (m *migrator) ensureTable(ctx context.Context) error {
	const ddl = `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version    INTEGER PRIMARY KEY,
		name       TEXT    NOT NULL,
		applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := m.db.ExecContext(ctx, ddl)
	return err
}

func (m *migrator) appliedVersions(ctx context.Context) (map[int]bool, error) {
	rows, err := m.db.QueryContext(ctx, "SELECT version FROM schema_migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make(map[int]bool)
	for rows.Next() {
		var v int
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		out[v] = true
	}
	return out, rows.Err()
}

// applyOne ejecuta una migración en una transacción.
func (m *migrator) applyOne(ctx context.Context, version int, sqlText string) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:errcheck // Rollback tras Commit es no-op.

	if _, err := tx.ExecContext(ctx, sqlText); err != nil {
		return fmt.Errorf("exec SQL: %w", err)
	}

	if _, err := tx.ExecContext(ctx,
		"INSERT INTO schema_migrations (version, name) VALUES (?, ?)",
		version, fmt.Sprintf("%03d", version),
	); err != nil {
		return fmt.Errorf("registrando versión: %w", err)
	}

	return tx.Commit()
}

// parseVersion extrae el número de "migrations/001_initial.sql" → 1.
func parseVersion(path string) (int, error) {
	base := path[strings.LastIndex(path, "/")+1:]
	idx := strings.Index(base, "_")
	if idx <= 0 {
		return 0, fmt.Errorf("falta prefijo numérico")
	}
	return strconv.Atoi(base[:idx])
}
