package screens

import (
	"github.com/NeRo0128/brain-cli/internal/core/execution"
	coretask "github.com/NeRo0128/brain-cli/internal/core/task"
	tea "github.com/charmbracelet/bubbletea"
)

// ActionMsg lo emite el Model al traducir un KeyMsg que matchea
// un ID declarado por la screen activa.
type ActionMsg struct{ ID string }

// BackMsg pide al Model hacer pop de la screen superior.
type BackMsg struct{}

// ExecuteTaskMsg pide al Model lanzar la ejecución de una task.
type ExecuteTaskMsg struct {
	TaskID   string
	TaskName string
}

// OpenDetailMsg pide al Model pushear el DetailScreen.
type OpenDetailMsg struct{ Task *coretask.Task }

// OpenHistoryMsg pide al Model pushear el HistoryScreen.
type OpenHistoryMsg struct{}

// OpenResultMsg pide al Model reemplazar el top por un ResultScreen.
type OpenResultMsg struct {
	TaskName string
	Exec     *execution.Execution
}

// OpenHelpMsg pide al Model pushear el HelpScreen.
type OpenHelpMsg struct{}

// --- Helpers como tea.Cmd ---

func Back() tea.Cmd { return func() tea.Msg { return BackMsg{} } }

func ExecuteTask(id, name string) tea.Cmd {
	return func() tea.Msg { return ExecuteTaskMsg{TaskID: id, TaskName: name} }
}

func OpenDetail(t *coretask.Task) tea.Cmd {
	return func() tea.Msg { return OpenDetailMsg{Task: t} }
}

func OpenHistory() tea.Cmd { return func() tea.Msg { return OpenHistoryMsg{} } }

func OpenResult(name string, e *execution.Execution) tea.Cmd {
	return func() tea.Msg { return OpenResultMsg{TaskName: name, Exec: e} }
}

func OpenHelp() tea.Cmd { return func() tea.Msg { return OpenHelpMsg{} } }

// OpenFormMsg pide al Model pushear el FormScreen.
// Task nil = crear, no-nil = editar.
type OpenFormMsg struct{ Task *coretask.Task }

// ReloadMsg pide a la screen activa recargar sus datos.
// Lo usa el FormScreen tras guardar, para que MainScreen refresque.
type ReloadMsg struct{}

func OpenForm(t *coretask.Task) tea.Cmd {
	return func() tea.Msg { return OpenFormMsg{Task: t} }
}

func Reload() tea.Cmd {
	return func() tea.Msg { return ReloadMsg{} }
}
