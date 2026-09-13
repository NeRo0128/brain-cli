package ui

import (
	core "github.com/NeRo0128/brain-cli/internal/core/config"
	"github.com/NeRo0128/brain-cli/internal/core/execution"
	coretask "github.com/NeRo0128/brain-cli/internal/core/task"
	"github.com/NeRo0128/brain-cli/internal/core/tool"
	"github.com/NeRo0128/brain-cli/internal/ui/keys"
	taskuc "github.com/NeRo0128/brain-cli/internal/usecases/task"
	tooluc "github.com/NeRo0128/brain-cli/internal/usecases/tool"
	"github.com/rs/zerolog"
)

// Deps agrupa las dependencias compartidas por el Model.
// Se pasa por constructor, no se guarda como global.
type Deps struct {
	Version     string
	Cfg         *core.Config
	ExecUC      *taskuc.Executor
	TaskRepo    coretask.Repository
	ToolRepo    tool.Repository
	Manager     *taskuc.Manager
	ToolManager *tooluc.Manager
	ExecRepo    execution.Repository
	Keys        *keys.Registry
	Log         zerolog.Logger

	Interpreter []tool.Interpreter
}
