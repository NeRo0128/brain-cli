package tool_test

import (
	"runtime"
	"testing"

	"github.com/NeRo0128/brain-cli/internal/core/tool"
)

func TestAvailable_ExcludesMissing(t *testing.T) {
	list := []tool.Interpreter{
		{Name: "bash", ScriptType: tool.ScriptTypeBash, Path: "/bin/bash"},
		{Name: "go", ScriptType: tool.ScriptTypeGo, Path: ""},
	}
	got := tool.Available(list)
	if len(got) != 1 {
		t.Fatalf("esperaba 1, dio %d", len(got))
	}
	if got[0].Name != "bash" {
		t.Errorf("esperaba bash, dio %q", got[0].Name)
	}
}

func TestAvailable_DeduplicatesByScriptType(t *testing.T) {
	list := []tool.Interpreter{
		{Name: "python3", ScriptType: tool.ScriptTypePython, Path: "/usr/bin/python3"},
		{Name: "python", ScriptType: tool.ScriptTypePython, Path: "/usr/bin/python"},
	}
	got := tool.Available(list)
	if len(got) != 1 {
		t.Fatalf("esperaba 1 (dedup), dio %d", len(got))
	}
	if got[0].Name != "python3" {
		t.Errorf("esperaba python3, dio %q", got[0].Name)
	}
}

func TestAvailable_EmptyInput(t *testing.T) {
	if got := tool.Available(nil); len(got) != 0 {
		t.Errorf("esperaba 0, dio %d", len(got))
	}
}

func TestMissing(t *testing.T) {
	list := []tool.Interpreter{
		{Name: "bash", Path: "/bin/bash"},
		{Name: "go", Path: ""},
		{Name: "python3", Path: ""},
	}
	got := tool.Missing(list)
	if len(got) != 2 {
		t.Fatalf("esperaba 2, dio %d", len(got))
	}
}

func TestDetect_FindsBash(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("bash no está en Windows por defecto")
	}
	list := tool.Detect()
	found := false
	for _, i := range list {
		if i.Name == "bash" && i.Available() {
			found = true
			break
		}
	}
	if !found {
		t.Error("bash debería estar disponible en este sistema")
	}
}

func TestInterpreter_Available(t *testing.T) {
	if !(tool.Interpreter{Path: "/bin/bash"}).Available() {
		t.Error("con path debería ser available")
	}
	if (tool.Interpreter{Path: ""}).Available() {
		t.Error("sin path no debería ser available")
	}
}
