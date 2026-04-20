package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// CacheService provides Redis-backed caching with TTL support.
type CacheService struct {
	rdb redis.Cmdable
}

func NewCacheService(rdb redis.Cmdable) *CacheService {
	return &CacheService{rdb: rdb}
}

// Get retrieves a cached value. Returns nil if not found or expired.
func (s *CacheService) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := s.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil // cache miss
	}
	if err != nil {
		return fmt.Errorf("cache get error: %w", err)
	}
	if err := json.Unmarshal([]byte(val), dest); err != nil {
		return fmt.Errorf("cache unmarshal error: %w", err)
	}
	return nil
}

// Set stores a value in cache with the given TTL.
func (s *CacheService) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("cache marshal error: %w", err)
	}
	return s.rdb.Set(ctx, key, data, ttl).Err()
}

// Delete removes a key from cache.
func (s *CacheService) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	return s.rdb.Del(ctx, keys...).Err()
}

// InvalidateByPattern removes all keys matching a pattern.
func (s *CacheService) InvalidateByPattern(ctx context.Context, pattern string) error {
	iter := s.rdb.Scan(ctx, 0, pattern, 100).Iterator()
	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if len(keys) > 0 {
		return s.rdb.Del(ctx, keys...).Err()
	}
	return nil
}

// CacheKey builds a namespaced cache key.
func CacheKey(parts ...string) string {
	result := "gofin:"
	for _, p := range parts {
		result += p + ":"
	}
	return result[:len(result)-1]
}
