package screens

import (
	"context"
	"strings"
	"time"

	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/rs/zerolog"

	coretask "github.com/NeRo0128/brain-cli/internal/core/task"
	"github.com/NeRo0128/brain-cli/internal/core/tool"
	"github.com/NeRo0128/brain-cli/internal/ui/keys"
	"github.com/NeRo0128/brain-cli/internal/ui/styles"
	taskuc "github.com/NeRo0128/brain-cli/internal/usecases/task"
)

type saveDoneMsg struct {
	created bool
	err     error
}

type toolLoadedForFormMsg struct {
	tool *tool.Tool
	err  error
}

type fieldID string

const (
	fID          fieldID = "id"
	fName        fieldID = "name"
	fDesc        fieldID = "desc"
	fKind        fieldID = "kind"        // reemplaza fType
	fInterpreter fieldID = "interpreter" // solo si kind=script
	fTool        fieldID = "tool"        // solo si kind=script
	fCommand     fieldID = "command"     // solo si kind=command
	fPrompt      fieldID = "prompt"      // solo si kind=ai
	fPriority    fieldID = "priority"
	fActive      fieldID = "active"
	fFavorite    fieldID = "favorite"
)

type fieldKind int

const (
	kindText fieldKind = iota
	kindSelect
	kindToggle
	kindButton
)

var fieldOrder = []fieldID{
	fID, fName, fDesc,
	fKind, fInterpreter, fTool, fCommand, fPrompt,
	fPriority, fActive, fFavorite,
}

// FormScreen edita o crea una Task.
type FormScreen struct {
	editing *coretask.Task

	idInput      textinput.Model
	nameInput    textinput.Model
	descInput    textinput.Model
	commandInput textinput.Model //
	promptInput  textarea.Model   // [S4a] prompt IA con múltiples líneas

	kind           coretask.TaskKind  //
	interpreters   []tool.Interpreter // solo disponibles
	interpreterIdx int                //
	selectedTool   *tool.Tool
	priority       coretask.Priority
	isActive       bool
	isFav          bool

	focus int

	err    error
	saving bool

	manager  *taskuc.Manager
	toolRepo tool.Repository
	log      zerolog.Logger

	width, height int
	styles        *styles.Styles
}

// constructor recibe la lista de intérpretes
func NewFormScreen(
	tk *coretask.Task,
	manager *taskuc.Manager,
	toolRepo tool.Repository,
	interpreters []tool.Interpreter,
	log zerolog.Logger,
	s *styles.Styles,
) FormScreen {
	action := "crear"
	if tk != nil {
		action = "editar"
	}
	screenLog := log.With().Str("screen", "form").Str("action", action).Logger()
	p := s.Theme.Resolve(s.Dark)

	f := FormScreen{
		editing:      tk,
		manager:      manager,
		toolRepo:     toolRepo,
		interpreters: interpreters,
		log:          screenLog,
		kind:         coretask.KindAI,
		priority:     coretask.PriorityMedium,
		isActive:     true,
		isFav:        false,
		styles:       s,
	}

	f.idInput = newInput("mi-task", 40)
	f.nameInput = newInput("Nombre de la task", 60)
	f.descInput = newInput("Descripción opcional...", 80)
	f.commandInput = newInput("docker ps -a", 80)

	// [S4a] Textarea para el prompt IA: bloque multi-línea
	ta := textarea.New()
	ta.Placeholder = "Prompt para la IA..."
	ta.SetWidth(60)
	ta.SetHeight(5)
	ta.ShowLineNumbers = false
	ta.CharLimit = 8192
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
	f.promptInput = ta

	if tk != nil {
		f.idInput.SetValue(tk.ID)
		f.nameInput.SetValue(tk.Name)
		f.descInput.SetValue(tk.Description)
		f.promptInput.SetValue(tk.AIPrompt)
		f.priority = tk.Priority
		f.isActive = tk.IsActive
		f.isFav = tk.IsFavorite
		f.kind = coretask.TaskKindFromType(tk.Type)

		// Si es script, buscar el intérprete que matchee el tool actual
		if f.kind == coretask.KindScript && len(interpreters) > 0 {
			f.interpreterIdx = 0 // se ajusta al cargar el tool
		}
	}

	f.applyFocus()
	return f
}

