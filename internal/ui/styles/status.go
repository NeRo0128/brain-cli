package styles

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// StatusIcon devuelve el símbolo y el color para un estado.
// Acepta tanto estados de Execution como strings cortos.
func StatusIcon(status string) (string, lipgloss.TerminalColor) {
	switch status {
	case "completed", "success", "ok":
		return "✓", Success
	case "failed", "error":
		return "✗", Error
	case "cancelled", "canceled", "warn":
		return "⊘", Warning
	case "running", "pending":
		return "◐", Primary
	}
	return "•", Muted
}

// Badge renderiza un texto con estilo discreto para metadatos.
// Uso: Badge("bash", Secondary)
func Badge(text string, color lipgloss.TerminalColor) string {
	return lipgloss.NewStyle().
		Foreground(color).
		Faint(true).
		Render("[" + text + "]")
}

// HumanTime formatea un timestamp de forma relativa si es reciente.
//
//	<1m  → "ahora"
//	<1h  → "hace 5m"
//	<24h → "hace 3h"
//	<7d  → "hace 2d"
//	≥7d  → "2026-09-05"
func HumanTime(t time.Time) string {
	if t.IsZero() {
		return "—"
	}
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "ahora"
	case d < time.Hour:
		return fmt.Sprintf("hace %dm", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("hace %dh", int(d.Hours()))
	case d < 7*24*time.Hour:
		return fmt.Sprintf("hace %dd", int(d.Hours()/24))
	default:
		return t.Local().Format("2006-01-02")
	}
}

// HumanDuration formatea una duración de forma compacta.
//
//	<1s  → "850ms"
//	<1m  → "1.2s"
//	≥1m  → "1m23s"
func HumanDuration(d time.Duration) string {
	switch {
	case d < time.Second:
		return fmt.Sprintf("%dms", d.Milliseconds())
	case d < time.Minute:
		return fmt.Sprintf("%.1fs", d.Seconds())
	default:
		m := int(d.Minutes())
		s := int(d.Seconds()) % 60
		return fmt.Sprintf("%dm%02ds", m, s)
	}
}

// Styles base para reusar en components (evita recrear estilo por render).

var (
	IconSuccess = lipgloss.NewStyle().Foreground(Success)
	IconWarning = lipgloss.NewStyle().Foreground(Warning)
	IconError   = lipgloss.NewStyle().Foreground(Error)
	IconInfo    = lipgloss.NewStyle().Foreground(Primary)
)
