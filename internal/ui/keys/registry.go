package keys

import "slices"

// Overrides mapea ID → keys custom. Viene de config.yaml.
// Ej: {"action.execute": ["ctrl+r", "enter"]}
type Overrides map[string][]string

// Registry resuelve bindings aplicando overrides sobre Defaults.
// Es inmutable una vez construido (seguro para leer desde múltiples goroutines).
type Registry struct {
	bindings map[string]Binding
}

// New construye un Registry con overrides opcionales.
// Si overrides es nil o un ID no existe, se ignora silenciosamente.
func New(overrides Overrides) *Registry {
	r := &Registry{bindings: make(map[string]Binding, len(Defaults))}
	for _, b := range Defaults {
		r.bindings[b.ID] = b
	}
	for id, newKeys := range overrides {
		if b, ok := r.bindings[id]; ok && len(newKeys) > 0 {
			b.Keys = newKeys
			r.bindings[b.ID] = b
		}
	}
	return r
}

// Get devuelve el binding por ID.
func (r *Registry) Get(id string) (Binding, bool) {
	b, ok := r.bindings[id]
	return b, ok
}

// Matches indica si una tecla dispara el binding con ese ID.
func (r *Registry) Matches(id, key string) bool {
	b, ok := r.bindings[id]
	if !ok {
		return false
	}
	return slices.Contains(b.Keys, key)
}

// Grouped devuelve los bindings agrupados por categoría,
// incluyendo SOLO los IDs pasados (los que maneja una screen).
// Respeta el orden de OrderedGroups.
//
// Uso: reg.Grouped(keys.ActionExecute, keys.NavUp, keys.NavDown)
func (r *Registry) Grouped(ids ...string) []GroupBindings {
	// Agrupa por grupo
	byGroup := make(map[Group][]Binding)
	for _, id := range ids {
		if b, ok := r.bindings[id]; ok {
			byGroup[b.Group] = append(byGroup[b.Group], b)
		}
	}

	// Ordena según OrderedGroups
	out := make([]GroupBindings, 0, len(byGroup))
	for _, g := range OrderedGroups {
		if bs, ok := byGroup[g]; ok {
			out = append(out, GroupBindings{Group: g, Bindings: bs})
		}
	}
	return out
}

// GroupBindings agrupa bindings bajo un título.
// Lo consume el help screen y Settings.
type GroupBindings struct {
	Group    Group
	Bindings []Binding
}
