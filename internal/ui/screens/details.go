package screens

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/NeRo0128/brain-cli/internal/core/task"
	"github.com/NeRo0128/brain-cli/internal/core/tool"
	"github.com/NeRo0128/brain-cli/internal/ui/keys"
	"github.com/NeRo0128/brain-cli/internal/ui/styles"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rs/zerolog"
)

// toolLoadedMsg transporta el Tool de la task (o error si no se pudo cargar).
type toolLoadedMsg struct {
	tool *tool.Tool
	err  error
}

// DetailScreen muestra el detalle de una task.
type DetailScreen struct {
	task    *task.Task
	tool    *tool.Tool
	toolErr error
	loading bool

	viewport viewport.Model
	ready    bool

	toolRepo tool.Repository
	log      zerolog.Logger
}

// NewDetailScreen construye la pantalla.
func NewDetailScreen(tk *task.Task, toolRepo tool.Repository, log zerolog.Logger) DetailScreen {
	screenLog := log.With().Str("screen", "detail").Str("task_id", tk.ID).Logger()
	screenLog.Debug().Bool("has_tool_id", tk.ToolID != nil).Msg("creando DetailScreen")
	return DetailScreen{
		task:     tk,
		toolRepo: toolRepo,
		loading:  tk.ToolID != nil,
		log:      screenLog,
	}
}

// Init dispara la carga del Tool asociado (si aplica).
func (m DetailScreen) Init() tea.Cmd {

	cmds := []tea.Cmd{tea.WindowSize()}

	if m.task.ToolID != nil {
		repo := m.toolRepo
		toolID := *m.task.ToolID
		log := m.log

		loadTool := func() tea.Msg {
			log.Debug().Int("tool_id", toolID).Msg("Init: cargando tool...")
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			tl, err := repo.GetByID(ctx, toolID)
			log.Debug().Err(err).Bool("has_tool", tl != nil).Msg("Init: tool cargado")
			return toolLoadedMsg{tool: tl, err: err}
		}
		cmds = append(cmds, loadTool)
	}

	return tea.Batch(cmds...)
}
func (m DetailScreen) Keys() []string {
	return []string{
		keys.ActionExecute,
		keys.EditUpdate,
		keys.NavBack,
		keys.ViewHelp,
	}
}

// REEMPLAZA Update:
func (m DetailScreen) Update(msg tea.Msg) (ScreenI, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport = viewport.New(msg.Width, msg.Height-8)
		m.viewport.SetContent(m.renderContent())
		m.ready = true
		return m, nil

	case toolLoadedMsg:
		m.loading = false
		m.tool = msg.tool
		m.toolErr = msg.err
		if m.ready {
			m.viewport.SetContent(m.renderContent())
		}
		return m, nil

	case ActionMsg:
		return m.handleAction(msg)
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m DetailScreen) handleAction(msg ActionMsg) (ScreenI, tea.Cmd) {
	switch msg.ID {
	case keys.ActionExecute:
		tk := m.Task()
		if tk == nil {
			return m, nil
		}
		return m, ExecuteTask(tk.ID, tk.Name)
	case keys.NavBack:
		return m, Back()
	case keys.ViewHelp:
		return m, OpenHelp()
	case keys.EditUpdate:
		return m, OpenForm(m.task)
	}
	return m, nil
}

func (m DetailScreen) Task() *task.Task { return m.task }
func (m DetailScreen) Ready() bool      { return m.ready }
func (m DetailScreen) HasTool() bool    { return m.tool != nil }

// View renderiza la pantalla.
func (m DetailScreen) View() string {
	header := styles.Title.Render("🧠 "+m.task.Name) + "\n"
	header += styles.Subtitle.Render("Task · "+m.task.ID) + "\n\n"

	if !m.ready {
		return header + styles.Subtitle.Render("Cargando...") + "\n"
	}

	footer := "\n" + styles.Help.Render(
		styles.Key.Render("↑↓")+" scroll  ·  "+
			styles.Key.Render("Enter/e")+" ejecutar  ·  "+
			styles.Key.Render("Esc")+" volver  ·  "+
			styles.Key.Render("q")+" salir",
	)

	return header + m.viewport.View() + footer
}

