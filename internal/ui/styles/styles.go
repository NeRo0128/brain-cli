package styles

import "github.com/charmbracelet/lipgloss"

// Paleta de colores centralizada.
// Cambiar aquí afecta toda la UI.
var (
	Primary   = lipgloss.Color("#7C3AED") // violeta
	Secondary = lipgloss.Color("#06B6D4") // cyan
	Muted     = lipgloss.Color("#6B7280") // gris
	Text      = lipgloss.Color("#E5E7EB") // gris claro
	Success   = lipgloss.Color("#10B981") // verde
	Warning   = lipgloss.Color("#F59E0B") // ámbar
	Error     = lipgloss.Color("#EF4444") // rojo
)

var (
	Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(Primary).
		Padding(0, 1)

	Subtitle = lipgloss.NewStyle().
			Foreground(Muted).
			Italic(true)

	Help = lipgloss.NewStyle().
		Foreground(Muted).
		Padding(1, 2)

	Key = lipgloss.NewStyle().
		Bold(true).
		Foreground(Secondary)
)
