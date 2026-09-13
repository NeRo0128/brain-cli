package screens

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/NeRo0128/brain-cli/internal/ui/keys"
	"github.com/NeRo0128/brain-cli/internal/ui/styles"
)

// ConfirmScreen es un modal de confirmación para acciones destructivas.
//
// Contrato:
//   - Se pushea con NewConfirm(title, message, action).
//   - Al confirmar (s/y): emite ConfirmYesMsg{Action}.
//   - Al cancelar (n/Esc): emite BackMsg.
//   - El Model raíz poppea este modal y despacha Action al top nuevo.
type ConfirmScreen struct {
	title   string
	message string
	action  tea.Msg

	width  int
	height int
}

func NewConfirm(title, message string, action tea.Msg) ConfirmScreen {
	return ConfirmScreen{
		title:   title,
		message: message,
		action:  action,
	}
}

func (m ConfirmScreen) Init() tea.Cmd {
	return tea.WindowSize()
}

// Keys declara solo NavBack (Esc). Las teclas s/y/n se manejan
// directamente en Update sin pasar por el Registry, porque son
// propias de este modal y no deben ser configurables globalmente.
func (m ConfirmScreen) Keys() []string {
	return []string{keys.NavBack}
}

func (m ConfirmScreen) Update(msg tea.Msg) (ScreenI, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "s", "y":
			action := m.action
			return m, func() tea.Msg {
				return ConfirmYesMsg{Action: action}
			}
		case "n":
			return m, Back()
		}

	case ActionMsg:
		if msg.ID == keys.NavBack {
			return m, Back()
		}
	}
	return m, nil
}

func (m ConfirmScreen) View() string {
	// Ancho de la caja: 60 por defecto, ajustado a terminales chicas.
	boxW := 60
	if m.width > 0 && m.width-10 < boxW {
		boxW = m.width - 10
	}
	if boxW < 30 {
		boxW = 30
	}

	var content strings.Builder
	content.WriteString(styles.ErrorStyle.Render("⚠  " + m.title))
	content.WriteString("\n\n")
	content.WriteString(m.message)
	content.WriteString("\n\n")
	content.WriteString(styles.Key.Render("s") + styles.Subtitle.Render(" confirmar"))
	content.WriteString("   ")
	content.WriteString(styles.Key.Render("n") + styles.Subtitle.Render(" cancelar"))

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Error).
		Padding(1, 3).
		Width(boxW).
		Render(content.String())

	// Centrar en el área de contenido (descontando el chrome del frame).
	if m.width > 0 && m.height > 0 {
		h := m.height - 4
		if h < 10 {
			h = 10
		}
		return lipgloss.Place(m.width, h, lipgloss.Center, lipgloss.Center, box)
	}
	return "\n" + box + "\n"
}
