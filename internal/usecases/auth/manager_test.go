package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rs/zerolog"

	coreauth "github.com/NeRo0128/brain-cli/internal/core/auth"
	"github.com/NeRo0128/brain-cli/internal/usecases/auth"
)

// --- Mocks ---

type mockRepo struct {
	session *coreauth.Session
	err     error
}

func (m *mockRepo) Save(_ context.Context, s *coreauth.Session) error {
	if m.err != nil {
		return m.err
	}
	m.session = s
	return nil
}
func (m *mockRepo) Load(_ context.Context) (*coreauth.Session, error) {
	if m.session == nil {
		return nil, coreauth.ErrNotAuthenticated
	}
	return m.session, nil
}
func (m *mockRepo) Delete(_ context.Context) error {
	m.session = nil
	return nil
}
func (m *mockRepo) Exists(_ context.Context) bool { return m.session != nil }

type mockOAuth struct {
	code       *coreauth.AuthCode
	session    *coreauth.Session
	requestErr error
	waitErr    error
}

func (m *mockOAuth) RequestCode(_ context.Context) (*coreauth.AuthCode, error) {
	if m.requestErr != nil {
		return nil, m.requestErr
	}
	if m.code == nil {
		m.code = &coreauth.AuthCode{
			UserCode:        "ABCD-1234",
			DeviceCode:      "device-code-123",
			VerificationURL: "https://github.com/login/device",
			ExpiresAt:       time.Now().Add(10 * time.Minute),
			Interval:        5 * time.Second,
		}
	}
	return m.code, nil
}

func (m *mockOAuth) WaitForToken(_ context.Context, code *coreauth.AuthCode) (*coreauth.Session, error) {
	if m.waitErr != nil {
		return nil, m.waitErr
	}
	if code == nil {
		return nil, errors.New("código nil")
	}
	return m.session, nil
}

// --- Tests ---

func newTestManager(repo coreauth.TokenRepository, oauth auth.OAuthClient) *auth.Manager {
	return auth.NewManager(repo, oauth, zerolog.Nop())
}

func sampleSession() *coreauth.Session {
	return &coreauth.Session{
		User:  coreauth.User{Login: "testuser"},
		Token: coreauth.Token{AccessToken: "gho_xxx", TokenType: "bearer"},
	}
}

func TestStartLogin_Success(t *testing.T) {
	repo := &mockRepo{}
	oauth := &mockOAuth{session: sampleSession()}
	mgr := newTestManager(repo, oauth)

	code, err := mgr.StartLogin(context.Background())
	if err != nil {
		t.Fatalf("StartLogin: %v", err)
	}
	if code.UserCode != "ABCD-1234" {
		t.Errorf("code = %q", code.UserCode)
	}
	if code.VerificationURL == "" {
		t.Error("url vacía")
	}
}

func TestStartLogin_ReturnsSameCode(t *testing.T) {
	repo := &mockRepo{}
	oauth := &mockOAuth{session: sampleSession()}
	mgr := newTestManager(repo, oauth)

	code1, err := mgr.StartLogin(context.Background())
	if err != nil {
		t.Fatalf("StartLogin 1: %v", err)
	}
	code2, err := mgr.StartLogin(context.Background())
	if err != nil {
		t.Fatalf("StartLogin 2: %v", err)
	}
	if code1 != code2 {
		t.Error("debería devolver el mismo código pendiente")
	}
}

func TestCompleteLogin_Success(t *testing.T) {
	repo := &mockRepo{}
	oauth := &mockOAuth{session: sampleSession()}
	mgr := newTestManager(repo, oauth)

	// Primero iniciar login
	_, err := mgr.StartLogin(context.Background())
	if err != nil {
		t.Fatalf("StartLogin: %v", err)
	}

	// Luego completarlo
	session, err := mgr.CompleteLogin(context.Background())
	if err != nil {
		t.Fatalf("CompleteLogin: %v", err)
	}
	if session.User.Login != "testuser" {
		t.Errorf("user = %q", session.User.Login)
	}
	if repo.session == nil {
		t.Error("sesión no persistida")
	}
}

func TestCompleteLogin_NoPending(t *testing.T) {
	repo := &mockRepo{}
	oauth := &mockOAuth{session: sampleSession()}
	mgr := newTestManager(repo, oauth)

	// Sin llamar a StartLogin
	_, err := mgr.CompleteLogin(context.Background())
	if !errors.Is(err, auth.ErrNoPendingLogin) {
		t.Fatalf("esperaba ErrNoPendingLogin, dio: %v", err)
	}
}

func TestCompleteLogin_OAuthFails(t *testing.T) {
	repo := &mockRepo{}
	oauth := &mockOAuth{session: sampleSession(), waitErr: coreauth.ErrUserDenied}
	mgr := newTestManager(repo, oauth)

	_, err := mgr.StartLogin(context.Background())
	if err != nil {
		t.Fatalf("StartLogin: %v", err)
	}

	_, err = mgr.CompleteLogin(context.Background())
	if !errors.Is(err, coreauth.ErrUserDenied) {
		t.Fatalf("esperaba ErrUserDenied, dio: %v", err)
	}
}

func TestCancelLogin(t *testing.T) {
	repo := &mockRepo{}
	oauth := &mockOAuth{session: sampleSession()}
	mgr := newTestManager(repo, oauth)

	_, err := mgr.StartLogin(context.Background())
	if err != nil {
		t.Fatalf("StartLogin: %v", err)
	}

	mgr.CancelLogin()

	// Después de cancelar, CompleteLogin debería fallar
	_, err = mgr.CompleteLogin(context.Background())
	if !errors.Is(err, auth.ErrNoPendingLogin) {
		t.Fatalf("esperaba ErrNoPendingLogin después de cancelar, dio: %v", err)
	}
}

func TestLogout(t *testing.T) {
	repo := &mockRepo{session: sampleSession()}
	mgr := newTestManager(repo, &mockOAuth{})

	if err := mgr.Logout(context.Background()); err != nil {
		t.Fatalf("Logout: %v", err)
	}
	if repo.session != nil {
		t.Error("sesión no eliminada")
	}
}

func TestIsAuthenticated(t *testing.T) {
	repo := &mockRepo{}
	mgr := newTestManager(repo, &mockOAuth{})

	if mgr.IsAuthenticated(context.Background()) {
		t.Error("no debería estar autenticado")
	}

	repo.session = sampleSession()
	if !mgr.IsAuthenticated(context.Background()) {
		t.Error("debería estar autenticado")
	}
}

func TestCurrentUser_NotAuthenticated(t *testing.T) {
	mgr := newTestManager(&mockRepo{}, &mockOAuth{})

	_, err := mgr.CurrentUser(context.Background())
	if !errors.Is(err, coreauth.ErrNotAuthenticated) {
		t.Fatalf("esperaba ErrNotAuthenticated, dio: %v", err)
	}
}

func TestCurrentUser_Success(t *testing.T) {
	repo := &mockRepo{session: sampleSession()}
	mgr := newTestManager(repo, &mockOAuth{})

	user, err := mgr.CurrentUser(context.Background())
	if err != nil {
		t.Fatalf("CurrentUser: %v", err)
	}
	if user.Login != "testuser" {
		t.Errorf("user = %q", user.Login)
	}
}
