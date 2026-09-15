package toast

import (
	"strings"

	"github.com/charmbracelet/x/ansi"
)

// Overlay superpone `layer` (los toasts) sobre `base` (la vista actual).
//
// `startLine` indica desde qué línea del base empezar a pegar.
//
// Alineación:
//   - width >= 80: a la derecha (con margen de 1 col)
//   - width <  80: centrado
//
// **Overlay real**: si el contenido de esa línea no cabe junto al
// toast, se trunca el contenido para hacer sitio. Los toasts NUNCA
// se salen del ancho del terminal.
func Overlay(base, layer string, width, startLine int) string {
	if layer == "" || width <= 0 {
		return base
	}
	if startLine < 0 {
		startLine = 0
	}

	layerLines := strings.Split(layer, "\n")
	baseLines := strings.Split(base, "\n")

	for i, ll := range layerLines {
		target := startLine + i
		if target >= len(baseLines) {
			break
		}
		baseLines[target] = placeAligned(baseLines[target], ll, width)
	}
	return strings.Join(baseLines, "\n")
}

// placeAligned coloca `right` dentro de `width`, respetando `left`
// como contenido previo. Si no caben juntos, `left` se trunca.
func placeAligned(left, right string, width int) string {
	rightW := ansi.StringWidth(right)

	// Si el toast es más ancho que el terminal, solo mostramos el toast.
	if rightW >= width {
		return ansi.Truncate(right, width, "")
	}

	// Espacio máximo que puede ocupar el contenido de la izquierda.
	maxLeftW := width - rightW - 1
	if maxLeftW < 0 {
		maxLeftW = 0
	}

	// Truncar el contenido si es más ancho que el hueco disponible.
	leftW := ansi.StringWidth(left)
	if leftW > maxLeftW {
		left = ansi.Truncate(left, maxLeftW, "")
		leftW = maxLeftW
	}

	if width < 80 {
		// Centrado: repartir el espacio libre a ambos lados.
		total := width - leftW - rightW
		if total < 0 {
			total = 0
		}
		leftPad := total / 2
		return left + strings.Repeat(" ", leftPad) + right
	}

	// Derecha: rellenar hasta que el toast toque el borde con 1 col.
	gap := width - leftW - rightW - 1
	if gap < 0 {
		gap = 0
	}
	return left + strings.Repeat(" ", gap) + right
}
