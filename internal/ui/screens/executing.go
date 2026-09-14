package screens

import (
	"strings"
	"time"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"

	"github.com/NeRo0128/brain-cli/internal/ui/components/progress"
	"github.com/NeRo0128/brain-cli/internal/ui/keys"
	"github.com/NeRo0128/brain-cli/internal/ui/styles"
)

type executingTickMsg struct{}

type ExecutingScreen struct {
	spinner   spinner.Model
	taskName  string
	startedAt time.Time
	canceling bool

	bar       progress.Model
	barOffset int
	styles    *styles.Styles
}

func NewExecutingScreen(taskName string, s *styles.Styles) ExecutingScreen {
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = s.SpinnerStyle

	return ExecutingScreen{
		spinner:   sp,
		taskName:  taskName,
		startedAt: time.Now(),
		bar:       progress.New(40, s),
		styles:    s,
	}
}

func (m ExecutingScreen) Init() tea.Cmd {
	return tea.Batch(func() tea.Msg { return m.spinner.Tick() }, m.tick(), func() tea.Msg { return tea.RequestWindowSize() })
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

func (m ExecutingScreen) MarkCanceling() ExecutingScreen {
	m.canceling = true
	return m
}

func (m ExecutingScreen) tick() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return executingTickMsg{}
	})
}

func (m ExecutingScreen) View() tea.View {
	elapsed := time.Since(m.startedAt).Round(time.Second)

	var status string
	if m.canceling {
		status = m.styles.WarningStyle.Render("⊘ Cancelando")
	} else {
		status = m.spinner.View() + " Ejecutando"
	}

	var b strings.Builder
	b.WriteString("  " + status + "  ")
	b.WriteString(m.styles.TimerStyle.Render(elapsed.String()))
	b.WriteString("\n\n")

	b.WriteString("  " + m.styles.Key.Render("Tarea:") + " " + m.taskName + "\n\n")

	b.WriteString("  " + m.bar.WithOffset(m.barOffset).View() + "\n")

	if m.canceling {
		b.WriteString("\n  " +
			m.styles.Subtitle.Render("Esperando a que termine el proceso...") +
			"\n")
	}

	return tea.NewView(b.String())
}
