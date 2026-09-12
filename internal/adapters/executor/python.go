package executor

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/NeRo0128/brain-cli/internal/core/tool"
)

// PythonExecutor ejecuta scripts python con python3 (o python).
type PythonExecutor struct{}

func NewPythonExecutor() *PythonExecutor { return &PythonExecutor{} }

func (e *PythonExecutor) Execute(ctx context.Context, t *tool.Tool, params map[string]string) (tool.Result, error) {
	scriptPath, cleanup, err := resolveScript(t, ".py")
	if err != nil {
		return tool.Result{}, err
	}
	defer cleanup()

	// python3 preferido; fallback a python.
	py, err := exec.LookPath("python3")
	if err != nil {
		py, err = exec.LookPath("python")
		if err != nil {
			return tool.Result{}, fmt.Errorf("python3 o python no encontrados en PATH")
		}
	}

	cmd := exec.CommandContext(ctx, py, scriptPath)
	cmd.Env = buildEnv(params)
	return run(ctx, cmd)
}
