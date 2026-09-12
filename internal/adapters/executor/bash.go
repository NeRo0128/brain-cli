package executor

import (
	"context"
	"os/exec"

	"github.com/NeRo0128/brain-cli/internal/core/tool"
)

// BashExecutor ejecuta scripts bash con /bin/bash.
type BashExecutor struct{}

func NewBashExecutor() *BashExecutor { return &BashExecutor{} }

func (e *BashExecutor) Execute(ctx context.Context, t *tool.Tool, params map[string]string) (tool.Result, error) {
	scriptPath, cleanup, err := resolveScript(t, ".sh")
	if err != nil {
		return tool.Result{}, err
	}
	defer cleanup()

	cmd := exec.CommandContext(ctx, "bash", scriptPath)
	cmd.Env = buildEnv(params)
	return run(ctx, cmd)
}
