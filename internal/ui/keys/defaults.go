package keys

// Defaults es el catálogo base.
// Los IDs son estables; las Keys pueden sobreescribirse desde config.
//
// Nota de diseño: "e" es EditUpdate (main) y NO ActionExecute.
// Para ejecutar desde main usa enter. Desde detalle, enter.
// Si quieres ejecutar con "e" desde detalle, lo manejas en el Update de DetailScreen.
var Defaults = []Binding{
	// --- Navegación ---
	{ID: NavUp, Keys: []string{"up", "k"}, Help: "subir", Group: GroupNavigation},
	{ID: NavDown, Keys: []string{"down", "j"}, Help: "bajar", Group: GroupNavigation},
	{ID: NavPageUp, Keys: []string{"pgup"}, Help: "página arriba", Group: GroupNavigation},
	{ID: NavPageDown, Keys: []string{"pgdown"}, Help: "página abajo", Group: GroupNavigation},
	{ID: NavTop, Keys: []string{"g"}, Help: "ir al inicio", Group: GroupNavigation},
	{ID: NavBottom, Keys: []string{"G"}, Help: "ir al final", Group: GroupNavigation},
	{ID: NavConfirm, Keys: []string{"enter"}, Help: "confirmar", Group: GroupNavigation},
	{ID: NavBack, Keys: []string{"esc"}, Help: "volver", Group: GroupNavigation},
	{ID: NavFilter, Keys: []string{"/"}, Help: "filtrar", Group: GroupNavigation},

	// --- Acciones ---
	{ID: ActionExecute, Keys: []string{"enter"}, Help: "ejecutar", Group: GroupAction},
	{ID: ActionSave, Keys: []string{"ctrl+s"}, Help: "guardar", Group: GroupAction},
	{ID: ActionCancel, Keys: []string{"esc"}, Help: "cancelar", Group: GroupAction},
	{ID: ActionRefresh, Keys: []string{"ctrl+r"}, Help: "refrescar", Group: GroupAction},
	{ID: ActionRerun, Keys: []string{"r"}, Help: "re-ejecutar", Group: GroupAction},

	// --- Vistas ---
	{ID: ViewDetail, Keys: []string{"d"}, Help: "detalle", Group: GroupView},
	{ID: ViewHistory, Keys: []string{"h"}, Help: "historial", Group: GroupView},
	{ID: ViewSettings, Keys: []string{"s"}, Help: "ajustes", Group: GroupView},
	{ID: ViewHelp, Keys: []string{"?"}, Help: "ayuda", Group: GroupView},

	// --- Edición ---
	{ID: EditNew, Keys: []string{"n"}, Help: "nueva", Group: GroupEditing},
	{ID: EditUpdate, Keys: []string{"e"}, Help: "editar", Group: GroupEditing},
	{ID: EditDelete, Keys: []string{"ctrl+d"}, Help: "borrar", Group: GroupEditing},

	// --- App ---
	{ID: AppQuit, Keys: []string{"q", "ctrl+c"}, Help: "salir", Group: GroupApp},
	{ID: AppInterrupt, Keys: []string{"ctrl+c"}, Help: "interrumpir", Group: GroupApp},
	{ID: AppTheme, Keys: []string{"ctrl+t"}, Help: "tema", Group: GroupApp},
}
