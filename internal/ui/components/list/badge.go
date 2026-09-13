// Package list provee un delegate unificado y helpers de render
// para todas las listas de la app (tasks, executions, tools).
package list

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/NeRo0128/brain-cli/internal/ui/styles"
)

// BadgeStyle clasifica el color semántico de un badge.
type BadgeStyle int

const (
	BadgeNeutral BadgeStyle = iota
	BadgeInfo
	BadgeSuccess
	BadgeWarning
	BadgeDanger
)

// Badge es una etiqueta discreta: [texto]
type Badge struct {
	Text  string
	Style BadgeStyle
}

// Render devuelve el badge con su color.
func (b Badge) Render() string {
	return lipgloss.NewStyle().
		Foreground(colorForBadge(b.Style)).
		Faint(true).
		Render("[" + b.Text + "]")
}

func colorForBadge(s BadgeStyle) lipgloss.TerminalColor {
	switch s {
	case BadgeInfo:
		return styles.Secondary
	case BadgeSuccess:
		return styles.Success
	case BadgeWarning:
		return styles.Warning
	case BadgeDanger:
		return styles.Error
	default:
		return styles.Muted
	}
}

// --- Helpers semánticos ---

// TypeBadge clasifica el TaskType / ToolType.
func TypeBadge(t string) Badge {
	switch t {
	case "script", "bash":
		return Badge{Text: t, Style: BadgeInfo}
	case "command", "native":
		return Badge{Text: t, Style: BadgeNeutral}
	case "ai":
		return Badge{Text: t, Style: BadgeWarning}
	case "python", "go":
		return Badge{Text: t, Style: BadgeWarning}
	}
	return Badge{Text: t, Style: BadgeNeutral}
}

// PriorityBadge clasifica la prioridad de una Task.
func PriorityBadge(p string) Badge {
	switch p {
	case "high":
		return Badge{Text: p, Style: BadgeDanger}
	case "medium":
		return Badge{Text: p, Style: BadgeNeutral}
	case "low":
		return Badge{Text: p, Style: BadgeNeutral}
	}
	return Badge{Text: p, Style: BadgeNeutral}
}

// StatusBadge clasifica el estado de una Execution.
func StatusBadge(s string) Badge {
	switch s {
	case "completed", "success":
		return Badge{Text: "ok", Style: BadgeSuccess}
	case "failed":
		return Badge{Text: "fail", Style: BadgeDanger}
	case "cancelled":
		return Badge{Text: "cancel", Style: BadgeWarning}
	case "running":
		return Badge{Text: "running", Style: BadgeInfo}
	case "pending":
		return Badge{Text: "pending", Style: BadgeNeutral}
	}
	return Badge{Text: s, Style: BadgeNeutral}
}

// CategoryBadge clasifica la categoría de un Tool.
func CategoryBadge(c string) Badge {
	return Badge{Text: c, Style: BadgeNeutral}
}

// joinBadges renderiza una lista de badges separados por espacio.
func joinBadges(badges []Badge) string {
	if len(badges) == 0 {
		return ""
	}
	out := badges[0].Render()
	for _, b := range badges[1:] {
		out += " " + b.Render()
	}
	return out
}
