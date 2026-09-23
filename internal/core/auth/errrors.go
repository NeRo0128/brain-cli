package auth

import "errors"

var (
	ErrNotAuthenticated = errors.New("no hay sesión de GitHub")
	ErrExpired          = errors.New("token de GitHub expirado")
	ErrInvalidToken     = errors.New("token de GitHub inválido")
	ErrUserDenied       = errors.New("el usuario canceló la autorización")
	ErrCodeExpired      = errors.New("código de dispositivo expirado")
)
