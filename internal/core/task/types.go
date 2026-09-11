package task

type TaskType string

const (
	TaskTypeScript  TaskType = "script"
	TaskTypeCommand TaskType = "command"
	TaskTypeAI      TaskType = "ai"
)

func (t TaskType) IsValid() bool {
	switch t {
	case TaskTypeScript, TaskTypeCommand, TaskTypeAI:
		return true
	}
	return false
}

func (t TaskType) RequiresTool() bool {
	return t == TaskTypeScript || t == TaskTypeCommand
}

type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

// "" se trata como medium (default razonable).
func (p Priority) IsValid() bool {
	switch p {
	case "", PriorityLow, PriorityMedium, PriorityHigh:
		return true
	}
	return false
}
