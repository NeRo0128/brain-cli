package screens

import (
	"fmt"

	"github.com/NeRo0128/brain-cli/internal/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
)

// MainScreen es la pantalla principal de la aplicación.
// Por ahora solo muestra el título y la versión, pero aquí
// vivirán la lista de tareas y el menú principal.
type MainScreen struct {
	width   int
	height  int
	version string
	appName string
}

// NewMainScreen construye la pantalla inicial.
// Acepta la versión para poder mostrarla (inyección de dependencias).
func NewMainScreen(version, appName string) MainScreen {
	return MainScreen{
		version: version,
		appName: appName,
	}
}

// Init devuelve el comando inicial. Por ahora no hacemos nada.
func (m MainScreen) Init() tea.Cmd {
	return nil
}

// Update procesa mensajes y devuelve el nuevo estado.
func (m MainScreen) Update(msg tea.Msg) (MainScreen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Guardamos el tamaño para layout responsive
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

// View renderiza la pantalla como string.
func (m MainScreen) View() string {
	title := styles.Title.Render("🧠 Brain CLI")
	version := styles.Subtitle.Render(fmt.Sprintf("v%s", m.version))

	help := styles.Help.Render(
		styles.Key.Render("q") + " salir  •  " +
			styles.Key.Render("?") + " ayuda",
	)

	return fmt.Sprintf("%s\n%s\n\n%s", title, version, help)
}
