// Package progress provee una barra de progreso simple.
//
// Modos:
//   - Determinado: SetPercent(p) → barra llena proporcional.
//   - Indeterminado: WithOffset(n) → segmento brillante animado.
package progress

import (
	"strings"

	"charm.land/lipgloss/v2"

	"github.com/NeRo0128/brain-cli/internal/ui/styles"
)

// Model es una barra horizontal.
type Model struct {
	width      int
	percent    float64
	determined bool
	offset     int
	styles     *styles.Styles
}

// New construye una barra indeterminada de `width` columnas.
func New(width int, s *styles.Styles) Model {
	if width < 4 {
		width = 4
	}
	return Model{width: width, styles: s}
}

// WithWidth devuelve una copia con nuevo ancho.
func (m Model) WithWidth(w int) Model {
	if w < 4 {
		w = 4
	}
	m.width = w
	return m
}

// WithPercent marca la barra como determinada.
// p se clampa a [0,1].
func (m Model) WithPercent(p float64) Model {
	m.determined = true
	if p < 0 {
		p = 0
	}
	if p > 1 {
		p = 1
	}
	m.percent = p
	return m
}

// WithOffset cambia el offset de animación (modo indeterminado).
// Llamar repetidamente con valores crecientes produce el movimiento.
func (m Model) WithOffset(o int) Model {
	m.offset = o
	return m
}

// View renderiza la barra según el modo.
func (m Model) View() string {
	if m.width <= 0 {
		return ""
	}
	p := m.styles.Theme.Resolve(m.styles.Dark)

	fillStyle := lipgloss.NewStyle().Foreground(p.Primary)
	dimStyle := lipgloss.NewStyle().Foreground(p.Muted).Faint(true)
	brightStyle := lipgloss.NewStyle().Foreground(p.Primary).Bold(true)

	if m.determined {
		return renderDeterminate(m.width, m.percent, fillStyle, dimStyle)
	}
	return renderIndeterminate(m.width, m.offset, dimStyle, brightStyle)
}

// --- Implementación ---

func renderDeterminate(w int, p float64, fill, dim lipgloss.Style) string {
	filled := int(float64(w) * p)
	if filled > w {
		filled = w
	}
	if filled < 0 {
		filled = 0
	}
	return fill.Render(strings.Repeat("━", filled)) +
		dim.Render(strings.Repeat("━", w-filled))
}

// renderIndeterminate dibuja un segmento brillante que se mueve
// de izquierda a derecha en un ciclo continuo.
func renderIndeterminate(w, offset int, dim, bright lipgloss.Style) string {
	seg := w / 5
	if seg < 3 {
		seg = 3
	}
	if seg > 12 {
		seg = 12
	}

	cycle := w + seg
	pos := offset % cycle
	if pos < 0 {
		pos += cycle
	}
	start := pos - seg
	if start < 0 {
		start = 0
	}
	end := pos
	if end > w {
		end = w
	}
	if start > w {
		start = w
	}

	return dim.Render(strings.Repeat("━", start)) +
		bright.Render(strings.Repeat("━", end-start)) +
		dim.Render(strings.Repeat("━", w-end))
}
