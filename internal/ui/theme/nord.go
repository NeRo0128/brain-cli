package theme

import lipgloss "charm.land/lipgloss/v2"

var nord = Theme{
	Name:        "nord",
	Description: "Azul frío, corporativo",
	Dark: Palette{
		// Texto (Nord 4-6)
		Text:     lipgloss.Color("#D8DEE9"),
		Emphasis: lipgloss.Color("#ECEFF4"),
		Muted:    lipgloss.Color("#7B88A1"),
		Faint:    lipgloss.Color("#4C566A"),

		// Marca (Nord 8-9-15)
		Primary:   lipgloss.Color("#88C0D0"),
		Secondary: lipgloss.Color("#81A1C1"),
		Tertiary:  lipgloss.Color("#B48EAD"),

		// Superficies (Nord 0-3)
		SelectedRow: lipgloss.Color("#3B4252"),
		Border:      lipgloss.Color("#4C566A"),
		BorderFocus: lipgloss.Color("#88C0D0"),
		StatusBar:   lipgloss.Color("#2E3440"),

		// Estados (Nord 11-14)
		Success: lipgloss.Color("#A3BE8C"),
		Warning: lipgloss.Color("#EBCB8B"),
		Error:   lipgloss.Color("#BF616A"),
		Info:    lipgloss.Color("#5E81AC"),

		// Prioridad (Nord 12-13-3)
		PriorityHigh:   lipgloss.Color("#D08770"),
		PriorityMedium: lipgloss.Color("#EBCB8B"),
		PriorityLow:    lipgloss.Color("#4C566A"),

		// Tipos (Nord 8-14-15)
		TypeScript:  lipgloss.Color("#A3BE8C"),
		TypeCommand: lipgloss.Color("#88C0D0"),
		TypeAI:      lipgloss.Color("#B48EAD"),

		// Inputs
		InputFocused: lipgloss.Color("#88C0D0"),
		InputBlurred: lipgloss.Color("#4C566A"),
	},
}
