package ui

import (
	"context"
	"time"

	tea "charm.land/bubbletea/v2"
	toasts "github.com/NeRo0128/brain-cli/internal/ui/components/toast"
	"github.com/NeRo0128/brain-cli/internal/ui/icons"
	"github.com/NeRo0128/brain-cli/internal/ui/keys"
	"github.com/NeRo0128/brain-cli/internal/ui/screens"
	"github.com/NeRo0128/brain-cli/internal/ui/screens/system"
	"github.com/NeRo0128/brain-cli/internal/ui/screens/tasks"
	"github.com/NeRo0128/brain-cli/internal/ui/screens/tools"
	"github.com/NeRo0128/brain-cli/internal/ui/styles"
	"github.com/NeRo0128/brain-cli/internal/ui/theme"
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
		return m, tea.Batch(toastCmd, func() tea.Msg { return tea.RequestWindowSize() })

	case popMsg:
		if len(m.stack) > 1 {
			m.stack = m.stack[:len(m.stack)-1]
		}
		return m, toastCmd

	case replaceMsg:
		m.setTop(msg.screen)
		return m, tea.Batch(toastCmd, func() tea.Msg { return tea.RequestWindowSize() })
	}

	// --- 1. Navegación desde screens ---
	switch msg := msg.(type) {
	case screens.BackMsg:
		if len(m.stack) > 1 {
			m.stack = m.stack[:len(m.stack)-1]
		}
		return m, toastCmd

	case screens.OpenDetailMsg:
		s := msg.Styles
		if s == nil {
			s = m.deps.Styles
		}
		d := tasks.NewDetailScreen(
			msg.Task,
			m.deps.TaskRepo,
			m.deps.ToolRepo,
			m.deps.Log,
			s,
		)
		return m, tea.Batch(toastCmd, push(d), d.Init())

	case screens.OpenHistoryMsg:
		s := msg.Styles
		if s == nil {
			s = m.deps.Styles
		}
		h := tasks.NewHistoryScreen(m.deps.ExecRepo, m.deps.TaskRepo, s)
		return m, tea.Batch(toastCmd, push(h), h.Init())

	case screens.OpenResultMsg:
		r := tasks.NewResultScreen(msg.TaskName, msg.Exec, m.deps.Styles)
		return m, tea.Batch(toastCmd, replace(r))

	case screens.OpenFormMsg:
		s := msg.Styles
		if s == nil {
			s = m.deps.Styles
		}
		f := tools.NewFormScreen(
			msg.Task,
			m.deps.Manager,
			m.deps.ToolRepo,
			m.deps.Interpreter,
			m.deps.Log,
			s,
		)
		return m, tea.Batch(toastCmd, push(f), f.Init())

	case screens.OpenToolPickerMsg:
		s := msg.Styles
		if s == nil {
			s = m.deps.Styles
		}
		p := tools.NewToolPickerScreen(
			m.deps.ToolRepo,
			msg.CurrentID,
			msg.FilterType,
			s,
		)
		return m, tea.Batch(toastCmd, push(p), p.Init())

	case screens.OpenAuthMsg:
		s := msg.Styles
		if s == nil {
			s = m.deps.Styles
		}
		var a screens.ScreenI
		if m.deps.AuthManager != nil {
			a = system.NewAuthScreen(m.deps.AuthManager, m.deps.Log, s)
		} else {
			a = system.NewAuthScreen(nil, m.deps.Log, s)
		}
		return m, tea.Batch(toastCmd, push(a), a.Init())

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
		s := msg.Styles
		if s == nil {
			s = m.deps.Styles
		}
		c := system.NewConfirm(msg.Title, msg.Message, msg.Action, s)
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

	// case screens.OpenToolFormMsg:
	// 	s := msg.Styles
	// 	if s == nil {
	// 		s = m.deps.Styles
	// 	}
	// 	f := tasks.NewToolFormScreen(msg.Tool, m.deps.ToolManager, m.deps.Log, s)
	// 	return m, tea.Batch(toastCmd, push(f), f.Init())

	case screens.ToolFormSavedMsg:
		if len(m.stack) > 1 {
			m.stack = m.stack[:len(m.stack)-1]
		}
		newTop, cmd := m.top().Update(screens.ReloadMsg{})
		m.setTop(newTop)
		toastCmd2 := toasts.ShowSuccess("Tool guardado")
		created := screens.ToolCreatedMsg{ToolID: msg.Tool.ID}
		newTop2, cmd2 := m.top().Update(created)
		m.setTop(newTop2)
		return m, tea.Batch(toastCmd, cmd, cmd2, toastCmd2)

	case screens.DeleteToolMsg:
		mgr := m.deps.ToolManager
		return m, tea.Batch(toastCmd, func() tea.Msg {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			err := mgr.Delete(ctx, msg.ToolID)
			return screens.ToolDeletedMsg{ToolID: msg.ToolID, Err: err}
		})

	case screens.ToolDeletedMsg:
		if msg.Err != nil {
			return m, tea.Batch(toastCmd,
				toasts.ShowError("Error al borrar tool: "+msg.Err.Error()))
		}
		newTop, cmd := m.top().Update(screens.ReloadMsg{})
		m.setTop(newTop)
		return m, tea.Batch(toastCmd, cmd, toasts.ShowSuccess("Tool borrado"))

	case screens.ExecuteTaskMsg:
		return m.startExecution(msg, toastCmd)

	case screens.OpenSettingsMsg:
		s := msg.Styles
		if s == nil {
			s = m.deps.Styles
		}
		st := system.NewSettingsScreen(m.deps.SettingsManager, m.deps.AuthManager, m.deps.Log, s)
		return m, tea.Batch(toastCmd, push(st), st.Init())

	case screens.SettingsChangedMsg:
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		newCfg, err := m.deps.SettingsManager.Resolve(ctx)
		if err != nil {
			return m, tea.Batch(toastCmd,
				toasts.ShowError("Recargando config: "+err.Error()))
		}

		*m.deps.Cfg = *newCfg

		newTheme := theme.Get(newCfg.UI.Theme)
		newStyles := styles.New(
			newTheme,
			m.styles.Dark,
			icons.Get(newCfg.UI.Icons),
		)
		*m.deps.Styles = newStyles

		m.toast.SetPalette(newTheme.Resolve(newStyles.Dark))
		m.toast.SetIcons(newStyles.Icons)

		return m, tea.Batch(toastCmd,
			toasts.ShowSuccess("Ajustes aplicados"),
			func() tea.Msg { return tea.RequestWindowSize() })
	case executionFinishedMsg:
		return m.handleExecutionFinished(msg, toastCmd)
	}

	// --- 2. Teclas globales ---
	if keyMsg, ok := msg.(tea.KeyPressMsg); ok {
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
				if ex, ok := m.top().(tasks.ExecutingScreen); ok {
					m.setTop(ex.MarkCanceling())
				}
				return m, toastCmd
			}
		case "ctrl+t":
			if !m.executing {
				nextName := theme.Next(m.styles.Theme.Name)
				next := theme.Get(nextName)

				s := styles.New(next, m.styles.Dark, m.styles.Icons)
				*m.deps.Styles = s
				m.toast.SetPalette(next.Resolve(s.Dark))
				m.toast.SetIcons(s.Icons) // [NUEVO]

				return m, tea.Batch(toastCmd,
					toasts.ShowSuccess("Tema: "+nextName),
					func() tea.Msg { return tea.RequestWindowSize() })
			}
			return m, toastCmd
		}
	}

	// --- 3. Traducción KeyMsg → ActionMsg ---
	if keyMsg, ok := msg.(tea.KeyPressMsg); ok {
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
