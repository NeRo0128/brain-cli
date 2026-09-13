package layout

// Breakpoint clasifica el ancho de la terminal.
type Breakpoint int

const (
	// SizeSmall: < 80 columnas. Layout compacto.
	SizeSmall Breakpoint = iota
	// SizeMedium: 80-119 columnas. Layout normal.
	SizeMedium
	// SizeLarge: >= 120 columnas. Layout aireado, 2 columnas.
	SizeLarge
)

// BreakpointFor devuelve el breakpoint para un ancho dado.
func BreakpointFor(width int) Breakpoint {
	switch {
	case width < 80:
		return SizeSmall
	case width < 120:
		return SizeMedium
	default:
		return SizeLarge
	}
}

// IsCompact es atajo para el caso común.
func IsCompact(width int) bool {
	return width < 80
}

// IsWide es atajo para 2-columnas.
func IsWide(width int) bool {
	return width >= 120
}

// ContentSize devuelve el área útil de una screen, restando
// header (2 líneas) + footer (2 líneas) + padding.
//
// Header: línea de brand + línea en blanco.
// Footer: línea en blanco + línea de help.
func ContentSize(width, height int) (w, h int) {
	w = width - 2 // 1 col de padding a cada lado
	h = height - 5
	if w < 20 {
		w = 20
	}
	if h < 5 {
		h = 5
	}
	return
}

// ShouldGate indica si la terminal es demasiado chica para
// renderizar la TUI con dignidad.
//
// Regla del plan: < 70 ancho || < 20 alto.
func ShouldGate(width, height int) bool {
	return width < 70 || height < 20
}
