package screens

import (
	"context"
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
	editing *coretask.Task // nil = crear

	idInput     textinput.Model
	nameInput   textinput.Model
	descInput   textinput.Model
	toolInput   textinput.Model
	promptInput textinput.Model

	typ      coretask.TaskType
	priority coretask.Priority
	isActive bool
	isFav    bool

	focus int // índice en visibleFields()

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
		editing: tk,
		manager: manager,
		log:     screenLog,
	}

	f.idInput = newInput("wifi-vpn", 40)
	f.nameInput = newInput("Conectar WiFi", 60)
	f.descInput = newInput("Descripción opcional...", 80)
	f.promptInput = newInput("Prompt para la IA...", 100)

	if tk != nil {
		f.idInput.SetValue(tk.ID)
		f.idInput.Blur() // el ID no se edita
		f.nameInput.SetValue(tk.Name)
		f.descInput.SetValue(tk.Description)
		f.promptInput.SetValue(tk.AIPrompt)
		f.focus = 1 // empezar en Name
	}

	f.updateFocus()
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
		return m, tea.Batch(Back(), Reload())

	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			m.focus = (m.focus + 1) % len(m.inputs())
			m.updateFocus()
			return m, nil
		case "shift+tab":
			m.focus = (m.focus - 1 + len(m.inputs())) % len(m.inputs())
			m.updateFocus()
			return m, nil
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
	}

	// Delegar al input activo
	var cmd tea.Cmd
	inputs := m.inputsPtr()
	*inputs[m.focus], cmd = inputs[m.focus].Update(msg)
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

	labels := []string{"ID", "Nombre", "Descripción", "Prompt IA"}
	inputs := m.inputs()
	for i, in := range inputs {
		label := labels[i]
		cursor := "  "
		if i == m.focus {
			cursor = "▶ "
			label = styles.Key.Render(label)
		} else {
			label = styles.Subtitle.Render(label)
		}
		b.WriteString(cursor)
		b.WriteString(label)
		b.WriteString("\n")
		b.WriteString("   ")
		b.WriteString(in.View())
		b.WriteString("\n\n")
	}

	if m.saving {
		b.WriteString(styles.Subtitle.Render("  Guardando..."))
		b.WriteString("\n\n")
	}
	if m.err != nil {
		b.WriteString("  ")
		b.WriteString(styles.ErrorStyle.Render("✗ "))
		b.WriteString(m.err.Error())
		b.WriteString("\n\n")
	}

	b.WriteString(styles.Help.Render(
		styles.Key.Render("Tab") + " siguiente  ·  " +
			styles.Key.Render("Shift+Tab") + " anterior  ·  " +
			styles.Key.Render("Ctrl+S") + " guardar  ·  " +
			styles.Key.Render("Esc") + " cancelar",
	))
	return b.String()
}

// --- helpers ---

func (m FormScreen) inputs() []textinput.Model {
	return []textinput.Model{m.idInput, m.nameInput, m.descInput, m.promptInput}
}

func (m FormScreen) inputsPtr() []*textinput.Model {
	return []*textinput.Model{&m.idInput, &m.nameInput, &m.descInput, &m.promptInput}
}

// updateFocus aplica Focus() al input activo y Blur() al resto.
func (m *FormScreen) updateFocus() {
	for i, in := range m.inputsPtr() {
		if i == m.focus {
			in.Focus()
		} else {
			in.Blur()
		}
	}
}

func (m FormScreen) save() (ScreenI, tea.Cmd) {
	if m.saving {
		return m, nil
	}
	m.saving = true
	m.err = nil

	in := taskuc.TaskInput{
		ID:          strings.TrimSpace(m.idInput.Value()),
		Name:        strings.TrimSpace(m.nameInput.Value()),
		Description: strings.TrimSpace(m.descInput.Value()),
		Type:        coretask.TaskTypeAI,
		RequiresAI:  true,
		AIPrompt:    strings.TrimSpace(m.promptInput.Value()),
		Priority:    coretask.PriorityMedium,
		IsActive:    true,
	}

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
