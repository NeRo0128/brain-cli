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
	Text      color.Color
	Muted     color.Color
	Primary   color.Color
	Secondary color.Color

	SelectedRow color.Color
	Border      color.Color
	StatusBar   color.Color

	Success color.Color
	Warning color.Color
	Error   color.Color

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
