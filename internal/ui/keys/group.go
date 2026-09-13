package keys

// Group categoriza bindings para el help screen y Settings.
// El ID técnico es en inglés (estable para config);
// el título se traduce en GroupTitles (capa de presentación).
type Group string

const (
	GroupNavigation Group = "navigation"
	GroupAction     Group = "action"
	GroupView       Group = "view"
	GroupEditing    Group = "editing"
	GroupApp        Group = "app"
)

// GroupTitles mapea IDs de grupo a títulos legibles.
// Cambiar esto cambia lo que se muestra; los IDs no.
var GroupTitles = map[Group]string{
	GroupNavigation: "Navegación",
	GroupAction:     "Acciones",
	GroupView:       "Vistas",
	GroupEditing:    "Edición",
	GroupApp:        "Aplicación",
}

// OrderedGroups define el orden de renderizado.
var OrderedGroups = []Group{
	GroupNavigation,
	GroupAction,
	GroupView,
	GroupEditing,
	GroupApp,
}

// Title devuelve el título legible del grupo.
// Si no está en el mapa, devuelve el ID crudo (nunca vacío).
func (g Group) Title() string {
	if t, ok := GroupTitles[g]; ok {
		return t
	}
	return string(g)
}

// Binding describe un atajo y su metadata.
// Keys son las teclas ACTIVAS (ya con overrides aplicados).
type Binding struct {
	ID    string
	Keys  []string
	Help  string
	Group Group
}
