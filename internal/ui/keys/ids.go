package keys

// IDs de bindings. Son strings estables y namespaced.
// El usuario los remapea desde config usando estos IDs.
//
// Convención: "<dominio>.<acción>"
// Añadir un ID aquí es la ÚNICA forma de exponer un atajo a config.
const (
	// --- Navegación ---
	NavUp       = "nav.up"
	NavDown     = "nav.down"
	NavPageUp   = "nav.page-up"
	NavPageDown = "nav.page-down"
	NavTop      = "nav.top"
	NavBottom   = "nav.bottom"
	NavConfirm  = "nav.confirm"
	NavBack     = "nav.back"
	NavFilter   = "nav.filter"

	// --- Acciones ---
	ActionExecute = "action.execute"
	ActionSave    = "action.save"
	ActionCancel  = "action.cancel"
	ActionRefresh = "action.refresh"
	ActionRerun   = "action.rerun"      // re-ejecutar task desde resultado

	// --- Vistas ---
	ViewDetail   = "view.detail"
	ViewHistory  = "view.history"
	ViewSettings = "view.settings"
	ViewHelp     = "view.help"

	// --- Edición ---
	EditNew    = "edit.new"
	EditUpdate = "edit.update"
	EditDelete = "edit.delete"

	// --- App ---
	AppQuit      = "app.quit"
	AppInterrupt = "app.interrupt"
)
