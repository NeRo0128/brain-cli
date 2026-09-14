package theme

import lipgloss "charm.land/lipgloss/v2"

var tokyoNight = Theme{
	Name:        "tokyo-night",
	Description: "Azul neón, synthwave",
	Dark: Palette{
		Text:      lipgloss.Color("#A9B1D6"),
		Muted:     lipgloss.Color("#565F89"),
		Primary:   lipgloss.Color("#BB9AF7"),
		Secondary: lipgloss.Color("#7AA2F7"),

		SelectedRow: lipgloss.Color("#292E42"),
		Border:      lipgloss.Color("#3B4261"),
		StatusBar:   lipgloss.Color("#1A1B26"),

		Success: lipgloss.Color("#9ECE6A"),
		Warning: lipgloss.Color("#E0AF68"),
		Error:   lipgloss.Color("#F7768E"),

		InputFocused: lipgloss.Color("#BB9AF7"),
		InputBlurred: lipgloss.Color("#3B4261"),
	},
}
