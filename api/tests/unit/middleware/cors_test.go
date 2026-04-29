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
	app := setupTestApp(middleware.CORS("http://localhost:5173", "production", ""))
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

func TestCORS_LocalhostFallback(t *testing.T) {
	// Non-production environment falls back to localhost origins (not wildcard)
	app := setupTestApp(middleware.CORS("", "local", ""))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	t.Run("localhost_5173_allowed", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Origin", "http://localhost:5173")
		resp, err := app.Test(req, -1)
		require.NoError(t, err)
		assert.Equal(t, "http://localhost:5173", resp.Header.Get("Access-Control-Allow-Origin"))
	})

	t.Run("random_origin_blocked", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Origin", "https://any-random-site.com")
		resp, err := app.Test(req, -1)
		require.NoError(t, err)
		// Non-production should NOT allow random origins anymore
		assert.Empty(t, resp.Header.Get("Access-Control-Allow-Origin"))
	})
}

func TestCORS_EnvVarOverride(t *testing.T) {
	// CORS_ALLOWED_ORIGINS takes priority over everything
	app := setupTestApp(middleware.CORS("http://localhost:8080", "production", "https://myapp.com,https://admin.myapp.com"))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	t.Run("allowed_origin_from_env_var", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Origin", "https://myapp.com")
		resp, err := app.Test(req, -1)
		require.NoError(t, err)
		assert.Equal(t, "https://myapp.com", resp.Header.Get("Access-Control-Allow-Origin"))
	})

	t.Run("origin_not_in_env_var_blocked", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Origin", "https://evil.com")
		resp, err := app.Test(req, -1)
		require.NoError(t, err)
		assert.Empty(t, resp.Header.Get("Access-Control-Allow-Origin"))
	})
}

func TestCORS_PreflightWithCustomOrigins(t *testing.T) {
	// In production with a specific appURL, preflight requests should also be restricted
	app := setupTestApp(middleware.CORS("http://localhost:5173", "production", ""))
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
