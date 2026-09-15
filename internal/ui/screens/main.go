package screens

import (
	"context"
	"fmt"
	"time"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/NeRo0128/brain-cli/internal/core/task"
	coretask "github.com/NeRo0128/brain-cli/internal/core/task"
	uilist "github.com/NeRo0128/brain-cli/internal/ui/components/list"
	"github.com/NeRo0128/brain-cli/internal/ui/components/states"
	"github.com/NeRo0128/brain-cli/internal/ui/keys"
	"github.com/NeRo0128/brain-cli/internal/ui/styles"
)

type tasksLoadedMsg struct {
	tasks []*task.Task
	err   error
}

type taskItem struct {
	task *task.Task
}

func (i taskItem) Title() string       { return i.task.Name }
func (i taskItem) Description() string { return "" }
func (i taskItem) FilterValue() string { return i.task.Name + " " + i.task.ID }

func (i taskItem) Row() uilist.Row {
	prefix := uilist.PrefixNone
	if i.task.IsFavorite {
		prefix = uilist.PrefixFavorite
	}
	return uilist.Row{
		Prefix: prefix,
		Title:  i.task.Name,
		Badges: []uilist.Badge{
			uilist.TypeBadge(string(i.task.Type)),
			uilist.PriorityBadge(string(i.task.Priority)),
		},
		Subtitle: "[" + i.task.ID + "]",
	}
}

type MainScreen struct {
	width, height int
	loading       bool
	err           error
	list          list.Model
	repo          coretask.Repository
	styles        *styles.Styles
}

func NewMainScreen(repo coretask.Repository, s *styles.Styles) MainScreen {
	l := list.New(nil, uilist.New(s), 80, 20)
	l.Title = "Tareas"
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = s.Title
	l.Styles.HelpStyle = s.Help

	return MainScreen{
		list:    l,
		repo:    repo,
		styles:  s,
		loading: true,
	}
}

func (m MainScreen) Init() tea.Cmd {
	repo := m.repo
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		tasks, err := repo.List(ctx)
		return tasksLoadedMsg{tasks: tasks, err: err}
	}
}

func (m MainScreen) Update(msg tea.Msg) (ScreenI, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.list.SetSize(msg.Width-2, msg.Height-6)

	case tasksLoadedMsg:
		m.loading = false
		m.err = msg.err
		if msg.err == nil {
			m.setItems(msg.tasks)
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

func (m MainScreen) handleAction(msg ActionMsg) (ScreenI, tea.Cmd) {
	switch msg.ID {
	case keys.ActionExecute:
		if tk := m.SelectedTask(); tk != nil {
			return m, ExecuteTask(tk.ID, tk.Name)
		}
	case keys.EditNew:
		return m, OpenForm(nil, m.styles)
	case keys.ViewDetail:
		if tk := m.SelectedTask(); tk != nil {
			return m, OpenDetail(tk, m.styles)
		}
	case keys.ViewHistory:
		return m, OpenHistory(m.styles)
	case keys.ViewHelp:
		return m, OpenHelp()
	case keys.EditDelete:
		tk := m.SelectedTask()
		if tk == nil {
			return m, nil
		}
		return m, OpenConfirm(
			"Borrar task",
			fmt.Sprintf(
				"¿Borrar la task '%s'?\n\nEsta acción no se puede deshacer.",
				tk.Name,
			),
			DeleteTaskMsg{TaskID: tk.ID, TaskName: tk.Name},
			m.styles,
		)
	}
	return m, nil
}

func (m *MainScreen) setItems(tasks []*coretask.Task) {
	items := make([]list.Item, len(tasks))
	for i, t := range tasks {
		items[i] = taskItem{task: t}
	}
	m.list.SetItems(items)
}

func (m MainScreen) SelectedTask() *coretask.Task {
	it, ok := m.list.SelectedItem().(taskItem)
	if !ok {
		return nil
	}
	return it.task
}

func (m MainScreen) View() tea.View {
	switch {
	case m.loading:
		return tea.NewView(states.Loading(m.styles, "tareas"))
	case m.err != nil:
		return tea.NewView(states.Error(m.styles, m.err))
	case len(m.list.Items()) == 0:
		return tea.NewView(states.Empty(m.styles, "Sin tareas", "Pulsa n para crear la primera"))
	}
	return tea.NewView(m.list.View())
}

func (m MainScreen) Keys() []string {
	return []string{
		keys.ActionExecute,
		keys.ViewDetail,
		keys.ViewHistory,
		keys.EditNew,
		keys.EditDelete,
		keys.ViewHelp,
	}
}
