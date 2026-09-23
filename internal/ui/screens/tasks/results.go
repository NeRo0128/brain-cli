package tasks

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"

	"github.com/NeRo0128/brain-cli/internal/core/execution"
	"github.com/NeRo0128/brain-cli/internal/ui/components/states"
	"github.com/NeRo0128/brain-cli/internal/ui/keys"
	"github.com/NeRo0128/brain-cli/internal/ui/screens"
	"github.com/NeRo0128/brain-cli/internal/ui/styles"
)

type ResultScreen struct {
	exec     *execution.Execution
	taskName string

	viewport viewport.Model
	ready    bool
	styles   *styles.Styles
}

func NewResultScreen(taskName string, exec *execution.Execution, s *styles.Styles) ResultScreen {
	return ResultScreen{
		exec:     exec,
		taskName: taskName,
		styles:   s,
	}
}

func (m ResultScreen) Init() tea.Cmd {
	return func() tea.Msg { return tea.RequestWindowSize() }
}

func (m ResultScreen) Keys() []string {
	return []string{
		keys.NavBack,
		keys.ActionRerun,
		keys.ViewHelp,
	}
}

func (m ResultScreen) Update(msg tea.Msg) (screens.ScreenI, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		vpW := msg.Width - 2
		vpH := msg.Height - 8
		if vpH < 3 {
			vpH = 3
		}
		m.viewport = viewport.New(viewport.WithWidth(vpW), viewport.WithHeight(vpH))
		m.viewport.SetContent(m.renderOutput())
		m.ready = true
		return m, nil

	case screens.ActionMsg:
		switch msg.ID {
		case keys.NavBack:
			return m, screens.Back()
		case keys.ViewHelp:
			return m, screens.OpenHelp()
		case keys.ActionRerun:
			return m, screens.ExecuteTask(m.exec.TaskID, m.taskName)
		}
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m ResultScreen) View() tea.View {
	if !m.ready {
		return tea.NewView(states.Loading(m.styles, "resultado"))
	}

	status := renderStatus(m.exec, m.styles)
	meta := " · " + styles.HumanDuration(m.exec.Duration()) +
		" · exit " + formatExitCode(m.exec.ExitCode)

	var b strings.Builder
	b.WriteString("  ")
	b.WriteString(status)
	b.WriteString(m.styles.Subtitle.Render(meta))
	b.WriteString("\n\n")
	b.WriteString(m.viewport.View())
	return tea.NewView(b.String())
}

func (m ResultScreen) renderOutput() string {
	out := m.exec.Output
	if out == "" {
		out = "(sin output)"
	}
	return strings.TrimRight(out, "\n")
}

func renderStatus(e *execution.Execution, s *styles.Styles) string {
	switch e.Status {
	case execution.StatusCompleted:
		return s.ColoredIcons.Success() + " " + s.SuccessStyle.Render("completado")
	case execution.StatusFailed:
		return s.ColoredIcons.Failed() + " " + s.ErrorStyle.Render("falló")
	case execution.StatusCancelled:
		return s.ColoredIcons.Cancelled() + " " + s.WarningStyle.Render("cancelado")
	default:
		return s.Subtitle.Render(string(e.Status))
	}
}

func formatExitCode(code *int) string {
	if code == nil {
		return "-"
	}
	return fmt.Sprintf("%d", *code)
}
