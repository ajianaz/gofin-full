package middleware

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

// memRateEntry tracks in-memory rate limit state per key.
type memRateEntry struct {
	timestamps []int64
	mu         sync.Mutex
}

// memRateLimiter is a global in-memory rate limiter used as fallback when Redis is unavailable.
var memRateLimiter sync.Map

// memRateLimitCheck performs an in-memory sliding window rate limit check.
// Returns true if the request should be allowed, false if rate limited.
func memRateLimitCheck(key string, limit int, window time.Duration) bool {
	now := time.Now().UnixMilli()
	windowStart := now - window.Milliseconds()

	val, _ := memRateLimiter.LoadOrStore(key, &memRateEntry{})
	entry := val.(*memRateEntry)

	entry.mu.Lock()
	defer entry.mu.Unlock()

	// Remove expired timestamps
	valid := entry.timestamps[:0]
	for _, ts := range entry.timestamps {
		if ts > windowStart {
			valid = append(valid, ts)
		}
	}
	entry.timestamps = valid

	if len(entry.timestamps) >= limit {
		return false
	}

	entry.timestamps = append(entry.timestamps, now)
	return true
}

// CacheControl adds HTTP cache headers to responses.
func CacheControl(maxAgeSeconds int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()
		c.Set("Cache-Control", "public, max-age="+strconv.Itoa(maxAgeSeconds))
		return err
	}
}

// RateLimit provides a Redis-based sliding window rate limiter.
// Falls back to an in-memory limiter when Redis is unavailable.
func RateLimit(rdb redis.Cmdable, limit int, window time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ip := c.IP()
		path := c.Path()
		key := fmt.Sprintf("ratelimit:%s:%s", ip, path)

		now := time.Now().UnixMilli()
		windowMs := window.Milliseconds()
		windowStart := now - windowMs

		pipe := rdb.Pipeline()
		incr := pipe.ZAdd(c.Context(), key, redis.Z{Score: float64(now), Member: now})
		removeOld := pipe.ZRemRangeByScore(c.Context(), key, "0", strconv.FormatInt(windowStart, 10))
		count := pipe.ZCard(c.Context(), key)
		expire := pipe.Expire(c.Context(), key, window+time.Second)

		_, _ = pipe.Exec(c.Context())

		if incr.Err() != nil || removeOld.Err() != nil || count.Err() != nil || expire.Err() != nil {
			// Redis unavailable — fall back to in-memory rate limiter
			log.Printf("rate limiter: redis error, applying in-memory fallback for %s", key)
			if !memRateLimitCheck(key, limit, window) {
				c.Set("Retry-After", strconv.Itoa(int(window.Seconds())))
				return c.Status(429).JSON(fiber.Map{
					"error": "Too many requests. Please try again later.",
				})
			}
			return c.Next()
		}

		if count.Val() > int64(limit) {
			ttl, _ := rdb.TTL(c.Context(), key).Result()
			retryAfter := int(ttl.Seconds())
			if retryAfter < 1 {
				retryAfter = int(window.Seconds())
			}
			c.Set("Retry-After", strconv.Itoa(retryAfter))
			return c.Status(429).JSON(fiber.Map{
				"error": "Too many requests. Please try again later.",
			})
		}

		return c.Next()
	}
}
