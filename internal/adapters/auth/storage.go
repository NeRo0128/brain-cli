package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/zalando/go-keyring"

	coreauth "github.com/NeRo0128/brain-cli/internal/core/auth"
)

const (
	keyringService = "brain-cli"
	keyringUser    = "github-session"
)

// KeyringRepository implementa TokenRepository usando el llavero
// del sistema (Secret Service en Linux, Keychain en macOS).
type KeyringRepository struct{}

func NewKeyringRepository() *KeyringRepository {
	return &KeyringRepository{}
}

// Save serializa la sesión completa a JSON y la guarda cifrada.
func (r *KeyringRepository) Save(_ context.Context, session *coreauth.Session) error {
	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("serializando sesión: %w", err)
	}

	if err := keyring.Set(keyringService, keyringUser, string(data)); err != nil {
		return fmt.Errorf("guardando en keyring: %w", err)
	}
	return nil
}

// Load recupera la sesión del llavero.
func (r *KeyringRepository) Load(_ context.Context) (*coreauth.Session, error) {
	raw, err := keyring.Get(keyringService, keyringUser)
	if errors.Is(err, keyring.ErrNotFound) {
		return nil, coreauth.ErrNotAuthenticated
	}
	if err != nil {
		return nil, fmt.Errorf("leyendo keyring: %w", err)
	}

	var session coreauth.Session
	if err := json.Unmarshal([]byte(raw), &session); err != nil {
		return nil, fmt.Errorf("deserializando sesión: %w", err)
	}

	if session.Token.IsExpired() {
		return nil, coreauth.ErrExpired
	}

	return &session, nil
}

// Delete elimina la sesión del llavero. Idempotente.
func (r *KeyringRepository) Delete(_ context.Context) error {
	err := keyring.Delete(keyringService, keyringUser)
	if errors.Is(err, keyring.ErrNotFound) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("borrando del keyring: %w", err)
	}
	return nil
}

// Exists indica si hay sesión guardada.
func (r *KeyringRepository) Exists(_ context.Context) bool {
	_, err := keyring.Get(keyringService, keyringUser)
	return err == nil
}
