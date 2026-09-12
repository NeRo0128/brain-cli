package screens

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/NeRo0128/brain-cli/internal/core/task"
	"github.com/NeRo0128/brain-cli/internal/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
)

// tasksLoadedMsg transporta las tasks que llegan de la DB.
type tasksLoadedMsg struct {
	tasks []*task.Task
	err   error
}

// MainScreen es la pantalla principal de la aplicación.
// Por ahora solo muestra el título y la versión, pero aquí
// vivirán la lista de tareas y el menú principal.
type MainScreen struct {
	width   int
	height  int
	version string
	appName string

	taskRepo task.Repository

	loading bool
	tasks   []*task.Task
	err     error
}

// NewMainScreen construye la pantalla inicial.
// Acepta la versión para poder mostrarla (inyección de dependencias).
func NewMainScreen(version, appName string, repo task.Repository) MainScreen {
	return MainScreen{
		version:  version,
		appName:  appName,
		taskRepo: repo,
		loading:  true,
	}
}

// Init devuelve el comando inicial. Por ahora no hacemos nada.
func (m MainScreen) Init() tea.Cmd {

	repo := m.taskRepo
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		tasks, err := repo.List(ctx)
		return tasksLoadedMsg{tasks: tasks, err: err}
	}
}

// Update procesa mensajes y devuelve el nuevo estado.
func (m MainScreen) Update(msg tea.Msg) (MainScreen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Guardamos el tamaño para layout responsive
		m.width = msg.Width
		m.height = msg.Height
	case tasksLoadedMsg:
		m.loading = false
		m.tasks = msg.tasks
		m.err = msg.err
	}
	return m, nil
}

// View renderiza la pantalla como string.
func (m MainScreen) View() string {
	var b strings.Builder

	b.WriteString(styles.Title.Render("🧠 " + m.appName))
	b.WriteString("  ")
	b.WriteString(styles.Subtitle.Render("v" + m.version))
	b.WriteString("\n\n")

	switch {
	case m.loading:
		b.WriteString(styles.Subtitle.Render("Cargando tareas..."))
		b.WriteString("\n")
	case m.err != nil:
		b.WriteString(styles.Key.Render("Error: "))
		b.WriteString(m.err.Error())
		b.WriteString("\n")
	case len(m.tasks) == 0:
		b.WriteString(styles.Subtitle.Render("No hay tareas registradas."))
		b.WriteString("\n")
	default:
		b.WriteString(fmt.Sprintf("Tareas (%d):\n\n", len(m.tasks)))
		for _, t := range m.tasks {
			b.WriteString("  ")
			b.WriteString(styles.Subtitle.Render("[" + t.ID + "] "))
			b.WriteString(t.Name)
			if t.IsFavorite {
				b.WriteString(" ★")
			}
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(styles.Help.Render(
		styles.Key.Render("q") + " salir  •  " +
			styles.Key.Render("?") + " ayuda",
	))

	return b.String()
}
