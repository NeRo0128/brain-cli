package tool

import (
	"context"
	"time"
)

// Result es el resultado de ejecutar un Tool.
type Result struct {
	Output   string        // Output: stdout + stderr combinados, en orden de emisión.
	ExitCode int           // ExitCode: 0 = éxito. -1 = el proceso fue matado (timeout/cancel).
	Duration time.Duration // Duration: tiempo total de la ejecución.
}

// Executor ejecuta un Tool concreto.
//
// Contrato:
//   - error != nil SOLO si no se pudo lanzar el proceso
//     (intérprete no instalado, archivo temporal no creado, etc.).
//   - ExitCode != 0 NO es un error de Go — es una ejecución válida que falló.
//   - Cancelación por ctx → ExitCode = -1, error = nil, con nota en Output.
type Executor interface {
	Execute(ctx context.Context, t *Tool, params map[string]string) (Result, error)
}
