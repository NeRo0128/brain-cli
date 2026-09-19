package screens

import (
	"context"
	"fmt"
	"strings"
	"time"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/muesli/reflow/wrap"
	"github.com/rs/zerolog"

	"github.com/NeRo0128/brain-cli/internal/core/task"
	"github.com/NeRo0128/brain-cli/internal/core/tool"
	"github.com/NeRo0128/brain-cli/internal/ui/components/states"
	"github.com/NeRo0128/brain-cli/internal/ui/keys"
	"github.com/NeRo0128/brain-cli/internal/ui/styles"
)

type toolLoadedMsg struct {
	tool *tool.Tool
	err  error
}

type DetailScreen struct {
	task    *task.Task
	tool    *tool.Tool
	toolErr error
	loading bool

	viewport viewport.Model
	ready    bool

	toolRepo tool.Repository
	taskRepo task.Repository
	log      zerolog.Logger
	styles   *styles.Styles
}

func NewDetailScreen(
	tk *task.Task,
	taskRepo task.Repository,
	toolRepo tool.Repository,
	log zerolog.Logger,
	s *styles.Styles,
) DetailScreen {
	screenLog := log.With().Str("screen", "detail").Str("task_id", tk.ID).Logger()
	return DetailScreen{
		task:     tk,
		taskRepo: taskRepo,
		toolRepo: toolRepo,
		loading:  tk.ToolID != nil,
		log:      screenLog,
		styles:   s,
	}
}

// taskReloadedMsg transporta la task + tool recargados tras un ReloadMsg.
type taskReloadedMsg struct {
	task    *task.Task
	tool    *tool.Tool
	toolErr error
	err     error
}

func (m DetailScreen) Init() tea.Cmd {
	cmds := []tea.Cmd{func() tea.Msg { return tea.RequestWindowSize() }}
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

func (m DetailScreen) Update(msg tea.Msg) (ScreenI, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport = viewport.New(viewport.WithWidth(msg.Width-2), viewport.WithHeight(msg.Height-8))
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
	case ReloadMsg:
		return m, m.reload()
	case taskReloadedMsg:
		if msg.err != nil {
			// Task probablemente borrada: volver a main.
			m.log.Warn().Err(msg.err).Msg("reload: task no encontrada, volviendo")
			return m, Back()
		}
		m.task = msg.task
		m.tool = msg.tool
		m.toolErr = msg.toolErr
		m.loading = false
		if m.ready {
			m.viewport.SetContent(m.renderContent())
		}
		return m, nil
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
		return m, OpenForm(m.task, m.styles)
	}
	return m, nil
}

func (m DetailScreen) Task() *task.Task { return m.task }
func (m DetailScreen) Ready() bool      { return m.ready }
func (m DetailScreen) HasTool() bool    { return m.tool != nil }

func (m DetailScreen) View() tea.View {
	if !m.ready {
		return tea.NewView(states.Loading(m.styles, "detalle"))
	}
	p := m.styles.Theme.Resolve(m.styles.Dark)
	nameStyle := lipgloss.NewStyle().Bold(true).Foreground(p.Secondary)
	header := nameStyle.Render(m.task.Name) + " " + m.styles.Subtitle.Render("Task · "+m.task.ID) + "\n\n"
	return tea.NewView(header + m.viewport.View())
}

func (m DetailScreen) renderContent() string {
	if m.viewport.Width() >= twoColMinWidth {
		return m.renderTwoColumn()
	}
	return m.renderStacked()
}

func (m DetailScreen) renderStacked() string {
	// Ancho útil: viewport menos 4 cols de margen (2 izq + 2 der).
	w := max(m.viewport.Width()-4, 30)

	var b strings.Builder
	b.WriteString(m.renderMeta())
	b.WriteString("\n")
	b.WriteString(m.renderToolSection())
	if m.task.AIPrompt != "" {
		b.WriteString("\n")
		b.WriteString(m.renderPromptSection(w))
	}
	return b.String()
}
func (m DetailScreen) renderTwoColumn() string {
	totalW := m.viewport.Width()
	gap := 3
	leftW := totalW/2 - gap
	rightW := totalW - leftW - gap

	rightContentW := max(rightW-4, 20)

	left := lipgloss.NewStyle().Width(leftW).Render(m.renderLeftColumn())
	right := lipgloss.NewStyle().Width(rightW).Render(m.renderRightColumn(rightContentW))
	gapStr := strings.Repeat(" ", gap)

	return lipgloss.JoinHorizontal(lipgloss.Top, left, gapStr, right)
}

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

func (m DetailScreen) renderRightColumn(contentW int) string {
	if m.task.AIPrompt != "" {
		return m.renderPromptSection(contentW)
	}
	if m.tool != nil && m.tool.ScriptContent != "" {
		return m.renderScriptSection(contentW)
	}
	return m.styles.Subtitle.Render("(sin contenido adicional)")
}

func (m DetailScreen) renderMeta() string {
	var b strings.Builder
	b.WriteString(m.styles.SectionHeader.Render("METADATOS"))
	b.WriteString("\n")
	writeKV(&b, m.styles, "Tipo", string(m.task.Type))
	writeKV(&b, m.styles, "Prioridad", string(m.task.Priority))
	writeKV(&b, m.styles, "Activa", boolYesNo(m.task.IsActive))
	writeKV(&b, m.styles, "Favorita", boolYesNo(m.task.IsFavorite))
	writeKV(&b, m.styles, "Requiere IA", boolYesNo(m.task.RequiresAI))
	if m.task.Description != "" {
		writeKV(&b, m.styles, "Descripción", m.task.Description)
	}
	return b.String()
}

