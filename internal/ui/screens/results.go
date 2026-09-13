package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/NeRo0128/brain-cli/internal/core/execution"
	"github.com/NeRo0128/brain-cli/internal/ui/components/states"
	"github.com/NeRo0128/brain-cli/internal/ui/keys"
	"github.com/NeRo0128/brain-cli/internal/ui/styles"
)

// ResultScreen muestra el resultado de una ejecución.
//
// [ACTUALIZADO] migrado a viewport.Model: scroll real, wrapping
// automático, PageUp/PageDown nativos.
type ResultScreen struct {
	exec     *execution.Execution
	taskName string

	viewport viewport.Model
	ready    bool
}

func NewResultScreen(taskName string, exec *execution.Execution) ResultScreen {
	return ResultScreen{
		exec:     exec,
		taskName: taskName,
	}
}

func (m ResultScreen) Init() tea.Cmd {
	return tea.WindowSize()
}

func (m ResultScreen) Keys() []string {
	return []string{
		keys.NavBack,
		keys.ActionRerun,
		keys.ViewHelp,
	}
}

func (m ResultScreen) Update(msg tea.Msg) (ScreenI, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		vpW := msg.Width - 2
		vpH := msg.Height - 8
		if vpH < 3 {
			vpH = 3
		}
		m.viewport = viewport.New(vpW, vpH)
		m.viewport.SetContent(m.renderOutput())
		m.ready = true
		return m, nil

	case ActionMsg:
		switch msg.ID {
		case keys.NavBack:
			return m, Back()
		case keys.ViewHelp:
			return m, OpenHelp()
		case keys.ActionRerun:
			return m, ExecuteTask(m.exec.TaskID, m.taskName)
		}
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m ResultScreen) View() string {
	if !m.ready {
		return states.Loading("resultado")
	}

	status := renderStatus(m.exec)
	meta := " · " + styles.HumanDuration(m.exec.Duration()) +
		" · exit " + formatExitCode(m.exec.ExitCode)

	var b strings.Builder
	b.WriteString("  " + status + styles.Subtitle.Render(meta))
	b.WriteString("\n\n")
	b.WriteString(m.viewport.View())
	return b.String()
}

func (m ResultScreen) renderOutput() string {
	out := m.exec.Output
	if out == "" {
		out = "(sin output)"
	}
	return strings.TrimRight(out, "\n")
}

func renderStatus(e *execution.Execution) string {
	switch e.Status {
	case execution.StatusCompleted:
		return styles.SuccessStyle.Render("✓ completado")
	case execution.StatusFailed:
		return styles.ErrorStyle.Render("✗ falló")
	case execution.StatusCancelled:
		return styles.WarningStyle.Render("⊘ cancelado")
	default:
		return styles.Subtitle.Render(string(e.Status))
	}
}

func formatExitCode(code *int) string {
	if code == nil {
		return "-"
	}
	return fmt.Sprintf("%d", *code)
}
