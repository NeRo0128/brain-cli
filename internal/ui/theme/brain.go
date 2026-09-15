package theme

import lipgloss "charm.land/lipgloss/v2"

var brain = Theme{
	Name:        "brain",
	Description: "Violeta saturado, la identidad original",
	Dark: Palette{
		Text:     lipgloss.Color("#E5E7EB"),
		Emphasis: lipgloss.Color("#FFFFFF"),
		Muted:    lipgloss.Color("#6B7280"),
		Faint:    lipgloss.Color("#4B5563"),

		Primary:   lipgloss.Color("#7C3AED"),
		Secondary: lipgloss.Color("#06B6D4"),
		Tertiary:  lipgloss.Color("#EC4899"),

		SelectedRow: lipgloss.Color("#1F2937"),
		Border:      lipgloss.Color("#374151"),
		BorderFocus: lipgloss.Color("#7C3AED"),
		StatusBar:   lipgloss.Color("#111827"),

		Success: lipgloss.Color("#10B981"),
		Warning: lipgloss.Color("#F59E0B"),
		Error:   lipgloss.Color("#EF4444"),
		Info:    lipgloss.Color("#3B82F6"),

		PriorityHigh:   lipgloss.Color("#F97316"),
		PriorityMedium: lipgloss.Color("#FBBF24"),
		PriorityLow:    lipgloss.Color("#6B7280"),

		TypeScript:  lipgloss.Color("#10B981"),
		TypeCommand: lipgloss.Color("#06B6D4"),
		TypeAI:      lipgloss.Color("#A855F7"),

		InputFocused: lipgloss.Color("#7C3AED"),
		InputBlurred: lipgloss.Color("#374151"),
	},
	// [NUEVO] Light: familia Tailwind, todas las contrastes ≥ 4.5
	Light: &Palette{
		Text:     lipgloss.Color("#1E293B"), // slate-800
		Emphasis: lipgloss.Color("#000000"),
		Muted:    lipgloss.Color("#64748B"), // slate-500
		Faint:    lipgloss.Color("#CBD5E1"), // slate-300

		Primary:   lipgloss.Color("#5B21B6"), // violet-800, 8.59
		Secondary: lipgloss.Color("#0E7490"), // cyan-700, 5.12
		Tertiary:  lipgloss.Color("#BE185D"), // pink-700, 5.77

		SelectedRow: lipgloss.Color("#E2E8F0"), // slate-200
		Border:      lipgloss.Color("#CBD5E1"), // slate-300
		BorderFocus: lipgloss.Color("#5B21B6"),
		StatusBar:   lipgloss.Color("#F1F5F9"), // slate-100

		Success: lipgloss.Color("#059669"), // emerald-600
		Warning: lipgloss.Color("#D97706"), // amber-600
		Error:   lipgloss.Color("#DC2626"), // red-600
		Info:    lipgloss.Color("#2563EB"), // blue-600

		PriorityHigh:   lipgloss.Color("#EA580C"), // orange-600
		PriorityMedium: lipgloss.Color("#D97706"),
		PriorityLow:    lipgloss.Color("#64748B"),

		TypeScript:  lipgloss.Color("#059669"),
		TypeCommand: lipgloss.Color("#0E7490"),
		TypeAI:      lipgloss.Color("#7C3AED"),

		InputFocused: lipgloss.Color("#5B21B6"),
		InputBlurred: lipgloss.Color("#CBD5E1"),
	},
}
