package screens

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/NeRo0128/brain-cli/internal/ui/styles"
)

// executingTickMsg dispara re-render para actualizar el cronómetro.
type executingTickMsg struct{}

// ExecutingScreen muestra el progreso de una ejecución en curso.
type ExecutingScreen struct {
	spinner   spinner.Model
	taskName  string
	startedAt time.Time
	canceling bool
}

// NewExecutingScreen construye la pantalla.
func NewExecutingScreen(taskName string) ExecutingScreen {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = styles.SpinnerStyle

	return ExecutingScreen{
		spinner:   s,
		taskName:  taskName,
		startedAt: time.Now(),
	}
}

// Init arranca el spinner y el tick del cronómetro (cada 1s).
func (m ExecutingScreen) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.tick())
}

func (m ExecutingScreen) Update(msg tea.Msg) (ExecutingScreen, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case executingTickMsg:
		return m, m.tick()
	}
	return m, nil
}

// MarkCanceling indica que el usuario pidió cancelar.
// La pantalla sigue corriendo hasta que llega executionFinishedMsg.
func (m ExecutingScreen) MarkCanceling() ExecutingScreen {
	m.canceling = true
	return m
}

func (m ExecutingScreen) tick() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return executingTickMsg{}
	})
}

// View renderiza la pantalla.
func (m ExecutingScreen) View() string {
	elapsed := time.Since(m.startedAt).Round(time.Second)

	var b string
	if m.canceling {
		b = styles.WarningStyle.Render("⊘ Cancelando...")
	} else {
		b = m.spinner.View() + " Ejecutando"
	}

	header := styles.Title.Render(b) + "\n\n"
	header += "  " + styles.Key.Render("Tarea:") + " " + m.taskName + "\n"
	header += "  " + styles.TimerStyle.Render("Tiempo: "+elapsed.String()) + "\n\n"

	var help string
	if m.canceling {
		help = styles.Subtitle.Render("Esperando a que termine el proceso...")
	} else {
		help = styles.Help.Render(
			styles.Key.Render("Esc") + " cancelar  ·  " +
				styles.Key.Render("Ctrl+C") + " salir",
		)
	}

	return header + help
}
