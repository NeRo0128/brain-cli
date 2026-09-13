package screens

import (
	"context"
	"fmt"
	"time"

	"github.com/NeRo0128/brain-cli/internal/core/task"
	coretask "github.com/NeRo0128/brain-cli/internal/core/task"
	uilist "github.com/NeRo0128/brain-cli/internal/ui/components/list"
	"github.com/NeRo0128/brain-cli/internal/ui/components/states"
	"github.com/NeRo0128/brain-cli/internal/ui/keys"
	"github.com/NeRo0128/brain-cli/internal/ui/styles"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// --- Mensajes propios de la pantalla ---

// tasksLoadedMsg transporta las tasks que llegan de la DB.
type tasksLoadedMsg struct {
	tasks []*task.Task
	err   error
}

// --- Item de la lista ---

// taskItem adapta *task.Task a list.Item (Title + Description).
type taskItem struct {
	task *task.Task
}

func (i taskItem) Title() string       { return i.task.Name }
func (i taskItem) Description() string { return "" }
func (i taskItem) FilterValue() string { return i.task.Name + " " + i.task.ID }

func (i taskItem) Row() uilist.Row {
	prefix := ""
	if i.task.IsFavorite {
		prefix = "★"
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

// --- Pantalla ---

type MainScreen struct {
	width, height int
	loading       bool
	err           error
	list          list.Model
	repo          coretask.Repository
}

func NewMainScreen(repo coretask.Repository) MainScreen {
	l := list.New(nil, uilist.New(), 80, 20)
	l.Title = "Tareas"
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = styles.Title
	l.Styles.HelpStyle = styles.Help

	return MainScreen{
		list:    l,
		repo:    repo,
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
		return m, OpenForm(nil)
	case keys.ViewDetail:
		if tk := m.SelectedTask(); tk != nil {
			return m, OpenDetail(tk)
		}
	case keys.ViewHistory:
		return m, OpenHistory()
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

func (m MainScreen) View() string {
	switch {
	case m.loading:
		return states.Loading("tareas")
	case m.err != nil:
		return states.Error(m.err)
	case len(m.list.Items()) == 0:
		return states.Empty("Sin tareas", "Pulsa n para crear la primera")
	}
	return m.list.View()
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
