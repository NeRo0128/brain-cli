package screens

import tea "github.com/charmbracelet/bubbletea"

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
	View() string
	Keys() []string
}
