package ui

import (
	toasts "github.com/NeRo0128/brain-cli/internal/ui/components/toast"
	"github.com/NeRo0128/brain-cli/internal/ui/keys"
	"github.com/NeRo0128/brain-cli/internal/ui/screens"
	tea "github.com/charmbracelet/bubbletea"
)

// Update procesa todos los mensajes de la aplicación.
//
// Orden:
//  0. Navegación interna (push/pop/replace)
//  1. Navegación desde screens (BackMsg, OpenXxxMsg)
//  2. Teclas globales (Ctrl+C, q, Esc durante ejecución)
//  3. Traducción KeyMsg → ActionMsg
//  4. Delegación al top del stack
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var toastCmd tea.Cmd
	m.toast, toastCmd = m.toast.Update(msg)

	if size, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = size.Width
		m.height = size.Height
	}

	// --- 0. Navegación interna ---
	switch msg := msg.(type) {
	case pushMsg:
		m.stack = append(m.stack, msg.screen)
		return m, tea.Batch(toastCmd, tea.WindowSize())

	case popMsg:
		if len(m.stack) > 1 {
			m.stack = m.stack[:len(m.stack)-1]
		}
		return m, toastCmd

	case replaceMsg:
		m.setTop(msg.screen)
		return m, tea.Batch(toastCmd, tea.WindowSize())
	}

	// --- 1. Navegación desde screens ---
	switch msg := msg.(type) {
	case screens.BackMsg:
		if len(m.stack) > 1 {
			m.stack = m.stack[:len(m.stack)-1]
		}
		return m, toastCmd

	case screens.OpenDetailMsg:
		d := screens.NewDetailScreen(msg.Task, m.deps.ToolRepo, m.deps.Log)
		return m, tea.Batch(toastCmd, push(d), d.Init())

	case screens.OpenHistoryMsg:
		h := screens.NewHistoryScreen(m.deps.ExecRepo, m.deps.TaskRepo)
		return m, tea.Batch(toastCmd, push(h), h.Init())

	case screens.OpenResultMsg:
		r := screens.NewResultScreen(msg.TaskName, msg.Exec)
		return m, tea.Batch(toastCmd, replace(r))

	case screens.OpenFormMsg:
		f := screens.NewFormScreen(
			msg.Task,
			m.deps.Manager,
			m.deps.ToolRepo,
			m.deps.Interpreter,
			m.deps.Log,
		)
		return m, tea.Batch(toastCmd, push(f), f.Init())

	case screens.OpenToolPickerMsg:
		p := screens.NewToolPickerScreen(
			m.deps.ToolRepo,
			msg.CurrentID,
			msg.FilterType,
		)
		return m, tea.Batch(toastCmd, push(p), p.Init())

	case screens.ToolSelectedMsg:
		if len(m.stack) > 1 {
			m.stack = m.stack[:len(m.stack)-1]
		}
		top := m.top()
		newTop, cmd := top.Update(msg)
		m.setTop(newTop)
		return m, tea.Batch(toastCmd, cmd)

	case screens.FormSavedMsg:
		if len(m.stack) > 1 {
			m.stack = m.stack[:len(m.stack)-1]
		}
		newTop, cmd := m.top().Update(screens.ReloadMsg{})
		m.setTop(newTop)
		return m, tea.Batch(toastCmd, cmd, toasts.ShowSuccess("Task guardada"))

	case screens.OpenConfirmMsg:
		c := screens.NewConfirm(msg.Title, msg.Message, msg.Action)
		return m, tea.Batch(toastCmd, push(c), c.Init())

	case screens.ConfirmYesMsg:
		if len(m.stack) > 1 {
			m.stack = m.stack[:len(m.stack)-1]
		}
		return m.dispatchConfirmed(msg.Action, toastCmd)

	case screens.TaskDeletedMsg:
		if msg.Err != nil {
			return m, tea.Batch(toastCmd,
				toasts.ShowError("Error al borrar: "+msg.Err.Error()))
		}
		newTop, cmd := m.top().Update(screens.ReloadMsg{})
		m.setTop(newTop)
		return m, tea.Batch(toastCmd, cmd, toasts.ShowSuccess("Task borrada"))

	case screens.ExecuteTaskMsg:
		return m.startExecution(msg, toastCmd)

	case executionFinishedMsg:
		return m.handleExecutionFinished(msg, toastCmd)
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
			return m, toastCmd

		case "esc":
			if m.executing && m.cancelExec != nil {
				m.cancelExec()
				m.cancelExec = nil
				m.executing = false
				if ex, ok := m.top().(screens.ExecutingScreen); ok {
					m.setTop(ex.MarkCanceling())
				}
				return m, toastCmd
			}
		}
	}

	// --- 3. Traducción KeyMsg → ActionMsg ---
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		top := m.top()
		for _, id := range top.Keys() {
			if m.deps.Keys.Matches(id, keyMsg.String()) {
				msg = screens.ActionMsg{ID: id}
				break
			}
		}
	}

	// --- 4. Delegación al top ---
	top := m.top()
	newTop, cmd := top.Update(msg)
	m.setTop(newTop)
	return m, tea.Batch(toastCmd, cmd)
}

// currentKeyMap combina las teclas globales + las del top para el help.
func (m Model) currentKeyMap() keys.KeyMap {
	ids := append([]string{}, m.top().Keys()...)
	return keys.NewKeyMap(m.deps.Keys, ids)
}
