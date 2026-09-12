package ui

import (
	"context"
	"fmt"
	"time"

	coreConfig "github.com/NeRo0128/brain-cli/internal/core/config"
	"github.com/NeRo0128/brain-cli/internal/core/execution"
	coreTask "github.com/NeRo0128/brain-cli/internal/core/task"
	"github.com/NeRo0128/brain-cli/internal/core/tool"
	"github.com/NeRo0128/brain-cli/internal/ui/screens"
	usercTask "github.com/NeRo0128/brain-cli/internal/usecases/task"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rs/zerolog"
)

// view identifica la pantalla activa.
type view int

const devMinExecutingDisplay = 0 * time.Second

const (
	viewMain view = iota
	viewDetails
	viewExecuting
	viewHistory
	viewResult
)

// executionFinishedMsg transporta el resultado de la ejecución.
type executionFinishedMsg struct {
	taskName string
	exec     *execution.Execution
	err      error
}

// Model es el modelo raíz de la aplicación.
// Contiene las pantallas y el estado global.
type Model struct {
	version  string
	cfg      *coreConfig.Config
	execUC   *usercTask.Executor
	toolRepo tool.Repository
	execRepo execution.Repository
	taskRepo coreTask.Repository
	active   view
	log      zerolog.Logger

	// Pantallas
	main      screens.MainScreen
	details   screens.DetailScreen
	executing screens.ExecutingScreen
	result    screens.ResultScreen
	history   screens.HistoryScreen

	// Estado de ejecución en curso
	cancelExec context.CancelFunc
	execErr    error

	resultParent view // resultParent: a dónde vuelve Esc desde viewResult.
}

// NewModels construye el modelo raíz con todas sus dependencias.
func NewModels(
	version string,
	cfg *coreConfig.Config,
	execUC *usercTask.Executor,
	taskRepo coreTask.Repository,
	toolRepo tool.Repository,
	execRepo execution.Repository,
	mainScreen screens.MainScreen,
	log zerolog.Logger,
) Model {
	return Model{
		version:  version,
		cfg:      cfg,
		execUC:   execUC,
		toolRepo: toolRepo,
		taskRepo: taskRepo,
		execRepo: execRepo,
		active:   viewMain,
		main:     mainScreen,
		log:      log.With().Str("component", "ui").Logger(),
	}
}

// Init arranca la aplicación.
func (m Model) Init() tea.Cmd { return m.main.Init() }

// Update maneja mensajes globales y delega a la pantalla activa.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			if m.active != viewExecuting {
				return m, tea.Quit
			}
		case "ctrl+c":
			if m.cancelExec != nil {
				m.cancelExec()
			}
			return m, tea.Quit
		case "esc":
			switch m.active {
			case viewExecuting:
				if m.cancelExec != nil {
					m.log.Info().Msg("cancelando ejecución")
					m.cancelExec()
					m.cancelExec = nil
					m.executing = m.executing.MarkCanceling()
				}
				return m, nil
			default:
				m.log.Debug().Msg("esc → volviendo a main")
				m.active = viewMain
				return m, nil

			}
		case "d":
			if m.active == viewMain {
				tk := m.main.SelectedTask()
				if tk == nil {
					m.log.Warn().Msg("no hay task seleccionada")
					return m, nil
				}
				m.log.Debug().Str("task_id", tk.ID).Msg("abriendo detalle")
				m.details = screens.NewDetailScreen(tk, m.toolRepo, m.log)
				m.active = viewDetails
				return m, m.details.Init()
			}
		case "h":
			switch m.active {
			case viewMain:
				m.history = screens.NewHistoryScreen(m.execRepo, m.taskRepo)
				m.active = viewHistory
				return m, m.history.Init()
			case viewHistory:
				m.active = viewMain
				return m, nil
			}
		case "enter", "e":
			switch m.active {
			case viewMain, viewDetails, viewHistory:
				if cmd := m.tryExecute(); cmd != nil {
					return m, cmd
				}
			}
		}

	case executionFinishedMsg:
		m.log.Debug().
			Err(msg.err).
			Bool("has_exec", msg.exec != nil).
			Msg("executionFinishedMsg recibido")
		m.cancelExec = nil
		if msg.err != nil {
			m.execErr = msg.err
			m.active = viewMain
			return m, nil
		}
		m.result = screens.NewResultScreen(msg.taskName, msg.exec)
		m.active = viewResult
		return m, nil
	}

	var cmd tea.Cmd
	switch m.active {
	case viewMain:
		m.main, cmd = m.main.Update(msg)
	case viewResult:
		m.result, cmd = m.result.Update(msg)
	case viewExecuting:
		m.executing, cmd = m.executing.Update(msg)
	case viewDetails:
		m.details, cmd = m.details.Update(msg)
	case viewHistory:
		m.history, cmd = m.history.Update(msg)
	}
	return m, cmd
}

// tryExecute intenta lanzar la ejecución de la task seleccionada.
// Devuelve nil si no hay nada que ejecutar (o ya hay una en curso).
func (m *Model) tryExecute() tea.Cmd {
	var tk *coreTask.Task
	switch m.active {
	case viewMain:
		tk = m.main.SelectedTask()
	case viewDetails:
		tk = m.details.Task()
	}
	if tk == nil || m.active == viewExecuting {
		return nil
	}

	// Creamos el ctx y guardamos el cancel en el modelo.
	ctx, cancel := context.WithCancel(context.Background())
	m.cancelExec = cancel
	m.resultParent = m.active
	m.executing = screens.NewExecutingScreen(tk.Name)
	m.active = viewExecuting

	m.log.Debug().Str("task_id", tk.ID).Msg("ejecutando")

	return tea.Batch(
		m.executing.Init(),
		m.runTask(ctx, tk.ID, tk.Name),
	)
}

// runTask dispara la ejecución en background y devuelve el resultado.
func (m Model) runTask(ctx context.Context, taskID, taskName string) tea.Cmd {
	execUC := m.execUC
	log := m.log
	return func() tea.Msg {
		start := time.Now()
		log.Debug().Str("task_id", taskID).Msg("runTask: iniciando")

		exec, err := execUC.Execute(ctx, taskID)

		// FIXME: mantener el spinner visible al menos N segundos (solo dev).
		// Si el usuario canceló, salimos inmediatamente.
		if remaining := devMinExecutingDisplay - time.Since(start); remaining > 0 {
			select {
			case <-time.After(remaining):
			case <-ctx.Done():
				// ctx cancelado → no esperamos
			}
		}

		log.Debug().Err(err).Bool("has_exec", exec != nil).Msg("runTask: terminado")
		return executionFinishedMsg{
			taskName: taskName,
			exec:     exec,
			err:      err,
		}
	}
}

func (m Model) View() string {
	if m.execErr != nil {
		return m.renderError()
	}
	switch m.active {
	case viewDetails:
		return m.details.View()
	case viewExecuting:
		return m.executing.View()
	case viewResult:
		return m.result.View()
	case viewHistory:
		return m.history.View()
	default:
		return m.main.View()
	}
}

func (m Model) renderExecuting() string {
	return "\n\n  ⏳ Ejecutando... (Ctrl+C para cancelar)\n"
}

func (m Model) renderError() string {
	return fmt.Sprintf("\n\n  ✗ Error: %v\n\n  Esc para volver\n", m.execErr)
}

func viewName(v view) string {
	switch v {
	case viewMain:
		return "main"
	case viewDetails:
		return "details"
	case viewExecuting:
		return "executing"
	case viewResult:
		return "result"
	case viewHistory:
		return "history"
	}
	return "unknown"
}
