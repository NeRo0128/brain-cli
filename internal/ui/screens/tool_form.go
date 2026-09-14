package screens

import (
	"context"
	"strconv"
	"strings"
	"time"

	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/rs/zerolog"

	coretool "github.com/NeRo0128/brain-cli/internal/core/tool"
	"github.com/NeRo0128/brain-cli/internal/ui/keys"
	"github.com/NeRo0128/brain-cli/internal/ui/styles"
	tooluc "github.com/NeRo0128/brain-cli/internal/usecases/tool"
)

// toolSaveDoneMsg es el resultado de un Create/Update.
type toolSaveDoneMsg struct {
	saved *coretool.Tool
	err   error
}

type toolFieldID string

const (
	tfName     toolFieldID = "name"
	tfDesc     toolFieldID = "desc"
	tfType     toolFieldID = "type"
	tfCategory toolFieldID = "category"
	tfContent  toolFieldID = "content"
	tfCommand  toolFieldID = "command"
	tfTimeout  toolFieldID = "timeout"
	tfSudo     toolFieldID = "sudo"
)

type toolFieldKind int

const (
	tfKindText toolFieldKind = iota
	tfKindSelect
	tfKindToggle
)

var toolFieldOrder = []toolFieldID{
	tfName, tfDesc, tfType, tfCategory, tfContent, tfCommand, tfTimeout, tfSudo,
}

// ToolFormScreen crea o edita un Tool.
type ToolFormScreen struct {
	editing *coretool.Tool

	nameInput    textinput.Model
	descInput    textinput.Model
	timeoutInput textinput.Model
	scriptArea   textarea.Model
	commandInput textinput.Model

	scriptType coretool.ScriptType
	category   coretool.Category
	sudo       bool

	focus int

	err    error
	saving bool

	manager *tooluc.Manager
	log     zerolog.Logger

	width, height int
	styles        *styles.Styles
}

func NewToolFormScreen(
	tk *coretool.Tool,
	manager *tooluc.Manager,
	log zerolog.Logger,
	s *styles.Styles,
) ToolFormScreen {
	action := "crear"
	if tk != nil {
		action = "editar"
	}
	screenLog := log.With().Str("screen", "tool_form").Str("action", action).Logger()
	p := s.Theme.Resolve(s.Dark)

	f := ToolFormScreen{
		editing:    tk,
		manager:    manager,
		log:        screenLog,
		scriptType: coretool.ScriptTypeBash,
		category:   coretool.CategoryUtils,
		styles:     s,
	}

	f.nameInput = textinput.New()
	f.nameInput.Placeholder = "mi-tool"
	f.nameInput.CharLimit = 128
	f.nameInput.SetWidth(60)

	f.descInput = textinput.New()
	f.descInput.Placeholder = "Descripción breve..."
	f.descInput.CharLimit = 256
	f.descInput.SetWidth(80)

	f.timeoutInput = textinput.New()
	f.timeoutInput.Placeholder = "300"
	f.timeoutInput.CharLimit = 6
	f.timeoutInput.SetWidth(10)

	f.commandInput = textinput.New()
	f.commandInput.Placeholder = "docker ps -a"
	f.commandInput.CharLimit = 512
	f.commandInput.SetWidth(80)

	ta := textarea.New()
	ta.Placeholder = "#!/bin/bash\nset -e\n..."
	ta.SetWidth(70)
	ta.SetHeight(8)
	ta.ShowLineNumbers = false
	ta.CharLimit = 16384
	ta.SetStyles(textarea.Styles{
		Focused: textarea.StyleState{
			Base: lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(p.InputFocused).
				Padding(0, 1),
			CursorLine: lipgloss.NewStyle(),
		},
		Blurred: textarea.StyleState{
			Base: lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(p.InputBlurred).
				Padding(0, 1),
			CursorLine: lipgloss.NewStyle(),
		},
	})
	f.scriptArea = ta

	if tk != nil {
		f.nameInput.SetValue(tk.Name)
		f.descInput.SetValue(tk.Description)
		f.scriptType = tk.ScriptType
		f.category = tk.Category
		f.sudo = tk.RequiresSudo
		f.scriptArea.SetValue(tk.ScriptContent)
		f.commandInput.SetValue(tk.Command)
		f.timeoutInput.SetValue(strconv.Itoa(tk.TimeoutSeconds))
	} else {
		f.timeoutInput.SetValue("300")
	}

	f.applyFocus()
	return f
}

func (m ToolFormScreen) Init() tea.Cmd {
	return tea.Batch(func() tea.Msg { return tea.RequestWindowSize() }, textinput.Blink)
}

