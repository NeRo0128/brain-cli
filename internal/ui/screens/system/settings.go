package system

import (
	"context"
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"

	core "github.com/NeRo0128/brain-cli/internal/core/config"
	"github.com/NeRo0128/brain-cli/internal/ui/icons"
	"github.com/NeRo0128/brain-cli/internal/ui/keys"
	"github.com/NeRo0128/brain-cli/internal/ui/screens"
	"github.com/NeRo0128/brain-cli/internal/ui/styles"
	"github.com/NeRo0128/brain-cli/internal/ui/theme"
	settingsuc "github.com/NeRo0128/brain-cli/internal/usecases/settings"
	"github.com/rs/zerolog"
)

// --- Tipos ---

type tabID int

const (
	tabApariencia tabID = iota
	tabCuenta
	tabSistema
	tabInfo
)

var tabTitles = []string{"Apariencia", "Cuenta", "Sistema", "Info"}

type fieldKind int

const (
	fkEnum fieldKind = iota
	fkToggle
	fkAction
	fkInfo
)

type field struct {
	key    string
	label  string
	kind   fieldKind
	values []string
}

// --- Mensajes internos ---

type settingsLoadedMsg struct {
	overrides map[string]string
	err       error
}

type settingsSavedMsg struct {
	count int
	err   error
}

// --- Screen ---

type SettingsScreen struct {
	mgr    *settingsuc.Manager
	log    zerolog.Logger
	styles *styles.Styles

	base     core.Config
	original map[string]string
	draft    map[string]string

	activeTab tabID
	focus     int

	loading bool
	saving  bool
	err     error
	toast   string

	width, height int
}

func NewSettingsScreen(mgr *settingsuc.Manager, log zerolog.Logger, s *styles.Styles) SettingsScreen {
	return SettingsScreen{
		mgr:       mgr,
		log:       log.With().Str("screen", "settings").Logger(),
		styles:    s,
		base:      mgr.Base(),
		original:  map[string]string{},
		draft:     map[string]string{},
		activeTab: tabApariencia,
		loading:   true,
	}
}

func (m SettingsScreen) Init() tea.Cmd {
	return tea.Batch(
		func() tea.Msg { return tea.RequestWindowSize() },
		m.loadOverrides(),
	)
}

func (m SettingsScreen) loadOverrides() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		ov, err := m.mgr.Overrides(ctx)
		return settingsLoadedMsg{overrides: ov, err: err}
	}
}

func (m SettingsScreen) Keys() []string {
	return []string{
		keys.ViewSettings,
		keys.NavBack,
		keys.ViewHelp,
		keys.ActionSave,
		keys.ActionExecute,
	}
}

func (m SettingsScreen) Update(msg tea.Msg) (screens.ScreenI, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case settingsLoadedMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.original = msg.overrides
		m.draft = cloneMap(msg.overrides)
		return m, nil

	case settingsSavedMsg:
		m.saving = false
		if msg.err != nil {
			m.err = msg.err
			m.toast = ""
			return m, nil
		}
		m.original = cloneMap(m.draft)
		m.toast = fmt.Sprintf("%d ajuste(s) guardado(s)", msg.count)
		m.err = nil
		return m, func() tea.Msg { return screens.SettingsChangedMsg{} }

	case screens.ActionMsg:
		return m.handleAction(msg)

	case tea.KeyPressMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m SettingsScreen) handleAction(msg screens.ActionMsg) (screens.ScreenI, tea.Cmd) {
	switch msg.ID {
	case keys.NavBack:
		return m, screens.Back()
	case keys.ViewHelp:
		return m, screens.OpenHelp()
	case keys.ActionSave:
		return m.startSave()
	case keys.ActionExecute:
		return m.handleEnter()
	}
	return m, nil
}

func (m SettingsScreen) handleEnter() (screens.ScreenI, tea.Cmd) {
	f := m.currentField()
	switch f.label {
	case "Cuenta GitHub":
		return m, screens.OpenAuth(m.styles)
	}
	return m, nil
}

func (m SettingsScreen) handleKey(msg tea.KeyPressMsg) (screens.ScreenI, tea.Cmd) {
	switch msg.String() {
	case "tab":
		m.activeTab = (m.activeTab + 1) % tabID(len(tabTitles))
		m.focus = 0
	case "shift+tab":
		m.activeTab = (m.activeTab - 1 + tabID(len(tabTitles))) % tabID(len(tabTitles))
		m.focus = 0
	case "up", "k":
		if m.focus > 0 {
			m.focus--
		}
	case "down", "j":
		fields := m.fieldsFor(m.activeTab)
		if m.focus < len(fields)-1 {
			m.focus++
		}
	case "left":
		return m.cycleValue(-1)
	case "right":
		return m.cycleValue(+1)
	}
	return m, nil
}

func (m SettingsScreen) cycleValue(dir int) (screens.ScreenI, tea.Cmd) {
	f := m.currentField()
	if f.kind != fkEnum || len(f.values) == 0 {
		return m, nil
	}
	current := m.effectiveValue(f.key)
	idx := indexOf(f.values, current)
	if idx < 0 {
		idx = 0
	}
	next := (idx + dir + len(f.values)) % len(f.values)
	m.draft[f.key] = f.values[next]
	return m, nil
}

