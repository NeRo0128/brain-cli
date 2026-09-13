package ui

import (
	"context"
	"fmt"
	"time"

	"github.com/NeRo0128/brain-cli/internal/core/execution"
	"github.com/NeRo0128/brain-cli/internal/ui/keys"
	"github.com/NeRo0128/brain-cli/internal/ui/screens"

	tea "github.com/charmbracelet/bubbletea"
)

// devMinExecutingDisplay es un delay solo-dev para que el spinner sea visible.
const devMinExecutingDisplay = 0 * time.Second

type Model struct {
	deps  Deps
	stack []screens.ScreenI

	width, height int
	executing     bool
	cancelExec    context.CancelFunc
	execErr       error
}

func NewModels(deps Deps, initial screens.ScreenI) Model {
	return Model{
		deps:  deps,
		stack: []screens.ScreenI{initial},
	}
}

// Init arranca la aplicación.
func (m Model) Init() tea.Cmd { return m.top().Init() }

func (m *Model) top() screens.ScreenI     { return m.stack[len(m.stack)-1] }
func (m *Model) setTop(s screens.ScreenI) { m.stack[len(m.stack)-1] = s }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// --- 0. Navegación interna ---
	switch msg := msg.(type) {
	case pushMsg:
		m.stack = append(m.stack, msg.screen)
		return m, tea.WindowSize()
	case popMsg:
		if len(m.stack) > 1 {
			m.stack = m.stack[:len(m.stack)-1]
		}
		return m, nil
	case replaceMsg:
		m.setTop(msg.screen)
		return m, tea.WindowSize()
	}

	// --- 1. Navegación desde screens ---
	switch msg := msg.(type) {
	case screens.BackMsg:
		if len(m.stack) > 1 {
			m.stack = m.stack[:len(m.stack)-1]
		}
		return m, nil
	case screens.OpenDetailMsg:
		d := screens.NewDetailScreen(msg.Task, m.deps.ToolRepo, m.deps.Log)
		return m, tea.Batch(push(d), d.Init())
	case screens.OpenHistoryMsg:
		h := screens.NewHistoryScreen(m.deps.ExecRepo, m.deps.TaskRepo)
		return m, tea.Batch(push(h), h.Init())
	case screens.OpenResultMsg:
		r := screens.NewResultScreen(msg.TaskName, msg.Exec)
		return m, replace(r)
	case screens.OpenFormMsg:
		f := screens.NewFormScreen(msg.Task, m.deps.Manager, m.deps.Log)
		return m, tea.Batch(push(f), f.Init())
	case screens.ReloadMsg:
		top := m.top()
		newTop, cmd := top.Update(msg)
		m.setTop(newTop)
		return m, cmd
	case screens.OpenHelpMsg:
		h := screens.NewHelpScreen(m.currentKeyMap())
		return m, tea.Batch(push(h), h.Init())
	case screens.ExecuteTaskMsg:
		return m.startExecution(msg)
	case executionFinishedMsg:
		return m.handleExecutionFinished(msg)
	case screens.FormSavedMsg:
		// Pop del form y reload del nuevo top
		if len(m.stack) > 1 {
			m.stack = m.stack[:len(m.stack)-1]
		}
		newTop, cmd := m.top().Update(screens.ReloadMsg{})
		m.setTop(newTop)
		return m, cmd
	}

	// --- 2. Teclas globales ---
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "ctrl+c":
			if m.cancelExec != nil {
				m.cancelExec()
			}
			return m, tea.Quit
		case "q":
			if !m.executing {
				return m, tea.Quit
			}
			return m, nil
		case "esc":
			if m.executing && m.cancelExec != nil {
				m.cancelExec()
				m.cancelExec = nil
				m.executing = false
				if ex, ok := m.top().(screens.ExecutingScreen); ok {
					m.setTop(ex.MarkCanceling())
				}
				return m, nil
			}
		}
	}

	// --- 3. Traducir KeyMsg → ActionMsg ---
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		top := m.top()
		for _, id := range top.Keys() {
			if m.deps.Keys.Matches(id, keyMsg.String()) {
				msg = screens.ActionMsg{ID: id}
				break
			}
		}
	}

	// --- 4. Delegar al top ---
	top := m.top()
	newTop, cmd := top.Update(msg)
	m.setTop(newTop)
	return m, cmd
}

func (m Model) startExecution(msg screens.ExecuteTaskMsg) (tea.Model, tea.Cmd) {
	ctx, cancel := context.WithCancel(context.Background())
	m.cancelExec = cancel
	m.executing = true

	ex := screens.NewExecutingScreen(msg.TaskName)
	m.stack = append(m.stack, ex)

	return m, tea.Batch(
		tea.WindowSize(),
		ex.Init(),
		m.runTask(ctx, msg.TaskID, msg.TaskName),
	)
}

func (m Model) handleExecutionFinished(msg executionFinishedMsg) (tea.Model, tea.Cmd) {
	m.cancelExec = nil
	m.executing = false

	if msg.err != nil {
		m.execErr = msg.err
		if len(m.stack) > 1 {
			m.stack = m.stack[:len(m.stack)-1]
		}
		return m, nil
	}

	// Reemplazar executing por result
	r := screens.NewResultScreen(msg.taskName, msg.exec)
	m.setTop(r)
	return m, r.Init()
}

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

// currentKeyMap combina las teclas globales + las del top para el help.
func (m Model) currentKeyMap() keys.KeyMap {
	ids := append([]string{}, m.top().Keys()...)
	return keys.NewKeyMap(m.deps.Keys, ids)
}

func (m Model) View() string {
	if m.execErr != nil {
		return m.renderError()
	}
	return m.top().View()
}

func (m Model) renderError() string {
	return fmt.Sprintf("\n\n  ✗ Error: %v\n\n  Pulsa cualquier tecla para volver\n", m.execErr)
}

// executionFinishedMsg transporta el resultado de la ejecución.
type executionFinishedMsg struct {
	taskName string
	exec     *execution.Execution
	err      error
}

// --- helpers de push/pop/replace ---

func push(s screens.ScreenI) tea.Cmd {
	return func() tea.Msg { return pushMsg{screen: s} }
}
func pop() tea.Cmd { return func() tea.Msg { return popMsg{} } }
func replace(s screens.ScreenI) tea.Cmd {
	return func() tea.Msg { return replaceMsg{screen: s} }
}

type pushMsg struct{ screen screens.ScreenI }
type popMsg struct{}
type replaceMsg struct{ screen screens.ScreenI }
