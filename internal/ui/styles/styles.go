package styles

import "github.com/charmbracelet/lipgloss"

// Estilos compuestos.
//
// Los COLORES viven en tokens.go.
// Los HELPERS de estado viven en status.go.
// Aquí solo van estilos reusables por 2+ screens.
//
// Regla: si un estilo solo lo usa UNA screen, va en su archivo.

// ============================================================
// Encabezados
// ============================================================

// Title: título principal de una screen.
var Title = lipgloss.NewStyle().
	Bold(true).
	Foreground(Primary).
	Padding(0, 1)

// Subtitle: texto secundario, contexto, contadores.
var Subtitle = lipgloss.NewStyle().
	Foreground(Muted)

// SectionHeader: encabezado de sección dentro de una screen.
var SectionHeader = lipgloss.NewStyle().
	Bold(true).
	Foreground(Secondary).
	MarginTop(1)

// ============================================================
// Contenido
// ============================================================

// Key: teclas de atajo en el footer o en línea.
var Key = lipgloss.NewStyle().
	Bold(true).
	Foreground(Secondary)

// Help: bloque de ayuda (usado por el footer de screens).
var Help = lipgloss.NewStyle().
	Foreground(Muted).
	Padding(1, 2)

// ============================================================
// Estado (iconos / badges)
// ============================================================

var (
	SuccessStyle = lipgloss.NewStyle().Foreground(Success).Bold(true)
	ErrorStyle   = lipgloss.NewStyle().Foreground(Error).Bold(true)
	WarningStyle = lipgloss.NewStyle().Foreground(Warning).Bold(true)
)

// ============================================================
// Específicos de screens concretas
// ============================================================

// SpinnerStyle: usado por ExecutingScreen.
var SpinnerStyle = lipgloss.NewStyle().
	Foreground(Primary).
	Bold(true)

// TimerStyle: cronómetro de ExecutingScreen.
var TimerStyle = lipgloss.NewStyle().
	Foreground(Secondary)

// ============================================================
// Inputs — foco visible en ambos temas
// ============================================================

// InputFocusedStyle: borde inferior acentuado, texto normal.
// El fondo + borde hacen visible el foco incluso sin color.
var InputFocusedStyle = lipgloss.NewStyle().
	Foreground(Text).
	Border(lipgloss.NormalBorder(), false, false, true, false).
	BorderForeground(InputFocused)

// InputBlurredStyle: borde inferior tenue, texto atenuado.
var InputBlurredStyle = lipgloss.NewStyle().
	Foreground(Muted).
	Border(lipgloss.NormalBorder(), false, false, true, false).
	BorderForeground(InputBlurred)

// ============================================================
// Paneles / bordes
// ============================================================

// BorderStyle: caja redondeada con padding, para agrupar bloques.
var BorderStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(Border).
	Padding(0, 1)

// ============================================================
// Badges y metadatos
// ============================================================

// BadgeMuted: texto tenue entre corchetes para metadatos.
var BadgeMuted = lipgloss.NewStyle().
	Foreground(Muted).
	Faint(true)

// KeyHint: "Enter" resaltado, "guardar" atenuado — para footers ricos.
var KeyHintKey = lipgloss.NewStyle().
	Bold(true).
	Foreground(Secondary)

var KeyHintText = lipgloss.NewStyle().
	Foreground(Muted)
