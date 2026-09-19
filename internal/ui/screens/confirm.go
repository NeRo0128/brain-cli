package screens

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/NeRo0128/brain-cli/internal/ui/keys"
	"github.com/NeRo0128/brain-cli/internal/ui/styles"
)

type ConfirmScreen struct {
	title   string
	message string
	action  tea.Msg
	styles  *styles.Styles

	width  int
	height int
}

func NewConfirm(title, message string, action tea.Msg, s *styles.Styles) ConfirmScreen {
	return ConfirmScreen{
		title:   title,
		message: message,
		action:  action,
		styles:  s,
	}
}

func (m ConfirmScreen) Init() tea.Cmd {
	return func() tea.Msg { return tea.RequestWindowSize() }
}

func (m ConfirmScreen) Keys() []string {
	return []string{keys.NavBack}
}

func (m ConfirmScreen) Update(msg tea.Msg) (ScreenI, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyPressMsg:
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

func (m ConfirmScreen) View() tea.View {
	boxW := 60
	if m.width > 0 && m.width-10 < boxW {
		boxW = m.width - 10
	}
	if boxW < 30 {
		boxW = 30
	}

	p := m.styles.Theme.Resolve(m.styles.Dark)

	var content strings.Builder

	content.WriteString(m.styles.ColoredIcons.Danger())
	content.WriteString(" ")
	content.WriteString(m.styles.ErrorStyle.Render(m.title))
	content.WriteString("\n\n")
	content.WriteString(m.message)
	content.WriteString("\n\n")
	content.WriteString(m.styles.Key.Render("s"))
	content.WriteString(m.styles.Subtitle.Render(" confirmar"))
	content.WriteString("   ")
	content.WriteString(m.styles.Key.Render("n"))
	content.WriteString(m.styles.Subtitle.Render(" cancelar"))

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(p.Error).
		Padding(1, 3).
		Width(boxW).
		Render(content.String())

	if m.width > 0 && m.height > 0 {
		h := max(m.height-4, 10)
		return tea.NewView(lipgloss.Place(m.width, h, lipgloss.Center, lipgloss.Center, box))
	}
	return tea.NewView("\n" + box + "\n")
}
