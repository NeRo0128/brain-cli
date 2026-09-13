package screens

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/NeRo0128/brain-cli/internal/ui/components/progress"
	"github.com/NeRo0128/brain-cli/internal/ui/keys"
	"github.com/NeRo0128/brain-cli/internal/ui/styles"
)

// executingTickMsg dispara re-render del cronómetro (cada 1s).
type executingTickMsg struct{}

// ExecutingScreen muestra el progreso de una ejecución en curso.
//
// [ACTUALIZADO] añade barra de progreso animada (indeterminada)
// sincronizada con el spinner.
type ExecutingScreen struct {
	spinner   spinner.Model
	taskName  string
	startedAt time.Time
	canceling bool

	bar       progress.Model
	barOffset int
}

func NewExecutingScreen(taskName string) ExecutingScreen {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = styles.SpinnerStyle

	return ExecutingScreen{
		spinner:   s,
		taskName:  taskName,
		startedAt: time.Now(),
		bar:       progress.New(40),
	}
}

func (m ExecutingScreen) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.tick(), tea.WindowSize())
}

func (m ExecutingScreen) Keys() []string {
	return []string{keys.ActionCancel}
}

func (m ExecutingScreen) Update(msg tea.Msg) (ScreenI, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		m.barOffset++
		return m, cmd

	case executingTickMsg:
		return m, m.tick()

	case tea.WindowSizeMsg:
		w := msg.Width - 6
		if w < 20 {
			w = 20
		}
		if w > 70 {
			w = 70
		}
		m.bar = m.bar.WithWidth(w)
		return m, nil
	}
	return m, nil
}

// MarkCanceling indica que el usuario pidió cancelar.
func (m ExecutingScreen) MarkCanceling() ExecutingScreen {
	m.canceling = true
	return m
}

func (m ExecutingScreen) tick() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return executingTickMsg{}
	})
}

// View renderiza SOLO el contenido del medio.
func (m ExecutingScreen) View() string {
	elapsed := time.Since(m.startedAt).Round(time.Second)

	var status string
	if m.canceling {
		status = styles.WarningStyle.Render("⊘ Cancelando")
	} else {
		status = m.spinner.View() + " Ejecutando"
	}

	var b strings.Builder
	b.WriteString("  " + status + "  ")
	b.WriteString(styles.TimerStyle.Render(elapsed.String()))
	b.WriteString("\n\n")

	b.WriteString("  " + styles.Key.Render("Tarea:") + " " + m.taskName + "\n\n")

	b.WriteString("  " + m.bar.WithOffset(m.barOffset).View() + "\n")

	if m.canceling {
		b.WriteString("\n  " +
			styles.Subtitle.Render("Esperando a que termine el proceso...") +
			"\n")
	}

	return b.String()
}
