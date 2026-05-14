package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// CORS returns a CORS middleware.
// Priority: CORS_ALLOWED_ORIGINS env var > APP_URL (production) > localhost fallback.
func CORS(appURL, appEnv, corsAllowedOrigins string) fiber.Handler {
	cfg := cors.Config{
		AllowMethods:  "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:  "Origin,Content-Type,Accept,Authorization,X-Trace-Id,X-Request-ID,X-API-Key",
		ExposeHeaders: "X-Request-ID,X-Trace-Id",
		MaxAge:        86400,
	}

	switch {
	case corsAllowedOrigins != "":
		// Explicit CORS origins from env var (supports comma-separated list)
		cfg.AllowOrigins = corsAllowedOrigins
	case appEnv == "production" && appURL != "":
		// Production: restrict to APP_URL
		cfg.AllowOrigins = appURL
	default:
		// Development: only allow localhost origins
		cfg.AllowOrigins = "http://localhost:5173,http://localhost:8080"
	}

	return cors.New(cfg)
}
