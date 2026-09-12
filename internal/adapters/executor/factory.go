package executor

import (
	"fmt"

	"github.com/NeRo0128/brain-cli/internal/core/tool"
)

// New devuelve el Executor adecuado para el ScriptType del Tool.
func New(t *tool.Tool) (tool.Executor, error) {
	switch t.ScriptType {
	case tool.ScriptTypeBash:
		return NewBashExecutor(), nil
	case tool.ScriptTypePython:
		return NewPythonExecutor(), nil
	case tool.ScriptTypeNative:
		return NewNativeExecutor(), nil
	case tool.ScriptTypeGo:
		return nil, fmt.Errorf("script_type=go aún no está soportado")
	default:
		return nil, fmt.Errorf("script_type desconocido: %q", t.ScriptType)
	}
}
