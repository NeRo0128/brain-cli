package ui

import (
	"github.com/NeRo0128/brain-cli/internal/core/config"
	core "github.com/NeRo0128/brain-cli/internal/core/config"
	"github.com/NeRo0128/brain-cli/internal/core/task"
	"github.com/NeRo0128/brain-cli/internal/ui/screens"
	tea "github.com/charmbracelet/bubbletea"
)

// Model es el modelo raíz de la aplicación.
// Contiene las pantallas y el estado global.
type Model struct {
	version  string
	cfg      *core.Config
	main     screens.MainScreen
	taskRepo task.Repository
}

// NewModels construye el modelo raíz con todas sus dependencias.
func NewModels(version string, cfg *config.Config, taskRepo task.Repository) Model {
	return Model{
		version:  version,
		main:     screens.NewMainScreen(version, cfg.App.Name, taskRepo),
		cfg:      cfg,
		taskRepo: taskRepo,
	}
}

// Init arranca la aplicación.
func (m Model) Init() tea.Cmd { return m.main.Init() }

// Update maneja mensajes globales y delega a la pantalla activa.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Atajos globales: se manejan ANTES de delegar
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	// Delegamos el resto de mensajes a la pantalla activa
	var cmd tea.Cmd
	m.main, cmd = m.main.Update(msg)
	return m, cmd
}

// View delega el renderizado a la pantalla activa.
func (m Model) View() string { return m.main.View() }
