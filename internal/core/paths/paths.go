// Package paths resuelve rutas del sistema siguiendo XDG Base Directory.
package paths

import (
	"os"
	"path/filepath"
)

const appName = "brain-cli"

func configBase() string {
	if v := os.Getenv("XDG_CONFIG_HOME"); v != "" {
		return v
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config")
}

func dataBase() string {
	if v := os.Getenv("XDG_DATA_HOME"); v != "" {
		return v
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "share")
}

func cacheBase() string {
	if v := os.Getenv("XDG_CACHE_HOME"); v != "" {
		return v
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".cache")
}

// ConfigDir → ~/.config/brain-cli/
func ConfigDir() string { return filepath.Join(configBase(), appName) }

// DataDir → ~/.local/share/brain-cli/
func DataDir() string { return filepath.Join(dataBase(), appName) }

// CacheDir → ~/.cache/brain-cli/
func CacheDir() string { return filepath.Join(cacheBase(), appName) }

// ConfigFile → ~/.config/brain-cli/config.yaml
func ConfigFile() string { return filepath.Join(ConfigDir(), "config.yaml") }

// DBFile → ~/.local/share/brain-cli/brain.db
func DBFile() string { return filepath.Join(DataDir(), "brain.db") }

// EnsureDirs crea los directorios si faltan.
func EnsureDirs() error {
	for _, d := range []string{ConfigDir(), DataDir(), CacheDir()} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			return err
		}
	}
	return nil
}

// DevMode indica si corremos desde el repo (go.mod en CWD).
func DevMode() bool {
	_, err := os.Stat("go.mod")
	return err == nil
}
