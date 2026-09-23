package config

// Writer persiste la configuración a disco.
// Las implementaciones deben escribir atómicamente: nunca dejar el
// archivo de configuración en un estado parcial.
type Writer interface {
	Write(path string, cfg *Config) error
}
