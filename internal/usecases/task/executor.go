package task

import (
	"errors"

	"github.com/NeRo0128/brain-cli/internal/core/execution"
	"github.com/NeRo0128/brain-cli/internal/core/task"
	"github.com/NeRo0128/brain-cli/internal/core/tool"
)

// Errores de dominio del use case.
var (
	ErrTaskNotFound        = errors.New("task no encontrada")
	ErrToolNotFound        = errors.New("tool de la task no encontrada")
	ErrUnsupportedTaskType = errors.New("tipo de task no soportado aún")
	ErrTaskWithoutTool     = errors.New("task no tiene tool asociada")
	ErrExecutionCreate     = errors.New("no se pudo registrar la ejecución")
	ErrExecutionUpdate     = errors.New("no se pudo actualizar la ejecución")
)

type ExecutorFactory func(t *tool.Tool) (tool.Executor, error)

// Executor orquesta la ejecución de Tasks.
// Depende de interfaces del dominio, no de implementaciones concretas.
type Executor struct {
	tasks      task.Repository
	tools      tool.Repository
	executions execution.Repository
	runner     ExecutorFactory
}

// NewExecutor construye el orquestador.
func NewExecutor(
	tasks task.Repository,
	tools tool.Repository,
	executions execution.Repository,
	runner ExecutorFactory,
) *Executor {
	return &Executor{
		tasks:      tasks,
		tools:      tools,
		executions: executions,
		runner:     runner,
	}
}
