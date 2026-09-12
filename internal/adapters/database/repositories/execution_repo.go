package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/NeRo0128/brain-cli/internal/core/execution"
)

// ExecutionRepository implementa execution.Repository con SQLite.
type ExecutionRepository struct {
	db *sql.DB
}

// NewExecutionRepository construye el repositorio.
func NewExecutionRepository(db *sql.DB) *ExecutionRepository {
	return &ExecutionRepository{db: db}
}

var _ execution.Repository = (*ExecutionRepository)(nil)

const executionColumns = `id, task_id, status, triggered_by,
	output, error, exit_code, provider_id, params,
	started_at, finished_at`

// Create inserta la execution y rellena e.ID.
// No insertamos finished_at: una execution nace en pending/running.
func (r *ExecutionRepository) Create(ctx context.Context, e *execution.Execution) error {
	params, err := marshalJSONMap(e.Params)
	if err != nil {
		return err
	}

	// Preservamos StartedAt si viene del caller; si no, now.
	if e.StartedAt.IsZero() {
		e.StartedAt = time.Now().UTC()
	}

	const q = `
		INSERT INTO executions (
			task_id, status, triggered_by,
			output, error, exit_code, provider_id, params,
			started_at, finished_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	res, err := r.db.ExecContext(ctx, q,
		e.TaskID, string(e.Status), string(e.TriggeredBy),
		e.Output, e.Error, ptrIntToArg(e.ExitCode), ptrIntToArg(e.ProviderID), params,
		e.StartedAt, e.FinishedAt,
	)
	if err != nil {
		return fmt.Errorf("insertando execution: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("last insert id: %w", err)
	}
	e.ID = int(id)
	return nil
}

// GetByID devuelve ErrNotFound si no existe.
func (r *ExecutionRepository) GetByID(ctx context.Context, id int) (*execution.Execution, error) {
	q := "SELECT " + executionColumns + " FROM executions WHERE id = ?"
	e, err := scanExecution(r.db.QueryRowContext(ctx, q, id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, execution.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get execution %d: %w", id, err)
	}
	return e, nil
}

// ListByTask devuelve las ejecuciones de una task, más recientes primero.
func (r *ExecutionRepository) ListByTask(ctx context.Context, taskID string, limit int) ([]*execution.Execution, error) {
	q := "SELECT " + executionColumns + " FROM executions WHERE task_id = ? ORDER BY started_at DESC"
	args := []any{taskID}
	if limit > 0 {
		q += " LIMIT ?"
		args = append(args, limit)
	}
	return r.queryExecutions(ctx, q, args...)
}

// ListRecent devuelve las últimas N ejecuciones globales.
func (r *ExecutionRepository) ListRecent(ctx context.Context, limit int) ([]*execution.Execution, error) {
	if limit <= 0 {
		limit = 100
	}
	q := "SELECT " + executionColumns + " FROM executions ORDER BY started_at DESC LIMIT ?"
	return r.queryExecutions(ctx, q, limit)
}

// Update muta status/output/error/exit_code/finished_at de una execution.
// NO toca task_id, started_at, triggered_by ni params (inmutables).
func (r *ExecutionRepository) Update(ctx context.Context, e *execution.Execution) error {
	const q = `
		UPDATE executions SET
			status = ?, output = ?, error = ?,
			exit_code = ?, provider_id = ?, finished_at = ?
		WHERE id = ?
	`
	res, err := r.db.ExecContext(ctx, q,
		string(e.Status), e.Output, e.Error,
		ptrIntToArg(e.ExitCode), ptrIntToArg(e.ProviderID), e.FinishedAt,
		e.ID,
	)
	if err != nil {
		return fmt.Errorf("actualizando execution %d: %w", e.ID, err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if n == 0 {
		return execution.ErrNotFound
	}
	return nil
}

// DeleteBefore borra ejecuciones anteriores a cutoff y devuelve cuántas borró.
func (r *ExecutionRepository) DeleteBefore(ctx context.Context, cutoff time.Time) (int64, error) {
	res, err := r.db.ExecContext(ctx, "DELETE FROM executions WHERE started_at < ?", cutoff)
	if err != nil {
		return 0, fmt.Errorf("borrando executions previas a %v: %w", cutoff, err)
	}
	return res.RowsAffected()
}

// --- Helpers privados --------------------------------------------------

func (r *ExecutionRepository) queryExecutions(ctx context.Context, q string, args ...any) ([]*execution.Execution, error) {
	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("query executions: %w", err)
	}
	defer rows.Close()

	var out []*execution.Execution
	for rows.Next() {
		e, err := scanExecution(rows)
		if err != nil {
			return nil, fmt.Errorf("scan execution: %w", err)
		}
		out = append(out, e)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterando executions: %w", err)
	}
	return out, nil
}

// scanExecution convierte una fila a *execution.Execution. Orden = executionColumns.
func scanExecution(s rowScanner) (*execution.Execution, error) {
	var (
		e           execution.Execution
		status      string
		triggeredBy string
		exitCode    sql.NullInt64
		providerID  sql.NullInt64
		params      string
		finishedAt  sql.NullTime
	)
	err := s.Scan(
		&e.ID, &e.TaskID, &status, &triggeredBy,
		&e.Output, &e.Error, &exitCode, &providerID, &params,
		&e.StartedAt, &finishedAt,
	)
	if err != nil {
		return nil, err
	}

	e.Status = execution.Status(status)
	e.TriggeredBy = execution.TriggeredBy(triggeredBy)
	e.ExitCode = nullInt64ToPtr(exitCode)
	e.ProviderID = nullInt64ToPtr(providerID)
	e.FinishedAt = nullTimeToPtr(finishedAt)

	e.Params, err = unmarshalJSONMap(params)
	if err != nil {
		return nil, err
	}
	return &e, nil
}
