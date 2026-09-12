package executor

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/NeRo0128/brain-cli/internal/core/tool"
)

// NativeExecutor ejecuta comandos del sistema directamente vía sh -c.
// Usamos sh -c para soportar quotes y args ("docker ps -a --format ...").
type NativeExecutor struct{}

func NewNativeExecutor() *NativeExecutor { return &NativeExecutor{} }

func (e *NativeExecutor) Execute(ctx context.Context, t *tool.Tool, params map[string]string) (tool.Result, error) {
	if t.Command == "" {
		return tool.Result{}, fmt.Errorf("tool %q (native) no tiene command", t.Name)
	}

	cmd := exec.CommandContext(ctx, "sh", "-c", t.Command)
	cmd.Env = buildEnv(params)
	return run(ctx, cmd)
}
