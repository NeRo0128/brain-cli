package theme

import lipgloss "charm.land/lipgloss/v2"

var nord = Theme{
	Name:        "nord",
	Description: "Azul frío, corporativo",
	Dark: Palette{
		Text:      lipgloss.Color("#D8DEE9"),
		Muted:     lipgloss.Color("#7B88A1"),
		Primary:   lipgloss.Color("#88C0D0"),
		Secondary: lipgloss.Color("#81A1C1"),

		SelectedRow: lipgloss.Color("#3B4252"),
		Border:      lipgloss.Color("#4C566A"),
		StatusBar:   lipgloss.Color("#2E3440"),

		Success: lipgloss.Color("#A3BE8C"),
		Warning: lipgloss.Color("#EBCB8B"),
		Error:   lipgloss.Color("#BF616A"),

		InputFocused: lipgloss.Color("#88C0D0"),
		InputBlurred: lipgloss.Color("#4C566A"),
	},
}
