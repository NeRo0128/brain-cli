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

// --- Confirm modal ---

// OpenConfirmMsg pide al Model pushear el ConfirmScreen.
type OpenConfirmMsg struct {
	Title   string
	Message string
	Action  tea.Msg // msg a despachar tras confirmar
}

// ConfirmYesMsg lo emite el ConfirmScreen al confirmar.
// El Model raíz poppea el confirm y despacha Action.
type ConfirmYesMsg struct {
	Action tea.Msg
}

// DeleteTaskMsg indica al Model que borre una task.
// Se usa como Action dentro de OpenConfirmMsg.
type DeleteTaskMsg struct {
	TaskID   string
	TaskName string
}

// TaskDeletedMsg es el resultado del borrado.
type TaskDeletedMsg struct {
	TaskID string
	Err    error
}

// OpenConfirm construye el comando que pushea el modal.
func OpenConfirm(title, message string, action tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return OpenConfirmMsg{Title: title, Message: message, Action: action}
	}
}

// --- Tool form ---

// OpenToolFormMsg pide al Model pushear el ToolFormScreen.
// Tool nil = crear, no-nil = editar.
type OpenToolFormMsg struct{ Tool *tool.Tool }

// ToolFormSavedMsg se emite tras guardar un tool con éxito.
// Lleva el tool creado/actualizado para que el Model pueda refrescar.
type ToolFormSavedMsg struct{ Tool *tool.Tool }

// DeleteToolMsg indica al Model que borre un tool (viene del ConfirmScreen).
type DeleteToolMsg struct {
	ToolID   int
	ToolName string
}

// ToolDeletedMsg es el resultado del borrado.
type ToolDeletedMsg struct {
	ToolID int
	Err    error
}

// ToolCreatedMsg notifica al picker que se creó un tool nuevo.
// El Model lo envía tras poppear el ToolFormScreen.
type ToolCreatedMsg struct{ ToolID int }

// --- Helpers ---

func OpenToolForm(t *tool.Tool) tea.Cmd {
	return func() tea.Msg { return OpenToolFormMsg{Tool: t} }
}
