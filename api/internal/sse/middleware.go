package sse

import (
	"bufio"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/rs/zerolog"
)

// Middleware provides SSE-related Fiber middleware.
type Middleware struct {
	hub *Hub
	log zerolog.Logger
}

// NewMiddleware creates a new SSE middleware.
func NewMiddleware(hub *Hub, log zerolog.Logger) *Middleware {
	return &Middleware{hub: hub, log: log}
}

// Stream upgrades an HTTP connection to SSE.
func (m *Middleware) Stream() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")
		c.Set("X-Accel-Buffering", "no")

		// Get user from auth context
		userID, err := c.ParamsInt("user_id")
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid user id"})
		}

		client := &Client{
			ID:     time.Now().UnixNano(),
			UserID: int64(userID),
			Ch:     make(chan Event, 16),
			Done:   make(chan struct{}),
		}

		m.hub.Subscribe(client)
		defer m.hub.Unsubscribe(client)

		c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
			// Send connected event
			_, _ = w.Write([]byte("event: connected\ndata: {\"status\":\"connected\"}\n\n"))
			_ = w.Flush()

			// Heartbeat ticker (30s to keep connection alive)
			ticker := time.NewTicker(30 * time.Second)
			defer ticker.Stop()

			for {
				select {
				case event := <-client.Ch:
					_, _ = w.Write(MarshalEvent(event))
					_ = w.Flush()
				case <-ticker.C:
					_, _ = w.Write([]byte(": heartbeat\n\n"))
					_ = w.Flush()
				case <-c.Context().Done():
					return
				}
			}
		})

		return nil
	}
}

// RateLimiter returns a rate limiter middleware for SSE connections.
func (m *Middleware) RateLimiter(maxConnections int, window time.Duration) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        maxConnections,
		Expiration: window,
		KeyGenerator: func(c *fiber.Ctx) string {
			return "sse:" + c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(429).JSON(fiber.Map{
				"error": "too many SSE connections",
			})
		},
	})
}
