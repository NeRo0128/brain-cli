package toast

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Overlay superpone `layer` (los toasts) sobre `base` (la vista actual).
//
// Comportamiento:
//   - width >= 80: alinea a la derecha (arriba-derecha).
//   - width < 80: alinea al centro (arriba-centro) — más legible.
//
// No modifica el alto del base: pega línea por línea.
func Overlay(base, layer string, width int) string {
	if layer == "" || width <= 0 {
		return base
	}

	layerLines := strings.Split(layer, "\n")
	baseLines := strings.Split(base, "\n")

	for i, ll := range layerLines {
		if i >= len(baseLines) {
			break
		}
		baseLines[i] = placeAligned(baseLines[i], ll, width)
	}
	return strings.Join(baseLines, "\n")
}

// placeAligned coloca `right` a la derecha (o centro si width<80)
// dentro de `width`, respetando `left` como contenido previo.
func placeAligned(left, right string, width int) string {
	leftW := lipgloss.Width(left)
	rightW := lipgloss.Width(right)

	if width < 80 {
		// Centrado: calculamos el padding total y repartimos.
		total := width - leftW - rightW - 1
		if total < 0 {
			return left + " " + right
		}
		leftPad := total / 2
		return left + strings.Repeat(" ", leftPad) + right
	}

	// Derecha: 1 col de margen desde el borde.
	remaining := width - leftW - 1
	if remaining < rightW {
		return left + " " + right
	}
	return left + strings.Repeat(" ", remaining-rightW) + right
}
