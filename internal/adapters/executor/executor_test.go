package executor_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/NeRo0128/brain-cli/internal/adapters/executor"
	"github.com/NeRo0128/brain-cli/internal/core/tool"
)

func bashTool(content string) *tool.Tool {
	return &tool.Tool{
		Name:           "test-bash",
		ScriptType:     tool.ScriptTypeBash,
		Category:       tool.CategoryUtils,
		ScriptContent:  content,
		TimeoutSeconds: 10,
		Version:        1,
	}
}

func TestBash_SimpleOutput(t *testing.T) {
	e := executor.NewBashExecutor()
	res, err := e.Execute(context.Background(), bashTool("echo hello"), nil)
	if err != nil {
		t.Fatal(err)
	}
	if res.ExitCode != 0 {
		t.Errorf("exit code = %d, quiero 0", res.ExitCode)
	}
	if !strings.Contains(res.Output, "hello") {
		t.Errorf("output: %q", res.Output)
	}
}

func TestBash_NonZeroExit_IsNotGoError(t *testing.T) {
	e := executor.NewBashExecutor()
	res, err := e.Execute(context.Background(), bashTool("exit 3"), nil)
	if err != nil {
		t.Fatalf("non-zero exit NO debe ser error de Go: %v", err)
	}
	if res.ExitCode != 3 {
		t.Errorf("exit code = %d, quiero 3", res.ExitCode)
	}
}

func TestBash_Params(t *testing.T) {
	e := executor.NewBashExecutor()
	tl := bashTool(`echo SSID=$BRAIN_PARAM_SSID`)
	res, err := e.Execute(context.Background(), tl, map[string]string{"SSID": "MiCasa"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(res.Output, "SSID=MiCasa") {
		t.Errorf("params no inyectados: %q", res.Output)
	}
}

func TestBash_CombinesStdoutAndStderr(t *testing.T) {
	e := executor.NewBashExecutor()
	res, err := e.Execute(context.Background(), bashTool(`echo to-stdout; echo to-stderr >&2`), nil)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(res.Output, "to-stdout") || !strings.Contains(res.Output, "to-stderr") {
		t.Errorf("output incompleto: %q", res.Output)
	}
}

func TestBash_Timeout(t *testing.T) {
	e := executor.NewBashExecutor()
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	res, err := e.Execute(ctx, bashTool("sleep 5; echo never"), nil)
	if err != nil {
		t.Fatal(err)
	}
	if res.ExitCode != -1 {
		t.Errorf("exit code = %d, quiero -1", res.ExitCode)
	}
	if !strings.Contains(res.Output, "timeout") {
		t.Errorf("output no menciona timeout: %q", res.Output)
	}
}

func TestBash_Cancellation(t *testing.T) {
	e := executor.NewBashExecutor()
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	res, err := e.Execute(ctx, bashTool("sleep 5; echo never"), nil)
	if err != nil {
		t.Fatal(err)
	}
	if res.ExitCode != -1 {
		t.Errorf("exit code = %d, quiero -1", res.ExitCode)
	}
	if !strings.Contains(res.Output, "cancelado") {
		t.Errorf("output: %q", res.Output)
	}
}

func TestBash_MissingInterpreter_FailsFast(t *testing.T) {
	e := executor.NewBashExecutor()
	// Script content vacío → debe fallar en resolveScript
	res, err := e.Execute(context.Background(), &tool.Tool{
		Name: "empty", ScriptType: tool.ScriptTypeBash,
	}, nil)
	if err == nil {
		t.Fatal("esperaba error por script vacío")
	}
	if res.ExitCode != 0 {
		t.Errorf("exit code debería ser 0 en error de infraestructura")
	}
}

func TestNative_Echo(t *testing.T) {
	e := executor.NewNativeExecutor()
	tl := &tool.Tool{
		Name:       "native",
		ScriptType: tool.ScriptTypeNative,
		Command:    `echo native-works`,
	}
	res, err := e.Execute(context.Background(), tl, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(res.Output, "native-works") {
		t.Errorf("output: %q", res.Output)
	}
}

func TestNative_WithArgs(t *testing.T) {
	e := executor.NewNativeExecutor()
	tl := &tool.Tool{
		Name:       "native-args",
		ScriptType: tool.ScriptTypeNative,
		Command:    `printf '%s\n' one two three`,
	}
	res, err := e.Execute(context.Background(), tl, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, w := range []string{"one", "two", "three"} {
		if !strings.Contains(res.Output, w) {
			t.Errorf("falta %q en output: %q", w, res.Output)
		}
	}
}

func TestNative_EmptyCommand_FailsFast(t *testing.T) {
	e := executor.NewNativeExecutor()
	_, err := e.Execute(context.Background(), &tool.Tool{
		Name: "empty", ScriptType: tool.ScriptTypeNative,
	}, nil)
	if err == nil {
		t.Fatal("esperaba error por command vacío")
	}
}

func TestFactory_ReturnsCorrectExecutor(t *testing.T) {
	cases := map[tool.ScriptType]string{
		tool.ScriptTypeBash:   "*executor.BashExecutor",
		tool.ScriptTypePython: "*executor.PythonExecutor",
		tool.ScriptTypeNative: "*executor.NativeExecutor",
	}
	for st, want := range cases {
		tl := &tool.Tool{ScriptType: st, Name: "x"}
		ex, err := executor.New(tl)
		if err != nil {
			t.Errorf("%s: %v", st, err)
			continue
		}
		got := fmt.Sprintf("%T", ex)
		if got != want {
			t.Errorf("%s: got %T, want %s", st, ex, want)
		}
	}
}

func TestFactory_UnknownType_Errors(t *testing.T) {
	_, err := executor.New(&tool.Tool{ScriptType: "ruby"})
	if err == nil {
		t.Fatal("esperaba error para script_type desconocido")
	}
}
