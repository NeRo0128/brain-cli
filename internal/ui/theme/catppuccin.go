package theme

import lipgloss "charm.land/lipgloss/v2"

var catppuccin = Theme{
	Name:        "catppuccin",
	Description: "Pastel violeta-mauve, suave",
	Dark: Palette{
		Text:      lipgloss.Color("#CAD3F5"),
		Muted:     lipgloss.Color("#A5ADCB"),
		Primary:   lipgloss.Color("#C6A0F6"),
		Secondary: lipgloss.Color("#8AADF4"),

		SelectedRow: lipgloss.Color("#363A4F"),
		Border:      lipgloss.Color("#494D64"),
		StatusBar:   lipgloss.Color("#24273A"),

		Success: lipgloss.Color("#A6DA95"),
		Warning: lipgloss.Color("#EED49F"),
		Error:   lipgloss.Color("#ED8796"),

		InputFocused: lipgloss.Color("#C6A0F6"),
		InputBlurred: lipgloss.Color("#494D64"),
	},
	Light: &Palette{
		Text:      lipgloss.Color("#4C4F69"),
		Muted:     lipgloss.Color("#6C6F85"),
		Primary:   lipgloss.Color("#8839EF"),
		Secondary: lipgloss.Color("#1E66F5"),

		SelectedRow: lipgloss.Color("#CCD0DA"),
		Border:      lipgloss.Color("#BCBFC8"),
		StatusBar:   lipgloss.Color("#E6E9EF"),

		Success: lipgloss.Color("#40A02B"),
		Warning: lipgloss.Color("#DF8E1D"),
		Error:   lipgloss.Color("#D20F39"),

		InputFocused: lipgloss.Color("#8839EF"),
		InputBlurred: lipgloss.Color("#BCBFC8"),
	},
}
