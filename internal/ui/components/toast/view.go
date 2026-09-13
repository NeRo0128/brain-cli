package toast

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/NeRo0128/brain-cli/internal/ui/styles"
)

const (
	// boxContentWidth es el ancho del área de contenido (sin borde ni padding).
	// El ancho total renderizado = boxContentWidth + 2 (padding) + 2 (borde).
	boxContentWidth = 30

	// barWidth es cuánto ocupa la barra de progreso.
	barWidth = 20

	// timeWidth reserva columnas para el tiempo restante ("12s" = 3 + 1 espacio).
	timeWidth = 5
)

// View renderiza la pila de toasts (sin posicionar).
// El Overlay se encarga de colocarla en pantalla.
func (m Model) View() string {
	if len(m.toasts) == 0 {
		return ""
	}
	blocks := make([]string, 0, len(m.toasts))
	for _, t := range m.toasts {
		blocks = append(blocks, renderBox(t))
	}
	// Alineados a la derecha entre sí.
	return lipgloss.JoinVertical(lipgloss.Right, blocks...)
}

// renderBox construye una caja individual con:
//
//	╭────────────────────────────────╮
//	│ ✓ Task guardada                │
//	│ ━━━━━━━━━━━━━━╸━━━━━━  2s      │
//	╰────────────────────────────────╯
func renderBox(t toast) string {
	icon, color := iconAndColor(t.level)

	// Tiempo restante (0 si ya expiró).
	remaining := time.Until(t.expires)
	if remaining < 0 {
		remaining = 0
	}
	ratio := float64(remaining) / float64(t.duration)
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}

	// Barra proporcional.
	filled := int(float64(barWidth) * ratio)
	if filled > barWidth {
		filled = barWidth
	}
	bar := strings.Repeat("━", filled) + strings.Repeat("░", barWidth-filled)

	// Texto del tiempo: "0s".."5s", alineado a la derecha del bloque.
	secs := int(remaining.Seconds())
	if secs < 0 {
		secs = 0
	}
	timeStr := fmt.Sprintf("%ds", secs)

	// Estilos por nivel.
	iconStyle := lipgloss.NewStyle().Foreground(color).Bold(true)
	barStyle := lipgloss.NewStyle().Foreground(color)
	timeStyle := lipgloss.NewStyle().Foreground(styles.Muted)
	msgStyle := lipgloss.NewStyle().Foreground(styles.Text)

	// Línea 1: icono + mensaje (truncado si hace falta).
	// Reservamos 2 columnas para "icon + espacio".
	msg := truncate(t.message, boxContentWidth-2)
	line1 := iconStyle.Render(icon) + " " + msgStyle.Render(msg)

	// Línea 2: barra + espacio + tiempo alineado a la derecha.
	barCell := barStyle.Render(bar) + strings.Repeat(" ", timeWidth-len(timeStr)-1) + timeStyle.Render(timeStr)

	content := line1 + "\n" + barCell

	// Caja con borde redondeado, color por nivel, padding lateral 1.
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(color).
		Padding(0, 1).
		Width(boxContentWidth)

	return box.Render(content)
}

// iconAndColor devuelve el icono y color según el nivel.
func iconAndColor(l Level) (string, lipgloss.TerminalColor) {
	switch l {
	case LevelSuccess:
		return "✓", styles.Success
	case LevelWarning:
		return "⚠", styles.Warning
	case LevelError:
		return "✗", styles.Error
	default:
		return "•", styles.Primary
	}
}

// truncate corta un string a n runas, añadiendo "…" si sobra.
func truncate(s string, n int) string {
	if n <= 0 {
		return ""
	}
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	if n == 1 {
		return "…"
	}
	return string(r[:n-1]) + "…"
}
