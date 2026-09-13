package screens

import (
	"context"
	"time"

	"github.com/NeRo0128/brain-cli/internal/core/execution"
	coretask "github.com/NeRo0128/brain-cli/internal/core/task"
	uilist "github.com/NeRo0128/brain-cli/internal/ui/components/list"
	"github.com/NeRo0128/brain-cli/internal/ui/components/states"
	"github.com/NeRo0128/brain-cli/internal/ui/keys"
	"github.com/NeRo0128/brain-cli/internal/ui/styles"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type historyLoadedMsg struct {
	executions []*execution.Execution
	taskNames  map[string]string
	err        error
}

// --- Item ---

type executionItem struct {
	exec     *execution.Execution
	taskName string
}

func (i executionItem) Title() string       { return i.taskName }
func (i executionItem) Description() string { return "" }
func (i executionItem) FilterValue() string {
	return i.taskName + " " + string(i.exec.Status)
}

func (i executionItem) Row() uilist.Row {
	icon, color := statusIconAndColor(i.exec.Status)
	return uilist.Row{
		Prefix:      icon,
		PrefixColor: color,
		Title:       i.taskName,
		Badges:      []uilist.Badge{uilist.StatusBadge(string(i.exec.Status))},
		Subtitle:    styles.HumanTime(i.exec.StartedAt),
		Meta:        styles.HumanDuration(i.exec.Duration()),
	}
}

// statusIconAndColor centraliza iconos de estado.
func statusIconAndColor(s execution.Status) (string, lipgloss.TerminalColor) {
	switch s {
	case execution.StatusCompleted:
		return "✓", styles.Success
	case execution.StatusFailed:
		return "✗", styles.Error
	case execution.StatusCancelled:
		return "⊘", styles.Warning
	case execution.StatusRunning, execution.StatusPending:
		return "◐", styles.Primary
	}
	return "•", styles.Muted
}

// --- Pantalla ---

type HistoryScreen struct {
	list    list.Model
	loading bool
	err     error

	execRepo execution.Repository
	taskRepo coretask.Repository

	names map[string]string
}

func NewHistoryScreen(execRepo execution.Repository, taskRepo coretask.Repository) HistoryScreen {
	l := list.New(nil, uilist.New(), 80, 20)
	l.Title = "Historial"
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
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

func (m HistoryScreen) Keys() []string {
	return []string{keys.NavConfirm, keys.NavBack, keys.ViewHelp}
}

func (m HistoryScreen) Update(msg tea.Msg) (ScreenI, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width-2, msg.Height-6)
		return m, nil

	case historyLoadedMsg:
		m.loading = false
		m.err = msg.err
		if msg.err == nil {
			m.names = msg.taskNames
			m.setItems(msg.executions)
		}
		return m, nil

	case ActionMsg:
		return m.handleAction(msg)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m HistoryScreen) handleAction(msg ActionMsg) (ScreenI, tea.Cmd) {
	switch msg.ID {
	case keys.NavConfirm:
		exec, name := m.SelectedExecution()
		if exec == nil {
			return m, nil
		}
		return m, OpenResult(name, exec)
	case keys.NavBack:
		return m, Back()
	case keys.ViewHelp:
		return m, OpenHelp()
	}
	return m, nil
}

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

func (m HistoryScreen) View() string {
	switch {
	case m.loading:
		return states.Loading("historial")
	case m.err != nil:
		return states.Error(m.err)
	case len(m.list.Items()) == 0:
		return states.Empty("Sin ejecuciones registradas", "Ejecuta una task para empezar")
	}
	return m.list.View()
}
