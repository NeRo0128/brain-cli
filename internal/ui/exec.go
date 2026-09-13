package ui

import (
	"context"
	"time"

	"github.com/NeRo0128/brain-cli/internal/core/execution"
	toastComp "github.com/NeRo0128/brain-cli/internal/ui/components/toast"
	"github.com/NeRo0128/brain-cli/internal/ui/screens"

	tea "github.com/charmbracelet/bubbletea"
)

// executionFinishedMsg transporta el resultado de la ejecución.
type executionFinishedMsg struct {
	taskName string
	exec     *execution.Execution
	err      error
}

// startExecution arranca una ejecución y pushea el ExecutingScreen.
func (m Model) startExecution(msg screens.ExecuteTaskMsg, toastCmd tea.Cmd) (tea.Model, tea.Cmd) {
	// Si ya hay una ejecución en curso, ignorar para no abrir spinners duplicados.
	if m.executing {
		return m, toastCmd
	}
	ctx, cancel := context.WithCancel(context.Background())
	m.cancelExec = cancel
	m.executing = true

	ex := screens.NewExecutingScreen(msg.TaskName)
	m.stack = append(m.stack, ex)

	return m, tea.Batch(
		toastCmd,
		tea.WindowSize(),
		ex.Init(),
		m.runTask(ctx, msg.TaskID, msg.TaskName),
	)
}

// handleExecutionFinished reemplaza el ExecutingScreen por el ResultScreen.
func (m Model) handleExecutionFinished(msg executionFinishedMsg, toastCmd tea.Cmd) (tea.Model, tea.Cmd) {
	m.cancelExec = nil
	m.executing = false

	if msg.err != nil {
		m.execErr = msg.err
		if len(m.stack) > 1 {
			m.stack = m.stack[:len(m.stack)-1]
		}
		return m, toastCmd
	}

	r := screens.NewResultScreen(msg.taskName, msg.exec)
	m.setTop(r)

	var toast tea.Cmd
	switch msg.exec.Status {
	case execution.StatusCompleted:
		toast = toastComp.ShowSuccess(msg.taskName + " completada")
	case execution.StatusFailed:
		toast = toastComp.ShowError(msg.taskName + " falló")
	case execution.StatusCancelled:
		toast = toastComp.ShowWarning(msg.taskName + " cancelada")
	}
	return m, tea.Batch(toastCmd, r.Init(), toast)
}

// runTask ejecuta la Task en background y devuelve el resultado.
func (m Model) runTask(ctx context.Context, taskID, taskName string) tea.Cmd {
	execUC := m.deps.ExecUC
	log := m.deps.Log
	return func() tea.Msg {
		start := time.Now()
		log.Debug().Str("task_id", taskID).Msg("runTask: iniciando")

		exec, err := execUC.Execute(ctx, taskID)

		if remaining := devMinExecutingDisplay - time.Since(start); remaining > 0 {
			select {
			case <-time.After(remaining):
			case <-ctx.Done():
			}
		}

		log.Debug().Err(err).Bool("has_exec", exec != nil).Msg("runTask: terminado")
		return executionFinishedMsg{taskName: taskName, exec: exec, err: err}
	}
}
