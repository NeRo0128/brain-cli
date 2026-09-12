package screens

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/NeRo0128/brain-cli/internal/core/execution"
	"github.com/NeRo0128/brain-cli/internal/core/task"
	"github.com/NeRo0128/brain-cli/internal/ui/styles"
)

// historyLoadedMsg trae las ejecuciones + nombres de tasks.
type historyLoadedMsg struct {
	executions []*execution.Execution
	taskNames  map[string]string
	err        error
}

// executionItem adapta *execution.Execution a list.Item.
type executionItem struct {
	exec     *execution.Execution
	taskName string
}

func (i executionItem) Title() string {
	icon, _ := statusIcon(i.exec.Status)
	return icon + " " + i.taskName
}

func (i executionItem) Description() string {
	dur := i.exec.Duration().Round(time.Millisecond)
	ts := i.exec.StartedAt.Local().Format("2006-01-02 15:04:05")
	meta := fmt.Sprintf("%s · %s · %s", ts, dur, i.exec.Status)
	if i.exec.ExitCode != nil && *i.exec.ExitCode > 0 {
		meta += fmt.Sprintf(" (exit %d)", *i.exec.ExitCode)
	}
	return meta
}

func (i executionItem) FilterValue() string {
	return i.taskName + " " + string(i.exec.Status)
}

// --- delegate ---

type executionDelegate struct{}

func (d executionDelegate) Height() int                             { return 2 }
func (d executionDelegate) Spacing() int                            { return 1 }
func (d executionDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d executionDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	it, ok := item.(executionItem)
	if !ok {
		return
	}

	selected := index == m.Index()
	cursor := "  "
	titleStyle := lipgloss.NewStyle().Foreground(styles.Text)
	if selected {
		cursor = "▶ "
		titleStyle = titleStyle.Bold(true).Foreground(styles.Primary)
	}

	icon, iconStyle := statusIcon(it.exec.Status)
	metaStyle := lipgloss.NewStyle().Foreground(styles.Muted)

	fmt.Fprintf(w, "%s%s %s\n", cursor, iconStyle.Render(icon), titleStyle.Render(it.taskName))
	fmt.Fprintf(w, "   %s", metaStyle.Render(it.Description()))
}

// statusIcon devuelve el símbolo y su estilo según el estado.
func statusIcon(s execution.Status) (string, lipgloss.Style) {
	switch s {
	case execution.StatusCompleted:
		return "✓", styles.SuccessStyle
	case execution.StatusFailed:
		return "✗", styles.ErrorStyle
	case execution.StatusCancelled:
		return "⊘", styles.WarningStyle
	case execution.StatusRunning, execution.StatusPending:
		return "◐", styles.Subtitle
	}
	return "?", styles.Subtitle
}

// --- pantalla ---

// HistoryScreen muestra el historial de ejecuciones.
type HistoryScreen struct {
	list    list.Model
	loading bool
	err     error

	execRepo execution.Repository
	taskRepo task.Repository

	names map[string]string
}

// NewHistoryScreen construye la pantalla.
func NewHistoryScreen(execRepo execution.Repository, taskRepo task.Repository) HistoryScreen {
	l := list.New(nil, executionDelegate{}, 80, 20)
	l.Title = "Historial"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = styles.Title
	l.Styles.HelpStyle = styles.Help

	return HistoryScreen{
		list:     l,
		execRepo: execRepo,
		taskRepo: taskRepo,
		loading:  true,
		names:    make(map[string]string),
	}
}

// Init carga ejecuciones + nombres de tasks (batch).
func (m HistoryScreen) Init() tea.Cmd {
	execRepo, taskRepo := m.execRepo, m.taskRepo
	load := func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		execs, err := execRepo.ListRecent(ctx, 100)
		if err != nil {
			return historyLoadedMsg{err: err}
		}
		tasks, err := taskRepo.List(ctx)
		if err != nil {
			return historyLoadedMsg{err: err}
		}
		names := make(map[string]string, len(tasks))
		for _, t := range tasks {
			names[t.ID] = t.Name
		}
		return historyLoadedMsg{executions: execs, taskNames: names}
	}
	return tea.Batch(tea.WindowSize(), load)
}

// Update maneja mensajes.
func (m HistoryScreen) Update(msg tea.Msg) (HistoryScreen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height-6)
		return m, nil

	case historyLoadedMsg:
		m.loading = false
		m.err = msg.err
		if msg.err == nil {
			m.names = msg.taskNames
			m.setItems(msg.executions)
		}
		return m, nil
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// SelectedExecution devuelve la ejecución seleccionada + nombre de la task.
func (m HistoryScreen) SelectedExecution() (*execution.Execution, string) {
	it, ok := m.list.SelectedItem().(executionItem)
	if !ok {
		return nil, ""
	}
	return it.exec, it.taskName
}

func (m *HistoryScreen) setItems(execs []*execution.Execution) {
	items := make([]list.Item, len(execs))
	for i, e := range execs {
		name := m.names[e.TaskID]
		if name == "" {
			name = e.TaskID
		}
		items[i] = executionItem{exec: e, taskName: name}
	}
	m.list.SetItems(items)
}

// View renderiza la pantalla.
func (m HistoryScreen) View() string {
	header := styles.Title.Render("📜 Historial") + "\n\n"

	switch {
	case m.loading:
		return header + styles.Subtitle.Render("Cargando ejecuciones...") + "\n"
	case m.err != nil:
		return header + styles.ErrorStyle.Render("Error: ") + m.err.Error() + "\n"
	case len(m.list.Items()) == 0:
		return header + styles.Subtitle.Render("(sin ejecuciones registradas)") + "\n"
	}

	footer := "\n" + styles.Help.Render(
		styles.Key.Render("↑↓")+" navegar  ·  "+
			styles.Key.Render("Enter")+" ver  ·  "+
			styles.Key.Render("/")+" filtrar  ·  "+
			styles.Key.Render("h/Esc")+" volver  ·  "+
			styles.Key.Render("q")+" salir",
	)

	return header + m.list.View() + footer
}
