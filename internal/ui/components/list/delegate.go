package list

import (
	"fmt"
	"image/color"
	"io"
	"strings"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/NeRo0128/brain-cli/internal/ui/styles"
	"github.com/NeRo0128/brain-cli/internal/ui/theme"
)

// PrefixKind clasifica el prefix del row.
type PrefixKind int

const (
	PrefixNone PrefixKind = iota
	PrefixFavorite
	PrefixSuccess
	PrefixFailed
	PrefixCancelled
	PrefixRunning
	PrefixPending
)

// Row es la representación visual de un item.
type Row struct {
	Prefix   PrefixKind
	Title    string
	Badges   []Badge
	Subtitle string
	Meta     string
}

// RowProvider es lo que cada item debe implementar para
// ser renderizado por el Delegate.
type RowProvider interface {
	list.Item
	Row() Row
}

// Delegate es el renderizador unificado. Alto 2 líneas + 1 de spacing.
type Delegate struct {
	styles *styles.Styles
}

// New construye el delegate.
func New(s *styles.Styles) Delegate { return Delegate{styles: s} }

func (Delegate) Height() int                             { return 2 }
func (Delegate) Spacing() int                            { return 1 }
func (Delegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d Delegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	rp, ok := item.(RowProvider)
	if !ok {
		// Fallback: muestra el Title() del item.
		fmt.Fprint(w, item.FilterValue())
		return
	}
	row := rp.Row()
	selected := index == m.Index()

	p := d.styles.Theme.Resolve(d.styles.Dark)

	// --- Estilos ---
	// --- Cursor + prefix ---
	cursor := "  "
	titleStyle := lipgloss.NewStyle().Foreground(p.Text)
	if selected {
		cursor = d.styles.Icons.Selected + " "
		titleStyle = titleStyle.Bold(true).Foreground(p.Primary)
	}

	// --- Prefix ---
	prefix := ""
	if row.Prefix != PrefixNone {
		glyph, clr := d.resolvePrefix(row.Prefix, p)
		prefix = lipgloss.NewStyle().Foreground(clr).Render(glyph) + " "
	}

	// --- Línea 1 ---
	line1Left := cursor + prefix + titleStyle.Render(row.Title)
	line1Right := joinBadges(row.Badges, d.styles)

	width := m.Width() - 2 // margen lateral
	leftW := lipgloss.Width(line1Left)
	rightW := lipgloss.Width(line1Right)
	gap := width - leftW - rightW

	var line1 string
	switch {
	case rightW == 0 || gap < 2:
		line1 = line1Left
	default:
		line1 = line1Left + strings.Repeat(" ", gap) + line1Right
	}

	// --- Línea 2 ---
	indent := "   "
	line2 := indent
	if row.Subtitle != "" {
		line2 += d.styles.Subtitle.Render(row.Subtitle)
	}
	if row.Meta != "" {
		if row.Subtitle != "" {
			line2 += d.styles.Subtitle.Render(" · ")
		}
		line2 += d.styles.Subtitle.Render(row.Meta)
	}
	if row.Subtitle == "" && row.Meta == "" {
		line2 = "" // sin línea 2
	}

	if line2 == "" {
		fmt.Fprint(w, line1)
		return
	}
	fmt.Fprintf(w, "%s\n%s", line1, line2)
}
func (d Delegate) resolvePrefix(k PrefixKind, p theme.Palette) (string, color.Color) {
	ic := d.styles.Icons
	switch k {
	case PrefixFavorite:
		return ic.Favorite, p.Warning
	case PrefixSuccess:
		return ic.Success, p.Success
	case PrefixFailed:
		return ic.Failed, p.Error
	case PrefixCancelled:
		return ic.Cancelled, p.Warning
	case PrefixRunning:
		return ic.Running, p.Primary
	case PrefixPending:
		return ic.Pending, p.Muted
	}
	return "", p.Text
}
