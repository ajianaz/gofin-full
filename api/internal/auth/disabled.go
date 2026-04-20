package auth

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// disabledProvider skips authentication entirely.
// All requests are treated as a single superuser (ID=1).
// Useful for single-user self-hosted deployments.
type disabledProvider struct{}

// NewDisabledProvider creates a disabled auth provider.
func NewDisabledProvider() *disabledProvider {
	return &disabledProvider{}
}

func (p *disabledProvider) Name() string { return "disabled" }

func (p *disabledProvider) AuthURL(_ string) string { return "" }

func (p *disabledProvider) SetDB(_ *pgxpool.Pool) {}

func (p *disabledProvider) Authenticate(_ context.Context, _ Credentials) (*UserIdentity, error) {
	return &UserIdentity{
		ID:       1,
		Email:    "admin@local",
		Blocked:  false,
		DemoUser: false,
	}, nil
}
