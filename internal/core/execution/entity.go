package execution

import (
	"errors"
	"fmt"
	"time"
)

// Execution es el registro de UNA corrida de una Task.
type Execution struct {
	ID          int               `db:"id"           json:"id"`
	TaskID      string            `db:"task_id"      json:"task_id"`
	Status      Status            `db:"status"       json:"status"`
	TriggeredBy TriggeredBy       `db:"triggered_by" json:"triggered_by"`
	Output      string            `db:"output"       json:"output,omitempty"`
	Error       string            `db:"error"        json:"error,omitempty"`
	ExitCode    *int              `db:"exit_code"    json:"exit_code,omitempty"`
	ProviderID  *int              `db:"provider_id"  json:"provider_id,omitempty"`
	Params      map[string]string `db:"params"       json:"params,omitempty"`
	StartedAt   time.Time         `db:"started_at"   json:"started_at"`
	FinishedAt  *time.Time        `db:"finished_at"  json:"finished_at,omitempty"`
}

// Duration devuelve el tiempo transcurrido. Si aún corre, contra ahora.
func (e *Execution) Duration() time.Duration {
	if e.FinishedAt == nil {
		return time.Since(e.StartedAt)
	}
	return e.FinishedAt.Sub(e.StartedAt)
}

// Validate comprueba la coherencia de la Execution.
func (e *Execution) Validate() error {
	var errs []error

	if e.TaskID == "" {
		errs = append(errs, errors.New("execution.task_id no puede estar vacío"))
	}
	if !e.Status.IsValid() {
		errs = append(errs, fmt.Errorf("execution.status inválido: %q", e.Status))
	}
	if !e.TriggeredBy.IsValid() {
		errs = append(errs, fmt.Errorf("execution.triggered_by inválido: %q", e.TriggeredBy))
	}
	if e.StartedAt.IsZero() {
		errs = append(errs, errors.New("execution.started_at requerido"))
	}

	if e.Status.IsTerminal() && e.FinishedAt == nil {
		errs = append(errs, fmt.Errorf("execution status=%q requiere finished_at", e.Status))
	}
	if !e.Status.IsTerminal() && e.FinishedAt != nil {
		errs = append(errs, fmt.Errorf("execution status=%q no admite finished_at", e.Status))
	}

	switch e.Status {
	case StatusCompleted:
		if e.Error != "" {
			errs = append(errs, errors.New("execution status=completed no admite error"))
		}
	case StatusPending, StatusRunning:
		if e.Error != "" {
			errs = append(errs, fmt.Errorf("execution status=%q no admite error", e.Status))
		}
	}

	if e.ExitCode != nil && !e.Status.IsTerminal() {
		errs = append(errs, errors.New("execution no terminal no admite exit_code"))
	}

	return errors.Join(errs...)
}
