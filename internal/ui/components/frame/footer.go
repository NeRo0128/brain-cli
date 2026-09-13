package frame

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/NeRo0128/brain-cli/internal/ui/styles"
)

// Binding es un atajo ya resuelto para mostrar en el footer.
// El Model lo construye desde el Registry + los Keys() de la screen.
type Binding struct {
	Keys []string // ej: ["↑↓"], ["enter"], ["q", "ctrl+c"]
	Help string   // ej: "navegar", "ejecutar", "salir"
}

var (
	footerKeyStyle  = lipgloss.NewStyle().Bold(true).Foreground(styles.Secondary)
	footerHelpStyle = lipgloss.NewStyle().Foreground(styles.Muted)
	footerSepStyle  = lipgloss.NewStyle().Foreground(styles.Border)
)

// Footer devuelve UNA línea con los atajos activos.
//
//	↑↓ navegar · Enter ejecutar · d detalle · ? ayuda · q salir
//
// Si no caben todos, corta por el final y añade "…".
// Si width <= 0, muestra todo (útil para tests).
func Footer(bindings []Binding, width int) string {
	if len(bindings) == 0 {
		return ""
	}

	sep := footerSepStyle.Render(" · ")
	rendered := make([]string, 0, len(bindings))

	accumulated := 0
	for i, b := range bindings {
		keysStr := strings.Join(b.Keys, "/")
		piece := footerKeyStyle.Render(keysStr) + " " + footerHelpStyle.Render(b.Help)
		pieceW := lipgloss.Width(piece)

		// Reservar espacio para el separador (salvo el primero).
		sepW := 0
		if i > 0 {
			sepW = lipgloss.Width(sep)
		}

		if width > 0 && accumulated+sepW+pieceW+2 > width {
			// No cabe: añadir "…" y cortar.
			rendered = append(rendered, footerHelpStyle.Render("…"))
			break
		}

		rendered = append(rendered, piece)
		accumulated += sepW + pieceW
	}

	return " " + strings.Join(rendered, sep) + " "
}
