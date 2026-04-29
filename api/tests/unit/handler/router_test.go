package handler_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/ajianaz/gofin-full/api/internal/auth"
	"github.com/ajianaz/gofin-full/api/internal/config"
	"github.com/ajianaz/gofin-full/api/internal/handler"
	"github.com/ajianaz/gofin-full/api/internal/router"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestRouter(healthHandler *handler.HealthHandler) *router.RouterConfig {
	return &router.RouterConfig{
		HealthHandler: healthHandler,
		JWTManager:    auth.NewJWTManager("test-secret-for-unit-tests", 60, 30),
	}
}

func TestRouter_HealthEndpoint(t *testing.T) {
	healthHandler := handler.NewHealthHandler(nil, nil)
	cfg := newTestRouter(healthHandler)
	app := router.New(*cfg)

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	assert.Equal(t, 503, resp.StatusCode) // Degraded since no DB/Redis
}

func TestRouter_APIv1Endpoint(t *testing.T) {
	healthHandler := handler.NewHealthHandler(nil, nil)
	cfg := newTestRouter(healthHandler)
	app := router.New(*cfg)

	req := httptest.NewRequest("GET", "/api/v1/", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var body map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Equal(t, "Gofin API v1", body["message"])
	assert.Equal(t, "1.0.0", body["version"])
}

func TestRouter_NotFound(t *testing.T) {
	healthHandler := handler.NewHealthHandler(nil, nil)
	cfg := newTestRouter(healthHandler)
	app := router.New(*cfg)

	req := httptest.NewRequest("GET", "/nonexistent", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)
}

func TestRouter_CORSHeaders(t *testing.T) {
	healthHandler := handler.NewHealthHandler(nil, nil)
	cfg := newTestRouter(healthHandler)
	app := router.New(*cfg)

	req := httptest.NewRequest("GET", "/health", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	resp, err := app.Test(req, -1)
	require.NoError(t, err)

	// CORS header should be set for localhost origins
	assert.NotEmpty(t, resp.Header.Get("Access-Control-Allow-Origin"))
}

func TestRouter_RequestIDHeader(t *testing.T) {
	healthHandler := handler.NewHealthHandler(nil, nil)
	cfg := newTestRouter(healthHandler)
	app := router.New(*cfg)

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)

	// Request ID should be set
	assert.NotEmpty(t, resp.Header.Get("X-Request-ID"))
}

func TestRouter_XTraceIDPassthrough(t *testing.T) {
	healthHandler := handler.NewHealthHandler(nil, nil)
	cfg := newTestRouter(healthHandler)
	app := router.New(*cfg)

	req := httptest.NewRequest("GET", "/health", nil)
	req.Header.Set("X-Trace-Id", "my-custom-trace-id")
	resp, err := app.Test(req, -1)
	require.NoError(t, err)

	// Should use the provided trace ID
	assert.Equal(t, "my-custom-trace-id", resp.Header.Get("X-Request-ID"))
}

func TestRouter_AuthProviderEndpoint(t *testing.T) {
	healthHandler := handler.NewHealthHandler(nil, nil)
	jwtMgr := auth.NewJWTManager("test-secret", 60, 30)
	provider := auth.NewDisabledProvider()

	cfg := &router.RouterConfig{
		HealthHandler: healthHandler,
		AuthHandler:   handler.NewAuthHandler(jwtMgr, provider, &config.Config{AuthProvider: "disabled"}, nil, nil, nil),
		JWTManager:    jwtMgr,
	}
	app := router.New(*cfg)

	req := httptest.NewRequest("GET", "/api/v1/auth/provider", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var body map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Equal(t, "disabled", body["provider"])
}

func TestRouter_ProtectedRouteUnauthorized(t *testing.T) {
	healthHandler := handler.NewHealthHandler(nil, nil)
	jwtMgr := auth.NewJWTManager("test-secret", 60, 30)

	cfg := &router.RouterConfig{
		HealthHandler: healthHandler,
		JWTManager:    jwtMgr,
	}
	app := router.New(*cfg)

	req := httptest.NewRequest("GET", "/api/v1/users/me", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}
