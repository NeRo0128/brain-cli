package execution

import (
	"context"
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("execution no encontrada")
)

// Repository define el contrato de persistencia para Executions.
type Repository interface {
	Create(ctx context.Context, e *Execution) error                                 // Create inserta y rellena e.ID.
	GetByID(ctx context.Context, id int) (*Execution, error)                        // GetByID devuelve ErrNotFound si no existe.
	ListByTask(ctx context.Context, taskID string, limit int) ([]*Execution, error) // ListByTask devuelve ejecuciones de una task (DESC por started_at).
	ListRecent(ctx context.Context, limit int) ([]*Execution, error)                // ListRecent devuelve las últimas N ejecuciones globales.
	Update(ctx context.Context, e *Execution) error                                 // Update muta estado/output/error/exit_code/finished_at.
	DeleteBefore(ctx context.Context, cutoff time.Time) (int64, error)              // DeleteBefore borra ejecuciones previas a cutoff. Devuelve filas afectadas.
}
