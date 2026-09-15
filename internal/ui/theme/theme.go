// Package theme define la paleta de colores de la aplicación.
//
// Cada tema tiene una paleta Dark obligatoria y una Light opcional.
// El constructor de estilos decide cuál usar según el terminal.
package theme

import (
	"image/color"
)

// Palette es un conjunto de colores para UN modo (dark o light).
type Palette struct {
	// ============================================================
	// Jerarquía de texto — 4 niveles
	// ============================================================
	Text     color.Color // texto normal
	Emphasis color.Color // texto destacado (más brillante que Text)
	Muted    color.Color // texto secundario (subtítulos, keys)
	Faint    color.Color // texto muy tenue (hints, placeholders, bordes suaves)

	// ============================================================
	// Marca — 3 acentos
	// ============================================================
	Primary   color.Color // accent principal (selección, brand)
	Secondary color.Color // accent secundario (keys, links)
	Tertiary  color.Color // tercer acento (tags, highlights, decoración)

	// ============================================================
	// Superficies — 4 tonos
	// ============================================================
	SelectedRow color.Color // fila seleccionada (bg)
	Border      color.Color // bordes normales
	BorderFocus color.Color // borde con foco (accent)
	StatusBar   color.Color // barra de estado (bg)

	// ============================================================
	// Estados semánticos — 4
	// ============================================================
	Success color.Color
	Warning color.Color
	Error   color.Color
	Info    color.Color // azul informativo, distinto de Secondary

	// ============================================================
	// Prioridad (badges) — 3
	// ============================================================
	PriorityHigh   color.Color
	PriorityMedium color.Color
	PriorityLow    color.Color

	// ============================================================
	// Tipos de task (badges) — 3
	// ============================================================
	TypeScript  color.Color
	TypeCommand color.Color
	TypeAI      color.Color

	// ============================================================
	// Inputs — 2
	// ============================================================
	InputFocused color.Color
	InputBlurred color.Color
}

// Theme agrupa una paleta dark + opcional light.
type Theme struct {
	Name        string
	Description string

	Dark  Palette
	Light *Palette
}

// IsDarkOnly indica si el tema no tiene variante light.
func (t Theme) IsDarkOnly() bool { return t.Light == nil }

// Resolve devuelve la paleta adecuada para el modo dado.
func (t Theme) Resolve(isDark bool) Palette {
	if isDark || t.Light == nil {
		return t.Dark
	}
	return *t.Light
}
