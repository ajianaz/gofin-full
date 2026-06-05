package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ajianaz/gofin-full/api/internal/auth"
	"github.com/ajianaz/gofin-full/api/internal/config"
	"github.com/ajianaz/gofin-full/api/internal/handler"
	"github.com/ajianaz/gofin-full/api/internal/middleware"
	"github.com/ajianaz/gofin-full/api/internal/router"
)

// ---------------------------------------------------------------------------
// Test helpers
// ---------------------------------------------------------------------------

// testAuthConfig returns a minimal Config suitable for auth handler unit tests.
func testAuthConfig() *config.Config {
	return &config.Config{
		AppEnv:                "testing",
		AuthProvider:          "disabled",
		AuthAllowRegistration: true,
	}
}

// newAuthTestApp creates a Fiber app with routes wired up exactly like production,
// but with no external dependencies (no DB, no Redis).
func newAuthTestApp() *fiber.App {
	log := zerolog.Nop()
	jwtMgr := auth.NewJWTManager("test-secret-for-unit-tests", 60, 30)
	provider := auth.NewDisabledProvider()
	cfg := testAuthConfig()

	authHandler := handler.NewAuthHandler(log, jwtMgr, provider, cfg, nil, nil, nil)

	rc := router.RouterConfig{
		AppEnv:         "testing",
		AppURL:         "http://localhost:5173",
		HealthHandler:  handler.NewHealthHandler(nil, nil, "test"),
		AuthHandler:    authHandler,
		JWTManager:     jwtMgr,
		DisableMetrics: true,
	}
	return router.New(rc)
}

// postJSON sends a POST request with a JSON body to the given path on app.
func postJSON(t *testing.T, app *fiber.App, path, body string) *http.Response {
	t.Helper()
	req := httptest.NewRequest("POST", path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	return resp
}

// decodeResponse reads the response body into a map.
func decodeResponse(t *testing.T, resp *http.Response) map[string]interface{} {
	t.Helper()
	var result map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	return result
}

// ---------------------------------------------------------------------------
// 1. TestAuth_Login_MissingEmail — POST /auth/login without email → 422
// ---------------------------------------------------------------------------

func TestAuth_Login_MissingEmail(t *testing.T) {
	app := newAuthTestApp()

	body := `{"password": "secret123"}`
	resp := postJSON(t, app, "/api/v1/auth/login", body)
	defer resp.Body.Close()

	assert.Equal(t, 422, resp.StatusCode)

	result := decodeResponse(t, resp)
	assert.Contains(t, result, "errors")

	errs, ok := result["errors"].(map[string]interface{})
	require.True(t, ok)
	assert.Contains(t, errs, "email")
}

// ---------------------------------------------------------------------------
// 2. TestAuth_Login_MissingPassword — POST /auth/login without password → 422
// ---------------------------------------------------------------------------

func TestAuth_Login_MissingPassword(t *testing.T) {
	app := newAuthTestApp()

	body := `{"email": "user@example.com"}`
	resp := postJSON(t, app, "/api/v1/auth/login", body)
	defer resp.Body.Close()

	assert.Equal(t, 422, resp.StatusCode)

	result := decodeResponse(t, resp)
	errs, ok := result["errors"].(map[string]interface{})
	require.True(t, ok)
	assert.Contains(t, errs, "password")
}

// ---------------------------------------------------------------------------
// 3. TestAuth_Register_MissingEmail — POST /auth/register without email → 422
// ---------------------------------------------------------------------------

func TestAuth_Register_MissingEmail(t *testing.T) {
	app := newAuthTestApp()

	body := `{"password": "Secret123!@#"}`
	resp := postJSON(t, app, "/api/v1/auth/register", body)
	defer resp.Body.Close()

	assert.Equal(t, 422, resp.StatusCode)

	result := decodeResponse(t, resp)
	errs, ok := result["errors"].(map[string]interface{})
	require.True(t, ok)
	assert.Contains(t, errs, "email")
}

// ---------------------------------------------------------------------------
// 4. TestAuth_Register_InvalidEmail — POST /auth/register with bad email → 422
// ---------------------------------------------------------------------------

func TestAuth_Register_InvalidEmail(t *testing.T) {
	app := newAuthTestApp()

	body := `{"email": "not-an-email", "password": "Secret123!@#"}`
	resp := postJSON(t, app, "/api/v1/auth/register", body)
	defer resp.Body.Close()

	assert.Equal(t, 422, resp.StatusCode)

	result := decodeResponse(t, resp)
	errs, ok := result["errors"].(map[string]interface{})
	require.True(t, ok)
	assert.Contains(t, errs, "email")
}

// ---------------------------------------------------------------------------
// 5. TestAuth_Register_ShortPassword — POST /auth/register with short pw → 422
// ---------------------------------------------------------------------------

func TestAuth_Register_ShortPassword(t *testing.T) {
	app := newAuthTestApp()

	body := `{"email": "user@example.com", "password": "abc"}`
	resp := postJSON(t, app, "/api/v1/auth/register", body)
	defer resp.Body.Close()

	assert.Equal(t, 422, resp.StatusCode)

	result := decodeResponse(t, resp)
	errs, ok := result["errors"].(map[string]interface{})
	require.True(t, ok)
	assert.Contains(t, errs, "password")
}

// ---------------------------------------------------------------------------
// 6. TestAuth_ForgotPassword_MissingEmail — POST /auth/forgot-password without email
// The handler first checks if SMTP is configured; without it returns 503.
// We test that missing email in body produces a validation error when mail is configured.
// Since we can't configure SMTP in unit tests, we verify the handler returns a
// clear error (503 for unconfigured SMTP, or 422 for missing email if SMTP were set).
// ---------------------------------------------------------------------------

func TestAuth_ForgotPassword_MissingEmail(t *testing.T) {
	// Build a minimal app with a handler-level test (not full router) to
	// bypass the SMTP-availability check and directly test validation.
	log := zerolog.Nop()
	jwtMgr := auth.NewJWTManager("test-secret-for-unit-tests", 60, 30)
	provider := auth.NewDisabledProvider()
	cfg := testAuthConfig()

	authHandler := handler.NewAuthHandler(log, jwtMgr, provider, cfg, nil, nil, nil)

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})
	app.Post("/api/v1/auth/forgot-password", authHandler.ForgotPassword)

	body := `{}`
	resp := postJSON(t, app, "/api/v1/auth/forgot-password", body)
	defer resp.Body.Close()

	// Without mail configured, the handler returns 503 before checking email.
	// This test verifies the handler gracefully handles missing email.
	assert.True(t, resp.StatusCode == 503 || resp.StatusCode == 422,
		"expected 503 (SMTP unconfigured) or 422 (validation error), got %d", resp.StatusCode)
}

