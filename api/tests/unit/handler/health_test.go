package handler_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/azfirazka/gofin-full/api/internal/handler"
	"github.com/azfirazka/gofin-full/api/internal/middleware"
	response "github.com/azfirazka/gofin-full/api/internal/dto/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthCheck_NoDependencies(t *testing.T) {
	h := handler.NewHealthHandler(nil, nil)

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})
	app.Get("/health", h.Check)

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	assert.Equal(t, 503, resp.StatusCode)

	var body response.HealthResponse
	err = json.NewDecoder(resp.Body).Decode(&body)
	require.NoError(t, err)

	assert.Equal(t, "degraded", body.Status)
	assert.Len(t, body.Services, 2)

	// Find service statuses
	serviceMap := make(map[string]response.ServiceHealth)
	for _, s := range body.Services {
		serviceMap[s.Name] = s
	}

	assert.Equal(t, "error", serviceMap["postgresql"].Status)
	assert.Contains(t, serviceMap["postgresql"].Error, "not initialized")

	assert.Equal(t, "error", serviceMap["redis"].Status)
	assert.Contains(t, serviceMap["redis"].Error, "not initialized")
}

func TestHealthCheck_ResponseFormat(t *testing.T) {
	h := handler.NewHealthHandler(nil, nil)

	app := fiber.New()
	app.Get("/health", h.Check)

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)

	// Verify response is valid JSON
	var body map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	require.NoError(t, err)

	assert.Contains(t, body, "status")
	assert.Contains(t, body, "services")

	services, ok := body["services"].([]interface{})
	require.True(t, ok)
	assert.Len(t, services, 2)
}

func TestHealthCheck_ServicesOrder(t *testing.T) {
	h := handler.NewHealthHandler(nil, nil)

	app := fiber.New()
	app.Get("/health", h.Check)

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)

	var body response.HealthResponse
	err = json.NewDecoder(resp.Body).Decode(&body)
	require.NoError(t, err)

	// PostgreSQL should be first, Redis second
	assert.Equal(t, "postgresql", body.Services[0].Name)
	assert.Equal(t, "redis", body.Services[1].Name)
}

func TestHealthCheck_ContentLength(t *testing.T) {
	h := handler.NewHealthHandler(nil, nil)

	app := fiber.New()
	app.Get("/health", h.Check)

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)

	// Verify response body has content
	buf := make([]byte, 1024)
	n, _ := resp.Body.Read(buf)
	assert.Greater(t, n, 0)
}
