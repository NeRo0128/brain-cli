package theme

import lipgloss "charm.land/lipgloss/v2"

var catppuccin = Theme{
	Name:        "catppuccin",
	Description: "Pastel violeta-mauve, suave",
	Dark: Palette{ // Macchiato
		// Texto
		Text:     lipgloss.Color("#CAD3F5"),
		Emphasis: lipgloss.Color("#FFFFFF"),
		Muted:    lipgloss.Color("#A5ADCB"),
		Faint:    lipgloss.Color("#8087A2"),

		// Marca
		Primary:   lipgloss.Color("#C6A0F6"), // mauve
		Secondary: lipgloss.Color("#8AADF4"), // blue
		Tertiary:  lipgloss.Color("#F5A97F"), // peach — Δhue 246° vs mauve

		// Superficies
		SelectedRow: lipgloss.Color("#363A4F"),
		Border:      lipgloss.Color("#494D64"),
		BorderFocus: lipgloss.Color("#C6A0F6"),
		StatusBar:   lipgloss.Color("#24273A"),

		// Estados
		Success: lipgloss.Color("#A6DA95"), // green
		Warning: lipgloss.Color("#EED49F"), // yellow
		Error:   lipgloss.Color("#ED8796"), // red
		Info:    lipgloss.Color("#8AADF4"), // blue

		// Prioridad
		PriorityHigh:   lipgloss.Color("#F5A97F"), // peach
		PriorityMedium: lipgloss.Color("#EED49F"), // yellow
		PriorityLow:    lipgloss.Color("#8087A2"), // overlay0

		// Tipos
		TypeScript:  lipgloss.Color("#A6DA95"), // green
		TypeCommand: lipgloss.Color("#8BD5CA"), // teal
		TypeAI:      lipgloss.Color("#C6A0F6"), // mauve

		// Inputs
		InputFocused: lipgloss.Color("#C6A0F6"),
		InputBlurred: lipgloss.Color("#494D64"),
	},
	Light: &Palette{ // Latte
		// Texto
		Text:     lipgloss.Color("#4C4F69"),
		Emphasis: lipgloss.Color("#000000"),
		Muted:    lipgloss.Color("#6C6F85"),
		Faint:    lipgloss.Color("#9CA0B0"),

		// Marca
		Primary:   lipgloss.Color("#8839EF"), // mauve
		Secondary: lipgloss.Color("#1E66F5"), // blue
		Tertiary:  lipgloss.Color("#EA76CB"), // pink

		// Superficies
		SelectedRow: lipgloss.Color("#CCD0DA"),
		Border:      lipgloss.Color("#BCBFC8"),
		BorderFocus: lipgloss.Color("#8839EF"),
		StatusBar:   lipgloss.Color("#E6E9EF"),

		// Estados
		Success: lipgloss.Color("#40A02B"),
		Warning: lipgloss.Color("#DF8E1D"),
		Error:   lipgloss.Color("#D20F39"),
		Info:    lipgloss.Color("#1E66F5"),

		// Prioridad
		PriorityHigh:   lipgloss.Color("#FE640B"), // peach
		PriorityMedium: lipgloss.Color("#DF8E1D"), // yellow
		PriorityLow:    lipgloss.Color("#9CA0B0"),

		// Tipos
		TypeScript:  lipgloss.Color("#40A02B"),
		TypeCommand: lipgloss.Color("#179299"), // teal
		TypeAI:      lipgloss.Color("#8839EF"),

		// Inputs
		InputFocused: lipgloss.Color("#8839EF"),
		InputBlurred: lipgloss.Color("#BCBFC8"),
	},
}
