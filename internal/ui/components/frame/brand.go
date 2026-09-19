package frame

import "strings"

// BrandStyle define el estilo del título ASCII.
type BrandStyle int

const (
	BrandMinimal BrandStyle = iota // "🧠 brain-cli" tradicional
	BrandSlim                      // 3 líneas, half-blocks
	BrandBig                       // 6 líneas, ANSI Shadow
)

// brandSlimASCII: 19 cols x 3 líneas. Unicode half-blocks.
// No depende de Nerd Font.
const brandSlimASCII = `█▀█ █▀█ ▄▀▄ ▀█▀ █▀█
█▀▄ ██▀ ███  █  ███
█▄█ █▀█ █▀█ ▄█▄ █▄█`

// brandBigASCII: 39 cols x 6 líneas. ANSI Shadow clásico.
const brandBigASCII = `██████╗ ██████╗  █████╗ ██╗███╗   ██╗    ██████╗██╗     ██╗
██╔══██╗██╔══██╗██╔══██╗██║████╗  ██║   ██╔════╝██║     ██║
██████╔╝██████╔╝███████║██║██╔██╗ ██║   ██║     ██║     ██║
██╔══██╗██╔══██╗██╔══██║██║██║╚██╗██║── ██║     ██║     ██║
██████╔╝██║  ██║██║  ██║██║██║ ╚████║   ╚██████╗███████╗██║
╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝╚═╝╚═╝  ╚═══╝    ╚═════╝╚══════╝╚═╝`

// BrandLines devuelve el ASCII art según el estilo.
// Devuelve nil para BrandMinimal (el header usa texto plano).
func BrandLines(style BrandStyle) []string {
	switch style {
	case BrandSlim:
		return splitLines(brandSlimASCII)
	case BrandBig:
		return splitLines(brandBigASCII)
	}
	return nil
}

// ParseBrandStyle convierte un string de config a BrandStyle.
func ParseBrandStyle(s string) BrandStyle {
	switch s {
	case "ascii-slim", "slim":
		return BrandSlim
	case "ascii-big", "big", "ascii":
		return BrandBig
	default:
		return BrandMinimal
	}
}

func splitLines(s string) []string {
	return strings.Split(s, "\n")
}
