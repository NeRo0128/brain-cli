package theme

import lipgloss "charm.land/lipgloss/v2"

var kanagawa = Theme{
	Name:        "kanagawa",
	Description: "Beige-tierra, cálido",
	Dark: Palette{
		Text:     lipgloss.Color("#DCD7BA"),
		Emphasis: lipgloss.Color("#FFFFFF"),
		Muted:    lipgloss.Color("#727169"),
		Faint:    lipgloss.Color("#54546D"),

		Primary:   lipgloss.Color("#957FB8"), // [SWAP] oniViolet (marca)
		Secondary: lipgloss.Color("#7E9CD8"), // [SWAP] crystalBlue
		Tertiary:  lipgloss.Color("#FFA066"), // surimiOrange

		SelectedRow: lipgloss.Color("#2A2A37"),
		Border:      lipgloss.Color("#54546D"),
		BorderFocus: lipgloss.Color("#957FB8"), // [ACTUALIZADO] sigue primary
		StatusBar:   lipgloss.Color("#1F1F28"),

		Success: lipgloss.Color("#98BB6C"),
		Warning: lipgloss.Color("#FF9E3B"), // [ACTUALIZADO] roninYellow, 7.94
		Error:   lipgloss.Color("#FF5D62"),
		Info:    lipgloss.Color("#7FB4CA"),

		PriorityHigh:   lipgloss.Color("#FF5D62"),
		PriorityMedium: lipgloss.Color("#FF9E3B"), // [ACTUALIZADO] mismo que warning
		PriorityLow:    lipgloss.Color("#727169"),

		TypeScript:  lipgloss.Color("#98BB6C"),
		TypeCommand: lipgloss.Color("#7FB4CA"),
		TypeAI:      lipgloss.Color("#957FB8"),

		InputFocused: lipgloss.Color("#957FB8"),
		InputBlurred: lipgloss.Color("#54546D"),
	},
	// [NUEVO] Light: Lotus con 3 fixes de contraste.
	Light: &Palette{
		Text:     lipgloss.Color("#545464"),
		Emphasis: lipgloss.Color("#1F1F28"),
		Muted:    lipgloss.Color("#8A8980"),
		Faint:    lipgloss.Color("#A09CAC"),

		Primary:   lipgloss.Color("#624C83"), // oniViolet Lotus
		Secondary: lipgloss.Color("#4D699B"), // crystalBlue Lotus
		Tertiary:  lipgloss.Color("#A34D68"), // [FIX] contraste 4.58

		SelectedRow: lipgloss.Color("#DCD5AC"),
		Border:      lipgloss.Color("#C9C1A0"),
		BorderFocus: lipgloss.Color("#624C83"),
		StatusBar:   lipgloss.Color("#E8E1B5"),

		Success: lipgloss.Color("#54703B"), // [FIX] contraste 4.66
		Warning: lipgloss.Color("#8F6428"), // [FIX] contraste 4.35 bold
		Error:   lipgloss.Color("#C84053"),
		Info:    lipgloss.Color("#4D699B"),

		PriorityHigh:   lipgloss.Color("#C84053"),
		PriorityMedium: lipgloss.Color("#8F6428"),
		PriorityLow:    lipgloss.Color("#8A8980"),

		TypeScript:  lipgloss.Color("#54703B"),
		TypeCommand: lipgloss.Color("#4D699B"),
		TypeAI:      lipgloss.Color("#624C83"),

		InputFocused: lipgloss.Color("#624C83"),
		InputBlurred: lipgloss.Color("#C9C1A0"),
	},
}
