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
func Empty(title, hint string) string {
	out := "\n\n  " + styles.Subtitle.Render(title)
	if hint != "" {
		out += "\n\n  " + styles.Key.Render(hint)
	}
	return out + "\n"
}

// Loading muestra un mensaje de carga.
func Loading(what string) string {
	return "\n\n  " + styles.Subtitle.Render("Cargando "+what+"...") + "\n"
}

// Error muestra un mensaje de error.
func Error(err error) string {
	return "\n\n  " + styles.ErrorStyle.Render("✗ ") + err.Error() + "\n"
}
