package execution_test

import (
	"strings"
	"testing"
	"time"

	"github.com/NeRo0128/brain-cli/internal/core/execution"
)

func intPtr(i int) *int              { return &i }
func timePtr(t time.Time) *time.Time { return &t }

func TestValidate_Running_OK(t *testing.T) {
	e := &execution.Execution{
		TaskID:      "task-x",
		Status:      execution.StatusRunning,
		TriggeredBy: execution.TriggerManual,
		StartedAt:   time.Now(),
	}
	if err := e.Validate(); err != nil {
		t.Fatalf("esperaba válido, dio: %v", err)
	}
}

func TestValidate_Completed_RequiresFinishedAt(t *testing.T) {
	e := &execution.Execution{
		TaskID:      "task-x",
		Status:      execution.StatusCompleted,
		TriggeredBy: execution.TriggerManual,
		StartedAt:   time.Now(),
	}
	err := e.Validate()
	if err == nil || !strings.Contains(err.Error(), "finished_at") {
		t.Fatalf("esperaba error de finished_at, dio: %v", err)
	}
}

func TestValidate_Running_RejectsFinishedAt(t *testing.T) {
	now := time.Now()
	e := &execution.Execution{
		TaskID:      "task-x",
		Status:      execution.StatusRunning,
		TriggeredBy: execution.TriggerManual,
		StartedAt:   now,
		FinishedAt:  &now,
	}
	err := e.Validate()
	if err == nil || !strings.Contains(err.Error(), "no admite finished_at") {
		t.Fatalf("esperaba error, dio: %v", err)
	}
}

func TestValidate_Completed_RejectsError(t *testing.T) {
	now := time.Now()
	e := &execution.Execution{
		TaskID:      "task-x",
		Status:      execution.StatusCompleted,
		TriggeredBy: execution.TriggerManual,
		StartedAt:   now,
		FinishedAt:  &now,
		Error:       "algo salió mal",
	}
	err := e.Validate()
	if err == nil || !strings.Contains(err.Error(), "no admite error") {
		t.Fatalf("esperaba error, dio: %v", err)
	}
}

func TestValidate_Failed_OK(t *testing.T) {
	now := time.Now()
	e := &execution.Execution{
		TaskID:      "task-x",
		Status:      execution.StatusFailed,
		TriggeredBy: execution.TriggerManual,
		StartedAt:   now,
		FinishedAt:  &now,
		Error:       "exit status 1",
		ExitCode:    intPtr(1),
	}
	if err := e.Validate(); err != nil {
		t.Fatalf("esperaba válido, dio: %v", err)
	}
}

func TestValidate_Pending_RejectsExitCode(t *testing.T) {
	e := &execution.Execution{
		TaskID:      "task-x",
		Status:      execution.StatusPending,
		TriggeredBy: execution.TriggerManual,
		StartedAt:   time.Now(),
		ExitCode:    intPtr(0),
	}
	err := e.Validate()
	if err == nil || !strings.Contains(err.Error(), "exit_code") {
		t.Fatalf("esperaba error de exit_code, dio: %v", err)
	}
}

func TestStatus_IsTerminal(t *testing.T) {
	cases := map[execution.Status]bool{
		execution.StatusPending:   false,
		execution.StatusRunning:   false,
		execution.StatusCompleted: true,
		execution.StatusFailed:    true,
		execution.StatusCancelled: true,
	}
	for status, want := range cases {
		if got := status.IsTerminal(); got != want {
			t.Errorf("%q.IsTerminal() = %v, quiero %v", status, got, want)
		}
	}
}

func TestDuration_NotFinished(t *testing.T) {
	e := &execution.Execution{StartedAt: time.Now().Add(-2 * time.Second)}
	if d := e.Duration(); d < 2*time.Second {
		t.Errorf("Duration = %v, quiero >=2s", d)
	}
}

func TestDuration_Finished(t *testing.T) {
	start := time.Now()
	end := start.Add(5 * time.Second)
	e := &execution.Execution{StartedAt: start, FinishedAt: &end}
	if d := e.Duration(); d != 5*time.Second {
		t.Errorf("Duration = %v, quiero 5s", d)
	}
}
