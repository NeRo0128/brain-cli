// Package toast implementa notificaciones efímeras tipo "toast"
// que aparecen arriba a la derecha y se desvanecen solas.
//
// Tres piezas:
//   - model.go: estado, mensajes, ticker, API pública.
//   - view.go: render de una caja y de la pila.
//   - overlay.go: superposición sobre la vista base.
package toast

import (
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/NeRo0128/brain-cli/internal/ui/theme"
)

// --- Mensajes ---

// Level clasifica visualmente el toast.
type Level int

const (
	LevelInfo Level = iota
	LevelSuccess
	LevelWarning
	LevelError
)

// Msg pide mostrar un toast. Cualquier screen puede emitirlo.
type Msg struct {
	Message  string
	Level    Level
	Duration time.Duration // 0 = default por nivel
}

// TickMsg barre los expirados. Se re-agenda mientras haya toasts vivos.
type TickMsg struct{}

// --- Constantes ---

const (
	tickInterval  = 100 * time.Millisecond
	durationInfo  = 3 * time.Second
	durationOk    = 3 * time.Second
	durationWarn  = 3 * time.Second
	durationError = 5 * time.Second

	maxToasts = 3
)

// --- Estado ---

// Model mantiene la cola de toasts visibles.
type Model struct {
	toasts  []toast
	nextID  int
	palette theme.Palette
}

type toast struct {
	id       int
	message  string
	level    Level
	duration time.Duration
	expires  time.Time
}

// NewModel construye un modelo vacío.
func NewModel() Model { return Model{} }

// SetPalette establece la paleta de colores para el renderizado.
func (m *Model) SetPalette(p theme.Palette) { m.palette = p }

// --- API ---

// Tick es el comando que dispara el barrido periódico.
// Se llama UNA vez desde el Init del Model raíz.
func Tick() tea.Cmd {
	return tea.Tick(tickInterval, func(time.Time) tea.Msg { return TickMsg{} })
}

// Update procesa mensajes. Devuelve el modelo mutado y el próximo comando.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case Msg:
		dur := msg.Duration
		if dur <= 0 {
			dur = durationForLevel(msg.Level)
		}
		wasEmpty := len(m.toasts) == 0

		m.toasts = append(m.toasts, toast{
			id:       m.nextID,
			message:  msg.Message,
			level:    msg.Level,
			duration: dur,
			expires:  time.Now().Add(dur),
		})
		m.nextID++

		// Cap: si pasamos el máximo, tiramos los más viejos.
		if len(m.toasts) > maxToasts {
			m.toasts = m.toasts[len(m.toasts)-maxToasts:]
		}

		// Arrancar el ticker solo si no había ninguno vivo.
		if wasEmpty {
			return m, Tick()
		}
		return m, nil

	case TickMsg:
		m.cleanup()
		if len(m.toasts) > 0 {
			return m, Tick()
		}
		return m, nil
	}
	return m, nil
}

// cleanup elimina toasts expirados in-place.
func (m *Model) cleanup() {
	now := time.Now()
	alive := m.toasts[:0]
	for _, t := range m.toasts {
		if t.expires.After(now) {
			alive = append(alive, t)
		}
	}
	m.toasts = alive
}

// HasToasts indica si hay alguno visible (útil para tests).
func (m Model) HasToasts() bool { return len(m.toasts) > 0 }

// --- Helpers de emisión ---

func ShowInfo(msg string) tea.Cmd {
	return func() tea.Msg { return Msg{Message: msg, Level: LevelInfo} }
}
func ShowSuccess(msg string) tea.Cmd {
	return func() tea.Msg { return Msg{Message: msg, Level: LevelSuccess} }
}
func ShowWarning(msg string) tea.Cmd {
	return func() tea.Msg { return Msg{Message: msg, Level: LevelWarning} }
}
func ShowError(msg string) tea.Cmd {
	return func() tea.Msg { return Msg{Message: msg, Level: LevelError} }
}

func durationForLevel(l Level) time.Duration {
	switch l {
	case LevelError:
		return durationError
	case LevelWarning:
		return durationWarn
	case LevelSuccess:
		return durationOk
	}
	return durationInfo
}
