package auth

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/rs/zerolog"

	coreauth "github.com/NeRo0128/brain-cli/internal/core/auth"
)

var ErrNoPendingLogin = errors.New("no hay login en curso")

// OAuthClient define el contrato para autenticar contra GitHub usando Device Flow de 2 pasos.
type OAuthClient interface {
	// RequestCode solicita un código de dispositivo al proveedor.
	RequestCode(ctx context.Context) (*coreauth.AuthCode, error)
	// WaitForToken espera a que el usuario autorice y devuelve la sesión.
	WaitForToken(ctx context.Context, code *coreauth.AuthCode) (*coreauth.Session, error)
}

// Manager orquesta la autenticación.
type Manager struct {
	repo    coreauth.TokenRepository
	oauth   OAuthClient
	log     zerolog.Logger
	mu      sync.Mutex
	pending *coreauth.AuthCode
}

func NewManager(repo coreauth.TokenRepository, oauth OAuthClient, log zerolog.Logger) *Manager {
	return &Manager{
		repo:  repo,
		oauth: oauth,
		log:   log.With().Str("component", "auth").Logger(),
	}
}

// StartLogin inicia el flujo OAuth solicitando un código de dispositivo.
// Guarda el código internamente para CompleteLogin.
func (m *Manager) StartLogin(ctx context.Context) (*coreauth.AuthCode, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.pending != nil {
		return m.pending, nil
	}

	code, err := m.oauth.RequestCode(ctx)
	if err != nil {
		m.log.Warn().Err(err).Msg("solicitando código de dispositivo")
		return nil, err
	}

	m.pending = code
	m.log.Info().Str("user_code", code.UserCode).Msg("código de dispositivo solicitado")
	return code, nil
}

// CompleteLogin espera a que el usuario autorice y persiste la sesión.
// Debe llamarse después de StartLogin.
func (m *Manager) CompleteLogin(ctx context.Context) (*coreauth.Session, error) {
	m.mu.Lock()
	code := m.pending
	m.mu.Unlock()

	if code == nil {
		return nil, ErrNoPendingLogin
	}

	session, err := m.oauth.WaitForToken(ctx, code)
	if err != nil {
		m.log.Warn().Err(err).Msg("esperando token")
		return nil, err
	}

	if err := m.repo.Save(ctx, session); err != nil {
		return nil, fmt.Errorf("guardando sesión: %w", err)
	}

	m.mu.Lock()
	m.pending = nil
	m.mu.Unlock()

	m.log.Info().Str("user", session.User.Login).Msg("login exitoso")
	return session, nil
}

// CancelLogin aborta un login en curso y limpia el estado pendiente.
func (m *Manager) CancelLogin() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.pending = nil
	m.log.Info().Msg("login cancelado")
}

// Logout elimina la sesión guardada.
func (m *Manager) Logout(ctx context.Context) error {
	if err := m.repo.Delete(ctx); err != nil {
		return fmt.Errorf("cerrando sesión: %w", err)
	}
	m.log.Info().Msg("logout exitoso")
	return nil
}

// IsAuthenticated indica si hay una sesión válida guardada.
func (m *Manager) IsAuthenticated(ctx context.Context) bool {
	return m.repo.Exists(ctx)
}

// CurrentUser devuelve el usuario autenticado.
// Devuelve ErrNotAuthenticated si no hay sesión.
func (m *Manager) CurrentUser(ctx context.Context) (*coreauth.User, error) {
	session, err := m.repo.Load(ctx)
	if err != nil {
		return nil, err
	}
	return &session.User, nil
}

// CurrentSession devuelve la sesión completa.
func (m *Manager) CurrentSession(ctx context.Context) (*coreauth.Session, error) {
	return m.repo.Load(ctx)
}
