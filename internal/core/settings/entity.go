package settings

import "time"

// Setting es un override de configuración persistido en la DB.
// La key sigue el formato "<sección>.<campo>" (ej: "ui.theme").
type Setting struct {
	Key         string
	Value       string
	Category    string
	IsEncrypted bool
	UpdatedAt   time.Time
}
