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
//	🧠 Brain CLI v2.0               ● IA ready  12 tasks  3★
//
// Si width < 80, oculta la parte derecha para no apretar la marca.
func Header(d HeaderData, s *styles.Styles, width int) string {
	p := s.Theme.Resolve(s.Dark)

	// --- Izquierda: marca + versión ---
	var left string
	switch d.BrandStyle {
	case "ascii":
		left = lipgloss.NewStyle().Bold(true).Foreground(p.Primary).Render("BrainCLI")
		if d.Version != "" {
			left += " " + s.Subtitle.Render("v"+d.Version)
		}
	case "none":
		left = ""
	default: // "minimal" or unset
		left = lipgloss.NewStyle().Bold(true).Foreground(p.Primary).Render("🧠 " + d.AppName)
		if d.Version != "" {
			left += " " + s.Subtitle.Render("v" + d.Version)
		}
	}

	// --- Derecha: estado IA + contadores ---
	right := ""
	if width >= 80 {
		right = renderAIStatus(d.AIStatus, s)
		if d.TaskCount > 0 {
			right += "  " + s.Subtitle.Render(
				fmt.Sprintf("%d tasks", d.TaskCount))
		}
		if d.FavoriteCount > 0 {
			right += "  " + lipgloss.NewStyle().Bold(true).Foreground(p.Warning).Render(
				fmt.Sprintf("%d★", d.FavoriteCount))
		}
	}

	// --- Ensamblar con gap flexible ---
	leftW := lipgloss.Width(left)
	rightW := lipgloss.Width(right)

	// Espacio total disponible (deja 1 col a cada lado).
	avail := width - leftW - rightW - 2
	if avail < 1 {
		// No cabe todo: priorizar izquierda, truncar derecha.
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

// renderAIStatus compone el indicador de IA con icono + color.
func renderAIStatus(status AIStatus, s *styles.Styles) string {
	p := s.Theme.Resolve(s.Dark)

	var icon string
	var clr color.Color
	switch status {
	case AIReady:
		icon, clr = "●", p.Success
	case AIOffline:
		icon, clr = "○", p.Warning
	default:
		icon, clr = "○", p.Muted
	}
	dot := lipgloss.NewStyle().Foreground(clr).Render(icon)
	label := s.Subtitle.Render(" IA ready")
	if status == AIOffline {
		label = s.Subtitle.Render(" IA offline")
	} else if status == AIDisabled {
		label = s.Subtitle.Render(" IA off")
	}
	return dot + label
}
