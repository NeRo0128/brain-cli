package config

// Manager es el contrato para cargar la configuración.
// El resto de la app depende de esta interfaz, no del loader concreto.
type Manager interface {
	// Load lee la configuración desde path, resuelve variables de
	// entorno y valida el resultado. Devuelve error si algo falla.
	Load(path string) (*Config, error)
}
