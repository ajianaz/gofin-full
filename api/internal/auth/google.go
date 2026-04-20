package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/azfirazka/gofin-full/api/internal/config"
)

// googleProvider implements Google OAuth2 authentication.
type googleProvider struct {
	oauthConfig *oauth2.Config
}

func newGoogleProvider(cfg *config.Config) *googleProvider {
	return &googleProvider{
		oauthConfig: &oauth2.Config{
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleClientSecret,
			RedirectURL:  cfg.AppURL + "/api/v1/auth/callback/google",
			Scopes:       []string{"openid", "email", "profile"},
			Endpoint:     google.Endpoint,
		},
	}
}

func (p *googleProvider) Name() string { return "google" }

func (p *googleProvider) SetDB(_ *pgxpool.Pool) {}

func (p *googleProvider) Authenticate(ctx context.Context, creds Credentials) (*UserIdentity, error) {
	if creds.Code == "" {
		return nil, fmt.Errorf("authorization code required for Google auth")
	}

	token, err := p.oauthConfig.Exchange(ctx, creds.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	client := p.oauthConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("google userinfo returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read userinfo response: %w", err)
	}

	var userInfo struct {
		Email string `json:"email"`
	}
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse userinfo: %w", err)
	}

	if userInfo.Email == "" {
		return nil, fmt.Errorf("google account has no email")
	}

	return &UserIdentity{
		Email:   userInfo.Email,
		Blocked: false,
	}, nil
}

// AuthURL returns the Google OAuth2 authorization URL.
func (p *googleProvider) AuthURL(state string) string {
	return p.oauthConfig.AuthCodeURL(state)
}
