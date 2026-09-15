package theme

import lipgloss "charm.land/lipgloss/v2"

var rosePine = Theme{
	Name:        "rose-pine",
	Description: "Rosa-morado, literario",
	Dark: Palette{ // Moon
		// Texto
		Text:     lipgloss.Color("#E0DEF4"),
		Emphasis: lipgloss.Color("#FFFFFF"),
		Muted:    lipgloss.Color("#908CAA"),
		Faint:    lipgloss.Color("#6E6A86"),

		// Marca
		Primary:   lipgloss.Color("#7E5A94"), // [ACTUALIZADO] iris oscurecido, 5.09
		Secondary: lipgloss.Color("#9CCFD8"), // foam
		Tertiary:  lipgloss.Color("#EBBCBA"), // rose

		// Superficies
		SelectedRow: lipgloss.Color("#2A273F"),
		Border:      lipgloss.Color("#393552"),
		BorderFocus: lipgloss.Color("#C4A7E7"),
		StatusBar:   lipgloss.Color("#232136"),

		// Estados
		Success: lipgloss.Color("#4CA183"), // [ACTUALIZADO] verde real (hue 159°)
		Warning: lipgloss.Color("#F6C177"), // gold
		Error:   lipgloss.Color("#EB6F92"), // love
		Info:    lipgloss.Color("#9CCFD8"), // foam

		// Prioridad
		PriorityHigh:   lipgloss.Color("#EB6F92"), // love
		PriorityMedium: lipgloss.Color("#F6C177"), // gold
		PriorityLow:    lipgloss.Color("#6E6A86"),

		// Tipos
		TypeScript:  lipgloss.Color("#31748F"), // pine
		TypeCommand: lipgloss.Color("#9CCFD8"), // foam
		TypeAI:      lipgloss.Color("#C4A7E7"), // iris

		// Inputs
		InputFocused: lipgloss.Color("#C4A7E7"),
		InputBlurred: lipgloss.Color("#393552"),
	},
	Light: &Palette{ // Dawn
		// Texto
		Text:     lipgloss.Color("#575279"),
		Emphasis: lipgloss.Color("#000000"),
		Muted:    lipgloss.Color("#797593"),
		Faint:    lipgloss.Color("#9893A5"),

		// Marca
		Primary:   lipgloss.Color("#907AA9"), // iris
		Secondary: lipgloss.Color("#56949F"), // foam
		Tertiary:  lipgloss.Color("#D7827E"), // rose

		// Superficies
		SelectedRow: lipgloss.Color("#F2E9E1"),
		Border:      lipgloss.Color("#DFDAD9"),
		BorderFocus: lipgloss.Color("#907AA9"),
		StatusBar:   lipgloss.Color("#FAF4ED"),

		// Estados
		Success: lipgloss.Color("#286983"), // pine
		Warning: lipgloss.Color("#EA9D34"), // gold
		Error:   lipgloss.Color("#B4637A"), // love
		Info:    lipgloss.Color("#56949F"), // foam

		// Prioridad
		PriorityHigh:   lipgloss.Color("#B4637A"),
		PriorityMedium: lipgloss.Color("#EA9D34"),
		PriorityLow:    lipgloss.Color("#9893A5"),

		// Tipos
		TypeScript:  lipgloss.Color("#286983"),
		TypeCommand: lipgloss.Color("#56949F"),
		TypeAI:      lipgloss.Color("#907AA9"),

		// Inputs
		InputFocused: lipgloss.Color("#907AA9"),
		InputBlurred: lipgloss.Color("#DFDAD9"),
	},
}