func newInput(placeholder string, width int) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.CharLimit = 512
	ti.SetWidth(width)
	return ti
}

func (m FormScreen) Init() tea.Cmd {
	cmds := []tea.Cmd{func() tea.Msg { return tea.RequestWindowSize() }, textinput.Blink}

	if m.editing != nil && m.editing.ToolID != nil {
		repo := m.toolRepo
		id := *m.editing.ToolID
		log := m.log
		loadTool := func() tea.Msg {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			tl, err := repo.GetByID(ctx, id)
			log.Debug().Err(err).Int("tool_id", id).Msg("Init: tool cargado para form")
			return toolLoadedForFormMsg{tool: tl, err: err}
		}
		cmds = append(cmds, loadTool)
	}

	return tea.Batch(cmds...)
}

func (m FormScreen) Keys() []string {
	return []string{keys.ActionSave, keys.NavBack, keys.ViewHelp}
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

	case toolLoadedForFormMsg:
		if msg.err == nil {
			m.selectedTool = msg.tool
			m.syncInterpreterToTool()
		}
		return m, nil

	case ToolSelectedMsg:
		m.selectedTool = msg.Tool
		m.syncInterpreterToTool()
		m.log.Debug().Str("tool", msg.Tool.Name).Msg("tool seleccionado")
		return m, nil

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

// syncInterpreterToTool alinea interpreterIdx con el tool actual.
func (m *FormScreen) syncInterpreterToTool() {
	if m.selectedTool == nil {
		return
	}
	for i, it := range m.interpreters {
		if it.ScriptType == m.selectedTool.ScriptType {
			m.interpreterIdx = i
			return
		}
	}
}

func (m FormScreen) handleKey(msg tea.KeyPressMsg) (ScreenI, tea.Cmd) {
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
	case "enter":
		// Enter en fTool abre el picker filtrado
		if m.currentField() == fTool {
			var current *int
			if m.selectedTool != nil {
				current = &m.selectedTool.ID
			}
			return m, OpenToolPicker(current, m.currentScriptType(), m.styles)
		}
	case "left":
		switch m.currentKind() {
		case kindSelect:
			m.cycleSelect(-1)
			return m, nil
		}
	case "right":
		switch m.currentKind() {
		case kindSelect:
			m.cycleSelect(+1)
			return m, nil
		}
	case "space":
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

func (m FormScreen) delegateToInput(msg tea.Msg) (ScreenI, tea.Cmd) {
	var cmd tea.Cmd
	switch m.currentField() {
	case fID:
		m.idInput, cmd = m.idInput.Update(msg)
	case fName:
		m.nameInput, cmd = m.nameInput.Update(msg)
	case fDesc:
		m.descInput, cmd = m.descInput.Update(msg)
	case fCommand:
		m.commandInput, cmd = m.commandInput.Update(msg)
	case fPrompt:
		m.promptInput, cmd = m.promptInput.Update(msg)
	}
	return m, cmd
}

// View renderiza SOLO el contenido del medio.
// [ACTUALIZADO] sin title ni footer propios.
func (m FormScreen) View() tea.View {
	var b strings.Builder

	// En modo edición, mostramos el ID como referencia (no editable).
	if m.editing != nil {
		b.WriteString("  ")
		b.WriteString(m.styles.Subtitle.Render("ID: "))
		b.WriteString(m.editing.ID)
		b.WriteString("\n\n")
	}

	// [S4c] section dividers entre grupos de campos
	currentSection := ""
	fields := m.visibleFields()
	for i, id := range fields {
		sec := sectionFor(id)
		if sec != "" && sec != currentSection {
			b.WriteString(renderSectionDivider(sec, m.styles))
			currentSection = sec
		}
		b.WriteString(m.renderField(id, i == m.focus))
	}

	if m.saving {
		b.WriteString("\n  ")
		b.WriteString(m.styles.Subtitle.Render("Guardando..."))
	}
	if m.err != nil {
		b.WriteString("\n  ")
		b.WriteString(m.styles.ErrorStyle.Render("✗ "))
		b.WriteString(m.err.Error())
	}

	return tea.NewView(b.String())
}

func (m FormScreen) renderField(id fieldID, focused bool) string {
	cursor := "  "
	if focused {
		cursor = "▶ "
	}

	label := m.labelFor(id)
	labelStyle := m.styles.Subtitle
	if focused {
		labelStyle = m.styles.Key
	}

	// El textarea del prompt se renderiza en bloque multi-línea.
	if id == fPrompt {
		return cursor + labelStyle.Render(label) + "\n   " +
			m.promptInput.View() + "\n\n"
	}

	var value string
	switch id {
	case fID:
		value = m.idInput.View()
	case fName:
		value = m.nameInput.View()
	case fDesc:
		value = m.descInput.View()
	case fCommand:
		value = m.commandInput.View()
	case fPrompt:
		value = m.promptInput.View()
	case fKind:
		value = renderSelect(string(m.kind), focused, m.styles)
	case fInterpreter:
		value = m.renderInterpreterField(focused)
	case fPriority:
		value = renderSelect(string(m.priority), focused, m.styles)
	case fActive:
		value = renderToggle(m.isActive, focused, m.styles)
	case fFavorite:
		value = renderToggle(m.isFav, focused, m.styles)
	case fTool:
		value = m.renderToolField(focused)
	}

	return cursor + labelStyle.Render(label) + "\n   " + value + "\n\n"
}

// renderInterpreterField: muestra el intérprete o warning si no hay.
func (m FormScreen) renderInterpreterField(focused bool) string {
	if len(m.interpreters) == 0 {
		return m.styles.ErrorStyle.Render("⚠ No hay intérpretes instalados")
	}
	if focused {
		return m.styles.Key.Render("◀ ") + m.interpreters[m.interpreterIdx].Display +
			m.styles.Key.Render(" ▶")
	}
	return m.styles.Subtitle.Render("  " + m.interpreters[m.interpreterIdx].Display + "  ")
}

func (m FormScreen) renderToolField(focused bool) string {
	if m.selectedTool == nil {
		if focused {
			return m.styles.Key.Render("[ Elegir tool... ]")
		}
		return m.styles.Subtitle.Render("[ Elegir tool... ]")
	}
	txt := "[" + m.selectedTool.Name + "]"
	if focused {
		return m.styles.Key.Render(txt)
	}
	return m.styles.Subtitle.Render(txt)
}

func (m FormScreen) labelFor(id fieldID) string {
	switch id {
	case fID:
		return "ID"
	case fName:
		return "Nombre"
	case fDesc:
		return "Descripción"
	case fKind:
		return "¿Qué quieres hacer?"
	case fInterpreter:
		return "Intérprete"
	case fTool:
		return "Script"
	case fCommand:
		return "Comando"
	case fPrompt:
		return "Prompt IA"
	case fPriority:
		return "Prioridad"
	case fActive:
		return "Activa"
	case fFavorite:
		return "Favorita"
	}
	return string(id)
}

func renderSelect(value string, focused bool, s *styles.Styles) string {
	if focused {
		return s.Key.Render("◀ ") + value + s.Key.Render(" ▶")
	}
	return s.Subtitle.Render("  " + value + "  ")
}

func renderToggle(on bool, focused bool, s *styles.Styles) string {
	mark := "[ ]"
	if on {
		mark = "[x]"
	}
	if focused {
		return s.Key.Render(mark)
	}
	return s.Subtitle.Render(mark)
}

// visibleFields según kind
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
	if id == fID && m.editing != nil {
		return false
	}
	switch id {
	case fInterpreter, fTool:
		return m.kind == coretask.KindScript
	case fCommand:
		return m.kind == coretask.KindCommand
	case fPrompt:
		return m.kind == coretask.KindAI
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
	case fKind, fInterpreter, fPriority:
		return kindSelect
	case fActive, fFavorite:
		return kindToggle
	case fTool:
		return kindButton
	default:
		return kindText
	}
}

// currentScriptType: ScriptType del intérprete elegido.
func (m FormScreen) currentScriptType() tool.ScriptType {
	if len(m.interpreters) == 0 {
		return ""
	}
	return m.interpreters[m.interpreterIdx].ScriptType
}

func (m *FormScreen) applyFocus() {
	m.idInput.Blur()
	m.nameInput.Blur()
	m.descInput.Blur()
	m.commandInput.Blur()
	m.promptInput.Blur()

	switch m.currentField() {
	case fID:
		m.idInput.Focus()
	case fName:
		m.nameInput.Focus()
	case fDesc:
		m.descInput.Focus()
	case fCommand:
		m.commandInput.Focus()
	case fPrompt:
		m.promptInput.Focus()
	}
}

func (m *FormScreen) cycleSelect(dir int) {
	switch m.currentField() {
	case fKind:
		opts := []coretask.TaskKind{
			coretask.KindScript,
			coretask.KindCommand,
			coretask.KindAI,
		}
		m.kind = cycleKind(opts, m.kind, dir)
		// Al cambiar de kind, resetear tool (puede no aplicar)
		m.selectedTool = nil
		m.clampFocus()

	case fInterpreter:
		if len(m.interpreters) == 0 {
			return
		}
		n := len(m.interpreters)
		m.interpreterIdx = (m.interpreterIdx + dir + n) % n
		// Si cambia el intérprete, resetear tool (tipo distinto)
		m.selectedTool = nil

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
	m.applyFocus()
}

func (m FormScreen) save() (ScreenI, tea.Cmd) {
	if m.saving {
		return m, nil
	}
	// validar que haya intérprete si kind=script
	if m.kind == coretask.KindScript && len(m.interpreters) == 0 {
		m.err = errNoInterpreters
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

// [REESCRITO] buildInput usa Kind → Type + Tool + Command/AIPrompt
func (m FormScreen) buildInput() taskuc.TaskInput {
	in := taskuc.TaskInput{
		ID:          strings.TrimSpace(m.idInput.Value()),
		Name:        strings.TrimSpace(m.nameInput.Value()),
		Description: strings.TrimSpace(m.descInput.Value()),
		Type:        m.kind.ToTaskType(),
		Priority:    m.priority,
		IsActive:    m.isActive,
		IsFavorite:  m.isFav,
	}
	if m.editing != nil {
		in.ID = m.editing.ID
	}

	switch m.kind {
	case coretask.KindScript:
		if m.selectedTool != nil {
			in.ToolID = &m.selectedTool.ID
		}
	case coretask.KindCommand:
		in.Command = strings.TrimSpace(m.commandInput.Value())
		if m.selectedTool != nil {
			in.ToolID = &m.selectedTool.ID
		}
	case coretask.KindAI:
		in.RequiresAI = true
		in.AIPrompt = strings.TrimSpace(m.promptInput.Value())
	}
	return in
}

func cycleKind(opts []coretask.TaskKind, current coretask.TaskKind, dir int) coretask.TaskKind {
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
	return opts[1]
}

// error local
var errNoInterpreters = errNoInterpretersT("no hay intérpretes instalados, elige otro tipo")

type errNoInterpretersT string

func (e errNoInterpretersT) Error() string { return string(e) }

// [S4c] Section dividers: agrupan campos en IDENTIDAD / ACCIÓN / METADATOS.

// sectionFor devuelve el nombre de la sección a la que pertenece un campo.
// "" significa "sin sección".
func sectionFor(id fieldID) string {
	switch id {
	case fID, fName, fDesc:
		return "Identidad"
	case fKind, fInterpreter, fTool, fCommand, fPrompt:
		return "Acción"
	case fPriority, fActive, fFavorite:
		return "Metadatos"
	}
	return ""
}

// renderSectionDivider dibuja una línea tipo "─── IDENTIDAD ───".
func renderSectionDivider(name string, s *styles.Styles) string {
	line := strings.Repeat("─", 3)
	header := strings.ToUpper(name)
	return "\n" + s.Subtitle.Render(line+" "+header+" "+line) + "\n\n"
}
