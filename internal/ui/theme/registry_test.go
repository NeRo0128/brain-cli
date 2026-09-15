package theme_test

import (
	"testing"

	"github.com/NeRo0128/brain-cli/internal/ui/theme"
)

func TestNames_Stable(t *testing.T) {
	names := theme.Names()
	if len(names) != 6 {
		t.Fatalf("esperaba 6 temas, dio %d", len(names))
	}
	if names[0] != "brain" {
		t.Errorf("primer tema = %q", names[0])
	}
}

func TestGet_Known(t *testing.T) {
	tm := theme.Get("catppuccin")
	if tm.Name != "catppuccin" {
		t.Errorf("Name = %q", tm.Name)
	}
}

func TestGet_Unknown_FallsBackToDefault(t *testing.T) {
	tm := theme.Get("no-existe")
	if tm.Name != "brain" {
		t.Errorf("debería caer a brain, dio %q", tm.Name)
	}
}

func TestExists(t *testing.T) {
	if !theme.Exists("nord") {
		t.Error("nord debería existir")
	}
	if theme.Exists("inventado") {
		t.Error("inventado no debería existir")
	}
}

func TestNext_CyclesInOrder(t *testing.T) {
	cases := map[string]string{
		"brain":      "catppuccin",
		"catppuccin": "tokyo-night",
		"kanagawa":   "brain",
	}
	for from, want := range cases {
		if got := theme.Next(from); got != want {
			t.Errorf("Next(%q) = %q, quiero %q", from, got, want)
		}
	}
}

func TestResolve_DarkOnly(t *testing.T) {
	// Usar un tema que sí sea dark-only.
	tm := theme.Get("nord")
	if !tm.IsDarkOnly() {
		t.Fatal("nord debería ser dark-only")
	}
	p := tm.Resolve(false)
	if p.Text != tm.Dark.Text {
		t.Error("dark-only en light debería devolver Dark")
	}
}

func TestResolve_LightAvailable(t *testing.T) {
	tm := theme.Get("catppuccin")
	if tm.IsDarkOnly() {
		t.Fatal("catppuccin NO debería ser dark-only")
	}
	pLight := tm.Resolve(false)
	pDark := tm.Resolve(true)
	if pLight.Text == pDark.Text {
		t.Error("light y dark deberían tener colores distintos")
	}
}