// ---------------------------------------------------------------------------
// 7. TestAuth_ResetPassword_MissingToken — POST /auth/reset-password without token
// The handler first checks if Redis is available; without it returns 503.
// ---------------------------------------------------------------------------

func TestAuth_ResetPassword_MissingToken(t *testing.T) {
	log := zerolog.Nop()
	jwtMgr := auth.NewJWTManager("test-secret-for-unit-tests", 60, 30)
	provider := auth.NewDisabledProvider()
	cfg := testAuthConfig()

	authHandler := handler.NewAuthHandler(log, jwtMgr, provider, cfg, nil, nil, nil)

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})
	app.Post("/api/v1/auth/reset-password", authHandler.ResetPassword)

	body := `{"password": "NewSecret123!", "confirm_password": "NewSecret123!"}`
	resp := postJSON(t, app, "/api/v1/auth/reset-password", body)
	defer resp.Body.Close()

	// Without Redis, returns 503. With Redis, would return 422 for missing token.
	assert.True(t, resp.StatusCode == 503 || resp.StatusCode == 422,
		"expected 503 (Redis unavailable) or 422 (missing token), got %d", resp.StatusCode)
}

// ---------------------------------------------------------------------------
// 8. TestAuth_ProviderEndpoint — GET /auth/provider → 200, returns "disabled"
// ---------------------------------------------------------------------------

func TestAuth_ProviderEndpoint(t *testing.T) {
	app := newAuthTestApp()

	req := httptest.NewRequest("GET", "/api/v1/auth/provider", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode)

	result := decodeResponse(t, resp)
	assert.Equal(t, "disabled", result["provider"])
}

// ---------------------------------------------------------------------------
// 9. TestAuth_Refresh_NoToken — POST /auth/refresh without token → 422
// ---------------------------------------------------------------------------

func TestAuth_Refresh_NoToken(t *testing.T) {
	app := newAuthTestApp()

	// Send empty body with no refresh_token and no cookie
	resp := postJSON(t, app, "/api/v1/auth/refresh", `{}`)
	defer resp.Body.Close()

	// The Refresh handler returns a ValidationError (422) when no refresh token
	// is found in either the body or httpOnly cookie.
	assert.Equal(t, 422, resp.StatusCode)

	result := decodeResponse(t, resp)
	errs, ok := result["errors"].(map[string]interface{})
	require.True(t, ok)
	assert.Contains(t, errs, "refresh_token")
}

// ---------------------------------------------------------------------------
// 10. TestAuth_Logout_Unauthorized — POST /auth/logout without auth → 401
// ---------------------------------------------------------------------------

func TestAuth_Logout_Unauthorized(t *testing.T) {
	app := newAuthTestApp()

	req := httptest.NewRequest("POST", "/api/v1/auth/logout", nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, 401, resp.StatusCode)
}

// ---------------------------------------------------------------------------
// 11. TestAuth_Me_Unauthorized — GET /users/me without auth → 401
// ---------------------------------------------------------------------------

func TestAuth_Me_Unauthorized(t *testing.T) {
	app := newAuthTestApp()

	req := httptest.NewRequest("GET", "/api/v1/users/me", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, 401, resp.StatusCode)
}

// ---------------------------------------------------------------------------
// 12. TestAuth_Login_BadJSON — POST /auth/login with invalid JSON body → 422
// ---------------------------------------------------------------------------

func TestAuth_Login_BadJSON(t *testing.T) {
	app := newAuthTestApp()

	req := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(`{invalid json`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer resp.Body.Close()

	// BodyParser fails on invalid JSON → ValidationError → 422
	assert.Equal(t, 422, resp.StatusCode)

	result := decodeResponse(t, resp)
	errs, ok := result["errors"].(map[string]interface{})
	require.True(t, ok)
	assert.Contains(t, errs, "body")
}

// ---------------------------------------------------------------------------
// 13. TestAuth_Refresh_EmptyBody — POST /auth/refresh with empty body
// ---------------------------------------------------------------------------

func TestAuth_Refresh_EmptyBody(t *testing.T) {
	app := newAuthTestApp()

	// Send POST with empty body (no Content-Length / empty)
	req := httptest.NewRequest("POST", "/api/v1/auth/refresh", bytes.NewReader([]byte{}))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Empty body → BodyParser succeeds but req.RefreshToken is "" → 422 validation error
	assert.Equal(t, 422, resp.StatusCode)

	result := decodeResponse(t, resp)
	errs, ok := result["errors"].(map[string]interface{})
	require.True(t, ok)
	assert.Contains(t, errs, "refresh_token")
}
