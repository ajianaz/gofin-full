package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	response "github.com/ajianaz/gofin-full/api/internal/dto/response"
)

// HealthHandler handles health check requests.
type HealthHandler struct {
	db    *pgxpool.Pool
	redis redis.Cmdable
}

// NewHealthHandler creates a new health handler.
func NewHealthHandler(db *pgxpool.Pool, rdb redis.Cmdable) *HealthHandler {
	return &HealthHandler{db: db, redis: rdb}
}

// Check returns the health status of all services.
func (h *HealthHandler) Check(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(healthTimeout)*time.Second)
	defer cancel()

	health := response.HealthResponse{
		Status:   "ok",
		Services: []response.ServiceHealth{},
	}

	// Check PostgreSQL
	pgStatus := response.ServiceHealth{Name: "postgresql"}
	if err := h.checkPostgres(ctx); err != nil {
		pgStatus.Status = "error"
		pgStatus.Error = err.Error()
		health.Status = "degraded"
	} else {
		pgStatus.Status = "ok"
	}
	health.Services = append(health.Services, pgStatus)

	// Check Redis
	redisStatus := response.ServiceHealth{Name: "redis"}
	if err := h.checkRedis(ctx); err != nil {
		redisStatus.Status = "error"
		redisStatus.Error = err.Error()
		health.Status = "degraded"
	} else {
		redisStatus.Status = "ok"
	}
	health.Services = append(health.Services, redisStatus)

	if health.Status == "ok" {
		return c.JSON(health)
	}
	return c.Status(503).JSON(health)
}

func (h *HealthHandler) checkPostgres(ctx context.Context) error {
	if h.db == nil {
		return fmt.Errorf("database connection not initialized")
	}
	conn, err := h.db.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()
	return conn.Ping(ctx)
}

func (h *HealthHandler) checkRedis(ctx context.Context) error {
	if h.redis == nil {
		return fmt.Errorf("redis connection not initialized")
	}
	return h.redis.Ping(ctx).Err()
}

const healthTimeout = 10
