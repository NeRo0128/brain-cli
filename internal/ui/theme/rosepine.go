package theme

import lipgloss "charm.land/lipgloss/v2"

var rosePine = Theme{
	Name:        "rose-pine",
	Description: "Rosa-morado, literario",
	Dark: Palette{
		Text:      lipgloss.Color("#E0DEF4"),
		Muted:     lipgloss.Color("#908CAA"),
		Primary:   lipgloss.Color("#C4A7E7"),
		Secondary: lipgloss.Color("#9CCFD8"),

		SelectedRow: lipgloss.Color("#2A273F"),
		Border:      lipgloss.Color("#393552"),
		StatusBar:   lipgloss.Color("#232136"),

		Success: lipgloss.Color("#3E8FB0"),
		Warning: lipgloss.Color("#F6C177"),
		Error:   lipgloss.Color("#EB6F92"),

		InputFocused: lipgloss.Color("#C4A7E7"),
		InputBlurred: lipgloss.Color("#393552"),
	},
	Light: &Palette{
		Text:      lipgloss.Color("#575279"),
		Muted:     lipgloss.Color("#797593"),
		Primary:   lipgloss.Color("#907AA9"),
		Secondary: lipgloss.Color("#56949F"),

		SelectedRow: lipgloss.Color("#F2E9E1"),
		Border:      lipgloss.Color("#DFDAD9"),
		StatusBar:   lipgloss.Color("#FAF4ED"),

		Success: lipgloss.Color("#286983"),
		Warning: lipgloss.Color("#EA9D34"),
		Error:   lipgloss.Color("#B4637A"),

		InputFocused: lipgloss.Color("#907AA9"),
		InputBlurred: lipgloss.Color("#DFDAD9"),
	},
}
