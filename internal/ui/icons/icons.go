package icons

// Set agrupa los glyphs de un estilo de iconografía.
// Cada glyph tiene un comentario con su codepoint para buscarlo en el cheat sheet.
type Set struct {
	Name string

	// === Estados de ejecución ===
	SuccessA string // ✓ U+2713
	FailedA  string // ✗ U+2717
	CancelA  string // ⊘ U+2298
	RunA     string // ● U+25CF
	PendA    string // ○ U+25CB

	SuccessB string // nf-fa-check_circle \uf058
	FailedB  string // nf-fa-times_circle \uf057
	CancelB  string // nf-fa-ban \uf05e
	RunB     string // nf-fa-spinner \uf110
	PendB    string // nf-fa-circle_o \uf10c

	SuccessC string // nf-fa-check_square \uf14a
	FailedC  string // nf-fa-times_square \uf2d3
	CancelC  string // nf-fa-square \uf0c8
	RunC     string // nf-fa-play_circle \uf144
	PendC    string // nf-fa-square_o \uf096

	// === Navegación ===
	SelA string // ▸ U+25B8
	SelB string // nf-fa-caret_right \uf0da
	SelC string // nf-oct-chevron_right \uf460

	FavA string // ★ U+2605
	FavB string // nf-fa-star \uf005
	FavC string // nf-md-star \uf4ce

	// === Tipos de tarea ===
	TypeScriptA string // ⚙ U+2699
	TypeScriptB string // nf-fa-gear \uf013
	TypeScriptC string // nf-md-script \uf3a8

	TypeCommandA string // ❯ U+276F
	TypeCommandB string // nf-fa-terminal \uf120
	TypeCommandC string // nf-md-console \uf18d

	TypeAIA string // ✦ U+2726
	TypeAIB string // nf-fa-magic \uf0d0
	TypeAIC string // nf-md-robot \uf48a

	// === Branding ===
	BrandA string // ◆ U+25C6
	BrandB string // nf-fa-brain \uf5dc
	BrandC string // nf-md-brain \uf4d4

	// === Modales ===
	DangerA string // ⚠ U+26A0
	DangerB string // nf-fa-exclamation_triangle \uf071
	DangerC string // nf-md-alert \uf026

	// === Aliases (los que usa el código existente) ===
	Success     string
	Failed      string
	Cancelled   string
	Running     string
	Pending     string
	Warning     string // [FIX] ahora poblado
	Info        string // [FIX] ahora poblado
	Selected    string
	Favorite    string
	Unfavorite  string // [FIX] ahora poblado
	Bullet      string // [FIX] ahora poblado
	TypeScript  string
	TypeCommand string
	TypeAI      string
	Brand       string
	Danger      string
}

// Unicode es el set por defecto. Solo símbolos geométricos.
var Unicode = Set{
	Name: "unicode",

	SuccessA: "✓", FailedA: "✗", CancelA: "⊘", RunA: "●", PendA: "○",
	SelA: "▸", FavA: "★",
	TypeScriptA: "⚙", TypeCommandA: "❯", TypeAIA: "✦",
	BrandA:  "◆",
	DangerA: "⚠",

	// Aliases
	Success:    "✓",
	Failed:     "✗",
	Cancelled:  "⊘",
	Running:    "●",
	Pending:    "○",
	Warning:    "▲", // U+25B2
	Info:       "•", // U+2022
	Selected:   "▸",
	Favorite:   "★",
	Unfavorite: "☆", // U+2606
	Bullet:     "•",
	TypeScript: "⚙", TypeCommand: "❯", TypeAI: "✦",
	Brand: "◆", Danger: "⚠",
}

// NerdB es el set balanceado. Círculos rellenos + Font Awesome.
var NerdB = Set{
	Name: "nerd-b",

	SuccessB: "\uf058", FailedB: "\uf057", CancelB: "\uf05e",
	RunB: "\uf110", PendB: "\uf10c",
	SelB: "\uf0da", FavB: "\uf005",
	TypeScriptB: "\uf013", TypeCommandB: "\uf120", TypeAIB: "\uf0d0",
	BrandB:  "\uf5dc",
	DangerB: "\uf071",

	// Aliases
	Success:    "\uf058",
	Failed:     "\uf057",
	Cancelled:  "\uf05e",
	Running:    "\uf110",
	Pending:    "\uf10c",
	Warning:    "\uf071", // nf-fa-exclamation_triangle
	Info:       "\uf05a", // nf-fa-info_circle
	Selected:   "\uf0da",
	Favorite:   "\uf005",
	Unfavorite: "\uf006", // nf-fa-star_o
	Bullet:     "\uf111", // nf-fa-circle (sólido)
	TypeScript: "\uf013", TypeCommand: "\uf120", TypeAI: "\uf0d0",
	Brand: "\uf5dc", Danger: "\uf071",
}

// NerdC es el set "techy". Cuadrados + Material Design.
var NerdC = Set{
	Name: "nerd-c",

	SuccessC: "\uf14a", FailedC: "\uf2d3", CancelC: "\uf0c8",
	RunC: "\uf144", PendC: "\uf096",
	SelC: "\uf460", FavC: "\uf4ce",
	TypeScriptC: "\uf3a8", TypeCommandC: "\uf18d", TypeAIC: "\uf48a",
	BrandC:  "\uf4d4",
	DangerC: "\uf026",

	// Aliases
	Success:    "\uf14a",
	Failed:     "\uf2d3",
	Cancelled:  "\uf0c8",
	Running:    "\uf144",
	Pending:    "\uf096",
	Warning:    "\uf026", // nf-md-alert
	Info:       "\uf02a", // nf-md-information
	Selected:   "\uf460",
	Favorite:   "\uf4ce",
	Unfavorite: "\uf4d2", // nf-md-star_outline
	Bullet:     "\uf444", // nf-md-circle_small
	TypeScript: "\uf3a8", TypeCommand: "\uf18d", TypeAI: "\uf48a",
	Brand: "\uf4d4", Danger: "\uf026",
}

// --- Registry ---

var registry = map[string]Set{
	"unicode": Unicode,
	"nerd":    NerdB,
	"nerd-b":  NerdB,
	"nerd-c":  NerdC,
}

func Get(name string) Set {
	if s, ok := registry[name]; ok {
		return s
	}
	return Unicode
}

func Names() []string {
	return []string{"unicode", "nerd-b", "nerd-c"}
}

func Exists(name string) bool {
	_, ok := registry[name]
	return ok
}

const DefaultName = "unicode"
