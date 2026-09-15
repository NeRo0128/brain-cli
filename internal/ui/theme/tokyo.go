package theme

import lipgloss "charm.land/lipgloss/v2"

var tokyoNight = Theme{
	Name:        "tokyo-night",
	Description: "Azul neón, synthwave",
	Dark: Palette{
		Text:     lipgloss.Color("#A9B1D6"),
		Emphasis: lipgloss.Color("#C0CAF5"),
		Muted:    lipgloss.Color("#565F89"),
		Faint:    lipgloss.Color("#414868"),

		Primary:   lipgloss.Color("#BB9AF7"),
		Secondary: lipgloss.Color("#7AA2F7"),
		Tertiary:  lipgloss.Color("#FF9E64"),

		SelectedRow: lipgloss.Color("#292E42"),
		Border:      lipgloss.Color("#3B4261"),
		BorderFocus: lipgloss.Color("#BB9AF7"),
		StatusBar:   lipgloss.Color("#1A1B26"),

		Success: lipgloss.Color("#9ECE6A"),
		Warning: lipgloss.Color("#E0AF68"),
		Error:   lipgloss.Color("#E65C5C"), // [ACTUALIZADO] rojo neutro 4.92
		Info:    lipgloss.Color("#7AA2F7"),

		PriorityHigh:   lipgloss.Color("#FF9E64"),
		PriorityMedium: lipgloss.Color("#E0AF68"),
		PriorityLow:    lipgloss.Color("#565F89"),

		TypeScript:  lipgloss.Color("#9ECE6A"),
		TypeCommand: lipgloss.Color("#7DCFFF"),
		TypeAI:      lipgloss.Color("#BB9AF7"),

		InputFocused: lipgloss.Color("#BB9AF7"),
		InputBlurred: lipgloss.Color("#3B4261"),
	},
	// [NUEVO] Light: Tokyo Night Day oficial.
	Light: &Palette{
		Text:     lipgloss.Color("#3760BF"),
		Emphasis: lipgloss.Color("#1A1B26"),
		Muted:    lipgloss.Color("#6172B0"),
		Faint:    lipgloss.Color("#A8AECB"),

		Primary:   lipgloss.Color("#9854F1"), // purple
		Secondary: lipgloss.Color("#2E7DE9"), // blue
		Tertiary:  lipgloss.Color("#B15C00"), // orange

		SelectedRow: lipgloss.Color("#D0D5E3"),
		Border:      lipgloss.Color("#A8AECB"),
		BorderFocus: lipgloss.Color("#9854F1"),
		StatusBar:   lipgloss.Color("#E9EAF0"),

		Success: lipgloss.Color("#587539"),
		Warning: lipgloss.Color("#8C6C3E"),
		Error:   lipgloss.Color("#F52A65"),
		Info:    lipgloss.Color("#2E7DE9"),

		PriorityHigh:   lipgloss.Color("#B15C00"),
		PriorityMedium: lipgloss.Color("#8C6C3E"),
		PriorityLow:    lipgloss.Color("#6172B0"),

		TypeScript:  lipgloss.Color("#587539"),
		TypeCommand: lipgloss.Color("#007197"),
		TypeAI:      lipgloss.Color("#9854F1"),

		InputFocused: lipgloss.Color("#9854F1"),
		InputBlurred: lipgloss.Color("#A8AECB"),
	},
}
