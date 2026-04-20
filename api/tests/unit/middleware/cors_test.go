package middleware_test

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ajianaz/gofin-full/api/internal/middleware"
)

func TestCORS_CustomOrigins(t *testing.T) {
	// In production with a specific appURL, CORS restricts to that origin
	app := setupTestApp(middleware.CORS("http://localhost:5173", "production"))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	// The configured app URL origin should get the origin back
	t.Run("allowed_origin_matches_app_url", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Origin", "http://localhost:5173")
		resp, err := app.Test(req, -1)
		require.NoError(t, err)
		assert.Equal(t, "http://localhost:5173", resp.Header.Get("Access-Control-Allow-Origin"))
	})

	t.Run("different_origin_blocked_in_production", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Origin", "https://evil.com")
		resp, err := app.Test(req, -1)
		require.NoError(t, err)
		// In production with a specific appURL, non-matching origins should not get CORS headers
		assert.Empty(t, resp.Header.Get("Access-Control-Allow-Origin"))
	})
}

func TestCORS_WildcardFallback(t *testing.T) {
	// Non-production environment falls back to wildcard
	app := setupTestApp(middleware.CORS("", "local"))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Origin", "https://any-random-site.com")
	resp, err := app.Test(req, -1)
	require.NoError(t, err)

	// Non-production should allow all origins
	assert.Equal(t, "*", resp.Header.Get("Access-Control-Allow-Origin"))
}

func TestCORS_PreflightWithCustomOrigins(t *testing.T) {
	// In production with a specific appURL, preflight requests should also be restricted
	app := setupTestApp(middleware.CORS("http://localhost:5173", "production"))
	app.Options("/", func(c *fiber.Ctx) error {
		return c.SendStatus(204)
	})

	req := httptest.NewRequest("OPTIONS", "/", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	req.Header.Set("Access-Control-Request-Method", "POST")
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	assert.Equal(t, 204, resp.StatusCode)
	assert.Equal(t, "http://localhost:5173", resp.Header.Get("Access-Control-Allow-Origin"))
}
