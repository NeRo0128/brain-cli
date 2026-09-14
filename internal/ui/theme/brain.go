package theme

import lipgloss "charm.land/lipgloss/v2"

var brain = Theme{
	Name:        "brain",
	Description: "Violeta saturado, la identidad original",
	Dark: Palette{
		Text:      lipgloss.Color("#E5E7EB"),
		Muted:     lipgloss.Color("#6B7280"),
		Primary:   lipgloss.Color("#7C3AED"),
		Secondary: lipgloss.Color("#06B6D4"),

		SelectedRow: lipgloss.Color("#1F2937"),
		Border:      lipgloss.Color("#374151"),
		StatusBar:   lipgloss.Color("#111827"),

		Success: lipgloss.Color("#10B981"),
		Warning: lipgloss.Color("#F59E0B"),
		Error:   lipgloss.Color("#EF4444"),

		InputFocused: lipgloss.Color("#7C3AED"),
		InputBlurred: lipgloss.Color("#374151"),
	},
}
