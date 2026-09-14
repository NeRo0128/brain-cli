package theme

var registry = map[string]Theme{
	"brain":       brain,
	"catppuccin":  catppuccin,
	"tokyo-night": tokyoNight,
	"nord":        nord,
	"rose-pine":   rosePine,
	"kanagawa":    kanagawa,
}

var orderedNames = []string{
	"brain",
	"catppuccin",
	"tokyo-night",
	"nord",
	"rose-pine",
	"kanagawa",
}

const DefaultName = "brain"

func Names() []string {
	out := make([]string, len(orderedNames))
	copy(out, orderedNames)
	return out
}

func Get(name string) Theme {
	if t, ok := registry[name]; ok {
		return t
	}
	return registry[DefaultName]
}

func Default() Theme { return registry[DefaultName] }

func Exists(name string) bool {
	_, ok := registry[name]
	return ok
}

func Next(current string) string {
	for i, n := range orderedNames {
		if n == current {
			return orderedNames[(i+1)%len(orderedNames)]
		}
	}
	return DefaultName
}
