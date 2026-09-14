package styles

import (
	lipgloss "charm.land/lipgloss/v2"

	"github.com/NeRo0128/brain-cli/internal/ui/theme"
)

// Styles agrupa todos los estilos visuales construidos desde un Theme.
//
// Uso:
//
//	s := styles.New(theme.Get("brain"), true)
//	s.Title.Render("Hola")
type Styles struct {
	Theme theme.Theme
	Dark  bool

	Title         lipgloss.Style
	Subtitle      lipgloss.Style
	SectionHeader lipgloss.Style
	Key           lipgloss.Style
	Help          lipgloss.Style

	SuccessStyle lipgloss.Style
	ErrorStyle   lipgloss.Style
	WarningStyle lipgloss.Style

	SpinnerStyle lipgloss.Style
	TimerStyle   lipgloss.Style

	InputFocusedStyle lipgloss.Style
	InputBlurredStyle lipgloss.Style

	BorderStyle lipgloss.Style

	BadgeMuted  lipgloss.Style
	KeyHintKey  lipgloss.Style
	KeyHintText lipgloss.Style

	IconSuccess lipgloss.Style
	IconWarning lipgloss.Style
	IconError   lipgloss.Style
	IconInfo    lipgloss.Style
}

// New construye el set de estilos desde un tema y el modo del terminal.
func New(t theme.Theme, isDark bool) Styles {
	p := t.Resolve(isDark)
	return Styles{
		Theme: t,
		Dark:  isDark,

		Title:         lipgloss.NewStyle().Bold(true).Foreground(p.Primary).Padding(0, 1),
		Subtitle:      lipgloss.NewStyle().Foreground(p.Muted),
		SectionHeader: lipgloss.NewStyle().Bold(true).Foreground(p.Secondary).MarginTop(1),
		Key:           lipgloss.NewStyle().Bold(true).Foreground(p.Secondary),
		Help:          lipgloss.NewStyle().Foreground(p.Muted).Padding(1, 2),

		SuccessStyle: lipgloss.NewStyle().Foreground(p.Success).Bold(true),
		ErrorStyle:   lipgloss.NewStyle().Foreground(p.Error).Bold(true),
		WarningStyle: lipgloss.NewStyle().Foreground(p.Warning).Bold(true),

		SpinnerStyle: lipgloss.NewStyle().Foreground(p.Primary).Bold(true),
		TimerStyle:   lipgloss.NewStyle().Foreground(p.Secondary),

		InputFocusedStyle: lipgloss.NewStyle().
			Foreground(p.Text).
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(p.InputFocused),
		InputBlurredStyle: lipgloss.NewStyle().
			Foreground(p.Muted).
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(p.InputBlurred),

		BorderStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(p.Border).
			Padding(0, 1),

		BadgeMuted:  lipgloss.NewStyle().Foreground(p.Muted).Faint(true),
		KeyHintKey:  lipgloss.NewStyle().Bold(true).Foreground(p.Secondary),
		KeyHintText: lipgloss.NewStyle().Foreground(p.Muted),

		IconSuccess: lipgloss.NewStyle().Foreground(p.Success),
		IconWarning: lipgloss.NewStyle().Foreground(p.Warning),
		IconError:   lipgloss.NewStyle().Foreground(p.Error),
		IconInfo:    lipgloss.NewStyle().Foreground(p.Primary),
	}
}
