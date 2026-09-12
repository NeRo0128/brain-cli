package execution

// Status representa el estado de una Execution.
type Status string

const (
	StatusPending   Status = "pending"
	StatusRunning   Status = "running"
	StatusCompleted Status = "completed"
	StatusFailed    Status = "failed"
	StatusCancelled Status = "cancelled"
)

func (s Status) IsValid() bool {
	switch s {
	case StatusPending, StatusRunning, StatusCompleted, StatusFailed, StatusCancelled:
		return true
	}
	return false
}

// IsTerminal indica si la Execution ya no puede cambiar de estado.
func (s Status) IsTerminal() bool {
	switch s {
	case StatusCompleted, StatusFailed, StatusCancelled:
		return true
	}
	return false
}

// TriggeredBy indica qué inició la Execution.
type TriggeredBy string

const (
	TriggerManual   TriggeredBy = "manual"
	TriggerSchedule TriggeredBy = "schedule"
	TriggerChain    TriggeredBy = "chain"
)

func (t TriggeredBy) IsValid() bool {
	switch t {
	case "", TriggerManual, TriggerSchedule, TriggerChain:
		return true
	}
	return false
}
