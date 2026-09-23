package settings

import (
	"context"
	"fmt"
	"slices"

	"github.com/rs/zerolog"

	core "github.com/NeRo0128/brain-cli/internal/core/config"
	coresettings "github.com/NeRo0128/brain-cli/internal/core/settings"
)

// Manager fusiona defaults + YAML + overrides de DB y persiste cambios.
type Manager struct {
	repo   coresettings.Repository
	schema *coresettings.Schema
	base   *core.Config // config del YAML, inmutable
	log    zerolog.Logger
}

func NewManager(
	repo coresettings.Repository,
	schema *coresettings.Schema,
	base *core.Config,
	log zerolog.Logger,
) *Manager {
	return &Manager{
		repo:   repo,
		schema: schema,
		base:   base,
		log:    log.With().Str("component", "settings").Logger(),
	}
}

// Resolve devuelve una copia del config base con los overrides aplicados.
// El base nunca se muta.
func (m *Manager) Resolve(ctx context.Context) (*core.Config, error) {
	out := *m.base

	overrides, err := m.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("cargando overrides: %w", err)
	}

	for _, o := range overrides {
		if err := applyKey(&out, o.Key, o.Value); err != nil {
			m.log.Warn().Err(err).Str("key", o.Key).Msg("override ignorado")
		}
	}

	return &out, nil
}

// Set valida y persiste un override.
func (m *Manager) Set(ctx context.Context, key, value string) error {
	return m.SetMany(ctx, map[string]string{key: value})
}

// SetMany valida y persiste varios overrides atómicamente.
func (m *Manager) SetMany(ctx context.Context, values map[string]string) error {
	for key, value := range values {
		if err := m.validate(key, value); err != nil {
			return err
		}
	}
	if err := m.repo.SetMany(ctx, values); err != nil {
		return fmt.Errorf("persistiendo: %w", err)
	}
	m.log.Debug().Int("count", len(values)).Msg("settings persistidos")
	return nil
}

// Reset borra todos los overrides.
func (m *Manager) Reset(ctx context.Context) error {
	if err := m.repo.Reset(ctx); err != nil {
		return fmt.Errorf("reseteando: %w", err)
	}
	m.log.Info().Msg("overrides reseteados")
	return nil
}

// ResetKey borra un override individual.
func (m *Manager) ResetKey(ctx context.Context, key string) error {
	return m.repo.Delete(ctx, key)
}

// Overrides devuelve el mapa actual de overrides (key → value).
func (m *Manager) Overrides(ctx context.Context) (map[string]string, error) {
	list, err := m.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	out := make(map[string]string, len(list))
	for _, s := range list {
		out[s.Key] = s.Value
	}
	return out, nil
}

// validate comprueba que la key existe en el schema y que el valor
// es válido si el campo es un enum.
func (m *Manager) validate(key, value string) error {
	field, ok := m.schema.Get(key)
	if !ok {
		return fmt.Errorf("setting desconocido: %q", key)
	}
	if len(field.Values) > 0 && !slices.Contains(field.Values, value) {
		return fmt.Errorf("valor inválido para %q: %q", key, value)
	}
	return nil
}

// applyKey muta cfg con el override. Las keys no reconocidas devuelven error.
func applyKey(cfg *core.Config, key, value string) error {
	switch key {
	case "ui.theme":
		cfg.UI.Theme = value
	case "ui.icons":
		cfg.UI.Icons = value
	case "ui.brand_style":
		cfg.UI.BrandStyle = value
	case "logging.level":
		cfg.Logging.Level = value
	case "logging.format":
		cfg.Logging.Format = value
	default:
		return fmt.Errorf("key no aplicable: %q", key)
	}
	return nil
}

func (m *Manager) Base() core.Config {
	return *m.base
}
