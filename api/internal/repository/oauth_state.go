package repository

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// OAuthStateRepository manages OAuth state tokens for CSRF protection.
type OAuthStateRepository struct {
	db *pgxpool.Pool
}

// NewOAuthStateRepository creates a new OAuth state repository.
func NewOAuthStateRepository(db *pgxpool.Pool) *OAuthStateRepository {
	return &OAuthStateRepository{db: db}
}

// Generate creates a new OAuth state token and stores it.
// Returns the state string or an error.
func (r *OAuthStateRepository) Generate(ctx context.Context, provider, redirect string) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate state: %w", err)
	}
	state := hex.EncodeToString(b)
	expiresAt := time.Now().UTC().Add(10 * time.Minute)

	_, err := r.db.Exec(ctx,
		`INSERT INTO oauth_states (state, provider, redirect, created_at, expires_at)
		 VALUES ($1, $2, $3, NOW(), $4)`,
		state, provider, redirect, expiresAt,
	)
	if err != nil {
		return "", fmt.Errorf("failed to store oauth state: %w", err)
	}

	return state, nil
}

// Validate checks if a state token is valid and not expired.
// Returns the provider and redirect URL, or an error.
// Deletes the state after validation (one-time use).
func (r *OAuthStateRepository) Validate(ctx context.Context, state string) (provider, redirect string, err error) {
	err = r.db.QueryRow(ctx,
		`SELECT provider, redirect FROM oauth_states
		 WHERE state = $1 AND expires_at > NOW()`,
		state,
	).Scan(&provider, &redirect)
	if err != nil {
		return "", "", fmt.Errorf("invalid or expired oauth state")
	}

	// Delete the used state
	_, _ = r.db.Exec(ctx, `DELETE FROM oauth_states WHERE state = $1`, state)

	return provider, redirect, nil
}

// CleanupExpired removes expired OAuth states. Call periodically.
func (r *OAuthStateRepository) CleanupExpired(ctx context.Context) error {
	_, err := r.db.Exec(ctx, `DELETE FROM oauth_states WHERE expires_at < NOW()`)
	return err
}
