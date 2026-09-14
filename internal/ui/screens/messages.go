package screens

import (
	"github.com/NeRo0128/brain-cli/internal/core/execution"
	coretask "github.com/NeRo0128/brain-cli/internal/core/task"
	"github.com/NeRo0128/brain-cli/internal/core/tool"
	"github.com/NeRo0128/brain-cli/internal/ui/styles"
	tea "charm.land/bubbletea/v2"
)

type ActionMsg struct{ ID string }

type BackMsg struct{}

type ExecuteTaskMsg struct {
	TaskID   string
	TaskName string
}

type OpenDetailMsg struct {
	Task   *coretask.Task
	Styles *styles.Styles
}

type OpenHistoryMsg struct {
	Styles *styles.Styles
}

type OpenResultMsg struct {
	TaskName string
	Exec     *execution.Execution
}

type OpenHelpMsg struct{}

type OpenFormMsg struct {
	Task   *coretask.Task
	Styles *styles.Styles
}

type ReloadMsg struct{}

type FormSavedMsg struct{}

type OpenToolPickerMsg struct {
	CurrentID  *int
	FilterType tool.ScriptType
	Styles     *styles.Styles
}

type ToolSelectedMsg struct{ Tool *tool.Tool }

type OpenConfirmMsg struct {
	Title   string
	Message string
	Action  tea.Msg
	Styles  *styles.Styles
}

type ConfirmYesMsg struct {
	Action tea.Msg
}

type DeleteTaskMsg struct {
	TaskID   string
	TaskName string
}

type TaskDeletedMsg struct {
	TaskID string
	Err    error
}

type OpenToolFormMsg struct {
	Tool   *tool.Tool
	Styles *styles.Styles
}

type ToolFormSavedMsg struct{ Tool *tool.Tool }

type DeleteToolMsg struct {
	ToolID   int
	ToolName string
}

type ToolDeletedMsg struct {
	ToolID int
	Err    error
}

type ToolCreatedMsg struct{ ToolID int }

// --- Helpers ---

func Back() tea.Cmd { return func() tea.Msg { return BackMsg{} } }

func ExecuteTask(id, name string) tea.Cmd {
	return func() tea.Msg { return ExecuteTaskMsg{TaskID: id, TaskName: name} }
}

func OpenDetail(t *coretask.Task, s *styles.Styles) tea.Cmd {
	return func() tea.Msg { return OpenDetailMsg{Task: t, Styles: s} }
}

func OpenHistory(s *styles.Styles) tea.Cmd {
	return func() tea.Msg { return OpenHistoryMsg{Styles: s} }
}

func OpenResult(name string, e *execution.Execution) tea.Cmd {
	return func() tea.Msg { return OpenResultMsg{TaskName: name, Exec: e} }
}

func OpenHelp() tea.Cmd { return func() tea.Msg { return OpenHelpMsg{} } }

func OpenForm(t *coretask.Task, s *styles.Styles) tea.Cmd {
	return func() tea.Msg { return OpenFormMsg{Task: t, Styles: s} }
}

func Reload() tea.Cmd {
	return func() tea.Msg { return ReloadMsg{} }
}

func FormSaved() tea.Cmd {
	return func() tea.Msg { return FormSavedMsg{} }
}

func OpenToolPicker(current *int, filter tool.ScriptType, s *styles.Styles) tea.Cmd {
	return func() tea.Msg {
		return OpenToolPickerMsg{CurrentID: current, FilterType: filter, Styles: s}
	}
}

func ToolSelected(t *tool.Tool) tea.Cmd {
	return func() tea.Msg { return ToolSelectedMsg{Tool: t} }
}

func OpenConfirm(title, message string, action tea.Msg, s *styles.Styles) tea.Cmd {
	return func() tea.Msg {
		return OpenConfirmMsg{Title: title, Message: message, Action: action, Styles: s}
	}
}

func OpenToolForm(t *tool.Tool, s *styles.Styles) tea.Cmd {
	return func() tea.Msg { return OpenToolFormMsg{Tool: t, Styles: s} }
}
