package screens

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rs/zerolog"

	coretask "github.com/NeRo0128/brain-cli/internal/core/task"
	"github.com/NeRo0128/brain-cli/internal/ui/keys"
	"github.com/NeRo0128/brain-cli/internal/ui/styles"
	taskuc "github.com/NeRo0128/brain-cli/internal/usecases/task"
)

// saveDoneMsg se emite al terminar de guardar.
type saveDoneMsg struct {
	created bool
	err     error
}

// fieldID identifica un campo del form.
type fieldID string

const (
	fID       fieldID = "id"
	fName     fieldID = "name"
	fDesc     fieldID = "desc"
	fType     fieldID = "type"
	fPriority fieldID = "priority"
	fTool     fieldID = "tool"
	fPrompt   fieldID = "prompt"
	fActive   fieldID = "active"
	fFavorite fieldID = "favorite"
)

type fieldKind int

const (
	kindText fieldKind = iota
	kindSelect
	kindToggle
)

// fieldOrder es el orden canónico de todos los campos.
// La visibilidad se decide en fieldVisible().
var fieldOrder = []fieldID{
	fID, fName, fDesc, fType, fPriority, fTool, fPrompt, fActive, fFavorite,
}

// FormScreen edita o crea una Task.
type FormScreen struct {
	editing *coretask.Task

	idInput     textinput.Model
	nameInput   textinput.Model
	descInput   textinput.Model
	toolInput   textinput.Model
	promptInput textinput.Model

	typ      coretask.TaskType
	priority coretask.Priority
	isActive bool
	isFav    bool

	focus int

	err    error
	saving bool

	manager *taskuc.Manager
	log     zerolog.Logger

	width, height int
}

// NewFormScreen construye el form.
// Si tk != nil, edita. Si tk == nil, crea.
func NewFormScreen(tk *coretask.Task, manager *taskuc.Manager, log zerolog.Logger) FormScreen {
	action := "crear"
	if tk != nil {
		action = "editar"
	}
	screenLog := log.With().Str("screen", "form").Str("action", action).Logger()

	f := FormScreen{
		editing:  tk,
		manager:  manager,
		log:      screenLog,
		typ:      coretask.TaskTypeAI,
		priority: coretask.PriorityMedium,
		isActive: true,
		isFav:    false,
	}

	f.idInput = newInput("mi-task", 40)
	f.nameInput = newInput("Nombre de la task", 60)
	f.descInput = newInput("Descripción opcional...", 80)
	f.toolInput = newInput("ID del tool (ej: 1)", 10)
	f.promptInput = newInput("Prompt para la IA...", 100)

	if tk != nil {
		f.idInput.SetValue(tk.ID)
		f.nameInput.SetValue(tk.Name)
		f.descInput.SetValue(tk.Description)
		f.promptInput.SetValue(tk.AIPrompt)
		f.typ = tk.Type
		f.priority = tk.Priority
		f.isActive = tk.IsActive
		f.isFav = tk.IsFavorite
		if tk.ToolID != nil {
			f.toolInput.SetValue(strconv.Itoa(*tk.ToolID))
		}
	}

	f.applyFocus()
	return f
}

// newInput construye un textinput con estilo consistente.
func newInput(placeholder string, width int) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.CharLimit = 512
	ti.Width = width
	return ti
}

func (m FormScreen) Init() tea.Cmd {
	return tea.Batch(
		tea.WindowSize(),
		textinput.Blink,
	)
}

func (m FormScreen) Keys() []string {
	return []string{
		keys.ActionSave,
		keys.NavBack,
		keys.ViewHelp,
	}
}
func (m FormScreen) Update(msg tea.Msg) (ScreenI, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil

	case saveDoneMsg:
		m.saving = false
		if msg.err != nil {
			m.err = msg.err
			m.log.Warn().Err(msg.err).Msg("guardar falló")
			return m, nil
		}
		m.log.Info().Bool("created", msg.created).Msg("task guardada")
		return m, FormSaved()

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

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	// Otros mensajes (BlinkMsg, etc.) → delegar al textinput activo
	return m.delegateToInput(msg)
}
func (m FormScreen) handleKey(msg tea.KeyMsg) (ScreenI, tea.Cmd) {
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
		if m.currentKind() == kindSelect {
			m.cycleSelect(-1)
			return m, nil
		}
	case "right":
		if m.currentKind() == kindSelect {
			m.cycleSelect(+1)
			return m, nil
		}
	case " ":
		if m.currentKind() == kindToggle {
			m.toggleCurrent()
			return m, nil
		}
	}

	if m.currentKind() == kindText {
		return m.delegateToInput(msg)
	}
	return m, nil
}

