package middleware_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	apperrors "github.com/ajianaz/gofin-full/api/pkg/errors"
	"github.com/ajianaz/gofin-full/api/internal/middleware"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestApp(handlers ...fiber.Handler) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})
	for _, h := range handlers {
		app.Use(h)
	}
	return app
}

// --- RequestID Middleware ---

func TestRequestID_GeneratesNewID(t *testing.T) {
	app := setupTestApp(middleware.RequestID())
	app.Get("/", func(c *fiber.Ctx) error {
		id := c.Locals("request_id")
		return c.JSON(fiber.Map{"request_id": id})
	})

	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)

	var body map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&body)
	assert.NotEmpty(t, body["request_id"])
	assert.Len(t, body["request_id"], 36) // UUID format
}

func TestRequestID_UsesXTraceID(t *testing.T) {
	app := setupTestApp(middleware.RequestID())
	app.Get("/", func(c *fiber.Ctx) error {
		id := c.Locals("request_id")
		return c.JSON(fiber.Map{"request_id": id})
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Trace-Id", "custom-trace-id")
	resp, err := app.Test(req, -1)
	require.NoError(t, err)

	var body map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Equal(t, "custom-trace-id", body["request_id"])
}

func TestRequestID_SetsHeader(t *testing.T) {
	app := setupTestApp(middleware.RequestID())
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	assert.NotEmpty(t, resp.Header.Get("X-Request-ID"))
}

// --- Recovery Middleware ---

func TestRecovery_RecoversFromPanic(t *testing.T) {
	var buf bytes.Buffer
	log := zerolog.New(&buf)

	app := setupTestApp(middleware.Recovery(log))
	app.Get("/", func(c *fiber.Ctx) error {
		panic("test panic")
	})

	req := httptest.NewRequest("GET", "/", nil)
	// Use -1 timeout since panic recovery adds latency
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	assert.Equal(t, 500, resp.StatusCode)

	var body map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Contains(t, body["message"], "Internal Server Error")
}

func TestRecovery_PassesThroughNormally(t *testing.T) {
	var buf bytes.Buffer
	log := zerolog.New(&buf)

	app := setupTestApp(middleware.Recovery(log))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

// --- Accept Headers Middleware ---

func TestAcceptHeaders_ValidAccept(t *testing.T) {
	app := setupTestApp(middleware.AcceptHeaders())
	app.Post("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	for _, accept := range []string{
		"application/json",
		"application/vnd.api+json",
		"*/*",
	} {
		t.Run(accept, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/", nil)
			req.Header.Set("Accept", accept)
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			assert.Equal(t, 200, resp.StatusCode)
		})
	}
}

func TestAcceptHeaders_InvalidAccept(t *testing.T) {
	app := setupTestApp(middleware.AcceptHeaders())
	app.Post("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	req := httptest.NewRequest("POST", "/", nil)
	req.Header.Set("Accept", "text/xml")
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	assert.Equal(t, 406, resp.StatusCode)
}

func TestAcceptHeaders_ValidContentType(t *testing.T) {
	app := setupTestApp(middleware.AcceptHeaders())
	app.Post("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	for _, ct := range []string{
		"application/json",
		"application/vnd.api+json",
	} {
		t.Run(ct, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/", nil)
			req.Header.Set("Content-Type", ct)
			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			assert.Equal(t, 200, resp.StatusCode)
		})
	}
}

func TestAcceptHeaders_EmptyContentType(t *testing.T) {
	app := setupTestApp(middleware.AcceptHeaders())
	app.Post("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	req := httptest.NewRequest("POST", "/", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	assert.Equal(t, 415, resp.StatusCode)
}

func TestAcceptHeaders_InvalidContentType(t *testing.T) {
	app := setupTestApp(middleware.AcceptHeaders())
	app.Post("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	req := httptest.NewRequest("POST", "/", strings.NewReader("{}"))
	req.Header.Set("Content-Type", "text/xml")
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	assert.Equal(t, 415, resp.StatusCode)
}

func TestAcceptHeaders_NoContentTypeForGET(t *testing.T) {
	app := setupTestApp(middleware.AcceptHeaders())
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

// --- Error Handler ---

func TestErrorHandler_AppError(t *testing.T) {
	app := setupTestApp()
	app.Get("/", func(c *fiber.Ctx) error {
		return apperrors.NotFoundResource("wallet", 42)
	})

	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)

	var body map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Equal(t, "Resource not found", body["title"])
}

func TestErrorHandler_ValidationError(t *testing.T) {
	app := setupTestApp()
	app.Post("/", func(c *fiber.Ctx) error {
		return apperrors.NewValidationError(map[string][]string{
			"name": {"The name field is required."},
		})
	})

	req := httptest.NewRequest("POST", "/", strings.NewReader("{}"))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	assert.Equal(t, 422, resp.StatusCode)

	var body map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Equal(t, "The given data was invalid.", body["message"])
}

func TestErrorHandler_InternalError(t *testing.T) {
	app := setupTestApp()
	app.Get("/", func(c *fiber.Ctx) error {
		return apperrors.ErrInternal
	})

	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	assert.Equal(t, 500, resp.StatusCode)

	var body map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Equal(t, "Internal Server Error", body["message"])
}

func TestErrorHandler_FiberError(t *testing.T) {
	app := setupTestApp()
	app.Get("/", func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusBadRequest, "Bad request from fiber")
	})

	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestErrorHandler_UnknownError(t *testing.T) {
	app := setupTestApp()
	app.Get("/", func(c *fiber.Ctx) error {
		return io.EOF
	})

	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	assert.Equal(t, 500, resp.StatusCode)
}

// --- Logger Middleware ---

func TestLogger_LogsRequest(t *testing.T) {
	var buf bytes.Buffer
	log := zerolog.New(&buf)

	app := setupTestApp(middleware.Logger(log))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/", nil)
	_, err := app.Test(req, -1)
	require.NoError(t, err)

	logOutput := buf.String()
	assert.Contains(t, logOutput, "GET")
	assert.Contains(t, logOutput, "\"status\":200")
	assert.Contains(t, logOutput, "request")
}

func TestLogger_CapturesRequestID(t *testing.T) {
	var buf bytes.Buffer
	log := zerolog.New(&buf)

	app := setupTestApp(middleware.RequestID(), middleware.Logger(log))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	req := httptest.NewRequest("GET", "/", nil)
	_, err := app.Test(req, -1)
	require.NoError(t, err)

	logOutput := buf.String()
	assert.Contains(t, logOutput, "request_id")
}

// --- CORS Middleware ---

func TestCORS_SetsHeaders(t *testing.T) {
	app := setupTestApp(middleware.CORS("", "local"))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Origin", "http://example.com")
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// CORS headers should be present
	assert.NotEmpty(t, resp.Header.Get("Access-Control-Allow-Origin"))
}

func TestCORS_Preflight(t *testing.T) {
	app := setupTestApp(middleware.CORS("", "local"))
	app.Options("/", func(c *fiber.Ctx) error {
		return c.SendStatus(204)
	})

	req := httptest.NewRequest("OPTIONS", "/", nil)
	req.Header.Set("Origin", "http://example.com")
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	assert.Equal(t, 204, resp.StatusCode)
}
