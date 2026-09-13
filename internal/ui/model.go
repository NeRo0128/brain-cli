package ui

import (
	"context"
	"time"

	toastComp "github.com/NeRo0128/brain-cli/internal/ui/components/toast"
	"github.com/NeRo0128/brain-cli/internal/ui/screens"

	tea "github.com/charmbracelet/bubbletea"
)

// devMinExecutingDisplay es un delay solo-dev para que el spinner sea visible.
const devMinExecutingDisplay = 0 * time.Second

// Model es el modelo raíz de la aplicación.
type Model struct {
	deps  Deps
	stack []screens.ScreenI

	width, height int
	executing     bool
	cancelExec    context.CancelFunc
	execErr       error
	toast         toastComp.Model // [ACTUALIZADO]
}

func NewModels(deps Deps, initial screens.ScreenI) Model {
	return Model{
		deps:  deps,
		stack: []screens.ScreenI{initial},
		toast: toastComp.NewModel(), // [ACTUALIZADO]
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.top().Init(),
		toastComp.Tick(), // [ACTUALIZADO]
	)
}

// --- Stack helpers ---

func (m *Model) top() screens.ScreenI     { return m.stack[len(m.stack)-1] }
func (m *Model) setTop(s screens.ScreenI) { m.stack[len(m.stack)-1] = s }

// --- Mensajes internos de navegación ---

type pushMsg struct{ screen screens.ScreenI }
type popMsg struct{}
type replaceMsg struct{ screen screens.ScreenI }

func push(s screens.ScreenI) tea.Cmd {
	return func() tea.Msg { return pushMsg{screen: s} }
}

func pop() tea.Cmd {
	return func() tea.Msg { return popMsg{} }
}

func replace(s screens.ScreenI) tea.Cmd {
	return func() tea.Msg { return replaceMsg{screen: s} }
}
