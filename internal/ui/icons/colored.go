package icons

import (
	lipgloss "charm.land/lipgloss/v2"

	"github.com/NeRo0128/brain-cli/internal/ui/theme"
)

// Colored agrupa iconos con color pre-aplicado según el tema.
type Colored struct {
	set Set
	pal theme.Palette
}

// NewColored construye el set de iconos coloreados.
func NewColored(s Set, isDark bool, t theme.Theme) Colored {
	return Colored{
		set: s,
		pal: t.Resolve(isDark),
	}
}

// --- Estados (cada uno con su color semántico) ---

func (c Colored) Success() string {
	return lipgloss.NewStyle().
		Foreground(c.pal.Success).
		Bold(true).
		Render(c.set.Success)
}

func (c Colored) Failed() string {
	return lipgloss.NewStyle().
		Foreground(c.pal.Error).
		Bold(true).
		Render(c.set.Failed)
}

func (c Colored) Cancelled() string {
	return lipgloss.NewStyle().
		Foreground(c.pal.Warning).
		Bold(true).
		Render(c.set.Cancelled)
}

func (c Colored) Running() string {
	return lipgloss.NewStyle().
		Foreground(c.pal.Primary).
		Bold(true).
		Render(c.set.Running)
}

func (c Colored) Pending() string {
	return lipgloss.NewStyle().
		Foreground(c.pal.Muted).
		Render(c.set.Pending)
}

// --- Navegación ---

func (c Colored) Selected() string {
	return lipgloss.NewStyle().
		Foreground(c.pal.Primary).
		Bold(true).
		Render(c.set.Selected)
}

func (c Colored) Favorite() string {
	return lipgloss.NewStyle().
		Foreground(c.pal.Warning).
		Render(c.set.Favorite)
}

// --- Tipos de tarea (cada uno con su color de tema) ---

func (c Colored) TypeScript() string {
	return lipgloss.NewStyle().
		Foreground(c.pal.TypeScript).
		Render(c.set.TypeScript)
}

func (c Colored) TypeCommand() string {
	return lipgloss.NewStyle().
		Foreground(c.pal.TypeCommand).
		Render(c.set.TypeCommand)
}

func (c Colored) TypeAI() string {
	return lipgloss.NewStyle().
		Foreground(c.pal.TypeAI).
		Render(c.set.TypeAI)
}

// --- Branding ---

func (c Colored) Brand() string {
	return lipgloss.NewStyle().
		Foreground(c.pal.Primary).
		Bold(true).
		Render(c.set.Brand)
}

// --- Modales ---

func (c Colored) Danger() string {
	return lipgloss.NewStyle().
		Foreground(c.pal.Error).
		Bold(true).
		Render(c.set.Danger)
}

// --- Helper: icono por estado de ejecución ---

// StatusIcon devuelve el icono coloreado según el estado de una ejecución.
// status: "completed" | "failed" | "cancelled" | "running" | "pending"
func (c Colored) StatusIcon(status string) string {
	switch status {
	case "completed", "success":
		return c.Success()
	case "failed", "error":
		return c.Failed()
	case "cancelled":
		return c.Cancelled()
	case "running":
		return c.Running()
	case "pending":
		return c.Pending()
	}
	return lipgloss.NewStyle().Foreground(c.pal.Muted).Render("•")
}

// --- Extras que faltaban ---

func (c Colored) Warning() string {
	return lipgloss.NewStyle().Foreground(c.pal.Warning).Bold(true).Render(c.set.Warning)
}

func (c Colored) Info() string {
	return lipgloss.NewStyle().Foreground(c.pal.Info).Render(c.set.Info)
}

func (c Colored) Unfavorite() string {
	return lipgloss.NewStyle().Foreground(c.pal.Muted).Render(c.set.Unfavorite)
}

func (c Colored) Bullet() string {
	return lipgloss.NewStyle().Foreground(c.pal.Muted).Render(c.set.Bullet)
}
