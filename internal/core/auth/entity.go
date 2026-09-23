package auth

import "time"

// User representa al usuario autenticado.
type User struct {
	Login     string
	Name      string
	Email     string
	AvatarURL string
}

// AuthCode representa el código de dispositivo para Device Flow.
type AuthCode struct {
	UserCode        string
	DeviceCode      string
	VerificationURL string
	ExpiresAt       time.Time
	Interval        time.Duration
}

// Token es una credencial emitida por el proveedor.
// AccessToken es el único campo sensible.
type Token struct {
	AccessToken string
	TokenType   string
	Scopes      []string
	ExpiresAt   time.Time
}

// IsExpired indica si el token ya expiró. Zero = sin expiración.
func (t Token) IsExpired() bool {
	if t.ExpiresAt.IsZero() {
		return false
	}
	return time.Now().After(t.ExpiresAt)
}

// Session agrupa user + token. Su existencia indica autenticación.
type Session struct {
	User      User
	Token     Token
	CreatedAt time.Time
}
