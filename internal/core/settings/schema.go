package settings

// Field describe un setting editable.
type Field struct {
	Key         string
	Category    string
	Description string
	Values      []string // vacío = texto libre, con valores = enum
}

// Schema define el conjunto de settings editables.
// Se construye en main.go con los valores válidos de cada dominio.
type Schema struct {
	fields map[string]Field
	order  []string
}

func NewSchema(fields []Field) *Schema {
	s := &Schema{
		fields: make(map[string]Field, len(fields)),
		order:  make([]string, 0, len(fields)),
	}
	for _, f := range fields {
		s.fields[f.Key] = f
		s.order = append(s.order, f.Key)
	}
	return s
}

func (s *Schema) Get(key string) (Field, bool) {
	f, ok := s.fields[key]
	return f, ok
}

func (s *Schema) Has(key string) bool {
	_, ok := s.fields[key]
	return ok
}

// All devuelve los campos en el orden de inserción.
func (s *Schema) All() []Field {
	out := make([]Field, 0, len(s.order))
	for _, k := range s.order {
		out = append(out, s.fields[k])
	}
	return out
}
