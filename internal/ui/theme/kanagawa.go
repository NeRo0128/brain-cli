package theme

import lipgloss "charm.land/lipgloss/v2"

var kanagawa = Theme{
	Name:        "kanagawa",
	Description: "Beige-tierra, cálido",
	Dark: Palette{
		Text:      lipgloss.Color("#DCD7BA"),
		Muted:     lipgloss.Color("#727169"),
		Primary:   lipgloss.Color("#7E9CD8"),
		Secondary: lipgloss.Color("#957FB8"),

		SelectedRow: lipgloss.Color("#2A2A37"),
		Border:      lipgloss.Color("#54546D"),
		StatusBar:   lipgloss.Color("#1F1F28"),

		Success: lipgloss.Color("#76946A"),
		Warning: lipgloss.Color("#C0A36E"),
		Error:   lipgloss.Color("#C34043"),

		InputFocused: lipgloss.Color("#7E9CD8"),
		InputBlurred: lipgloss.Color("#54546D"),
	},
}
