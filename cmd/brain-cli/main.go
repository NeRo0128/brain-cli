package main

import (
	"fmt"
	"os"
)

// Version will be set during build via ldflags
var Version = "dev"

func main() {
	fmt.Println("🧠 Brain CLI - Asistente Personal TUI")
	fmt.Printf("Versión: %s\n", Version)
	fmt.Println()
	fmt.Println("Estado: 🏗️ En desarrollo (Hito 1)")
	fmt.Println()
	fmt.Println("Próximos pasos:")
	fmt.Println("  1. Implementar configuración YAML")
	fmt.Println("  2. Crear entidades del dominio")
	fmt.Println("  3. Implementar UI básica con Bubble Tea")
	fmt.Println()
	fmt.Println("Ver AGENTS.md y docs/plans/development-roadmap.md para más detalles")
	
	os.Exit(0)
}
