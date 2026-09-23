package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cli/oauth/device"
	"github.com/google/go-github/v66/github"

	coreauth "github.com/NeRo0128/brain-cli/internal/core/auth"
)

// GitHubOAuth implementa el flujo Device Flow de GitHub.
type GitHubOAuth struct {
	clientID string
	scopes   []string
}

func NewGitHubOAuth(clientID string) *GitHubOAuth {
	return &GitHubOAuth{
		clientID: clientID,
		scopes:   []string{"read:user", "user:email"},
	}
}

// RequestCode solicita un código de dispositivo al proveedor.
func (g *GitHubOAuth) RequestCode(ctx context.Context) (*coreauth.AuthCode, error) {
	code, err := device.RequestCode(http.DefaultClient, "https://github.com/login/device/code", g.clientID, g.scopes)
	if err != nil {
		return nil, fmt.Errorf("solicitando código de dispositivo: %w", err)
	}

	return &coreauth.AuthCode{
		UserCode:        code.UserCode,
		DeviceCode:      code.DeviceCode,
		VerificationURL: code.VerificationURI,
		ExpiresAt:       time.Now().Add(time.Duration(code.ExpiresIn) * time.Second),
		Interval:        time.Duration(code.Interval) * time.Second,
	}, nil
}

// WaitForToken espera a que el usuario autorice y devuelve la sesión.
func (g *GitHubOAuth) WaitForToken(ctx context.Context, code *coreauth.AuthCode) (*coreauth.Session, error) {
	deviceCode := &device.CodeResponse{
		DeviceCode: code.DeviceCode,
		Interval:   int(code.Interval.Seconds()),
		ExpiresIn:  int(time.Until(code.ExpiresAt).Seconds()),
	}

	opts := device.WaitOptions{
		ClientID:   g.clientID,
		DeviceCode: deviceCode,
	}

	accessToken, err := device.Wait(ctx, http.DefaultClient, "https://github.com/login/oauth/access_token", opts)
	if err != nil {
		return nil, fmt.Errorf("esperando token: %w", err)
	}

	user, err := g.fetchUser(ctx, accessToken.Token)
	if err != nil {
		return nil, fmt.Errorf("fetching user: %w", err)
	}

	return &coreauth.Session{
		User:      *user,
		Token:     coreauth.Token{AccessToken: accessToken.Token, TokenType: accessToken.Type, Scopes: g.scopes},
		CreatedAt: time.Now(),
	}, nil
}

func (g *GitHubOAuth) fetchUser(ctx context.Context, token string) (*coreauth.User, error) {
	client := github.NewClient(nil).WithAuthToken(token)

	u, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return nil, err
	}

	return &coreauth.User{
		Login:     u.GetLogin(),
		Name:      u.GetName(),
		Email:     u.GetEmail(),
		AvatarURL: u.GetAvatarURL(),
	}, nil
}