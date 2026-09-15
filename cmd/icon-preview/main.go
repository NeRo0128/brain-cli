package main

import (
	"fmt"

	"github.com/NeRo0128/brain-cli/internal/ui/icons"
	"github.com/NeRo0128/brain-cli/internal/ui/theme"
)

func main() {
	sets := []icons.Set{
		icons.Unicode,
		icons.NerdB,
		icons.NerdC,
	}

	// Preview plano (sin color)
	fmt.Println("═══ SIN COLOR ═══")
	fmt.Printf("%-18s", "GLYPH")
	for _, s := range sets {
		fmt.Printf("  %-12s", s.Name)
	}
	fmt.Println()
	fmt.Println("─────────────────────────────────────────────────────")

	type row struct {
		name   string
		getter func(icons.Set) string
	}

	rows := []row{
		{"Success", func(s icons.Set) string { return s.Success }},
		{"Failed", func(s icons.Set) string { return s.Failed }},
		{"Cancelled", func(s icons.Set) string { return s.Cancelled }},
		{"Running", func(s icons.Set) string { return s.Running }},
		{"Pending", func(s icons.Set) string { return s.Pending }},
		{"Selected", func(s icons.Set) string { return s.Selected }},
		{"Favorite", func(s icons.Set) string { return s.Favorite }},
		{"TypeScript", func(s icons.Set) string { return s.TypeScript }},
		{"TypeCommand", func(s icons.Set) string { return s.TypeCommand }},
		{"TypeAI", func(s icons.Set) string { return s.TypeAI }},
		{"Brand", func(s icons.Set) string { return s.Brand }},
		{"Danger", func(s icons.Set) string { return s.Danger }},
	}

	for _, r := range rows {
		fmt.Printf("%-18s", r.name)
		for _, s := range sets {
			fmt.Printf("  %-12s", r.getter(s))
		}
		fmt.Println()
	}

	// Preview con color (usa el tema brain)

	for i := 0; i < len(theme.Names()); i++ {
		tm := theme.Get(theme.Names()[i])
		fmt.Println("\n═══ CON COLOR (tema ", tm.Name)
		colored := icons.NewColored(icons.NerdB, true, tm)

		fmt.Printf("  %s completado\n", colored.Success())
		fmt.Printf("  %s falló\n", colored.Failed())
		fmt.Printf("  %s cancelado\n", colored.Cancelled())
		fmt.Printf("  %s running\n", colored.Running())
		fmt.Printf("  %s pending\n", colored.Pending())
		fmt.Printf("  %s selected\n", colored.Selected())
		fmt.Printf("  %s favorite\n", colored.Favorite())
		fmt.Printf("  %s script\n", colored.TypeScript())
		fmt.Printf("  %s command\n", colored.TypeCommand())
		fmt.Printf("  %s ai\n", colored.TypeAI())
		fmt.Printf("  %s brand\n", colored.Brand())
		fmt.Printf("  %s danger\n", colored.Danger())
	}

}
