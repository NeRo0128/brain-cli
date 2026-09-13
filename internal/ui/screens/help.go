package screens

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/NeRo0128/brain-cli/internal/ui/keys"
	"github.com/NeRo0128/brain-cli/internal/ui/styles"
)

const helpTwoColMinWidth = 100

// HelpScreen muestra los atajos agrupados.
// [S4c] responsive: 2 columnas si width >= 100, 1 columna si no.
type HelpScreen struct {
	km     keys.KeyMap
	width  int
	height int
}

func NewHelpScreen(km keys.KeyMap) HelpScreen {
	return HelpScreen{km: km}
}

func (m HelpScreen) Init() tea.Cmd  { return tea.WindowSize() }
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

func (m HelpScreen) View() string {
	groups := m.km.Grouped()
	if len(groups) == 0 {
		return styles.Subtitle.Render("(sin atajos registrados)")
	}

	// 1 columna: terminal angosta o pocos grupos.
	if m.width < helpTwoColMinWidth || len(groups) < 2 {
		return m.renderGroups(groups)
	}

	// 2 columnas: partir grupos por la mitad.
	mid := (len(groups) + 1) / 2
	left := m.renderGroups(groups[:mid])
	right := m.renderGroups(groups[mid:])

	colW := (m.width - 4) / 2 // 2 cols de gap + 1 de margen c/lado
	leftCol := lipgloss.NewStyle().Width(colW).Render(left)
	rightCol := lipgloss.NewStyle().Width(colW).Render(right)

	return lipgloss.JoinHorizontal(lipgloss.Top, leftCol, "   ", rightCol)
}

// renderGroups aplana los grupos a un bloque de texto.
func (m HelpScreen) renderGroups(groups []keys.GroupBindings) string {
	var b strings.Builder
	for i, gb := range groups {
		if i > 0 {
			b.WriteString("\n") // separación entre grupos
		}
		b.WriteString(styles.SectionHeader.Render(strings.ToUpper(gb.Group.Title())))
		b.WriteString("\n")
		for _, bd := range gb.Bindings {
			keysStr := strings.Join(bd.Keys, "/")
			b.WriteString("  ")
			b.WriteString(styles.Key.Render(padRight(keysStr, 16)))
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
