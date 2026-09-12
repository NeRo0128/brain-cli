package executor

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/NeRo0128/brain-cli/internal/core/tool"
)

// paramPrefix es el prefijo de las env vars inyectadas.
// Un param {SSID: "MiCasa"} → BRAIN_PARAM_SSID=MiCasa
const paramPrefix = "BRAIN_PARAM_"

// syncBuffer es un io.Writer thread-safe.
// Necesario porque cmd.Stdout y cmd.Stderr escriben desde
// goroutines distintas cuando apuntan al mismo writer.
type syncBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *syncBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *syncBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}

// run ejecuta un exec.Cmd con manejo uniforme de timeout,
// cancelación y captura de output.
//
// Es la pieza compartida por bash/python/native.
func run(ctx context.Context, cmd *exec.Cmd) (tool.Result, error) {
	start := time.Now()

	var buf syncBuffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	cmd.Stdin = nil

	// WaitDelay: tras matar el proceso por timeout, esperar hasta
	// 5s a que libere los pipes antes de forzar cierre.
	// Reducido a 200ms para que los tests de timeout sean rápidos.
	cmd.WaitDelay = 200 * time.Millisecond

	err := cmd.Run()
	result := tool.Result{
		Output:   buf.String(),
		ExitCode: 0,
		Duration: time.Since(start),
	}

	// PRIORIDAD 1: chequear el ctx ANTES que err.
	// En Go 1.20+, exec.CommandContext puede devolver *exec.ExitError
	// (por SIGKILL) en vez de context.DeadlineExceeded. El ctx es
	// la fuente de verdad.
	if ctxErr := ctx.Err(); ctxErr != nil {
		result.ExitCode = -1
		switch {
		case errors.Is(ctxErr, context.DeadlineExceeded):
			result.Output += "\n[timeout: proceso terminado]\n"
		case errors.Is(ctxErr, context.Canceled):
			result.Output += "\n[cancelado por el usuario]\n"
		}
		return result, nil
	}

	if err == nil {
		return result, nil
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		// El proceso corrió y devolvió non-zero → resultado válido.
		result.ExitCode = exitErr.ExitCode()
		return result, nil
	}

	// No pudimos ni lanzar el proceso.
	return result, fmt.Errorf("ejecutando comando: %w", err)
}

// buildEnv devuelve os.Environ() + los params con prefijo.
// Las keys se uppercased: {ssid: x} → BRAIN_PARAM_SSID=x
func buildEnv(params map[string]string) []string {
	env := os.Environ()
	for k, v := range params {
		env = append(env, fmt.Sprintf("%s%s=%s", paramPrefix, strings.ToUpper(k), v))
	}
	return env
}

// resolveScript devuelve la ruta a un archivo de script ejecutable.
//   - Si t.ScriptPath está seteado, lo usa directamente.
//   - Si t.ScriptContent está seteado, lo escribe a un archivo temporal.
//
// Devuelve también un cleanup() que borra el temporal (o no-op).
func resolveScript(t *tool.Tool, ext string) (path string, cleanup func(), err error) {
	noop := func() {}

	if t.ScriptPath != "" {
		if _, err := os.Stat(t.ScriptPath); err != nil {
			return "", noop, fmt.Errorf("script_path %q: %w", t.ScriptPath, err)
		}
		return t.ScriptPath, noop, nil
	}

	if t.ScriptContent == "" {
		return "", noop, fmt.Errorf("tool %q no tiene script_path ni script_content", t.Name)
	}

	dir, err := os.MkdirTemp("", "brain-cli-*")
	if err != nil {
		return "", noop, fmt.Errorf("creando temp dir: %w", err)
	}

	path = filepath.Join(dir, "script"+ext)
	if err := os.WriteFile(path, []byte(t.ScriptContent), 0o700); err != nil {
		os.RemoveAll(dir)
		return "", noop, fmt.Errorf("escribiendo script temporal: %w", err)
	}

	return path, func() { _ = os.RemoveAll(dir) }, nil
}
