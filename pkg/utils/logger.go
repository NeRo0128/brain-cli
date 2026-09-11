package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LoggerOptions controla cómo se construye el logger.
type LoggerOptions struct {
	Level      string // Level: debug | info | warn | error
	Format     string // Format: json (producción) | pretty (desarrollo)
	OutputPath string // OutputPath vacío → stderr. Si se define, escribe a archivo.
	MaxSizeMB  int    // MaxSizeMB antes de rotar el archivo. 0 = sin rotación.
}

// NewLogger construye un logger zerolog según opts.
// Devuelve error si el nivel es inválido o si no se puede crear
// el archivo/directorio de salida.
func NewLogger(opts LoggerOptions) (zerolog.Logger, error) {
	level, err := zerolog.ParseLevel(opts.Level)
	if err != nil {
		return zerolog.Nop(), fmt.Errorf("nivel de log inválido %q: %w", opts.Level, err)
	}

	out, err := buildWriter(opts)
	if err != nil {
		return zerolog.Nop(), err
	}

	// En modo pretty envolvemos el writer con el ConsoleWriter de zerolog.
	if opts.Format == "pretty" {
		out = zerolog.ConsoleWriter{
			Out:        out,
			TimeFormat: time.Kitchen,
		}
	}

	logger := zerolog.New(out).Level(level).With().Timestamp().Logger()
	return logger, nil
}

// NopLogger devuelve un logger que descarta todo.
// Útil en tests y en código de librería que no debe ruidear.
func NopLogger() zerolog.Logger {
	return zerolog.Nop()
}

// buildWriter decide a dónde va el log y devuelve el io.Writer.
func buildWriter(opts LoggerOptions) (io.Writer, error) {
	// Sin OutputPath → stderr (NUNCA stdout, ver doc del struct).
	if opts.OutputPath == "" {
		return os.Stderr, nil
	}

	// Aseguramos que el directorio padre exista.
	dir := filepath.Dir(opts.OutputPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("creando directorio de logs %q: %w", dir, err)
	}

	// Sin rotación configurada: abrimos el archivo directo en append.
	if opts.MaxSizeMB <= 0 {
		f, err := os.OpenFile(opts.OutputPath,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return nil, fmt.Errorf("abriendo archivo de logs: %w", err)
		}
		return f, nil
	}

	// Con rotación: lumberjack se encarga de todo (crea el archivo on-demand).
	return &lumberjack.Logger{
		Filename:   opts.OutputPath,
		MaxSize:    opts.MaxSizeMB,
		MaxBackups: 3,
		MaxAge:     28, // días
		Compress:   true,
	}, nil
}
