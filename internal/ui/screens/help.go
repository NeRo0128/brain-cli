package screens

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/NeRo0128/brain-cli/internal/ui/keys"
	"github.com/NeRo0128/brain-cli/internal/ui/styles"
)

// HelpScreen muestra los atajos agrupados.
type HelpScreen struct {
	km keys.KeyMap
}

func NewHelpScreen(km keys.KeyMap) HelpScreen {
	return HelpScreen{km: km}
}

func (m HelpScreen) Init() tea.Cmd  { return tea.WindowSize() }
func (m HelpScreen) Keys() []string { return []string{keys.NavBack} }

func (m HelpScreen) Update(msg tea.Msg) (ScreenI, tea.Cmd) {
	if act, ok := msg.(ActionMsg); ok && act.ID == keys.NavBack {
		return m, Back()
	}
	return m, nil
}

func (m HelpScreen) View() string {
	var b strings.Builder
	b.WriteString(styles.Title.Render("⌨  Atajos de teclado"))
	b.WriteString("\n\n")

	for _, gb := range m.km.Grouped() {
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
		b.WriteString("\n")
	}

	b.WriteString(styles.Help.Render(styles.Key.Render("Esc") + " cerrar"))
	return b.String()
}

func padRight(s string, n int) string {
	if len(s) >= n {
		return s
	}
	return s + strings.Repeat(" ", n-len(s))
}
