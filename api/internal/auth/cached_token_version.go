package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// CachedTokenVersionLookup wraps a TokenVersionLookup with Redis caching.
// Falls back to the underlying DB lookup when Redis is unavailable.
type CachedTokenVersionLookup struct {
	db  TokenVersionLookup
	rdb redis.Cmdable
	ttl time.Duration
}

// NewCachedTokenVersionLookup creates a caching layer for token version checks.
// TTL defaults to 5 minutes if zero.
func NewCachedTokenVersionLookup(db TokenVersionLookup, rdb redis.Cmdable, ttl time.Duration) *CachedTokenVersionLookup {
	if ttl == 0 {
		ttl = 5 * time.Minute
	}
	return &CachedTokenVersionLookup{db: db, rdb: rdb, ttl: ttl}
}

func (c *CachedTokenVersionLookup) redisKey(userID uuid.UUID) string {
	return fmt.Sprintf("token_version:%s", userID)
}

// GetTokenVersion checks Redis first, falls back to DB, caches the result.
func (c *CachedTokenVersionLookup) GetTokenVersion(ctx context.Context, userID uuid.UUID) (int, error) {
	key := c.redisKey(userID)

	// Try Redis — distinguish cache miss (nil) from connection errors
	val, err := c.rdb.Get(ctx, key).Int()
	if err == nil {
		return val, nil // cache hit
	}
	if err != redis.Nil {
		// Redis error (connection issue, etc.) — fall back to DB silently
		// Log handled by caller if needed
	}

	// Cache miss or Redis unavailable — query DB
	version, dbErr := c.db.GetTokenVersion(ctx, userID)
	if dbErr != nil {
		return 0, dbErr
	}

	// Cache the result (fire-and-forget, ignore errors)
	_ = c.rdb.Set(ctx, key, version, c.ttl).Err()

	return version, nil
}

// InvalidateTokenVersion removes the cached token version for a user.
// Must be called on password change, logout, or any token invalidation event
// to ensure the cache stays consistent with the DB.
func (c *CachedTokenVersionLookup) InvalidateTokenVersion(ctx context.Context, userID uuid.UUID) {
	_ = c.rdb.Del(ctx, c.redisKey(userID)).Err()
}
