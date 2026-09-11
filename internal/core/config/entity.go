package config

// Config es la configuración global de la aplicación.
// Se carga desde un archivo YAML tras resolver variables de entorno.
type Config struct {
	App      AppConfig      `yaml:"app"`
	Database DatabaseConfig `yaml:"database"`
	Logging  LoggingConfig  `yaml:"logging"`
	UI       UIConfig       `yaml:"ui"`
	Scripts  ScriptsConfig  `yaml:"scripts"`
}

// AppConfig contiene metadatos de la aplicación.
type AppConfig struct {
	Name string `yaml:"name"`
}

// DatabaseConfig configura la capa de persistencia.
type DatabaseConfig struct {
	Path        string `yaml:"path"`         // Path al archivo SQLite (relativo o absoluto)
	AutoMigrate bool   `yaml:"auto_migrate"` // AutoMigrate ejecuta migraciones pendientes al arrancar.
}

// LoggingConfig controla el sistema de logs.
type LoggingConfig struct {
	Level      string `yaml:"level"`       // Level: debug | info | warn | error
	Format     string `yaml:"format"`      // Format: json | pretty
	OutputPath string `yaml:"output_path"` // OutputPath vacío significa stdout.
	MaxSizeMB  int    `yaml:"max_size_mb"` // MaxSizeMB antes de rotar (0 = sin rotación).
}

// UIConfig controla la interfaz de terminal.
type UIConfig struct {
	Theme string `yaml:"theme"` // Theme: dark | light
}

// ScriptsConfig define dónde viven los scripts de usuario.
type ScriptsConfig struct {
	CustomDir string `yaml:"custom_dir"` // CustomDir es la carpeta externa con scripts personalizados.
}
