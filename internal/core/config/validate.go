package config

import (
	"errors"
	"fmt"
)

var validLogLevels = map[string]bool{
	"debug": true, "info": true, "warn": true, "error": true,
}

var validLogFormats = map[string]bool{
	"json": true, "pretty": true,
}

var validThemes = map[string]bool{
	"dark": true, "light": true,
}

// Validate comprueba que la configuración es coherente.
// Devuelve un error combinado con TODOS los problemas encontrados,
// no solo el primero — más útil para el usuario.
func (c *Config) Validate() error {
	var errs []error

	if c.App.Name == "" {
		errs = append(errs, errors.New("app.name no puede estar vacío"))
	}
	if c.Database.Path == "" {
		errs = append(errs, errors.New("database.path no puede estar vacío"))
	}
	if !validLogLevels[c.Logging.Level] {
		errs = append(errs, fmt.Errorf("logging.level inválido: %q", c.Logging.Level))
	}
	if !validLogFormats[c.Logging.Format] {
		errs = append(errs, fmt.Errorf("logging.format inválido: %q", c.Logging.Format))
	}
	if c.Logging.MaxSizeMB < 0 {
		errs = append(errs, errors.New("logging.max_size_mb no puede ser negativo"))
	}
	if !validThemes[c.UI.Theme] {
		errs = append(errs, fmt.Errorf("ui.theme inválido: %q", c.UI.Theme))
	}

	return errors.Join(errs...)
}
