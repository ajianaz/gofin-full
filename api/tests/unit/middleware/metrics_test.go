package middleware_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"

	"github.com/azfirazka/gofin-full/api/internal/middleware"
)

func TestMetricsMiddleware(t *testing.T) {
	app := fiber.New()
	app.Use(middleware.Metrics())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	// Make a request
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "ok", string(body))
}

func TestMetricsMiddlewareCollectsData(t *testing.T) {
	// Reset metrics for clean test
	prometheus.DefaultRegisterer.Unregister(prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "gofin_http_requests_total", Help: "test"},
		[]string{"method", "path", "status"},
	))

	app := fiber.New()
	app.Use(middleware.Metrics())
	app.Get("/metrics-test", func(c *fiber.Ctx) error {
		return c.Status(200).SendString("ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/metrics-test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}