// delegateToInput pasa el mensaje al textinput activo.
func (m FormScreen) delegateToInput(msg tea.Msg) (ScreenI, tea.Cmd) {
	var cmd tea.Cmd
	switch m.currentField() {
	case fID:
		m.idInput, cmd = m.idInput.Update(msg)
	case fName:
		m.nameInput, cmd = m.nameInput.Update(msg)
	case fDesc:
		m.descInput, cmd = m.descInput.Update(msg)
	case fTool:
		m.toolInput, cmd = m.toolInput.Update(msg)
	case fPrompt:
		m.promptInput, cmd = m.promptInput.Update(msg)
	}
	return m, cmd
}
func (m FormScreen) View() string {
	title := "➕ Nueva tarea"
	if m.editing != nil {
		title = "✏️  Editar tarea"
	}
	var b strings.Builder
	b.WriteString(styles.Title.Render(title))
	b.WriteString("\n\n")

	// ID readonly en edición
	if m.editing != nil {
		b.WriteString("  ")
		b.WriteString(styles.Subtitle.Render("ID: "))
		b.WriteString(m.editing.ID)
		b.WriteString("\n\n")
	}

	for i, id := range m.visibleFields() {
		b.WriteString(m.renderField(id, i == m.focus))
	}

	if m.saving {
		b.WriteString("\n  ")
		b.WriteString(styles.Subtitle.Render("Guardando..."))
	}
	if m.err != nil {
		b.WriteString("\n  ")
		b.WriteString(styles.ErrorStyle.Render("✗ "))
		b.WriteString(m.err.Error())
	}

	b.WriteString("\n\n")
	b.WriteString(styles.Help.Render(
		styles.Key.Render("Tab") + " siguiente  ·  " +
			styles.Key.Render("←/→") + " cambiar  ·  " +
			styles.Key.Render("Space") + " alternar  ·  " +
			styles.Key.Render("Ctrl+S") + " guardar  ·  " +
			styles.Key.Render("Esc") + " cancelar",
	))
	return b.String()
}

func (m FormScreen) renderField(id fieldID, focused bool) string {
	cursor := "  "
	if focused {
		cursor = "▶ "
	}

	label := m.labelFor(id)
	labelStyle := styles.Subtitle
	if focused {
		labelStyle = styles.Key
	}

	var value string
	switch id {
	case fID:
		value = m.idInput.View()
	case fName:
		value = m.nameInput.View()
	case fDesc:
		value = m.descInput.View()
	case fTool:
		value = m.toolInput.View()
	case fPrompt:
		value = m.promptInput.View()
	case fType:
		value = renderSelect(string(m.typ), focused)
	case fPriority:
		value = renderSelect(string(m.priority), focused)
	case fActive:
		value = renderToggle(m.isActive, focused)
	case fFavorite:
		value = renderToggle(m.isFav, focused)
	}

	return cursor + labelStyle.Render(label) + "\n   " + value + "\n\n"
}

func (m FormScreen) labelFor(id fieldID) string {
	switch id {
	case fID:
		return "ID"
	case fName:
		return "Nombre"
	case fDesc:
		return "Descripción"
	case fType:
		return "Tipo"
	case fPriority:
		return "Prioridad"
	case fTool:
		return "Tool ID"
	case fPrompt:
		return "Prompt IA"
	case fActive:
		return "Activa"
	case fFavorite:
		return "Favorita"
	}
	return string(id)
}

func renderSelect(value string, focused bool) string {
	if focused {
		return styles.Key.Render("◀ ") + value + styles.Key.Render(" ▶")
	}
	return styles.Subtitle.Render("  " + value + "  ")
}

func renderToggle(on bool, focused bool) string {
	mark := "[ ]"
	if on {
		mark = "[x]"
	}
	if focused {
		return styles.Key.Render(mark)
	}
	return styles.Subtitle.Render(mark)
}

// --- field traversal ---

