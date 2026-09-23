package settings_test

import (
	"context"
	"errors"
	"testing"

	"github.com/rs/zerolog"

	core "github.com/NeRo0128/brain-cli/internal/core/config"
	coresettings "github.com/NeRo0128/brain-cli/internal/core/settings"
	uc "github.com/NeRo0128/brain-cli/internal/usecases/settings"
)

// --- mock repo ---

type mockRepo struct {
	data map[string]string
	err  error
}

func newMockRepo() *mockRepo { return &mockRepo{data: map[string]string{}} }

func (m *mockRepo) Get(_ context.Context, key string) (*coresettings.Setting, error) {
	v, ok := m.data[key]
	if !ok {
		return nil, coresettings.ErrNotFound
	}
	return &coresettings.Setting{Key: key, Value: v}, nil
}
func (m *mockRepo) Set(_ context.Context, key, value string) error {
	if m.err != nil {
		return m.err
	}
	m.data[key] = value
	return nil
}
func (m *mockRepo) SetMany(_ context.Context, values map[string]string) error {
	if m.err != nil {
		return m.err
	}
	for k, v := range values {
		m.data[k] = v
	}
	return nil
}
func (m *mockRepo) Delete(_ context.Context, key string) error {
	delete(m.data, key)
	return nil
}
func (m *mockRepo) Reset(_ context.Context) error {
	m.data = map[string]string{}
	return nil
}
func (m *mockRepo) List(_ context.Context) ([]coresettings.Setting, error) {
	out := make([]coresettings.Setting, 0, len(m.data))
	for k, v := range m.data {
		out = append(out, coresettings.Setting{Key: k, Value: v})
	}
	return out, nil
}
func (m *mockRepo) ListByCategory(_ context.Context, _ string) ([]coresettings.Setting, error) {
	return nil, nil
}

// --- helpers ---

func testSchema() *coresettings.Schema {
	return coresettings.NewSchema([]coresettings.Field{
		{Key: "ui.theme", Category: "ui", Values: []string{"brain", "nord", "catppuccin"}},
		{Key: "ui.icons", Category: "ui", Values: []string{"unicode", "nerd-b"}},
		{Key: "logging.level", Category: "logging", Values: []string{"debug", "info", "warn", "error"}},
	})
}

func newManager(repo coresettings.Repository) *uc.Manager {
	base := core.Default()
	base.App.Name = "test-app"
	return uc.NewManager(repo, testSchema(), base, zerolog.Nop())
}

// --- tests ---

func TestResolve_NoOverrides_ReturnsBase(t *testing.T) {
	mgr := newManager(newMockRepo())

	cfg, err := mgr.Resolve(context.Background())
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if cfg.App.Name != "test-app" {
		t.Errorf("App.Name = %q", cfg.App.Name)
	}
}

func TestResolve_AppliesOverrides(t *testing.T) {
	repo := newMockRepo()
	repo.data["ui.theme"] = "nord"
	repo.data["logging.level"] = "warn"

	mgr := newManager(repo)

	cfg, err := mgr.Resolve(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if cfg.UI.Theme != "nord" {
		t.Errorf("UI.Theme = %q, quiero 'nord'", cfg.UI.Theme)
	}
	if cfg.Logging.Level != "warn" {
		t.Errorf("Logging.Level = %q, quiero 'warn'", cfg.Logging.Level)
	}
}

func TestResolve_IgnoresUnknownKeys(t *testing.T) {
	repo := newMockRepo()
	repo.data["desconocida.key"] = "valor"

	mgr := newManager(repo)

	// No debe fallar; solo ignora la key.
	cfg, err := mgr.Resolve(context.Background())
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	// El resto de la config sigue igual.
	if cfg.App.Name != "test-app" {
		t.Errorf("config alterada inesperadamente")
	}
}

func TestSet_ValidValue_Persists(t *testing.T) {
	repo := newMockRepo()
	mgr := newManager(repo)

	if err := mgr.Set(context.Background(), "ui.theme", "nord"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if repo.data["ui.theme"] != "nord" {
		t.Errorf("no persistido: %v", repo.data)
	}
}

func TestSet_UnknownKey_Fails(t *testing.T) {
	mgr := newManager(newMockRepo())

	err := mgr.Set(context.Background(), "no.existe", "x")
	if err == nil {
		t.Fatal("esperaba error por key desconocida")
	}
}

func TestSet_InvalidEnumValue_Fails(t *testing.T) {
	mgr := newManager(newMockRepo())

	err := mgr.Set(context.Background(), "ui.theme", "inventado")
	if err == nil {
		t.Fatal("esperaba error por valor inválido")
	}
}

func TestSetMany_Atomic(t *testing.T) {
	repo := newMockRepo()
	mgr := newManager(repo)

	values := map[string]string{
		"ui.theme":      "nord",
		"ui.icons":      "nerd-b",
		"logging.level": "debug",
	}
	if err := mgr.SetMany(context.Background(), values); err != nil {
		t.Fatalf("SetMany: %v", err)
	}
	if len(repo.data) != 3 {
		t.Errorf("persistidos %d, quiero 3", len(repo.data))
	}
}

func TestSetMany_OneInvalid_Rollback(t *testing.T) {
	repo := newMockRepo()
	mgr := newManager(repo)

	values := map[string]string{
		"ui.theme": "nord",
		"ui.icons": "inventado", // inválido
	}
	err := mgr.SetMany(context.Background(), values)
	if err == nil {
		t.Fatal("esperaba error por valor inválido")
	}
	// Ningún valor debe haber sido persistido.
	if len(repo.data) != 0 {
		t.Errorf("datos persistidos a pesar del error: %v", repo.data)
	}
}

func TestReset_ClearsAll(t *testing.T) {
	repo := newMockRepo()
	repo.data["ui.theme"] = "nord"
	repo.data["logging.level"] = "warn"

	mgr := newManager(repo)

	if err := mgr.Reset(context.Background()); err != nil {
		t.Fatalf("Reset: %v", err)
	}
	if len(repo.data) != 0 {
		t.Errorf("quedan %d overrides", len(repo.data))
	}
}

func TestResetKey_RemovesOne(t *testing.T) {
	repo := newMockRepo()
	repo.data["ui.theme"] = "nord"
	repo.data["logging.level"] = "warn"

	mgr := newManager(repo)

	if err := mgr.ResetKey(context.Background(), "ui.theme"); err != nil {
		t.Fatal(err)
	}
	if _, ok := repo.data["ui.theme"]; ok {
		t.Error("key no eliminada")
	}
	if _, ok := repo.data["logging.level"]; !ok {
		t.Error("otra key fue eliminada")
	}
}

func TestOverrides_ReturnsMap(t *testing.T) {
	repo := newMockRepo()
	repo.data["ui.theme"] = "nord"
	repo.data["ui.icons"] = "nerd-b"

	mgr := newManager(repo)

	ov, err := mgr.Overrides(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(ov) != 2 || ov["ui.theme"] != "nord" {
		t.Errorf("mapa incorrecto: %v", ov)
	}
}

func TestResolve_BaseNotMutated(t *testing.T) {
	repo := newMockRepo()
	repo.data["ui.theme"] = "nord"

	base := core.Default()
	base.UI.Theme = "brain" // original

	mgr := uc.NewManager(repo, testSchema(), base, zerolog.Nop())

	_, err := mgr.Resolve(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	// El base original no debe haber cambiado.
	if base.UI.Theme != "brain" {
		t.Errorf("base fue mutado: %q", base.UI.Theme)
	}
}

// Referencia: silencia el import de errors si no se usa.
var _ = errors.Is
