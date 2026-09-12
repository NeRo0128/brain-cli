package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/NeRo0128/brain-cli/internal/core/tool"
)

// ToolRepository implementa tool.Repository con SQLite.
type ToolRepository struct {
	db *sql.DB
}

// NewToolRepository construye el repositorio.
func NewToolRepository(db *sql.DB) *ToolRepository {
	return &ToolRepository{db: db}
}

// Compile-time check: la struct cumple el contrato del dominio.
var _ tool.Repository = (*ToolRepository)(nil)

// toolColumns es el orden canónico de columnas para SELECTs.
// DEBE coincidir con el orden de scanTool.
const toolColumns = `id, name, description, script_type, category,
	script_content, script_path, command,
	requires_sudo, timeout_seconds, is_builtin, version,
	created_at, updated_at, deleted_at`

// Create inserta el tool y rellena t.ID y t.CreatedAt.
func (r *ToolRepository) Create(ctx context.Context, t *tool.Tool) error {
	now := time.Now().UTC()
	t.CreatedAt = now

	const q = `
		INSERT INTO tools (
			name, description, script_type, category,
			script_content, script_path, command,
			requires_sudo, timeout_seconds, is_builtin, version,
			created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	res, err := r.db.ExecContext(ctx, q,
		t.Name, t.Description, string(t.ScriptType), string(t.Category),
		t.ScriptContent, t.ScriptPath, t.Command,
		boolToInt(t.RequiresSudo), t.TimeoutSeconds, boolToInt(t.IsBuiltin), t.Version,
		now,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return tool.ErrDuplicateName
		}
		return fmt.Errorf("insertando tool %q: %w", t.Name, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("recuperando last insert id: %w", err)
	}
	t.ID = int(id)
	return nil
}

// GetByID devuelve tool.ErrNotFound si no existe o está soft-deleted.
func (r *ToolRepository) GetByID(ctx context.Context, id int) (*tool.Tool, error) {
	q := "SELECT " + toolColumns + " FROM tools WHERE id = ? AND deleted_at IS NULL"
	t, err := scanTool(r.db.QueryRowContext(ctx, q, id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, tool.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get tool by id %d: %w", id, err)
	}
	return t, nil
}

// GetByName devuelve tool.ErrNotFound si no existe.
func (r *ToolRepository) GetByName(ctx context.Context, name string) (*tool.Tool, error) {
	q := "SELECT " + toolColumns + " FROM tools WHERE name = ? AND deleted_at IS NULL"
	t, err := scanTool(r.db.QueryRowContext(ctx, q, name))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, tool.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get tool by name %q: %w", name, err)
	}
	return t, nil
}

// List devuelve todos los tools no borrados, ordenados por name.
func (r *ToolRepository) List(ctx context.Context) ([]*tool.Tool, error) {
	q := "SELECT " + toolColumns + " FROM tools WHERE deleted_at IS NULL ORDER BY name"
	return r.queryTools(ctx, q)
}

// ListByCategory filtra por categoría.
func (r *ToolRepository) ListByCategory(ctx context.Context, category tool.Category) ([]*tool.Tool, error) {
	q := "SELECT " + toolColumns + " FROM tools WHERE category = ? AND deleted_at IS NULL ORDER BY name"
	return r.queryTools(ctx, q, string(category))
}

// ListBuiltin devuelve solo tools con is_builtin=1.
func (r *ToolRepository) ListBuiltin(ctx context.Context) ([]*tool.Tool, error) {
	q := "SELECT " + toolColumns + " FROM tools WHERE is_builtin = 1 AND deleted_at IS NULL ORDER BY name"
	return r.queryTools(ctx, q)
}

// Update modifica un tool y rellena t.UpdatedAt.
func (r *ToolRepository) Update(ctx context.Context, t *tool.Tool) error {
	now := time.Now().UTC()
	t.UpdatedAt = &now

	const q = `
		UPDATE tools SET
			name = ?, description = ?, script_type = ?, category = ?,
			script_content = ?, script_path = ?, command = ?,
			requires_sudo = ?, timeout_seconds = ?, is_builtin = ?, version = ?,
			updated_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`
	res, err := r.db.ExecContext(ctx, q,
		t.Name, t.Description, string(t.ScriptType), string(t.Category),
		t.ScriptContent, t.ScriptPath, t.Command,
		boolToInt(t.RequiresSudo), t.TimeoutSeconds, boolToInt(t.IsBuiltin), t.Version,
		now,
		t.ID,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return tool.ErrDuplicateName
		}
		return fmt.Errorf("actualizando tool %d: %w", t.ID, err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if n == 0 {
		return tool.ErrNotFound
	}
	return nil
}

// Delete hace soft-delete (UPDATE deleted_at = now).
func (r *ToolRepository) Delete(ctx context.Context, id int) error {
	now := time.Now().UTC()
	const q = "UPDATE tools SET deleted_at = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL"
	res, err := r.db.ExecContext(ctx, q, now, now, id)
	if err != nil {
		return fmt.Errorf("soft-deleting tool %d: %w", id, err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if n == 0 {
		return tool.ErrNotFound
	}
	return nil
}

// queryTools ejecuta un SELECT y escanea todas las filas.
// Reduce el boilerplate de las 3 variantes de List.
func (r *ToolRepository) queryTools(ctx context.Context, q string, args ...any) ([]*tool.Tool, error) {
	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("query tools: %w", err)
	}
	defer rows.Close()

	var out []*tool.Tool
	for rows.Next() {
		t, err := scanTool(rows)
		if err != nil {
			return nil, fmt.Errorf("scan tool: %w", err)
		}
		out = append(out, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterando tools: %w", err)
	}
	return out, nil
}

// scanTool convierte una fila a *tool.Tool.
// El orden de Scan DEBE coincidir con toolColumns.
func scanTool(s rowScanner) (*tool.Tool, error) {
	var (
		t            tool.Tool
		scriptType   string
		category     string
		requiresSudo int
		isBuiltin    int
		updatedAt    sql.NullTime
		deletedAt    sql.NullTime
	)
	err := s.Scan(
		&t.ID, &t.Name, &t.Description, &scriptType, &category,
		&t.ScriptContent, &t.ScriptPath, &t.Command,
		&requiresSudo, &t.TimeoutSeconds, &isBuiltin, &t.Version,
		&t.CreatedAt, &updatedAt, &deletedAt,
	)
	if err != nil {
		return nil, err
	}
	t.ScriptType = tool.ScriptType(scriptType)
	t.Category = tool.Category(category)
	t.RequiresSudo = requiresSudo == 1
	t.IsBuiltin = isBuiltin == 1
	t.UpdatedAt = nullTimeToPtr(updatedAt)
	t.DeletedAt = nullTimeToPtr(deletedAt)
	return &t, nil
}
