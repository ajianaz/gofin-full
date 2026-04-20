package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// CORS returns a CORS middleware. In production, restrict to APP_URL.
// Falls back to wildcard for local/testing environments.
func CORS(appURL, appEnv string) fiber.Handler {
	cfg := cors.Config{
		AllowMethods:  "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:  "Origin,Content-Type,Accept,Authorization,X-Trace-Id,X-Request-ID,X-API-Key",
		ExposeHeaders: "X-Request-ID,X-Trace-Id",
		MaxAge:        86400,
	}

	// Restrict origins in production
	if appEnv == "production" && appURL != "" {
		cfg.AllowOrigins = appURL
	} else {
		cfg.AllowOrigins = "*"
	}

	return cors.New(cfg)
}
