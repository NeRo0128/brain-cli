package styles

import "github.com/charmbracelet/lipgloss"

// Tokens: única fuente de verdad de la paleta.
//
// Reglas:
//   - Máximo 4 colores base: Text, Muted, Primary, SelectedRow.
//   - Success/Warning/Error SOLO para iconos de estado.
//   - Todo es AdaptiveColor: se ve bien en terminal claro y oscuro.

var (
	// Base
	Text    = lipgloss.AdaptiveColor{Light: "#1F2937", Dark: "#E5E7EB"}
	Muted   = lipgloss.AdaptiveColor{Light: "#6B7280", Dark: "#9CA3AF"}
	Primary = lipgloss.AdaptiveColor{Light: "#4C1D95", Dark: "#A78BFA"}

	// Acento secundario (muy usado, se queda)
	Secondary = lipgloss.AdaptiveColor{Light: "#0E7490", Dark: "#67E8F9"}

	// Estado (solo iconos/badges)
	Success = lipgloss.AdaptiveColor{Light: "#047857", Dark: "#10B981"}
	Warning = lipgloss.AdaptiveColor{Light: "#B45309", Dark: "#F59E0B"}
	Error   = lipgloss.AdaptiveColor{Light: "#B91C1C", Dark: "#EF4444"}

	// Superficies
	SelectedRow = lipgloss.AdaptiveColor{Light: "#E5E7EB", Dark: "#1F2937"}
	Border      = lipgloss.AdaptiveColor{Light: "#D1D5DB", Dark: "#374151"}
	StatusBar   = lipgloss.AdaptiveColor{Light: "#F3F4F6", Dark: "#111827"}

	// Foco de inputs
	InputFocused = Primary
	InputBlurred = Border
)
