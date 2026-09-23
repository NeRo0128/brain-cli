package auth

import "context"

// TokenRepository define la persistencia de sesiones.
// Las implementaciones deben cifrar el AccessToken antes de guardarlo.
type TokenRepository interface {
	Save(ctx context.Context, session *Session) error
	Load(ctx context.Context) (*Session, error)
	Delete(ctx context.Context) error
	Exists(ctx context.Context) bool
}
