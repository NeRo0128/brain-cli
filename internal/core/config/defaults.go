package config

// Default devuelve una configuración con valores sensatos.
// Al hacer Unmarshal encima, solo se sobreescriben los campos
// que estén presentes en el YAML.
func Default() *Config {
	return &Config{
		App: AppConfig{
			Name: "brain-cli",
		},
		Database: DatabaseConfig{
			Path:        "data/brain.db",
			AutoMigrate: true,
		},
		Logging: LoggingConfig{
			Level:     "debug",
			Format:    "pretty",
			MaxSizeMB: 100,
		},
		UI: UIConfig{
			Theme: "dark",
		},
		Scripts: ScriptsConfig{
			CustomDir: "scripts",
		},
	}
}
