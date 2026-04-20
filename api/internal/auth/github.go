package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"

	"github.com/azfirazka/gofin-full/api/internal/config"
)

// gitHubProvider implements GitHub OAuth2 authentication.
type gitHubProvider struct {
	oauthConfig *oauth2.Config
}

func newGitHubProvider(cfg *config.Config) *gitHubProvider {
	return &gitHubProvider{
		oauthConfig: &oauth2.Config{
			ClientID:     cfg.GitHubClientID,
			ClientSecret: cfg.GitHubClientSecret,
			RedirectURL:  cfg.AppURL + "/api/v1/auth/callback/github",
			Scopes:       []string{"user:email", "read:user"},
			Endpoint:     github.Endpoint,
		},
	}
}

func (p *gitHubProvider) Name() string { return "github" }

func (p *gitHubProvider) SetDB(_ *pgxpool.Pool) {}

func (p *gitHubProvider) Authenticate(ctx context.Context, creds Credentials) (*UserIdentity, error) {
	if creds.Code == "" {
		return nil, fmt.Errorf("authorization code required for GitHub auth")
	}

	token, err := p.oauthConfig.Exchange(ctx, creds.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	client := p.oauthConfig.Client(ctx, token)

	// Fetch GitHub user profile
	profile, err := fetchGitHubProfile(client)
	if err != nil {
		return nil, err
	}

	// If email is public, we have it already.
	// Otherwise, fetch from the emails API endpoint.
	email := profile.Email
	if email == "" {
		email, err = fetchGitHubPrimaryEmail(client)
		if err != nil {
			return nil, err
		}
	}

	if email == "" {
		return nil, fmt.Errorf("github account has no verified email")
	}

	return &UserIdentity{
		Email:   email,
		Blocked: false,
	}, nil
}

// AuthURL returns the GitHub OAuth2 authorization URL.
func (p *gitHubProvider) AuthURL(state string) string {
	return p.oauthConfig.AuthCodeURL(state)
}

// githubProfile represents the relevant fields from GitHub's /user response.
type githubProfile struct {
	Login string `json:"login"`
	Email string `json:"email"`
	ID    int64  `json:"id"`
}

// githubEmail represents an entry from GitHub's /user/emails response.
type githubEmail struct {
	Email   string `json:"email"`
	Primary bool   `json:"primary"`
	Verified bool  `json:"verified"`
}

func fetchGitHubProfile(client *http.Client) (*githubProfile, error) {
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, fmt.Errorf("failed to get github user profile: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("github user profile returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read github profile response: %w", err)
	}

	var profile githubProfile
	if err := json.Unmarshal(body, &profile); err != nil {
		return nil, fmt.Errorf("failed to parse github profile: %w", err)
	}

	return &profile, nil
}

func fetchGitHubPrimaryEmail(client *http.Client) (string, error) {
	resp, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		return "", fmt.Errorf("failed to get github user emails: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("github user emails returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read github emails response: %w", err)
	}

	var emails []githubEmail
	if err := json.Unmarshal(body, &emails); err != nil {
		return "", fmt.Errorf("failed to parse github emails: %w", err)
	}

	// Find the primary verified email
	for _, e := range emails {
		if e.Primary && e.Verified {
			return e.Email, nil
		}
	}

	// Fallback: any verified email
	for _, e := range emails {
		if e.Verified {
			return e.Email, nil
		}
	}

	return "", fmt.Errorf("no verified email found on github account")
}
