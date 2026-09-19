package screens

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/NeRo0128/brain-cli/internal/ui/keys"
	"github.com/NeRo0128/brain-cli/internal/ui/styles"
)

type HelpScreen struct {
	km     keys.KeyMap
	width  int
	height int
	styles *styles.Styles
}

func NewHelpScreen(km keys.KeyMap, s *styles.Styles) HelpScreen {
	return HelpScreen{km: km, styles: s}
}

func (m HelpScreen) Init() tea.Cmd  { return func() tea.Msg { return tea.RequestWindowSize() } }
func (m HelpScreen) Keys() []string { return []string{keys.NavBack} }

func (m HelpScreen) Update(msg tea.Msg) (ScreenI, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case ActionMsg:
		if msg.ID == keys.NavBack {
			return m, Back()
		}
	}
	return m, nil
}

func (m HelpScreen) View() tea.View {
	groups := m.km.Grouped()
	if len(groups) == 0 {
		return tea.NewView(m.styles.Subtitle.Render("(sin atajos registrados)"))
	}

	if m.width < helpTwoColMinWidth || len(groups) < 2 {
		return tea.NewView(m.renderGroups(groups))
	}

	mid := (len(groups) + 1) / 2
	left := m.renderGroups(groups[:mid])
	right := m.renderGroups(groups[mid:])

	colW := (m.width - 4) / 2
	leftCol := lipgloss.NewStyle().Width(colW).Render(left)
	rightCol := lipgloss.NewStyle().Width(colW).Render(right)

	return tea.NewView(lipgloss.JoinHorizontal(lipgloss.Top, leftCol, "   ", rightCol))
}

func (m HelpScreen) renderGroups(groups []keys.GroupBindings) string {
	var b strings.Builder
	for i, gb := range groups {
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString(m.styles.SectionHeader.Render(strings.ToUpper(gb.Group.Title())))
		b.WriteString("\n")
		for _, bd := range gb.Bindings {
			keysStr := strings.Join(bd.Keys, "/")
			b.WriteString("  ")
			b.WriteString(m.styles.Key.Render(padRight(keysStr, 16)))
			b.WriteString(" ")
			b.WriteString(bd.Help)
			b.WriteString("\n")
		}
	}
	return b.String()
}

func padRight(s string, n int) string {
	if len(s) >= n {
		return s
	}
	return s + strings.Repeat(" ", n-len(s))
}
