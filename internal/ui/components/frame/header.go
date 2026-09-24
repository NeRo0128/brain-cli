// Package frame provee el chrome de la TUI: header arriba y footer abajo.
// Los componentes son puros: reciben datos ya resueltos y devuelven
// una línea (o bloque de líneas) listos para concatenar.
package frame

import (
	"fmt"
	"image/color"
	"strings"

	"charm.land/lipgloss/v2"

	"github.com/NeRo0128/brain-cli/internal/ui/styles"
)

// AIStatus indica el estado del proveedor de IA activo.
type AIStatus int

const (
	AIDisabled AIStatus = iota // sin provider configurado
	AIReady                    // provider activo y responde
	AIOffline                  // provider activo pero no responde
)

// HeaderData son los datos que muestra el header.
// El Model los provee ya resueltos; el header no consulta nada.
type HeaderData struct {
	AppName       string
	Version       string
	BrandStyle    string // "ascii" | "minimal" | "none"
	AIStatus      AIStatus
	TaskCount     int
	FavoriteCount int
}

// Header devuelve UNA línea con el chrome superior.
//
//	🧠 Brain CLI v1.0.0             ● IA ready  12 tasks  3★
//
// Si width < 80, oculta la parte derecha para no apretar la marca.
func Header(d HeaderData, s *styles.Styles, width int) string {
	p := s.Theme.Resolve(s.Dark)

	// [NUEVO] Parsear el estilo.
	style := ParseBrandStyle(d.BrandStyle)
	brandLines := BrandLines(style)

	// [NUEVO] Si no hay ASCII art, usar la versión minimal (1 línea).
	if len(brandLines) == 0 {
		return renderMinimalHeader(d, s, width)
	}

	// Renderizar ASCII + bloque derecho en paralelo.
	brandStyle := lipgloss.NewStyle().Bold(true).Foreground(p.Primary)

	brandW := 0
	for _, l := range brandLines {
		if w := lipgloss.Width(l); w > brandW {
			brandW = w
		}
	}

	// Bloque derecho: AI status + contadores + versión.
	rightLines := buildRightBlock(d, s, width)
	rightW := 0
	for _, l := range rightLines {
		if w := lipgloss.Width(l); w > rightW {
			rightW = w
		}
	}

	// Degradar a minimal si no cabe.
	if brandW+rightW+4 > width {
		return renderMinimalHeader(d, s, width)
	}

	// Ensamblar línea por línea.
	maxLines := max(len(brandLines), len(rightLines))
	var out strings.Builder
	for i := 0; i < maxLines; i++ {
		var left, right string
		if i < len(brandLines) {
			left = brandStyle.Render(brandLines[i])
		} else {
			left = strings.Repeat(" ", brandW)
		}
		if i < len(rightLines) {
			right = rightLines[i]
		}

		gap := width - brandW - lipgloss.Width(right) - 2
		if gap < 1 {
			gap = 1
		}

		out.WriteString(" " + left + strings.Repeat(" ", gap) + right)
		if i < maxLines-1 {
			out.WriteString("\n")
		}
	}
	return out.String()
}
func renderMinimalHeader(d HeaderData, s *styles.Styles, width int) string {
	p := s.Theme.Resolve(s.Dark)

	brand := s.ColoredIcons.Brand() + " " +
		lipgloss.NewStyle().Bold(true).Foreground(p.Primary).Render(d.AppName)
	left := brand
	if d.Version != "" {
		left += " " + s.Subtitle.Render("v"+d.Version)
	}

	right := buildRightBlockInline(d, s, width)

	leftW := lipgloss.Width(left)
	rightW := lipgloss.Width(right)
	avail := width - leftW - rightW - 2
	if avail < 1 {
		if leftW > 0 {
			return " " + left
		}
		return " " + right + " "
	}
	if leftW == 0 {
		return " " + right + " "
	}
	return " " + left + strings.Repeat(" ", avail) + right + " "
}

// buildRightBlock: bloque derecho multi-línea.
func buildRightBlock(d HeaderData, s *styles.Styles, width int) []string {
	if width < 80 {
		return nil
	}
	line1 := renderAIStatus(d.AIStatus, s)
	if d.TaskCount > 0 {
		line1 += "  " + s.Subtitle.Render(fmt.Sprintf("%d tasks", d.TaskCount))
	}
	if d.FavoriteCount > 0 {
		line1 += "  " + s.ColoredIcons.Favorite() + s.Subtitle.Render(fmt.Sprintf(" %d", d.FavoriteCount))
	}
	line2 := s.Subtitle.Render("v" + d.Version)
	return []string{line1, line2}
}

// buildRightBlockInline: bloque derecho en una sola línea.
func buildRightBlockInline(d HeaderData, s *styles.Styles, width int) string {
	if width < 80 {
		return ""
	}
	right := renderAIStatus(d.AIStatus, s)
	if d.TaskCount > 0 {
		right += "  " + s.Subtitle.Render(fmt.Sprintf("%d tasks", d.TaskCount))
	}
	if d.FavoriteCount > 0 {
		right += "  " + s.ColoredIcons.Favorite() + s.Subtitle.Render(fmt.Sprintf(" %d", d.FavoriteCount))
	}
	return right
}

// renderAIStatus compone el indicador de IA con icono + color.
func renderAIStatus(status AIStatus, s *styles.Styles) string {
	p := s.Theme.Resolve(s.Dark)

	var icon string
	var clr color.Color
	switch status {
	case AIReady:
		icon, clr = s.Icons.Running, p.Success
	case AIOffline:
		icon, clr = s.Icons.Pending, p.Warning
	default:
		icon, clr = s.Icons.Pending, p.Muted
	}
	dot := lipgloss.NewStyle().Foreground(clr).Render(icon)
	label := s.Subtitle.Render(" IA ready")
	switch status {
	case AIOffline:
		label = s.Subtitle.Render(" IA offline")
	case AIDisabled:
		label = s.Subtitle.Render(" IA off")
	}
	return dot + label
}