func (m ToolFormScreen) Keys() []string {
	return []string{keys.ActionSave, keys.NavBack, keys.ViewHelp}
}

func (m ToolFormScreen) Update(msg tea.Msg) (ScreenI, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil

	case toolSaveDoneMsg:
		m.saving = false
		if msg.err != nil {
			m.err = msg.err
			m.log.Warn().Err(msg.err).Msg("guardar tool falló")
			return m, nil
		}
		m.log.Info().Int("tool_id", msg.saved.ID).Msg("tool guardado")
		saved := msg.saved
		return m, func() tea.Msg {
			return ToolFormSavedMsg{Tool: saved}
		}

	case ActionMsg:
		switch msg.ID {
		case keys.ActionSave:
			return m.save()
		case keys.NavBack:
			return m, Back()
		case keys.ViewHelp:
			return m, OpenHelp()
		}
		return m, nil

	case tea.KeyPressMsg:
		return m.handleKey(msg)
	}

	return m.delegateToInput(msg)
}

func (m ToolFormScreen) handleKey(msg tea.KeyPressMsg) (ScreenI, tea.Cmd) {
	fields := m.visibleFields()
	if len(fields) == 0 {
		return m, nil
	}

	switch msg.String() {
	case "tab":
		m.focus = (m.focus + 1) % len(fields)
		m.applyFocus()
		return m, nil
	case "shift+tab":
		m.focus = (m.focus - 1 + len(fields)) % len(fields)
		m.applyFocus()
		return m, nil
	case "left":
		if m.currentKind() == tfKindSelect {
			m.cycleSelect(-1)
			return m, nil
		}
	case "right":
		if m.currentKind() == tfKindSelect {
			m.cycleSelect(+1)
			return m, nil
		}
	case "space":
		if m.currentKind() == tfKindToggle {
			m.sudo = !m.sudo
			return m, nil
		}
	}

	if m.currentKind() == tfKindText {
		return m.delegateToInput(msg)
	}
	return m, nil
}

func (m ToolFormScreen) delegateToInput(msg tea.Msg) (ScreenI, tea.Cmd) {
	var cmd tea.Cmd
	switch m.currentField() {
	case tfName:
		m.nameInput, cmd = m.nameInput.Update(msg)
	case tfDesc:
		m.descInput, cmd = m.descInput.Update(msg)
	case tfTimeout:
		m.timeoutInput, cmd = m.timeoutInput.Update(msg)
	case tfContent:
		m.scriptArea, cmd = m.scriptArea.Update(msg)
	case tfCommand:
		m.commandInput, cmd = m.commandInput.Update(msg)
	}
	return m, cmd
}

func (m ToolFormScreen) View() tea.View {
	var b strings.Builder

	if m.editing != nil {
		b.WriteString("  " + m.styles.Subtitle.Render("ID: ") +
			strconv.Itoa(m.editing.ID) + "\n\n")
	}

	for i, id := range m.visibleFields() {
		b.WriteString(m.renderField(id, i == m.focus))
	}

	if m.saving {
		b.WriteString("\n  " + m.styles.Subtitle.Render("Guardando..."))
	}
	if m.err != nil {
		b.WriteString("\n  " + m.styles.ErrorStyle.Render("✗ ") + m.err.Error())
	}

	return tea.NewView(b.String())
}

func (m ToolFormScreen) renderField(id toolFieldID, focused bool) string {
	cursor := "  "
	if focused {
		cursor = "▶ "
	}

	label := m.labelFor(id)
	labelStyle := m.styles.Subtitle
	if focused {
		labelStyle = m.styles.Key
	}

	// Textarea de script: bloque multi-línea.
	if id == tfContent {
		return cursor + labelStyle.Render(label) + "\n   " +
			m.scriptArea.View() + "\n\n"
	}

	var value string
	switch id {
	case tfName:
		value = m.nameInput.View()
	case tfDesc:
		value = m.descInput.View()
	case tfType:
		value = renderSelect(string(m.scriptType), focused, m.styles)
	case tfCategory:
		value = renderSelect(string(m.category), focused, m.styles)
	case tfCommand:
		value = m.commandInput.View()
	case tfTimeout:
		value = m.timeoutInput.View()
	case tfSudo:
		value = renderToggle(m.sudo, focused, m.styles)
	}

	return cursor + labelStyle.Render(label) + "\n   " + value + "\n\n"
}

func (m ToolFormScreen) labelFor(id toolFieldID) string {
	switch id {
	case tfName:
		return "Nombre"
	case tfDesc:
		return "Descripción"
	case tfType:
		return "Tipo"
	case tfCategory:
		return "Categoría"
	case tfContent:
		return "Script"
	case tfCommand:
		return "Comando"
	case tfTimeout:
		return "Timeout (segundos)"
	case tfSudo:
		return "Requiere sudo"
	}
	return string(id)
}

