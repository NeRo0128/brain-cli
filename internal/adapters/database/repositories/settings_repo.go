package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/NeRo0128/brain-cli/internal/core/settings"
)

type SettingsRepository struct {
	db *sql.DB
}

func NewSettingsRepository(db *sql.DB) *SettingsRepository {
	return &SettingsRepository{db: db}
}

var _ settings.Repository = (*SettingsRepository)(nil)

const settingsColumns = `key, value, category, is_encrypted, updated_at`

// Get devuelve ErrNotFound si la key no existe.
func (r *SettingsRepository) Get(ctx context.Context, key string) (*settings.Setting, error) {
	q := "SELECT " + settingsColumns + " FROM settings WHERE key = ?"
	s, err := scanSetting(r.db.QueryRowContext(ctx, q, key))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, settings.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get setting %q: %w", key, err)
	}
	return s, nil
}

// Set hace upsert. La categoría se deriva del prefijo de la key
// ("ui.theme" → "ui").
func (r *SettingsRepository) Set(ctx context.Context, key, value string) error {
	return r.SetMany(ctx, map[string]string{key: value})
}

// SetMany aplica varios upserts en una sola transacción.
func (r *SettingsRepository) SetMany(ctx context.Context, values map[string]string) error {
	if len(values) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	const q = `
        INSERT INTO settings (key, value, category, updated_at)
        VALUES (?, ?, ?, CURRENT_TIMESTAMP)
        ON CONFLICT(key) DO UPDATE SET
            value      = excluded.value,
            category   = excluded.category,
            updated_at = CURRENT_TIMESTAMP
    `

	for key, value := range values {
		cat := categoryFor(key)
		if _, err := tx.ExecContext(ctx, q, key, value, cat); err != nil {
			return fmt.Errorf("upsert setting %q: %w", key, err)
		}
	}

	return tx.Commit()
}

// Delete es idempotente: no falla si la key no existe.
func (r *SettingsRepository) Delete(ctx context.Context, key string) error {
	if _, err := r.db.ExecContext(ctx, "DELETE FROM settings WHERE key = ?", key); err != nil {
		return fmt.Errorf("delete setting %q: %w", key, err)
	}
	return nil
}

// Reset borra todos los overrides. No falla si la tabla está vacía.
func (r *SettingsRepository) Reset(ctx context.Context) error {
	if _, err := r.db.ExecContext(ctx, "DELETE FROM settings"); err != nil {
		return fmt.Errorf("reset settings: %w", err)
	}
	return nil
}

func (r *SettingsRepository) List(ctx context.Context) ([]settings.Setting, error) {
	q := "SELECT " + settingsColumns + " FROM settings ORDER BY key"
	return r.querySettings(ctx, q)
}

func (r *SettingsRepository) ListByCategory(ctx context.Context, category string) ([]settings.Setting, error) {
	q := "SELECT " + settingsColumns + " FROM settings WHERE category = ? ORDER BY key"
	return r.querySettings(ctx, q, category)
}

// --- helpers ---

func (r *SettingsRepository) querySettings(ctx context.Context, q string, args ...any) ([]settings.Setting, error) {
	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("query settings: %w", err)
	}
	defer rows.Close()

	var out []settings.Setting
	for rows.Next() {
		s, err := scanSetting(rows)
		if err != nil {
			return nil, fmt.Errorf("scan setting: %w", err)
		}
		out = append(out, *s)
	}
	return out, rows.Err()
}

// scanSetting convierte una fila en Setting. El orden debe coincidir
// con settingsColumns.
func scanSetting(s rowScanner) (*settings.Setting, error) {
	var (
		st          settings.Setting
		isEncrypted int
	)
	err := s.Scan(&st.Key, &st.Value, &st.Category, &isEncrypted, &st.UpdatedAt)
	if err != nil {
		return nil, err
	}
	st.IsEncrypted = isEncrypted == 1
	return &st, nil
}

// categoryFor extrae la categoría del prefijo de la key.
// "ui.theme" → "ui", "logging.level" → "logging".
// Sin prefijo devuelve "general".
func categoryFor(key string) string {
	for i := 0; i < len(key); i++ {
		if key[i] == '.' {
			return key[:i]
		}
	}
	return "general"
}
