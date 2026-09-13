package screens

import (
	"context"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/NeRo0128/brain-cli/internal/core/tool"
	uilist "github.com/NeRo0128/brain-cli/internal/ui/components/list"
	"github.com/NeRo0128/brain-cli/internal/ui/components/states"
	"github.com/NeRo0128/brain-cli/internal/ui/keys"
	"github.com/NeRo0128/brain-cli/internal/ui/styles"
)

type toolsLoadedMsg struct {
	tools []*tool.Tool
	err   error
}

// --- Item ---

type toolItem struct {
	tool    *tool.Tool
	current bool // [S4c] ★ actual badge si es el tool seleccionado
}

func (i toolItem) Title() string       { return i.tool.Name }
func (i toolItem) Description() string { return "" }
func (i toolItem) FilterValue() string { return i.tool.Name + " " + i.tool.Description }

func (i toolItem) Row() uilist.Row {
	badges := []uilist.Badge{uilist.TypeBadge(string(i.tool.ScriptType))}
	if i.tool.Category != "" {
		badges = append(badges, uilist.CategoryBadge(string(i.tool.Category)))
	}
	// [S4c] badge "actual" si es el tool seleccionado
	if i.current {
		badges = append(badges, uilist.Badge{
			Text:  "★ actual",
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

// --- Pantalla ---

type ToolPickerScreen struct {
	list    list.Model
	loading bool
	err     error

	repo       tool.Repository
	currentID  *int
	filterType tool.ScriptType
}

func NewToolPickerScreen(repo tool.Repository, currentID *int, filterType tool.ScriptType) ToolPickerScreen {
	l := list.New(nil, uilist.New(), 80, 20)
	title := "Elegir herramienta"
	if filterType != "" {
		title = "Elegir " + string(filterType)
	}
	l.Title = title
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = styles.Title
	l.Styles.HelpStyle = styles.Help

	return ToolPickerScreen{
		list:       l,
		repo:       repo,
		filterType: filterType,
		currentID:  currentID,
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
	return tea.Batch(tea.WindowSize(), load)
}

func (m ToolPickerScreen) Keys() []string {
	return []string{keys.NavConfirm, keys.NavBack, keys.NavFilter, keys.ViewHelp}
}

func (m ToolPickerScreen) Update(msg tea.Msg) (ScreenI, tea.Cmd) {
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

	case ActionMsg:
		switch msg.ID {
		case keys.NavConfirm:
			it, ok := m.list.SelectedItem().(toolItem)
			if !ok {
				return m, nil
			}
			return m, ToolSelected(it.tool)
		case keys.NavBack:
			return m, Back()
		case keys.ViewHelp:
			return m, OpenHelp()
		case keys.NavFilter:
			// dejar que la lista maneje "/"
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *ToolPickerScreen) setItems(tools []*tool.Tool) {
	items := make([]list.Item, len(tools))
	selectedIdx := -1
	for i, t := range tools {
		isCurrent := m.currentID != nil && t.ID == *m.currentID
		items[i] = toolItem{tool: t, current: isCurrent}
		if isCurrent {
			selectedIdx = i
		}
	}
	m.list.SetItems(items)
	if selectedIdx >= 0 {
		m.list.Select(selectedIdx)
	}
}

func (m ToolPickerScreen) View() string {
	switch {
	case m.loading:
		return states.Loading("tools")
	case m.err != nil:
		return states.Error(m.err)
	case len(m.list.Items()) == 0:
		return states.Empty("No hay tools disponibles", "Crea uno con el formulario")
	}
	return m.list.View()
}
