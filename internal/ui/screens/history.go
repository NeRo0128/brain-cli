package screens

import (
	"context"
	"time"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/NeRo0128/brain-cli/internal/core/execution"
	coretask "github.com/NeRo0128/brain-cli/internal/core/task"
	uilist "github.com/NeRo0128/brain-cli/internal/ui/components/list"
	"github.com/NeRo0128/brain-cli/internal/ui/components/states"
	"github.com/NeRo0128/brain-cli/internal/ui/keys"
	"github.com/NeRo0128/brain-cli/internal/ui/styles"
)

type historyLoadedMsg struct {
	executions []*execution.Execution
	taskNames  map[string]string
	err        error
}

type executionItem struct {
	exec     *execution.Execution
	taskName string
}

func (i executionItem) Title() string       { return i.taskName }
func (i executionItem) Description() string { return "" }
func (i executionItem) FilterValue() string {
	return i.taskName + " " + string(i.exec.Status)
}

type HistoryScreen struct {
	list    list.Model
	loading bool
	err     error

	execRepo execution.Repository
	taskRepo coretask.Repository

	names  map[string]string
	styles *styles.Styles
}

func NewHistoryScreen(
	execRepo execution.Repository,
	taskRepo coretask.Repository,
	s *styles.Styles,
) HistoryScreen {
	l := list.New(nil, uilist.New(s), 80, 20)
	l.Title = "Historial"
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = s.Title
	l.Styles.HelpStyle = s.Help

	return HistoryScreen{
		list:     l,
		execRepo: execRepo,
		taskRepo: taskRepo,
		styles:   s,
		loading:  true,
		names:    make(map[string]string),
	}
}

func (m HistoryScreen) Init() tea.Cmd {
	execRepo, taskRepo := m.execRepo, m.taskRepo
	load := func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		tasks, err := taskRepo.List(ctx)
		if err != nil {
			return historyLoadedMsg{err: err}
		}
		names := make(map[string]string, len(tasks))
		for _, t := range tasks {
			names[t.ID] = t.Name
		}

		execs, err := execRepo.ListRecent(ctx, 100)
		if err != nil {
			return historyLoadedMsg{err: err}
		}
		return historyLoadedMsg{executions: execs, taskNames: names}
	}
	return tea.Batch(func() tea.Msg { return tea.RequestWindowSize() }, load)
}

func (m HistoryScreen) Keys() []string {
	return []string{keys.NavConfirm, keys.NavBack, keys.ViewHelp}
}

func (m HistoryScreen) Update(msg tea.Msg) (ScreenI, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width-2, msg.Height-6)

	case historyLoadedMsg:
		m.loading = false
		m.err = msg.err
		if msg.err == nil {
			m.setItems(msg.executions, msg.taskNames)
		}
		return m, nil

	case ReloadMsg:
		m.loading = true
		return m, m.Init()

	case ActionMsg:
		return m.handleAction(msg)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m HistoryScreen) handleAction(msg ActionMsg) (ScreenI, tea.Cmd) {
	switch msg.ID {
	case keys.NavBack:
		return m, Back()
	case keys.ViewHelp:
		return m, OpenHelp()
	}
	return m, nil
}

func (m *HistoryScreen) setItems(execs []*execution.Execution, names map[string]string) {
	items := make([]list.Item, len(execs))
	for i, e := range execs {
		name := names[e.TaskID]
		if name == "" {
			name = e.TaskID
		}
		items[i] = executionItem{exec: e, taskName: name}
	}
	m.list.SetItems(items)
}

func (m HistoryScreen) View() tea.View {
	switch {
	case m.loading:
		return tea.NewView(states.Loading(m.styles, "historial"))
	case m.err != nil:
		return tea.NewView(states.Error(m.styles, m.err))
	case len(m.list.Items()) == 0:
		return tea.NewView(states.Empty(m.styles, "Sin ejecuciones registradas", "Ejecuta una task para empezar"))
	}
	return tea.NewView(m.list.View())
}
