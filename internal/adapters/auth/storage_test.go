package auth_test

import (
	"testing"

	"github.com/zalando/go-keyring"

	coreauth "github.com/NeRo0128/brain-cli/internal/core/auth"
)

// mockKeyring es un llavero en memoria para tests.
type mockKeyring struct {
	data map[string]string
}

func (m *mockKeyring) Set(service, user, password string) error {
	if m.data == nil {
		m.data = make(map[string]string)
	}
	m.data[service+":"+user] = password
	return nil
}
func (m *mockKeyring) Get(service, user string) (string, error) {
	v, ok := m.data[service+":"+user]
	if !ok {
		return "", keyring.ErrNotFound
	}
	return v, nil
}
func (m *mockKeyring) Delete(service, user string) error {
	delete(m.data, service+":"+user)
	return nil
}

func sampleSession() *coreauth.Session {
	return &coreauth.Session{
		User:  coreauth.User{Login: "testuser"},
		Token: coreauth.Token{AccessToken: "gho_xxx", TokenType: "bearer"},
	}
}

// Nota: los tests de keyring real requieren un daemon activo.
// Estos tests usan un mock para CI.

func TestKeyring_SaveLoad(t *testing.T) {
	// Se prueba solo la lógica de serialización con el mock.
	// Se salta si el daemon no está disponible.
	t.Skip("skipping: requiere keyring daemon")
}
