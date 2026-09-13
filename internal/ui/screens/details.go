package screens

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rs/zerolog"

	"github.com/NeRo0128/brain-cli/internal/core/task"
	"github.com/NeRo0128/brain-cli/internal/core/tool"
	"github.com/NeRo0128/brain-cli/internal/ui/components/states"
	"github.com/NeRo0128/brain-cli/internal/ui/keys"
	"github.com/NeRo0128/brain-cli/internal/ui/styles"
)

// detailTwoColMinWidth: umbral para activar layout de 2 columnas.
const detailTwoColMinWidth = 120

// toolLoadedMsg transporta el Tool de la task (o error si no se pudo cargar).
type toolLoadedMsg struct {
	tool *tool.Tool
	err  error
}

// DetailScreen muestra el detalle de una task.
//
// [S4b] responsive: 2 columnas si viewport >= 120 cols, apilado si no.
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

// Update procesa los mensajes de la pantalla.
func (m DetailScreen) Update(msg tea.Msg) (ScreenI, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Chrome (4) + header interno (2) + padding (2) = 8 líneas.
		m.viewport = viewport.New(msg.Width-2, msg.Height-8)
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
		if tk := m.Task(); tk != nil {
			return m, ExecuteTask(tk.ID, tk.Name)
		}
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

// View renderiza SOLO el contenido del medio.
func (m DetailScreen) View() string {
	if !m.ready {
		return states.Loading("detalle")
	}
	header := styles.Subtitle.Render("Task · "+m.task.ID) + "\n\n"
	return header + m.viewport.View()
}

// renderContent decide el layout según el ancho del viewport.
//
//	>= 120 cols: 2 columnas (metadatos | contenido)
//	<  120 cols: apilado vertical
func (m DetailScreen) renderContent() string {
	if m.viewport.Width >= detailTwoColMinWidth {
		return m.renderTwoColumn()
	}
	return m.renderStacked()
}

// renderStacked: layout apilado vertical (terminal normal).
func (m DetailScreen) renderStacked() string {
	var b strings.Builder
	b.WriteString(m.renderMeta())
	b.WriteString("\n")
	b.WriteString(m.renderToolSection())
	if m.task.AIPrompt != "" {
		b.WriteString("\n")
		b.WriteString(m.renderPromptSection())
	}
	return b.String()
}

// renderTwoColumn: layout de 2 columnas (terminal wide >= 120).
// Izquierda: metadatos + tags + params + tool.
// Derecha: prompt IA o script content.
func (m DetailScreen) renderTwoColumn() string {
	totalW := m.viewport.Width
	gap := 3
	leftW := totalW/2 - gap
	rightW := totalW - leftW - gap

	left := lipgloss.NewStyle().Width(leftW).Render(m.renderLeftColumn())
	right := lipgloss.NewStyle().Width(rightW).Render(m.renderRightColumn())
	gapStr := strings.Repeat(" ", gap)

	return lipgloss.JoinHorizontal(lipgloss.Top, left, gapStr, right)
}

// renderLeftColumn: metadatos + tags + params + tool (info compacta).
func (m DetailScreen) renderLeftColumn() string {
	var b strings.Builder
	b.WriteString(m.renderMeta())
	if len(m.task.Tags) > 0 {
		b.WriteString("\n")
		b.WriteString(m.renderTagsSection())
	}
	if len(m.task.Params) > 0 {
		b.WriteString("\n")
		b.WriteString(m.renderParamsSection())
	}
	b.WriteString("\n")
	b.WriteString(m.renderToolSection())
	return b.String()
}

// renderRightColumn: contenido pesado — prompt IA o script content.
func (m DetailScreen) renderRightColumn() string {
	if m.task.AIPrompt != "" {
		return m.renderPromptSection()
	}
	if m.tool != nil && m.tool.ScriptContent != "" {
		return m.renderScriptSection()
	}
	return styles.Subtitle.Render("(sin contenido adicional)")
}

// --- Secciones ---

func (m DetailScreen) renderMeta() string {
	var b strings.Builder
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
	return b.String()
}

func (m DetailScreen) renderTagsSection() string {
	var b strings.Builder
	b.WriteString(styles.SectionHeader.Render("TAGS"))
	b.WriteString("\n")
	for _, tg := range m.task.Tags {
		b.WriteString("  • ")
		b.WriteString(tg)
		b.WriteString("\n")
	}
	return b.String()
}

func (m DetailScreen) renderParamsSection() string {
	var b strings.Builder
	b.WriteString(styles.SectionHeader.Render("PARÁMETROS"))
	b.WriteString("\n")
	for k, v := range m.task.Params {
		writeKV(&b, k, v)
	}
	return b.String()
}

func (m DetailScreen) renderToolSection() string {
	var b strings.Builder
	b.WriteString(styles.SectionHeader.Render("TOOL ASOCIADO"))
	b.WriteString("\n")

	switch {
	case m.loading:
		b.WriteString("  " + styles.Subtitle.Render("Cargando...") + "\n")
	case m.toolErr != nil:
		b.WriteString("  " + styles.ErrorStyle.Render("Error: ") + m.toolErr.Error() + "\n")
	case m.tool == nil:
		b.WriteString("  " + styles.Subtitle.Render("(sin tool asociado)") + "\n")
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
	}
	return b.String()
}

func (m DetailScreen) renderPromptSection() string {
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Border).
		Padding(0, 1)

	return styles.SectionHeader.Render("PROMPT IA") + "\n" +
		box.Render(m.task.AIPrompt) + "\n"
}

func (m DetailScreen) renderScriptSection() string {
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Border).
		Padding(0, 1)

	return styles.SectionHeader.Render("CONTENIDO DEL SCRIPT") + "\n" +
		box.Render(m.tool.ScriptContent) + "\n"
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