func (m ToolFormScreen) visibleFields() []toolFieldID {
	out := make([]toolFieldID, 0, len(toolFieldOrder))
	for _, id := range toolFieldOrder {
		if m.fieldVisible(id) {
			out = append(out, id)
		}
	}
	return out
}

func (m ToolFormScreen) fieldVisible(id toolFieldID) bool {
	switch id {
	case tfContent:
		return m.scriptType != coretool.ScriptTypeNative
	case tfCommand:
		return m.scriptType == coretool.ScriptTypeNative
	}
	return true
}

func (m ToolFormScreen) currentField() toolFieldID {
	fields := m.visibleFields()
	if m.focus >= len(fields) {
		return ""
	}
	return fields[m.focus]
}

func (m ToolFormScreen) currentKind() toolFieldKind {
	switch m.currentField() {
	case tfType, tfCategory:
		return tfKindSelect
	case tfSudo:
		return tfKindToggle
	default:
		return tfKindText
	}
}

func (m *ToolFormScreen) applyFocus() {
	m.nameInput.Blur()
	m.descInput.Blur()
	m.timeoutInput.Blur()
	m.commandInput.Blur()
	m.scriptArea.Blur()

	switch m.currentField() {
	case tfName:
		m.nameInput.Focus()
	case tfDesc:
		m.descInput.Focus()
	case tfTimeout:
		m.timeoutInput.Focus()
	case tfContent:
		m.scriptArea.Focus()
	case tfCommand:
		m.commandInput.Focus()
	}
}

func (m *ToolFormScreen) cycleSelect(dir int) {
	switch m.currentField() {
	case tfType:
		opts := []coretool.ScriptType{
			coretool.ScriptTypeBash,
			coretool.ScriptTypePython,
			coretool.ScriptTypeGo,
			coretool.ScriptTypeNative,
		}
		m.scriptType = cycleScriptType(opts, m.scriptType, dir)
		m.clampFocus()
	case tfCategory:
		opts := []coretool.Category{
			coretool.CategorySystem,
			coretool.CategoryDev,
			coretool.CategoryAI,
			coretool.CategoryUtils,
			coretool.CategoryNetwork,
			coretool.CategoryMaintenance,
			coretool.CategoryCustom,
		}
		m.category = cycleCategory(opts, m.category, dir)
	}
}

func (m *ToolFormScreen) clampFocus() {
	n := len(m.visibleFields())
	if n == 0 {
		m.focus = 0
		return
	}
	if m.focus >= n {
		m.focus = n - 1
	}
	m.applyFocus()
}

func (m ToolFormScreen) save() (ScreenI, tea.Cmd) {
	if m.saving {
		return m, nil
	}
	m.saving = true
	m.err = nil

	in := m.buildInput()
	editing := m.editing
	mgr := m.manager
	log := m.log

	return m, func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var saved *coretool.Tool
		var err error
		if editing != nil {
			log.Debug().Int("id", in.ID).Msg("update tool")
			saved, err = mgr.Update(ctx, in)
		} else {
			log.Debug().Str("name", in.Name).Msg("create tool")
			saved, err = mgr.Create(ctx, in)
		}
		return toolSaveDoneMsg{saved: saved, err: err}
	}
}

func (m ToolFormScreen) buildInput() tooluc.ToolInput {
	in := tooluc.ToolInput{
		Name:         strings.TrimSpace(m.nameInput.Value()),
		Description:  strings.TrimSpace(m.descInput.Value()),
		ScriptType:   m.scriptType,
		Category:     m.category,
		RequiresSudo: m.sudo,
	}
	if m.editing != nil {
		in.ID = m.editing.ID
	}

	if m.scriptType == coretool.ScriptTypeNative {
		in.Command = strings.TrimSpace(m.commandInput.Value())
	} else {
		in.ScriptContent = m.scriptArea.Value()
	}

	if v := strings.TrimSpace(m.timeoutInput.Value()); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			in.TimeoutSeconds = n
		}
	}

	return in
}

// --- Helpers de ciclo ---

func cycleScriptType(opts []coretool.ScriptType, cur coretool.ScriptType, dir int) coretool.ScriptType {
	for i, o := range opts {
		if o == cur {
			return opts[(i+dir+len(opts))%len(opts)]
		}
	}
	return opts[0]
}

func cycleCategory(opts []coretool.Category, cur coretool.Category, dir int) coretool.Category {
	for i, o := range opts {
		if o == cur {
			return opts[(i+dir+len(opts))%len(opts)]
		}
	}
	return opts[0]
}
