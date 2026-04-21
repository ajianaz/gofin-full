package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// RefreshTokenRepository handles refresh token persistence.
type RefreshTokenRepository struct {
	db *pgxpool.Pool
}

func NewRefreshTokenRepository(db *pgxpool.Pool) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

// Store persists a hashed refresh token.
func (r *RefreshTokenRepository) Store(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO refresh_tokens (user_id, token_hash, expires_at, created_at) VALUES ($1, $2, $3, $4)`,
		userID, tokenHash, expiresAt, time.Now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("failed to store refresh token: %w", err)
	}
	return nil
}

// GetByHash returns the token row matching the hash.
func (r *RefreshTokenRepository) GetByHash(ctx context.Context, tokenHash string) (userID uuid.UUID, expiresAt time.Time, err error) {
	err = r.db.QueryRow(ctx,
		`SELECT user_id, expires_at FROM refresh_tokens WHERE token_hash = $1`,
		tokenHash,
	).Scan(&userID, &expiresAt)
	if err != nil {
		return uuid.Nil, time.Time{}, fmt.Errorf("refresh token not found: %w", err)
	}
	return userID, expiresAt, nil
}

// RevokeByHash deletes a single refresh token by hash.
func (r *RefreshTokenRepository) RevokeByHash(ctx context.Context, tokenHash string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM refresh_tokens WHERE token_hash = $1`, tokenHash)
	return err
}

// RevokeByUserID deletes all refresh tokens for a user.
func (r *RefreshTokenRepository) RevokeByUserID(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM refresh_tokens WHERE user_id = $1`, userID)
	return err
}

// RevokeExpired removes all expired refresh tokens.
func (r *RefreshTokenRepository) RevokeExpired(ctx context.Context) error {
	_, err := r.db.Exec(ctx, `DELETE FROM refresh_tokens WHERE expires_at < NOW()`)
	return err
}

// CountByUserID returns the number of active refresh tokens for a user.
func (r *RefreshTokenRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int, error) {
	var count int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM refresh_tokens WHERE user_id = $1 AND expires_at > NOW()`,
		userID,
	).Scan(&count)
	return count, err
}
