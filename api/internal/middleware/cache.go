package middleware

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

// CacheControl adds HTTP cache headers to responses.
func CacheControl(maxAgeSeconds int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()
		c.Set("Cache-Control", "public, max-age="+strconv.Itoa(maxAgeSeconds))
		return err
	}
}

// RateLimit provides a Redis-based sliding window rate limiter.
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
