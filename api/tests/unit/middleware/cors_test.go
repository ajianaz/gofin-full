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
	origins := []string{"http://localhost:5173", "https://app.gofin.io"}

	app := setupTestApp(middleware.CORS(origins...))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	// Allowed origin should get the origin back
	t.Run("allowed_origin_localhost", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Origin", "http://localhost:5173")
		resp, err := app.Test(req, -1)
		require.NoError(t, err)
		assert.Equal(t, "http://localhost:5173", resp.Header.Get("Access-Control-Allow-Origin"))
	})

	t.Run("allowed_origin_app", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Origin", "https://app.gofin.io")
		resp, err := app.Test(req, -1)
		require.NoError(t, err)
		assert.Equal(t, "https://app.gofin.io", resp.Header.Get("Access-Control-Allow-Origin"))
	})

	t.Run("disallowed_origin_blocked", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Origin", "https://evil.com")
		resp, err := app.Test(req, -1)
		require.NoError(t, err)
		// With specific origins, a non-matching origin should get an empty header
		assert.Empty(t, resp.Header.Get("Access-Control-Allow-Origin"))
	})
}

func TestCORS_WildcardFallback(t *testing.T) {
	app := setupTestApp(middleware.CORS())
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Origin", "https://any-random-site.com")
	resp, err := app.Test(req, -1)
	require.NoError(t, err)

	// With no origins specified, should fall back to wildcard
	assert.Equal(t, "*", resp.Header.Get("Access-Control-Allow-Origin"))
}

func TestCORS_PreflightWithCustomOrigins(t *testing.T) {
	origins := []string{"http://localhost:5173"}
	app := setupTestApp(middleware.CORS(origins...))
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
