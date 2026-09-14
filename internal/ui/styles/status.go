package styles

import (
	"fmt"
	"image/color"
	"time"

	"github.com/NeRo0128/brain-cli/internal/ui/theme"
)

// HumanTime formatea un timestamp de forma relativa si es reciente.
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

// ColorForStatus devuelve el color semántico para un estado.
func ColorForStatus(status string, p theme.Palette) color.Color {
	switch status {
	case "completed", "success", "ok":
		return p.Success
	case "failed", "error":
		return p.Error
	case "cancelled", "canceled", "warn":
		return p.Warning
	case "running", "pending":
		return p.Primary
	}
	return p.Muted
}

// IconForStatus devuelve el icono para un estado.
func IconForStatus(status string) string {
	switch status {
	case "completed", "success", "ok":
		return "✓"
	case "failed", "error":
		return "✗"
	case "cancelled", "canceled", "warn":
		return "⊘"
	case "running", "pending":
		return "◐"
	}
	return "•"
}
