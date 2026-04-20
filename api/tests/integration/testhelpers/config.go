package testhelpers

import (
	"fmt"
	"os"
)

// TestConfig holds connection details for integration test dependencies.
// All fields have safe defaults suitable for a local Docker-based test environment.
type TestConfig struct {
	DBHost     string
	DBPort     string
	DBDatabase string
	DBUser     string
	DBPassword string
	RedisHost  string
	RedisPort  string
	JWTSecret  string
}

// NewTestConfig builds a TestConfig using environment variables when set,
// otherwise falling back to defaults that match the project's Docker Compose setup.
func NewTestConfig() *TestConfig {
	return &TestConfig{
		DBHost:     envOr("DB_HOST", "localhost"),
		DBPort:     envOr("DB_PORT", "5433"),
		DBDatabase: envOr("DB_DATABASE", "gofin_test"),
		DBUser:     envOr("DB_USERNAME", "gofin_test"),
		DBPassword: envOr("DB_PASSWORD", "gofin_test"),
		RedisHost:  envOr("REDIS_HOST", "localhost"),
		RedisPort:  envOr("REDIS_PORT", "6380"),
		JWTSecret:  envOr("AUTH_JWT_SECRET", "test-jwt-secret-for-integration-tests-32ch"),
	}
}

// DSN returns a PostgreSQL connection string suitable for pgx.
func (c *TestConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBDatabase,
	)
}

// RedisAddr returns the host:port address for Redis.
func (c *TestConfig) RedisAddr() string {
	return c.RedisHost + ":" + c.RedisPort
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
