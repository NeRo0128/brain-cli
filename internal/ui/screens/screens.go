package screens

import (
	tea "charm.land/bubbletea/v2"
)

// ScreenI es una pantalla navegable de la TUI.
//
// Contrato:
//   - Init se llama UNA vez al pushear la screen.
//   - Update devuelve la screen mutada y un comando.
//   - View devuelve el string a renderizar.
//   - Keys devuelve los IDs del Registry que la screen maneja.
//     El Model los usa para traducir teclas a ActionMsg.
//
// Las screens NO conocen al Model. Emiten intenciones de navegación
// y acciones vía mensajes del paquete screens (BackMsg, OpenDetailMsg...).
type ScreenI interface {
	Init() tea.Cmd
	Update(tea.Msg) (ScreenI, tea.Cmd)
	View() tea.View
	Keys() []string
	// WantsTextInput() bool
}

const twoColMinWidth = 120

// maxTextAreaPreviewLines limita cuántas líneas del script se muestran
// en el detalle. Si el script tiene más, se añade un indicador.
// El usuario puede editarlo con `e` para ver el contenido completo.
const maxTextAreaPreviewLines = 25

const helpTwoColMinWidth = 100

func boolYesNo(v bool) string {
	if v {
		return "sí"
	}
	return "no"
}
