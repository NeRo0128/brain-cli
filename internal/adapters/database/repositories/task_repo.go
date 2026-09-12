package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/NeRo0128/brain-cli/internal/core/task"
)

// TaskRepository implementa task.Repository con SQLite.
type TaskRepository struct {
	db *sql.DB
}

// NewTaskRepository construye el repositorio.
func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

var _ task.Repository = (*TaskRepository)(nil)

const taskColumns = `id, name, description, type, tool_id,
	requires_ai, ai_prompt, params, priority,
	is_active, is_favorite,
	created_at, updated_at, deleted_at`

// Create inserta la task y sus tags en una sola transacción.
func (r *TaskRepository) Create(ctx context.Context, t *task.Task) error {
	params, err := marshalJSONMap(t.Params)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	t.CreatedAt = now

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck // no-op si commit fue exitoso

	const q = `
		INSERT INTO tasks (
			id, name, description, type, tool_id,
			requires_ai, ai_prompt, params, priority,
			is_active, is_favorite, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err = tx.ExecContext(ctx, q,
		t.ID, t.Name, t.Description, string(t.Type), ptrIntToArg(t.ToolID),
		boolToInt(t.RequiresAI), t.AIPrompt, params, string(t.Priority),
		boolToInt(t.IsActive), boolToInt(t.IsFavorite), now,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return task.ErrDuplicateID
		}
		return fmt.Errorf("insertando task %q: %w", t.ID, err)
	}

	if err := syncTags(ctx, tx, t.ID, t.Tags); err != nil {
		return err
	}

	return tx.Commit()
}

// GetByID devuelve task.ErrNotFound si no existe o está soft-deleted.
func (r *TaskRepository) GetByID(ctx context.Context, id string) (*task.Task, error) {
	q := "SELECT " + taskColumns + " FROM tasks WHERE id = ? AND deleted_at IS NULL"
	t, err := scanTask(r.db.QueryRowContext(ctx, q, id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, task.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get task %q: %w", id, err)
	}

	// Cargar tags
	tagsByID, err := loadTagsByTaskIDs(ctx, r.db, []string{t.ID})
	if err != nil {
		return nil, err
	}
	t.Tags = tagsByID[t.ID]
	return t, nil
}

// List devuelve todas las tasks activas (no soft-deleted).
func (r *TaskRepository) List(ctx context.Context) ([]*task.Task, error) {
	q := "SELECT " + taskColumns + " FROM tasks WHERE deleted_at IS NULL ORDER BY priority DESC, name"
	return r.queryTasksWithTags(ctx, q)
}

// ListByTag devuelve tasks que contengan TODAS las tags dadas.
func (r *TaskRepository) ListByTag(ctx context.Context, tags []string) ([]*task.Task, error) {
	if len(tags) == 0 {
		return r.List(ctx)
	}

	q := `SELECT ` + taskColumns + ` FROM tasks
		WHERE deleted_at IS NULL
		AND id IN (
			SELECT task_id FROM task_tags
			WHERE tag_id IN (SELECT id FROM tags WHERE name IN (` + placeholders(len(tags)) + `))
			GROUP BY task_id
			HAVING COUNT(DISTINCT tag_id) = ?
		)
		ORDER BY priority DESC, name`

	args := make([]any, 0, len(tags)+1)
	for _, tag := range tags {
		args = append(args, tag)
	}
	args = append(args, len(tags))

	return r.queryTasksWithTags(ctx, q, args...)
}

// ListByPriority filtra por prioridad.
func (r *TaskRepository) ListByPriority(ctx context.Context, p task.Priority) ([]*task.Task, error) {
	q := "SELECT " + taskColumns + " FROM tasks WHERE priority = ? AND deleted_at IS NULL ORDER BY name"
	return r.queryTasksWithTags(ctx, q, string(p))
}

// Update modifica la task y sincroniza sus tags.
func (r *TaskRepository) Update(ctx context.Context, t *task.Task) error {
	params, err := marshalJSONMap(t.Params)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	t.UpdatedAt = &now

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck

	const q = `
		UPDATE tasks SET
			name = ?, description = ?, type = ?, tool_id = ?,
			requires_ai = ?, ai_prompt = ?, params = ?, priority = ?,
			is_active = ?, is_favorite = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`
	res, err := tx.ExecContext(ctx, q,
		t.Name, t.Description, string(t.Type), ptrIntToArg(t.ToolID),
		boolToInt(t.RequiresAI), t.AIPrompt, params, string(t.Priority),
		boolToInt(t.IsActive), boolToInt(t.IsFavorite), now,
		t.ID,
	)
	if err != nil {
		return fmt.Errorf("actualizando task %q: %w", t.ID, err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if n == 0 {
		return task.ErrNotFound
	}

	if err := syncTags(ctx, tx, t.ID, t.Tags); err != nil {
		return err
	}

	return tx.Commit()
}

// Delete hace soft-delete. No toca task_tags (se preservan para auditoría).
func (r *TaskRepository) Delete(ctx context.Context, id string) error {
	now := time.Now().UTC()
	const q = "UPDATE tasks SET deleted_at = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL"
	res, err := r.db.ExecContext(ctx, q, now, now, id)
	if err != nil {
		return fmt.Errorf("soft-deleting task %q: %w", id, err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if n == 0 {
		return task.ErrNotFound
	}
	return nil
}

// --- Helpers privados --------------------------------------------------

// queryTasksWithTags ejecuta un SELECT y carga tags en batch (evita N+1).
func (r *TaskRepository) queryTasksWithTags(ctx context.Context, q string, args ...any) ([]*task.Task, error) {
	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("query tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*task.Task
	for rows.Next() {
		t, err := scanTask(rows)
		if err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}
		tasks = append(tasks, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterando tasks: %w", err)
	}

	if len(tasks) == 0 {
		return tasks, nil
	}

	ids := make([]string, len(tasks))
	for i, t := range tasks {
		ids[i] = t.ID
	}
	tagsByID, err := loadTagsByTaskIDs(ctx, r.db, ids)
	if err != nil {
		return nil, err
	}
	for _, t := range tasks {
		t.Tags = tagsByID[t.ID]
	}
	return tasks, nil
}

// scanTask convierte una fila a *task.Task. Orden = taskColumns.
func scanTask(s rowScanner) (*task.Task, error) {
	var (
		t          task.Task
		taskType   string
		toolID     sql.NullInt64
		requiresAI int
		params     string
		priority   string
		isActive   int
		isFavorite int
		updatedAt  sql.NullTime
		deletedAt  sql.NullTime
	)
	err := s.Scan(
		&t.ID, &t.Name, &t.Description, &taskType, &toolID,
		&requiresAI, &t.AIPrompt, &params, &priority,
		&isActive, &isFavorite,
		&t.CreatedAt, &updatedAt, &deletedAt,
	)
	if err != nil {
		return nil, err
	}

	t.Type = task.TaskType(taskType)
	t.ToolID = nullInt64ToPtr(toolID)
	t.RequiresAI = requiresAI == 1
	t.Priority = task.Priority(priority)
	t.IsActive = isActive == 1
	t.IsFavorite = isFavorite == 1
	t.UpdatedAt = nullTimeToPtr(updatedAt)
	t.DeletedAt = nullTimeToPtr(deletedAt)

	t.Params, err = unmarshalJSONMap(params)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// loadTagsByTaskIDs carga los tags de N tasks en una sola query.
// Devuelve map[taskID][]tagName. Tasks sin tags no aparecen en el map.
func loadTagsByTaskIDs(ctx context.Context, db *sql.DB, ids []string) (map[string][]string, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	q := `SELECT tt.task_id, t.name
		FROM task_tags tt
		JOIN tags t ON t.id = tt.tag_id
		WHERE tt.task_id IN (` + placeholders(len(ids)) + `)
		ORDER BY t.name`

	args := make([]any, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	rows, err := db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("query tags: %w", err)
	}
	defer rows.Close()

	out := make(map[string][]string)
	for rows.Next() {
		var taskID, tagName string
		if err := rows.Scan(&taskID, &tagName); err != nil {
			return nil, fmt.Errorf("scan tag: %w", err)
		}
		out[taskID] = append(out[taskID], tagName)
	}
	return out, rows.Err()
}

// syncTags reemplaza los tags de una task: DELETE + INSERT.
// Auto-crea tags nuevos con color vacío (INSERT OR IGNORE).
// Se ejecuta dentro de una transacción existente.
func syncTags(ctx context.Context, tx *sql.Tx, taskID string, tags []string) error {
	if _, err := tx.ExecContext(ctx, "DELETE FROM task_tags WHERE task_id = ?", taskID); err != nil {
		return fmt.Errorf("limpiando task_tags: %w", err)
	}
	if len(tags) == 0 {
		return nil
	}

	for _, name := range tags {
		// Auto-crear tag si no existe. Color vacío → UI usa default.
		if _, err := tx.ExecContext(ctx,
			"INSERT OR IGNORE INTO tags (name, color) VALUES (?, '')", name,
		); err != nil {
			return fmt.Errorf("creando tag %q: %w", name, err)
		}
		// Vincular (subquery por name para obtener el id)
		if _, err := tx.ExecContext(ctx,
			"INSERT INTO task_tags (task_id, tag_id) SELECT ?, id FROM tags WHERE name = ?",
			taskID, name,
		); err != nil {
			return fmt.Errorf("linkeando tag %q: %w", name, err)
		}
	}
	return nil
}

// Evita el warning de import no usado si el compilador se queja.
var _ = strings.TrimSpace
