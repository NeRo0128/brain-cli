package ui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/NeRo0128/brain-cli/internal/ui/components/frame"
	"github.com/NeRo0128/brain-cli/internal/ui/components/toast"
	"github.com/NeRo0128/brain-cli/internal/ui/keys"
	"github.com/NeRo0128/brain-cli/internal/ui/layout"
)

// View renderiza el estado completo de la aplicación.
func (m Model) View() tea.View {
	// --- Gate: terminal demasiado chica ---
	if layout.ShouldGate(m.width, m.height) {
		return tea.NewView(m.renderGate())
	}

	// --- Contenido base ---
	var content string
	if m.execErr != nil {
		content = m.renderError()
	} else {
		content = m.top().View().Content
	}

	// --- Frame ---
	header := frame.Header(m.headerData(), m.styles, m.width)
	footer := frame.Footer(m.footerBindings(), m.styles, m.width)

	// --- Ensamblado (footer pegado abajo) ---
	body := m.assembleBody(header, content, footer)

	// --- Overlay del toast (esquina superior derecha) ---
	v := tea.NewView(toast.Overlay(body, m.toast.View(), m.width))
	v.AltScreen = true
	return v
}

// --- Datos del header ---

// headerData arma los datos que consume el header.
//
// NOTA: hoy AIStatus=AIDisabled y los contadores son 0 porque
// todavía no cacheamos esa info en el Model. Cuando se añada el
// provider IA y se cacheen los contadores, se pueblan aquí.
func (m Model) headerData() frame.HeaderData {
	brandStyle := "minimal"
	if m.deps.Cfg != nil {
		brandStyle = m.deps.Cfg.UI.BrandStyle
	}
	return frame.HeaderData{
		AppName:       m.deps.Cfg.App.Name,
		Version:       m.deps.Version,
		BrandStyle:    brandStyle,
		AIStatus:      frame.AIDisabled,
		TaskCount:     0,
		FavoriteCount: 0,
	}
}

// --- Datos del footer ---

// footerBindings convierte los bindings activos (screen + globales
// imprescindibles) a la forma que espera frame.Footer.
//
// Añade AppQuit siempre porque el usuario debe saber cómo salir,
// aunque la screen no lo declare.
func (m Model) footerBindings() []frame.Binding {
	km := m.currentKeyMap()
	out := make([]frame.Binding, 0, len(km)+1)
	seen := make(map[string]bool, len(km)+1)

	for _, b := range km {
		if seen[b.ID] {
			continue
		}
		seen[b.ID] = true
		out = append(out, frame.Binding{Keys: b.Keys, Help: b.Help})
	}

	// Añadir q/ctrl+c si la screen no lo declaró.
	if !seen[keys.AppQuit] {
		if b, ok := m.deps.Keys.Get(keys.AppQuit); ok {
			out = append(out, frame.Binding{Keys: b.Keys, Help: b.Help})
		}
	}
	return out
}

// --- Ensamblado ---

// assembleBody compone header + contenido + filler + footer,
// rellenando líneas vacías para que el footer quede pegado abajo.
//
// Si el contenido ya es más alto que la terminal, no rellena
// y deja que el terminal haga scroll.
func (m Model) assembleBody(header, content, footer string) string {
	// Sin dimensiones conocidas: layout mínimo sin padding.
	if m.height <= 0 {
		return header + "\n\n" + content + "\n\n" + footer
	}

	// Reservamos 4 líneas de overhead:
	//   1 header
	//   1 blank después del header
	//   1 blank antes del footer
	//   1 footer
	contentHeight := m.height - 4
	if contentHeight < 1 {
		contentHeight = 1
	}

	// Convertimos el contenido a líneas exactas: truncar si sobra,
	// rellenar con vacías si falta.
	lines := strings.Split(content, "\n")
	if len(lines) > contentHeight {
		lines = lines[:contentHeight]
	}
	for len(lines) < contentHeight {
		lines = append(lines, "")
	}

	var b strings.Builder
	b.WriteString(header)
	b.WriteString("\n\n")
	b.WriteString(strings.Join(lines, "\n"))
	b.WriteString("\n\n")
	b.WriteString(footer)
	return b.String()
}

// --- Pantallas auxiliares ---

// renderError dibuja la pantalla de error global.
func (m Model) renderError() string {
	return fmt.Sprintf(
		"\n\n  ✗ Error: %v\n\n  Pulsa cualquier tecla para volver\n",
		m.execErr,
	)
}

// renderGate avisa cuando la terminal no cumple el mínimo.
func (m Model) renderGate() string {
	return fmt.Sprintf(
		"\n\n  ⚠  Terminal demasiado pequeña (%d×%d)\n\n"+
			"     Mínimo requerido: 70×20\n\n"+
			"     Redimensiona la terminal para continuar.\n",
		m.width, m.height,
	)
}
