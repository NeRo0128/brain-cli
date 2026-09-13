package keys_test

import (
	"testing"

	"github.com/NeRo0128/brain-cli/internal/ui/keys"
)

func TestNew_NoOverrides_UsesDefaults(t *testing.T) {
	r := keys.New(nil)
	b, ok := r.Get(keys.ActionExecute)
	if !ok {
		t.Fatal("action.execute no encontrado")
	}
	if len(b.Keys) != 1 || b.Keys[0] != "enter" {
		t.Errorf("keys = %v, quiero [enter]", b.Keys)
	}
}

func TestNew_WithOverrides_Replaces(t *testing.T) {
	r := keys.New(keys.Overrides{
		keys.ActionExecute: {"ctrl+r", "enter"},
	})
	b, _ := r.Get(keys.ActionExecute)
	if len(b.Keys) != 2 || b.Keys[0] != "ctrl+r" {
		t.Errorf("keys = %v", b.Keys)
	}
}

func TestNew_UnknownID_Ignored(t *testing.T) {
	r := keys.New(keys.Overrides{
		"no.existe": {"x"},
	})
	if _, ok := r.Get("no.existe"); ok {
		t.Error("ID desconocido no debería existir")
	}
}

func TestMatches(t *testing.T) {
	r := keys.New(nil)
	if !r.Matches(keys.ActionExecute, "enter") {
		t.Error("enter debería matchear action.execute")
	}
	if r.Matches(keys.ActionExecute, "e") {
		t.Error("e NO debería matchear action.execute (default)")
	}
}

func TestGrouped_PreservesOrder(t *testing.T) {
	r := keys.New(nil)
	groups := r.Grouped(
		keys.ActionExecute,
		keys.NavUp,
		keys.ViewDetail,
	)
	if len(groups) != 3 {
		t.Fatalf("esperaba 3 grupos, dio %d", len(groups))
	}
	// Orden: Navigation, Action, View
	if groups[0].Group != keys.GroupNavigation {
		t.Errorf("primer grupo = %q", groups[0].Group)
	}
	if groups[1].Group != keys.GroupAction {
		t.Errorf("segundo grupo = %q", groups[1].Group)
	}
	if groups[2].Group != keys.GroupView {
		t.Errorf("tercer grupo = %q", groups[2].Group)
	}
}

func TestGroup_Title(t *testing.T) {
	if keys.GroupNavigation.Title() != "Navegación" {
		t.Errorf("Title = %q", keys.GroupNavigation.Title())
	}
	if keys.Group("desconocido").Title() != "desconocido" {
		t.Error("grupo desconocido debería devolver su ID")
	}
}
