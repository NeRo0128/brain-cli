// Package states provee bloques visuales para estados comunes
// (vacío, cargando, error) que varias screens comparten.
package states

import (
	"github.com/NeRo0128/brain-cli/internal/ui/styles"
)

// Empty devuelve un bloque con un mensaje y una pista de acción.
//
//	Sin tareas
//
//	Pulsa n para crear la primera
func Empty(s *styles.Styles, title, hint string) string {
	out := "\n\n  " + s.Subtitle.Render(title)
	if hint != "" {
		out += "\n\n  " + s.Key.Render(hint)
	}
	return out + "\n"
}

// Loading muestra un mensaje de carga.
func Loading(s *styles.Styles, what string) string {
	return "\n\n  " + s.Subtitle.Render("Cargando "+what+"...") + "\n"
}

// Error muestra un mensaje de error.
func Error(s *styles.Styles, err error) string {
	return "\n\n  " + s.ColoredIcons.Failed() + " " + err.Error() + "\n"
}
