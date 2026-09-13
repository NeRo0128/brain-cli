// Package frame provee el chrome de la TUI: header arriba y footer abajo.
// Los componentes son puros: reciben datos ya resueltos y devuelven
// una línea (o bloque de líneas) listos para concatenar.
package frame

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

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
	AIStatus      AIStatus
	TaskCount     int
	FavoriteCount int
}

// Estilos locales del header (no contaminan styles global).
var (
	headerBrandStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(styles.Primary)

	headerMutedStyle = lipgloss.NewStyle().
				Foreground(styles.Muted)

	headerFavoriteStyle = lipgloss.NewStyle().
				Foreground(styles.Warning).
				Bold(true)
)

// Header devuelve UNA línea con el chrome superior.
//
//	🧠 Brain CLI v2.0               ● IA ready  12 tasks  3★
//
// Si width < 80, oculta la parte derecha para no apretar la marca.
func Header(d HeaderData, width int) string {
	// --- Izquierda: marca + versión ---
	left := headerBrandStyle.Render("🧠 " + d.AppName)
	if d.Version != "" {
		left += " " + headerMutedStyle.Render("v"+d.Version)
	}

	// --- Derecha: estado IA + contadores ---
	right := ""
	if width >= 80 {
		right = renderAIStatus(d.AIStatus)
		if d.TaskCount > 0 {
			right += "  " + headerMutedStyle.Render(
				fmt.Sprintf("%d tasks", d.TaskCount))
		}
		if d.FavoriteCount > 0 {
			right += "  " + headerFavoriteStyle.Render(
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
		return " " + left
	}

	return " " + left + strings.Repeat(" ", avail) + right + " "
}

// renderAIStatus compone el indicador de IA con icono + color.
func renderAIStatus(s AIStatus) string {
	var icon string
	var color lipgloss.TerminalColor
	switch s {
	case AIReady:
		icon, color = "●", styles.Success
	case AIOffline:
		icon, color = "○", styles.Warning
	default:
		icon, color = "○", styles.Muted
	}
	dot := lipgloss.NewStyle().Foreground(color).Render(icon)
	label := headerMutedStyle.Render(" IA ready")
	if s == AIOffline {
		label = headerMutedStyle.Render(" IA offline")
	} else if s == AIDisabled {
		label = headerMutedStyle.Render(" IA off")
	}
	return dot + label
}
