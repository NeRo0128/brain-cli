package executor

import (
	"context"
	"io"
)

// Executor define la interfaz para ejecutar herramientas/scripts
type Executor interface {
	// Execute ejecuta una herramienta con los parámetros dados
	Execute(ctx context.Context, scriptContent string, params map[string]string) (*Output, error)
	
	// IsAvailable verifica si el ejecutor está disponible en el sistema
	IsAvailable() bool
	
	// GetType retorna el tipo de ejecutor
	GetType() string
}

// Output representa la salida de una ejecución
type Output struct {
	Stdout   string
	Stderr   string
	ExitCode int
	Duration int64 // milisegundos
}

// StreamingExecutor extiende Executor para soportar streaming de output
type StreamingExecutor interface {
	Executor
	
	// ExecuteStreaming ejecuta con streaming de stdout/stderr
	ExecuteStreaming(ctx context.Context, scriptContent string, params map[string]string, stdout, stderr io.Writer) (*Output, error)
}

// TODO: Implementar ejecutores concretos en:
// - bash.go (BashExecutor)
// - python.go (PythonExecutor)
// - native.go (NativeExecutor)