func (m DetailScreen) renderTagsSection() string {
	var b strings.Builder
	b.WriteString(m.styles.SectionHeader.Render("TAGS"))
	b.WriteString("\n")
	for _, tg := range m.task.Tags {
		b.WriteString("  ")
		b.WriteString(m.styles.ColoredIcons.Bullet())
		b.WriteString(" ")
		b.WriteString(tg)
		b.WriteString("\n")
	}
	return b.String()
}

func (m DetailScreen) renderParamsSection() string {
	var b strings.Builder
	b.WriteString(m.styles.SectionHeader.Render("PARÁMETROS"))
	b.WriteString("\n")
	for k, v := range m.task.Params {
		writeKV(&b, m.styles, k, v)
	}
	return b.String()
}

func (m DetailScreen) renderToolSection() string {
	var b strings.Builder
	b.WriteString(m.styles.SectionHeader.Render("TOOL ASOCIADO"))
	b.WriteString("\n")

	switch {
	case m.loading:
		b.WriteString("  ")
		b.WriteString(m.styles.Subtitle.Render("Cargando..."))
		b.WriteString("\n")
	case m.toolErr != nil:
		b.WriteString("  ")
		b.WriteString(m.styles.ErrorStyle.Render("Error: "))
		b.WriteString(m.toolErr.Error())
		b.WriteString("\n")
	case m.tool == nil:
		b.WriteString("  ")
		b.WriteString(m.styles.Subtitle.Render("(sin tool asociado)"))
		b.WriteString("\n")
	default:
		writeKV(&b, m.styles, "Nombre", m.tool.Name)
		writeKV(&b, m.styles, "Tipo", string(m.tool.ScriptType))
		writeKV(&b, m.styles, "Categoría", string(m.tool.Category))
		writeKV(&b, m.styles, "Timeout", fmt.Sprintf("%ds", m.tool.TimeoutSeconds))
		writeKV(&b, m.styles, "Sudo", boolYesNo(m.tool.RequiresSudo))
		writeKV(&b, m.styles, "Builtin", boolYesNo(m.tool.IsBuiltin))
		if m.tool.ScriptPath != "" {
			writeKV(&b, m.styles, "Script path", m.tool.ScriptPath)
		}
		if m.tool.Command != "" {
			writeKV(&b, m.styles, "Comando", m.tool.Command)
		}
	}
	return b.String()
}

func (m DetailScreen) renderPromptSection(contentW int) string {
	p := m.styles.Theme.Resolve(m.styles.Dark)

	lines := strings.Split(m.task.AIPrompt, "\n")
	total := len(lines)
	truncated := false
	if total > maxTextAreaPreviewLines {
		lines = lines[:maxTextAreaPreviewLines]
		truncated = true
	}

	wrapped := wrap.String(strings.Join(lines, "\n"), contentW)

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(p.Border).
		Padding(0, 1).
		Width(contentW)

	out := m.styles.SectionHeader.Render("PROMPT IA") + "\n" +
		box.Render(wrapped) + "\n"

	if truncated {
		missing := total - maxTextAreaPreviewLines
		out += m.styles.Subtitle.Render(
			fmt.Sprintf("  … (+%d líneas, pulsa e para ver completo)", missing),
		) + "\n"
	}
	return out
}

func (m DetailScreen) renderScriptSection(contentW int) string {
	p := m.styles.Theme.Resolve(m.styles.Dark)

	// Truncar a N líneas.
	lines := strings.Split(m.tool.ScriptContent, "\n")
	total := len(lines)
	truncated := false
	if total > maxTextAreaPreviewLines {
		lines = lines[:maxTextAreaPreviewLines]
		truncated = true
	}

	body := strings.Join(lines, "\n")
	wrapped := wrap.String(body, contentW)

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(p.Border).
		Padding(0, 1).
		Width(contentW)

	out := m.styles.SectionHeader.Render("CONTENIDO DEL SCRIPT") + "\n" +
		box.Render(wrapped) + "\n"

	if truncated {
		missing := total - maxTextAreaPreviewLines
		out += m.styles.Subtitle.Render(
			fmt.Sprintf("  … (+%d líneas, pulsa e para ver completo)", missing),
		) + "\n"
	}
	return out
}
func writeKV(b *strings.Builder, s *styles.Styles, key, value string) {
	b.WriteString("  ")
	b.WriteString(s.Key.Render(key + ":"))
	b.WriteString(" ")
	b.WriteString(value)
	b.WriteString("\n")
}

// reload re-fetchea la task y su tool desde los repositorios.
// Se llama tras un ReloadMsg (ej: después de editar la task en el form).
func (m DetailScreen) reload() tea.Cmd {
	taskRepo := m.taskRepo
	toolRepo := m.toolRepo
	id := m.task.ID
	log := m.log

	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		tk, err := taskRepo.GetByID(ctx, id)
		if err != nil {
			log.Warn().Err(err).Str("task_id", id).Msg("reload: task no encontrada")
			return taskReloadedMsg{err: err}
		}

		var tl *tool.Tool
		var toolErr error
		if tk.ToolID != nil {
			tl, toolErr = toolRepo.GetByID(ctx, *tk.ToolID)
		}

		log.Debug().
			Str("task_id", id).
			Bool("has_tool", tl != nil).
			Msg("reload: task actualizada")

		return taskReloadedMsg{task: tk, tool: tl, toolErr: toolErr}
	}
}
