package screens

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/NeRo0128/brain-cli/internal/core/execution"
	"github.com/NeRo0128/brain-cli/internal/ui/styles"
)

// ResultScreen muestra el resultado de una ejecución.
type ResultScreen struct {
	exec     *execution.Execution
	taskName string
	scroll   int
	lines    []string
}

// NewResultScreen construye la pantalla a partir de una Execution.
func NewResultScreen(taskName string, exec *execution.Execution) ResultScreen {
	output := exec.Output
	if output == "" {
		output = "(sin output)"
	}
	return ResultScreen{
		exec:     exec,
		taskName: taskName,
		lines:    strings.Split(output, "\n"),
	}
}

func (m ResultScreen) Init() tea.Cmd { return nil }

func (m ResultScreen) Update(msg tea.Msg) (ResultScreen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.scroll > 0 {
				m.scroll--
			}
		case "down", "j":
			if m.scroll < len(m.lines)-1 {
				m.scroll++
			}
		case "g":
			m.scroll = 0
		case "G":
			m.scroll = len(m.lines) - 1
		}
	}
	return m, nil
}

func (m ResultScreen) View() string {
	var b strings.Builder

	// Encabezado con estado
	status := renderStatus(m.exec)
	b.WriteString(styles.Title.Render("🧠 " + m.taskName))
	b.WriteString("  ")
	b.WriteString(status)
	b.WriteString("\n")

	// Metadata
	meta := fmt.Sprintf("duración: %s  ·  exit: %s",
		m.exec.Duration().Round(10*1e6), formatExitCode(m.exec.ExitCode))
	b.WriteString(styles.Subtitle.Render(meta))
	b.WriteString("\n\n")

	// Output con ventana deslizante
	maxLines := 15
	start := m.scroll
	end := start + maxLines
	if end > len(m.lines) {
		end = len(m.lines)
	}
	for _, line := range m.lines[start:end] {
		b.WriteString(line)
		b.WriteString("\n")
	}

	if len(m.lines) > maxLines {
		b.WriteString(styles.Subtitle.Render(
			fmt.Sprintf("\n[%d/%d líneas]", end, len(m.lines))))
	}

	b.WriteString("\n\n")
	b.WriteString(styles.Help.Render(
		styles.Key.Render("↑↓") + " scroll  ·  " +
			styles.Key.Render("Esc") + " volver  ·  " +
			styles.Key.Render("q") + " salir",
	))

	return b.String()
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
