package tools

import (
	"context"
	"fmt"
	"time"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"

	"github.com/NeRo0128/brain-cli/internal/core/tool"
	uilist "github.com/NeRo0128/brain-cli/internal/ui/components/list"
	"github.com/NeRo0128/brain-cli/internal/ui/components/states"
	"github.com/NeRo0128/brain-cli/internal/ui/keys"
	"github.com/NeRo0128/brain-cli/internal/ui/screens"
	"github.com/NeRo0128/brain-cli/internal/ui/styles"
)

type toolsLoadedMsg struct {
	tools []*tool.Tool
	err   error
}

type toolItem struct {
	tool    *tool.Tool
	current bool
}

func (i toolItem) Title() string       { return i.tool.Name }
func (i toolItem) Description() string { return "" }
func (i toolItem) FilterValue() string { return i.tool.Name + " " + i.tool.Description }

func (i toolItem) Row() uilist.Row {
	badges := []uilist.Badge{uilist.TypeBadge(string(i.tool.ScriptType))}
	if i.tool.Category != "" {
		badges = append(badges, uilist.CategoryBadge(string(i.tool.Category)))
	}
	if i.current {
		badges = append(badges, uilist.Badge{
			Text:  "actual",
			Style: uilist.BadgeWarning,
		})
	}
	subtitle := i.tool.Description
	if subtitle == "" {
		subtitle = "(sin descripción)"
	}
	return uilist.Row{
		Title:    i.tool.Name,
		Badges:   badges,
		Subtitle: subtitle,
	}
}

type ToolPickerScreen struct {
	list    list.Model
	loading bool
	err     error

	repo            tool.Repository
	currentID       *int
	filterType      tool.ScriptType
	pendingSelectID *int
	styles          *styles.Styles
}

func NewToolPickerScreen(repo tool.Repository, currentID *int, filterType tool.ScriptType, s *styles.Styles) ToolPickerScreen {
	l := list.New(nil, uilist.New(s), 80, 20)
	title := "Elegir herramienta"
	if filterType != "" {
		title = "Elegir " + string(filterType)
	}
	l.Title = title
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = s.Title
	l.Styles.HelpStyle = s.Help

	return ToolPickerScreen{
		list:       l,
		repo:       repo,
		filterType: filterType,
		currentID:  currentID,
		styles:     s,
		loading:    true,
	}
}

func (m ToolPickerScreen) Init() tea.Cmd {
	repo := m.repo
	filter := m.filterType
	load := func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		tools, err := repo.List(ctx)
		if err != nil {
			return toolsLoadedMsg{err: err}
		}
		if filter != "" {
			filtered := tools[:0]
			for _, t := range tools {
				if t.ScriptType == filter {
					filtered = append(filtered, t)
				}
			}
			tools = filtered
		}
		return toolsLoadedMsg{tools: tools, err: nil}
	}
	return tea.Batch(func() tea.Msg { return tea.RequestWindowSize() }, load)
}

func (m ToolPickerScreen) Keys() []string {
	return []string{
		keys.NavConfirm,
		keys.NavBack,
		keys.NavFilter,
		keys.EditNew,
		keys.EditUpdate,
		keys.EditDelete,
		keys.ViewHelp,
	}
}

func (m ToolPickerScreen) Update(msg tea.Msg) (screens.ScreenI, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width-2, msg.Height-6)
		return m, nil

	case toolsLoadedMsg:
		m.loading = false
		m.err = msg.err
		if msg.err == nil {
			m.setItems(msg.tools)
		}
		return m, nil

	case screens.ReloadMsg:
		m.loading = true
		return m, m.Init()

	case screens.ToolCreatedMsg:
		m.pendingSelectID = &msg.ToolID
		return m, nil

	case screens.ActionMsg:
		return m.handleAction(msg)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m ToolPickerScreen) handleAction(msg screens.ActionMsg) (screens.ScreenI, tea.Cmd) {
	switch msg.ID {
	case keys.NavConfirm:
		it, ok := m.list.SelectedItem().(toolItem)
		if !ok {
			return m, nil
		}
		return m, screens.ToolSelected(it.tool)
	case keys.NavBack:
		return m, screens.Back()
	case keys.ViewHelp:
		return m, screens.OpenHelp()
	case keys.EditNew:
		return m, screens.OpenToolForm(nil, m.styles)
	case keys.EditUpdate:
		it, ok := m.list.SelectedItem().(toolItem)
		if !ok {
			return m, nil
		}
		return m, screens.OpenToolForm(it.tool, m.styles)
	case keys.EditDelete:
		it, ok := m.list.SelectedItem().(toolItem)
		if !ok {
			return m, nil
		}
		return m, screens.OpenConfirm(
			"Borrar tool",
			fmt.Sprintf("¿Borrar el tool '%s'?\n\nEsta acción no se puede deshacer.",
				it.tool.Name),
			screens.DeleteToolMsg{ToolID: it.tool.ID, ToolName: it.tool.Name},
			m.styles,
		)
	case keys.NavFilter:
	}
	return m, nil
}

func (m *ToolPickerScreen) setItems(tools []*tool.Tool) {
	selectID := m.currentID
	if m.pendingSelectID != nil {
		selectID = m.pendingSelectID
		m.pendingSelectID = nil
	}

	items := make([]list.Item, len(tools))
	selectedIdx := -1
	for i, t := range tools {
		isExisting := m.currentID != nil && t.ID == *m.currentID
		isCurrent := selectID != nil && t.ID == *selectID
		items[i] = toolItem{tool: t, current: isExisting}
		if isCurrent {
			selectedIdx = i
		}
	}
	m.list.SetItems(items)
	if selectedIdx >= 0 {
		m.list.Select(selectedIdx)
	}
}

func (m ToolPickerScreen) View() tea.View {
	switch {
	case m.loading:
		return tea.NewView(states.Loading(m.styles, "tools"))
	case m.err != nil:
		return tea.NewView(states.Error(m.styles, m.err))
	case len(m.list.Items()) == 0:
		return tea.NewView(states.Empty(m.styles, "No hay tools disponibles", "Pulsa n para crear uno"))
	}
	return tea.NewView(m.list.View())
}
