package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/oauth2"

	"github.com/ajianaz/gofin-full/api/internal/config"
)

// keycloakProvider implements Keycloak OIDC authentication.
type keycloakProvider struct {
	realmURL     string
	clientID     string
	clientSecret string
	redirectURL  string
}

func newKeycloakProvider(cfg *config.Config) *keycloakProvider {
	realmURL := cfg.KeycloakRealmURL()
	return &keycloakProvider{
		realmURL:     realmURL,
		clientID:     cfg.KeycloakClientID,
		clientSecret: cfg.KeycloakClientSecret,
		redirectURL:  cfg.AppURL + "/api/v1/auth/callback/keycloak",
	}
}

func (p *keycloakProvider) Name() string { return "keycloak" }

func (p *keycloakProvider) SetDB(_ *pgxpool.Pool) {}

func (p *keycloakProvider) Authenticate(ctx context.Context, creds Credentials) (*UserIdentity, error) {
	if creds.Code == "" {
		return nil, fmt.Errorf("authorization code required for Keycloak auth")
	}

	// Exchange authorization code for tokens via Keycloak's token endpoint
	tokenURL := p.realmURL + "/protocol/openid-connect/token"

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", p.clientID)
	data.Set("client_secret", p.clientSecret)
	data.Set("code", creds.Code)
	data.Set("redirect_uri", p.redirectURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code with keycloak: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("keycloak token endpoint returned status %d", resp.StatusCode)
	}

	var tokenResp struct {
		IDToken string `json:"id_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode keycloak token response: %w", err)
	}

	if tokenResp.IDToken == "" {
		return nil, fmt.Errorf("keycloak did not return an id_token")
	}

	// Decode the ID token (JWT) to extract user info
	// We parse the payload section (second part) of the JWT without signature verification.
	// Signature verification should be done with JWKS in production.
	email, err := extractEmailFromIDToken(tokenResp.IDToken)
	if err != nil {
		return nil, err
	}

	if email == "" {
		return nil, fmt.Errorf("keycloak account has no email")
	}

	return &UserIdentity{
		Email:   email,
		Blocked: false,
	}, nil
}

// AuthURL returns the Keycloak authorization URL.
func (p *keycloakProvider) AuthURL(state string) string {
	return fmt.Sprintf("%s/protocol/openid-connect/auth?client_id=%s&response_type=code&redirect_uri=%s&scope=openid+email+profile&state=%s",
		p.realmURL, url.QueryEscape(p.clientID), url.QueryEscape(p.redirectURL), url.QueryEscape(state))
}

// extractEmailFromIDToken parses the JWT payload to extract the email claim.
func extractEmailFromIDToken(idToken string) (string, error) {
	parts := strings.Split(idToken, ".")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid id_token format")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", fmt.Errorf("failed to decode id_token payload: %w", err)
	}

	var claims struct {
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return "", fmt.Errorf("failed to parse id_token claims: %w", err)
	}

	return claims.Email, nil
}

// Ensure oauth2.Config is available for potential future use.
var _ = oauth2.Config{}
