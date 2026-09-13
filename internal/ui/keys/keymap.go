package keys

// KeyMap es una lista plana de bindings resueltos (con overrides aplicados).
// Lo consume el HelpScreen.
type KeyMap []Binding

// NewKeyMap resuelve IDs contra el Registry y devuelve un KeyMap plano.
// Si un ID no existe, se ignora.
func NewKeyMap(reg *Registry, ids []string) KeyMap {
	out := make(KeyMap, 0, len(ids))
	for _, id := range ids {
		if b, ok := reg.Get(id); ok {
			out = append(out, b)
		}
	}
	return out
}

// Grouped agrupa por grupo preservando el orden de OrderedGroups.
func (km KeyMap) Grouped() []GroupBindings {
	byGroup := make(map[Group][]Binding)
	for _, b := range km {
		byGroup[b.Group] = append(byGroup[b.Group], b)
	}
	out := make([]GroupBindings, 0, len(byGroup))
	for _, g := range OrderedGroups {
		if bs, ok := byGroup[g]; ok {
			out = append(out, GroupBindings{Group: g, Bindings: bs})
		}
	}
	return out
}