func (m SettingsScreen) startSave() (screens.ScreenI, tea.Cmd) {
	if m.saving {
		return m, nil
	}
	diff := m.diff()
	if len(diff) == 0 {
		m.toast = "sin cambios"
		return m, nil
	}
	m.saving = true
	m.toast = ""
	m.err = nil

	mgr := m.mgr
	return m, func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := mgr.SetMany(ctx, diff)
		return settingsSavedMsg{count: len(diff), err: err}
	}
}

// --- Navegación y datos ---

func (m SettingsScreen) currentField() field {
	fields := m.fieldsFor(m.activeTab)
	if m.focus < 0 || m.focus >= len(fields) {
		return field{}
	}
	return fields[m.focus]
}

func (m SettingsScreen) fieldsFor(tab tabID) []field {
	switch tab {
	case tabApariencia:
		return []field{
			{key: "ui.theme", label: "Tema", kind: fkEnum, values: theme.Names()},
			{key: "ui.icons", label: "Icons", kind: fkEnum, values: icons.Names()},
			{key: "ui.brand_style", label: "Brand style", kind: fkEnum,
				values: []string{"minimal", "slim", "big"}},
		}
	case tabCuenta:
		return []field{
			{label: "Cuenta GitHub", kind: fkAction},
		}
	case tabSistema:
		return []field{
			{key: "logging.level", label: "Log level", kind: fkEnum,
				values: []string{"debug", "info", "warn", "error"}},
			{key: "logging.format", label: "Log format", kind: fkEnum,
				values: []string{"pretty", "json"}},
		}
	case tabInfo:
		return []field{
			{label: "Versión", kind: fkInfo},
		}
	}
	return nil
}

func (m SettingsScreen) effectiveValue(key string) string {
	if v, ok := m.draft[key]; ok && v != "" {
		return v
	}
	return valueFromConfig(&m.base, key)
}

func (m SettingsScreen) diff() map[string]string {
	out := map[string]string{}
	for k, v := range m.draft {
		if m.original[k] != v {
			out[k] = v
		}
	}
	return out
}

func (m SettingsScreen) dirty() bool {
	return len(m.diff()) > 0
}

// --- View ---

func (m SettingsScreen) View() tea.View {
	var b strings.Builder

	title := m.styles.Title.Render("Ajustes")
	if m.dirty() {
		title += " " + m.styles.WarningStyle.Render("*")
	}
	b.WriteString(title)
	b.WriteString("\n\n")

	b.WriteString(m.renderTabs())
	b.WriteString("\n\n")

	switch {
	case m.loading:
		b.WriteString(m.styles.Subtitle.Render("Cargando..."))
	case m.err != nil:
		b.WriteString(m.styles.ErrorStyle.Render(m.err.Error()))
	default:
		b.WriteString(m.renderFields())
	}

	if m.toast != "" {
		b.WriteString("\n")
		b.WriteString(m.styles.SuccessStyle.Render(m.toast))
	}

	return tea.NewView(b.String())
}

func (m SettingsScreen) renderTabs() string {
	parts := make([]string, 0, len(tabTitles))
	for i, title := range tabTitles {
		style := m.styles.Subtitle
		if tabID(i) == m.activeTab {
			style = m.styles.Key
		}
		parts = append(parts, style.Render(title))
	}
	return strings.Join(parts, "   ")
}

func (m SettingsScreen) renderFields() string {
	fields := m.fieldsFor(m.activeTab)
	var b strings.Builder
	for i, f := range fields {
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString(m.renderField(f, i == m.focus))
		b.WriteString("\n")
	}
	return b.String()
}

func (m SettingsScreen) renderField(f field, focused bool) string {
	cursor := "  "
	if focused {
		cursor = "▶ "
	}

	labelStyle := m.styles.Subtitle
	if focused {
		labelStyle = m.styles.Key
	}

	var value string
	switch f.kind {
	case fkEnum:
		current := m.effectiveValue(f.key)
		if focused {
			value = m.styles.Key.Render("◀ ") + current + m.styles.Key.Render(" ▶")
		} else {
			value = m.styles.Subtitle.Render(current)
		}
	case fkAction:
		if focused {
			value = m.styles.Key.Render("[ Enter ]")
		} else {
			value = m.styles.Subtitle.Render("[ Enter ]")
		}
	case fkInfo:
		value = m.styles.Subtitle.Render("vdev")
	}

	return cursor + labelStyle.Render(f.label) + "\n   " + value
}

// --- Helpers ---

func valueFromConfig(cfg *core.Config, key string) string {
	switch key {
	case "ui.theme":
		return cfg.UI.Theme
	case "ui.icons":
		return cfg.UI.Icons
	case "ui.brand_style":
		return cfg.UI.BrandStyle
	case "logging.level":
		return cfg.Logging.Level
	case "logging.format":
		return cfg.Logging.Format
	}
	return ""
}

func cloneMap(in map[string]string) map[string]string {
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func indexOf(s []string, v string) int {
	for i, x := range s {
		if x == v {
			return i
		}
	}
	return -1
}
