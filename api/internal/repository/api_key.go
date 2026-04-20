package repository

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ajianaz/gofin-full/api/internal/domain"
)

// APIKeyRepository handles API key data access.
type APIKeyRepository struct {
	db *pgxpool.Pool
}

// NewAPIKeyRepository creates a new API key repository.
func NewAPIKeyRepository(db *pgxpool.Pool) *APIKeyRepository {
	return &APIKeyRepository{db: db}
}

// Create generates a new API key, stores its hash, and returns the model + raw key.
// The raw key is only returned once — it cannot be retrieved later.
func (r *APIKeyRepository) Create(ctx context.Context, userID int64, name string) (*domain.APIKey, string, error) {
	// Generate raw key: gofin_<32 random bytes hex>
	rawBytes := make([]byte, 32)
	if _, err := rand.Read(rawBytes); err != nil {
		return nil, "", fmt.Errorf("failed to generate key: %w", err)
	}
	rawKey := "gofin_" + hex.EncodeToString(rawBytes)

	// Hash for storage
	hash := sha256.Sum256([]byte(rawKey))
	keyHash := hex.EncodeToString(hash[:])

	// Prefix for display: "gofin_" + first 4 hex chars = 10 chars total
	prefix := rawKey[:10]

	now := time.Now().UTC()
	var ak domain.APIKey
	err := r.db.QueryRow(ctx,
		`INSERT INTO api_keys (user_id, name, key_hash, key_prefix, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, user_id, name, key_hash, key_prefix, created_at, updated_at`,
		userID, name, keyHash, prefix, now, now,
	).Scan(&ak.ID, &ak.UserID, &ak.Name, &ak.KeyHash, &ak.KeyPrefix, &ak.CreatedAt, &ak.UpdatedAt)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create api key: %w", err)
	}

	return &ak, rawKey, nil
}

// FindByHash looks up an API key by its SHA-256 hash.
// Returns the userID and keyID, or an error if not found or soft-deleted.
func (r *APIKeyRepository) FindByHash(ctx context.Context, keyHash string) (int64, int64, error) {
	var userID, keyID int64
	err := r.db.QueryRow(ctx,
		`SELECT user_id, id FROM api_keys WHERE key_hash = $1 AND deleted_at IS NULL`,
		keyHash,
	).Scan(&userID, &keyID)
	if err != nil {
		return 0, 0, fmt.Errorf("api key not found")
	}
	return userID, keyID, nil
}

// ListByUser returns all non-deleted API keys for a user.
func (r *APIKeyRepository) ListByUser(ctx context.Context, userID int64) ([]domain.APIKey, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, name, key_prefix, last_used_at, created_at, updated_at
		 FROM api_keys WHERE user_id = $1 AND deleted_at IS NULL
		 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list api keys: %w", err)
	}
	defer rows.Close()

	var keys []domain.APIKey
	for rows.Next() {
		var ak domain.APIKey
		if err := rows.Scan(&ak.ID, &ak.UserID, &ak.Name, &ak.KeyPrefix, &ak.LastUsedAt, &ak.CreatedAt, &ak.UpdatedAt); err != nil {
			return nil, err
		}
		keys = append(keys, ak)
	}
	return keys, rows.Err()
}

// Delete soft-deletes an API key.
func (r *APIKeyRepository) Delete(ctx context.Context, id, userID int64) error {
	_, err := r.db.Exec(ctx,
		`UPDATE api_keys SET deleted_at = $1, updated_at = $2
		 WHERE id = $3 AND user_id = $4 AND deleted_at IS NULL`,
		time.Now().UTC(), time.Now().UTC(), id, userID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete api key: %w", err)
	}
	return nil
}

// UpdateLastUsed updates the last_used_at timestamp.
func (r *APIKeyRepository) UpdateLastUsed(ctx context.Context, keyID int64) error {
	_, err := r.db.Exec(ctx,
		`UPDATE api_keys SET last_used_at = $1 WHERE id = $2 AND deleted_at IS NULL`,
		time.Now().UTC(), keyID,
	)
	return err
}