func (m FormScreen) visibleFields() []fieldID {
	out := make([]fieldID, 0, len(fieldOrder))
	for _, id := range fieldOrder {
		if m.fieldVisible(id) {
			out = append(out, id)
		}
	}
	return out
}

func (m FormScreen) fieldVisible(id fieldID) bool {
	// ID no se edita: oculto en modo edición
	if id == fID && m.editing != nil {
		return false
	}
	switch id {
	case fTool:
		return m.typ == coretask.TaskTypeScript || m.typ == coretask.TaskTypeCommand
	case fPrompt:
		return m.typ == coretask.TaskTypeAI
	}
	return true
}

func (m FormScreen) currentField() fieldID {
	fields := m.visibleFields()
	if m.focus >= len(fields) {
		return ""
	}
	return fields[m.focus]
}

func (m FormScreen) currentKind() fieldKind {
	switch m.currentField() {
	case fType, fPriority:
		return kindSelect
	case fActive, fFavorite:
		return kindToggle
	default:
		return kindText
	}
}

func (m *FormScreen) applyFocus() {
	m.idInput.Blur()
	m.nameInput.Blur()
	m.descInput.Blur()
	m.toolInput.Blur()
	m.promptInput.Blur()

	switch m.currentField() {
	case fID:
		m.idInput.Focus()
	case fName:
		m.nameInput.Focus()
	case fDesc:
		m.descInput.Focus()
	case fTool:
		m.toolInput.Focus()
	case fPrompt:
		m.promptInput.Focus()
	}
}

func (m *FormScreen) cycleSelect(dir int) {
	switch m.currentField() {
	case fType:
		opts := []coretask.TaskType{
			coretask.TaskTypeAI,
			coretask.TaskTypeScript,
			coretask.TaskTypeCommand,
		}
		m.typ = cycleTaskType(opts, m.typ, dir)
		m.clampFocus()
	case fPriority:
		opts := []coretask.Priority{
			coretask.PriorityLow,
			coretask.PriorityMedium,
			coretask.PriorityHigh,
		}
		m.priority = cyclePriority(opts, m.priority, dir)
	}
}

func (m *FormScreen) toggleCurrent() {
	switch m.currentField() {
	case fActive:
		m.isActive = !m.isActive
	case fFavorite:
		m.isFav = !m.isFav
	}
}

func (m *FormScreen) clampFocus() {
	n := len(m.visibleFields())
	if n == 0 {
		m.focus = 0
		return
	}
	if m.focus >= n {
		m.focus = n - 1
	}
}

// --- save ---

func (m FormScreen) save() (ScreenI, tea.Cmd) {
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

		if editing != nil {
			log.Debug().Str("id", in.ID).Msg("update task")
			_, err := mgr.Update(ctx, in)
			return saveDoneMsg{created: false, err: err}
		}
		log.Debug().Str("id", in.ID).Msg("create task")
		_, err := mgr.Create(ctx, in)
		return saveDoneMsg{created: true, err: err}
	}
}

func (m FormScreen) buildInput() taskuc.TaskInput {
	in := taskuc.TaskInput{
		ID:          strings.TrimSpace(m.idInput.Value()),
		Name:        strings.TrimSpace(m.nameInput.Value()),
		Description: strings.TrimSpace(m.descInput.Value()),
		Type:        m.typ,
		Priority:    m.priority,
		IsActive:    m.isActive,
		IsFavorite:  m.isFav,
	}
	if m.editing != nil {
		in.ID = m.editing.ID
	}

	switch m.typ {
	case coretask.TaskTypeAI:
		in.RequiresAI = true
		in.AIPrompt = strings.TrimSpace(m.promptInput.Value())
	default:
		if v := strings.TrimSpace(m.toolInput.Value()); v != "" {
			if id, err := strconv.Atoi(v); err == nil {
				in.ToolID = &id
			}
		}
	}
	return in
}

// --- helpers de ciclo ---

func cycleTaskType(opts []coretask.TaskType, current coretask.TaskType, dir int) coretask.TaskType {
	for i, o := range opts {
		if o == current {
			return opts[(i+dir+len(opts))%len(opts)]
		}
	}
	return opts[0]
}

func cyclePriority(opts []coretask.Priority, current coretask.Priority, dir int) coretask.Priority {
	for i, o := range opts {
		if o == current {
			return opts[(i+dir+len(opts))%len(opts)]
		}
	}
	return opts[1] // medium por defecto
}
