package task

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/NeRo0128/brain-cli/internal/adapters/executor"
	"github.com/NeRo0128/brain-cli/internal/core/execution"
	coretask "github.com/NeRo0128/brain-cli/internal/core/task"
	"github.com/NeRo0128/brain-cli/internal/core/tool"
)

// defaultTimeout se usa si el Tool no define TimeoutSeconds.
const defaultTimeout = 5 * time.Minute

// Execute corre una Task end-to-end.
//
// Pasos:
//  1. Carga la Task de la DB.
//  2. Valida que sea ejecutable localmente (script/command).
//  3. Carga el Tool asociado.
//  4. Registra una Execution en estado "running".
//  5. Ejecuta el script.
//  6. Actualiza la Execution con el resultado.
//  7. Devuelve la Execution finalizada.
//
// La Execution se registra ANTES de ejecutar para que quede
// auditoría incluso si el proceso crashea la app.
func (e *Executor) Execute(ctx context.Context, taskID string) (*execution.Execution, error) {
	// 1. Cargar Task
	tk, err := e.tasks.GetByID(ctx, taskID)
	if errors.Is(err, coretask.ErrNotFound) {
		return nil, ErrTaskNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("cargando task %q: %w", taskID, err)
	}

	// 2. Validar tipo (ai no está soportado aún)
	if tk.Type == coretask.TaskTypeAI {
		return nil, fmt.Errorf("%w: %q (type=ai)", ErrUnsupportedTaskType, taskID)
	}
	if tk.ToolID == nil || *tk.ToolID <= 0 {
		return nil, fmt.Errorf("%w: %q", ErrTaskWithoutTool, taskID)
	}

	// 3. Cargar Tool
	tl, err := e.tools.GetByID(ctx, *tk.ToolID)
	if errors.Is(err, tool.ErrNotFound) {
		return nil, fmt.Errorf("%w: tool_id=%d", ErrToolNotFound, *tk.ToolID)
	}
	if err != nil {
		return nil, fmt.Errorf("cargando tool %d: %w", *tk.ToolID, err)
	}

	// 4. Crear Execution en "running" ANTES de ejecutar
	exec, err := e.startExecution(ctx, tk)
	if err != nil {
		return nil, err
	}

	// 5. Elegir executor concreto vía factory
	runner, err := executor.New(tl)
	if err != nil {
		// Registramos el fallo: el tool existe pero no se puede ejecutar
		return e.finishFailed(ctx, exec, fmt.Errorf("creando executor: %w", err))
	}

	// 6. Ejecutar con timeout del Tool (o default)
	runCtx, cancel := context.WithTimeout(ctx, toolTimeout(tl))
	defer cancel()

	result, runErr := runner.Execute(runCtx, tl, tk.Params)
	if runErr != nil {
		// Error de infraestructura (no se pudo lanzar el proceso)
		return e.finishFailed(ctx, exec, runErr)
	}

	// 7. Traducir Result → estado final
	return e.finishFromResult(ctx, exec, result)
}

// startExecution crea y persiste una Execution en estado "running".
func (e *Executor) startExecution(ctx context.Context, tk *coretask.Task) (*execution.Execution, error) {
	exec := &execution.Execution{
		TaskID:      tk.ID,
		Status:      execution.StatusRunning,
		TriggeredBy: execution.TriggerManual,
		Params:      cloneParams(tk.Params),
		StartedAt:   time.Now().UTC(),
	}
	if err := e.executions.Create(ctx, exec); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrExecutionCreate, err)
	}
	return exec, nil
}

// finishFromResult traduce un Result exitoso a una Execution terminada.
func (e *Executor) finishFromResult(
	ctx context.Context,
	exec *execution.Execution,
	result tool.Result,
) (*execution.Execution, error) {
	now := time.Now().UTC()
	exec.FinishedAt = &now
	exec.Output = result.Output

	switch {
	case result.ExitCode == 0:
		exec.Status = execution.StatusCompleted
		exit := 0
		exec.ExitCode = &exit

	case result.ExitCode == -1:
		// Timeout o cancelación: el output ya trae la nota del executor
		exec.Status = execution.StatusCancelled
		exit := -1
		exec.ExitCode = &exit

	default:
		exec.Status = execution.StatusFailed
		exec.Error = fmt.Sprintf("exit code %d", result.ExitCode)
		exit := result.ExitCode
		exec.ExitCode = &exit
	}

	// Contexto de persistencia: independiente del de ejecución.
	// Aunque el script fue cancelado/timeout, debemos poder guardar el resultado.
	saveCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := e.executions.Update(saveCtx, exec); err != nil {
		return exec, fmt.Errorf("%w: %v", ErrExecutionUpdate, err)
	}
	return exec, nil
}

// finishFailed registra una Execution fallida por error de infraestructura.
func (e *Executor) finishFailed(
	ctx context.Context,
	exec *execution.Execution,
	cause error,
) (*execution.Execution, error) {
	now := time.Now().UTC()
	exec.Status = execution.StatusFailed
	exec.Error = cause.Error()
	exec.FinishedAt = &now

	saveCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := e.executions.Update(saveCtx, exec); err != nil {
		return exec, fmt.Errorf("%w: %v", ErrExecutionUpdate, err)
	}
	return exec, cause
}

// toolTimeout devuelve el timeout efectivo del Tool.
func toolTimeout(tl *tool.Tool) time.Duration {
	if tl.TimeoutSeconds <= 0 {
		return defaultTimeout
	}
	return time.Duration(tl.TimeoutSeconds) * time.Second
}

// cloneParams copia el map para que la Execution no comparta
// memoria con la Task (evita mutaciones accidentales).
func cloneParams(in map[string]string) map[string]string {
	if len(in) == 0 {
		return nil
	}
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}
