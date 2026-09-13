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

//	TaskKind es la elección de alto nivel del usuario en el form.
//
// No se persiste — se traduce a TaskType al guardar y se deriva
// desde TaskType al editar.
type TaskKind string

const (
	KindScript  TaskKind = "script"
	KindCommand TaskKind = "command"
	KindAI      TaskKind = "ai"
)

// IsValid indica si el kind es reconocido.
func (k TaskKind) IsValid() bool {
	switch k {
	case KindScript, KindCommand, KindAI:
		return true
	}
	return false
}

// ToTaskType mapea la elección de UI a TaskType.
func (k TaskKind) ToTaskType() TaskType {
	switch k {
	case KindScript:
		return TaskTypeScript
	case KindCommand:
		return TaskTypeCommand
	case KindAI:
		return TaskTypeAI
	}
	return ""
}

//	TaskKindFromType deriva el kind a partir de un TaskType.
//
// Útil al abrir el form en modo edición.
func TaskKindFromType(t TaskType) TaskKind {
	switch t {
	case TaskTypeScript:
		return KindScript
	case TaskTypeCommand:
		return KindCommand
	case TaskTypeAI:
		return KindAI
	}
	return ""
}
