package screens

import (
	"github.com/NeRo0128/brain-cli/internal/core/execution"
	coretask "github.com/NeRo0128/brain-cli/internal/core/task"
	"github.com/NeRo0128/brain-cli/internal/core/tool"
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

// FormSavedMsg se emite tras guardar con éxito.
// El Model pop-ea el form y envía ReloadMsg al nuevo top.
type FormSavedMsg struct{}

func FormSaved() tea.Cmd {
	return func() tea.Msg { return FormSavedMsg{} }
}

// OpenToolPickerMsg pide al Model pushear el ToolPickerScreen.
// CurrentID permite resaltar el tool actualmente seleccionado.
type OpenToolPickerMsg struct {
	CurrentID  *int
	FilterType tool.ScriptType
}

// ToolSelectedMsg lo emite el ToolPickerScreen al elegir un tool.
// El Model lo delega al nuevo top (el form que lo abrió).
type ToolSelectedMsg struct{ Tool *tool.Tool }

func OpenToolPicker(current *int, filter tool.ScriptType) tea.Cmd {
	return func() tea.Msg {
		return OpenToolPickerMsg{CurrentID: current, FilterType: filter}
	}
}

func ToolSelected(t *tool.Tool) tea.Cmd {
	return func() tea.Msg { return ToolSelectedMsg{Tool: t} }
}