// renderContent construye el cuerpo del detalle.
func (m DetailScreen) renderContent() string {
	var b strings.Builder

	// --- Sección: metadatos de la task ---
	b.WriteString(styles.SectionHeader.Render("METADATOS"))
	b.WriteString("\n")
	writeKV(&b, "Tipo", string(m.task.Type))
	writeKV(&b, "Prioridad", string(m.task.Priority))
	writeKV(&b, "Activa", boolYesNo(m.task.IsActive))
	writeKV(&b, "Favorita", boolYesNo(m.task.IsFavorite))
	writeKV(&b, "Requiere IA", boolYesNo(m.task.RequiresAI))
	if m.task.Description != "" {
		writeKV(&b, "Descripción", m.task.Description)
	}
	b.WriteString("\n")

	// --- Sección: tags ---
	if len(m.task.Tags) > 0 {
		b.WriteString(styles.SectionHeader.Render("TAGS"))
		b.WriteString("\n")
		for _, tg := range m.task.Tags {
			b.WriteString("  • ")
			b.WriteString(tg)
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	// --- Sección: parámetros ---
	if len(m.task.Params) > 0 {
		b.WriteString(styles.SectionHeader.Render("PARÁMETROS"))
		b.WriteString("\n")
		for k, v := range m.task.Params {
			writeKV(&b, k, v)
		}
		b.WriteString("\n")
	}

	// --- Sección: prompt IA (si aplica) ---
	if m.task.AIPrompt != "" {
		b.WriteString(styles.SectionHeader.Render("PROMPT IA"))
		b.WriteString("\n")
		b.WriteString(indent(m.task.AIPrompt, "  "))
		b.WriteString("\n\n")
	}

	// --- Sección: Tool asociado ---
	b.WriteString(styles.SectionHeader.Render("TOOL ASOCIADO"))
	b.WriteString("\n")
	switch {
	case m.loading:
		b.WriteString("  ")
		b.WriteString(styles.Subtitle.Render("Cargando..."))
		b.WriteString("\n")
	case m.toolErr != nil:
		b.WriteString("  ")
		b.WriteString(styles.ErrorStyle.Render("Error: "))
		b.WriteString(m.toolErr.Error())
		b.WriteString("\n")
	case m.tool == nil:
		b.WriteString("  ")
		b.WriteString(styles.Subtitle.Render("(sin tool asociado)"))
		b.WriteString("\n")
	default:
		writeKV(&b, "Nombre", m.tool.Name)
		writeKV(&b, "Tipo", string(m.tool.ScriptType))
		writeKV(&b, "Categoría", string(m.tool.Category))
		writeKV(&b, "Timeout", fmt.Sprintf("%ds", m.tool.TimeoutSeconds))
		writeKV(&b, "Sudo", boolYesNo(m.tool.RequiresSudo))
		writeKV(&b, "Builtin", boolYesNo(m.tool.IsBuiltin))
		if m.tool.ScriptPath != "" {
			writeKV(&b, "Script path", m.tool.ScriptPath)
		}
		if m.tool.Command != "" {
			writeKV(&b, "Comando", m.tool.Command)
		}
		if m.tool.ScriptContent != "" {
			b.WriteString("\n  ")
			b.WriteString(styles.Subtitle.Render("Contenido del script:"))
			b.WriteString("\n")
			b.WriteString(indent(m.tool.ScriptContent, "  "))
			b.WriteString("\n")
		}
	}

	return b.String()
}

// --- helpers de formato ---

func writeKV(b *strings.Builder, key, value string) {
	b.WriteString("  ")
	b.WriteString(styles.Key.Render(key + ":"))
	b.WriteString(" ")
	b.WriteString(value)
	b.WriteString("\n")
}

func boolYesNo(v bool) string {
	if v {
		return "sí"
	}
	return "no"
}

// indent añade prefix a cada línea.
func indent(s, prefix string) string {
	lines := strings.Split(s, "\n")
	for i, l := range lines {
		lines[i] = prefix + l
	}
	return strings.Join(lines, "\n")
}
