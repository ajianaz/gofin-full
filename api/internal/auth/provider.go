package auth

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ajianaz/gofin-full/api/internal/config"
)

// UserIdentity represents an authenticated user, extracted from any auth provider.
type UserIdentity struct {
	ID          int64
	Email       string
	Blocked     bool
	DemoUser    bool
	UserGroupID *int64
}

// AuthProvider is the strategy interface for authentication.
// Each provider (local, google, github, keycloak, disabled) implements this.
type AuthProvider interface {
	// Name returns the provider identifier (e.g. "local", "google").
	Name() string

	// Authenticate validates credentials and returns the user identity.
	// For local: validates email + password.
	// For OAuth: exchanges authorization code for user info.
	Authenticate(ctx context.Context, creds Credentials) (*UserIdentity, error)

	// AuthURL returns the OAuth authorization URL for redirect-based flows.
	// Returns empty string for providers that don't support OAuth (e.g. local, disabled).
	AuthURL(state string) string

	// SetDB injects the database pool for providers that need it.
	SetDB(db *pgxpool.Pool)
}

// Credentials holds login input. Only relevant fields are populated per provider.
type Credentials struct {
	Email    string
	Password string
	Code     string // OAuth authorization code
}

// NewProvider creates the auth provider based on config.
func NewProvider(cfg *config.Config) AuthProvider {
	switch cfg.AuthProvider {
	case "google":
		return newGoogleProvider(cfg)
	case "github":
		return newGitHubProvider(cfg)
	case "keycloak":
		return newKeycloakProvider(cfg)
	case "disabled":
		return &disabledProvider{}
	default:
		return newLocalProvider(cfg)
	}
}
