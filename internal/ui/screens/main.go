package screens

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/charmbracelet/bubbles/list"

	"github.com/NeRo0128/brain-cli/internal/core/task"
	"github.com/NeRo0128/brain-cli/internal/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

func (i taskItem) Title() string {
	title := i.task.Name
	if i.task.IsFavorite {
		title = "★ " + title
	}
	return title
}

func (i taskItem) Description() string {
	return fmt.Sprintf("[%s] %s  ·  %s", i.task.ID, i.task.Type, i.task.Priority)
}

// FilterValue permite buscar por nombre e ID.
func (i taskItem) FilterValue() string {
	return i.task.Name + " " + i.task.ID
}

// --- Item delegate (cómo se renderiza cada fila) ---

type taskDelegate struct{}

func (d taskDelegate) Height() int                             { return 2 }
func (d taskDelegate) Spacing() int                            { return 1 }
func (d taskDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d taskDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	it, ok := item.(taskItem)
	if !ok {
		return
	}

	selected := index == m.Index()

	titleStyle := lipgloss.NewStyle().Foreground(styles.Text)
	descStyle := lipgloss.NewStyle().Foreground(styles.Muted)
	cursor := "  "

	if selected {
		titleStyle = titleStyle.Bold(true).Foreground(styles.Primary)
		cursor = "▶ "
	}

	fmt.Fprintf(w, "%s%s\n", cursor, titleStyle.Render(it.Title()))
	fmt.Fprintf(w, "   %s", descStyle.Render(it.Description()))
}

// --- Pantalla principal ---

// MainScreen es la pantalla de lista de tasks.
type MainScreen struct {
	width   int
	height  int
	appName string
	version string

	loading bool
	err     error
	list    list.Model
	repo    task.Repository
}

// NewMainScreen construye la pantalla.
func NewMainScreen(version, appName string, repo task.Repository) MainScreen {
	l := list.New(nil, taskDelegate{}, 80, 20)
	l.Title = "Tareas"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = styles.Title
	l.Styles.HelpStyle = styles.Help

	return MainScreen{
		appName: appName,
		version: version,
		list:    l,
		repo:    repo,
		loading: true,
	}
}

// Init dispara la carga asíncrona.
func (m MainScreen) Init() tea.Cmd {
	repo := m.repo
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		tasks, err := repo.List(ctx)
		return tasksLoadedMsg{tasks: tasks, err: err}
	}
}

// Update maneja los mensajes.
func (m MainScreen) Update(msg tea.Msg) (MainScreen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width, msg.Height-6)

	case tasksLoadedMsg:
		m.loading = false
		m.err = msg.err
		if msg.err == nil {
			m.setItems(msg.tasks)
		}
		return m, nil
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// setItems pobla la lista con las tasks.
func (m *MainScreen) setItems(tasks []*task.Task) {
	items := make([]list.Item, len(tasks))
	for i, t := range tasks {
		items[i] = taskItem{task: t}
	}
	m.list.SetItems(items)
}

// SelectedTask devuelve la task actualmente seleccionada (nil si no hay).
func (m MainScreen) SelectedTask() *task.Task {
	it, ok := m.list.SelectedItem().(taskItem)
	if !ok {
		return nil
	}
	return it.task
}

// View renderiza la pantalla.
func (m MainScreen) View() string {
	header := styles.Title.Render("🧠 "+m.appName) + "  " +
		styles.Subtitle.Render("v"+m.version)

	switch {
	case m.loading:
		return header + "\n\n" + styles.Subtitle.Render("Cargando tareas...") + "\n"
	case m.err != nil:
		return header + "\n\n" + styles.Key.Render("Error: ") + m.err.Error() + "\n"
	}

	help := styles.Help.Render(
		styles.Key.Render("↑↓") + " navegar  ·  " +
			styles.Key.Render("Enter") + " ejecutar  ·  " +
			styles.Key.Render("d") + " detalle  ·  " +
			styles.Key.Render("h") + " historial  ·  " +
			styles.Key.Render("/") + " filtrar  ·  " +
			styles.Key.Render("q") + " salir",
	)

	return header + "\n\n" + m.list.View() + "\n" + help
}
