package config_test

import (
	"os"
	"testing"

	"github.com/ajianaz/gofin-full/api/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_Defaults(t *testing.T) {
	// Ensure no .env file interferes
	_ = os.Remove(".env")

	// Clear env vars that docker-compose sets, so we test true defaults
	envVars := []string{
		"APP_ENV", "APP_DEBUG", "APP_URL", "APP_TIMEZONE",
		"HTTP_PORT", "HTTP_HOST",
		"DB_HOST", "DB_PORT", "DB_DATABASE", "DB_USERNAME", "DB_PASSWORD", "DB_SSL_MODE", "DB_SCHEMA", "DB_DSN",
		"REDIS_HOST", "REDIS_PORT", "REDIS_DB",
		"KEYCLOAK_URL", "KEYCLOAK_REALM", "KEYCLOAK_CLIENT_ID",
		"LOG_LEVEL", "LOG_FORMAT",
		"FEATURE_EXPORT", "FEATURE_WEBHOOKS", "FEATURE_HANDLE_DEBTS", "FEATURE_EXPRESSION_ENGINE", "FEATURE_RUNNING_BALANCE",
		"BUSINESS_MAX_UPLOAD_SIZE", "BUSINESS_ALLOW_WEBHOOKS", "BUSINESS_WEBHOOK_MAX_ATTEMPTS", "BUSINESS_ENABLE_EXTERNAL_RATES", "BUSINESS_ENABLE_EXCHANGE_RATES",
		"RATE_LIMIT_MAX", "RATE_LIMIT_WINDOW_SECONDS", "MAX_REQUEST_BODY_BYTES", "CORS_ALLOWED_ORIGINS", "ALLOW_2FA_BYPASS",
	}
	for _, k := range envVars {
		t.Setenv(k, "")
		os.Unsetenv(k)
	}

	cfg, err := config.Load()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// App defaults
	assert.Equal(t, "production", cfg.AppEnv)
	assert.False(t, cfg.AppDebug)
	assert.Equal(t, "http://localhost", cfg.AppURL)
	assert.Equal(t, "UTC", cfg.AppTimezone)

	// HTTP defaults
	assert.Equal(t, "8080", cfg.HTTPPort)
	assert.Equal(t, "0.0.0.0", cfg.HTTPHost)
	assert.Equal(t, "0.0.0.0:8080", cfg.HTTPAddr())

	// Database defaults
	assert.Equal(t, "localhost", cfg.DBHost)
	assert.Equal(t, 5432, cfg.DBPort)
	assert.Equal(t, "gofin", cfg.DBDatabase)
	assert.Equal(t, "gofin", cfg.DBUsername)
	assert.Equal(t, "prefer", cfg.DBSSLMode)
	assert.Equal(t, "public", cfg.DBSchema)
	assert.Equal(t, 5, cfg.DBMaxOpenConns)
	assert.Equal(t, 2, cfg.DBMaxIdleConns)

	// Redis defaults
	assert.Equal(t, "localhost", cfg.RedisHost)
	assert.Equal(t, 6379, cfg.RedisPort)
	assert.Equal(t, "localhost:6379", cfg.RedisAddr())
	assert.Equal(t, 0, cfg.RedisDB)
	assert.Equal(t, 1, cfg.RedisCacheDB)

	// Keycloak defaults
	assert.Equal(t, "http://localhost:8088", cfg.KeycloakURL)
	assert.Equal(t, "gofin", cfg.KeycloakRealm)
	assert.Equal(t, "gofin-api", cfg.KeycloakClientID)
	assert.Equal(t, "http://localhost:8088/realms/gofin", cfg.KeycloakRealmURL())
	assert.Contains(t, cfg.KeycloakJWKSURL(), "/protocol/openid-connect/certs")
	assert.Contains(t, cfg.KeycloakTokenURL(), "/protocol/openid-connect/token")

	// Log defaults
	assert.Equal(t, "info", cfg.LogLevel)
	assert.Equal(t, "json", cfg.LogFormat)

	// Feature defaults
	assert.True(t, cfg.FeatureExport)
	assert.True(t, cfg.FeatureWebhooks)
	assert.True(t, cfg.FeatureHandleDebts)
	assert.True(t, cfg.FeatureExpressionEngine)
	assert.True(t, cfg.FeatureRunningBalance)

	// Business defaults
	assert.Equal(t, int64(1073741824), cfg.BusinessMaxUploadSize)
	assert.False(t, cfg.BusinessAllowWebhooks)
	assert.Equal(t, 3, cfg.BusinessWebhookMaxAttempts)
	assert.False(t, cfg.BusinessEnableExternalRates)
	assert.False(t, cfg.BusinessEnableExchangeRates)

	// Security defaults
	assert.Equal(t, 100, cfg.RateLimitMax)
	assert.Equal(t, 60, cfg.RateLimitWindowSeconds)
	assert.Equal(t, int64(10485760), cfg.MaxRequestBodyBytes)
	assert.Equal(t, "http://localhost:5173", cfg.CORSAllowedOrigins)
	assert.False(t, cfg.Allow2FABypass)
}

func TestLoad_FromEnvVars(t *testing.T) {
	// Viper reads env vars at load time via AutomaticEnv
	// Set env vars before calling Load
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_DEBUG", "true")
	t.Setenv("HTTP_PORT", "9090")
	t.Setenv("DB_HOST", "db-host")
	t.Setenv("DB_PORT", "5433")
	t.Setenv("REDIS_HOST", "redis-host")
	t.Setenv("LOG_FORMAT", "console")
	t.Setenv("RATE_LIMIT_MAX", "50")
	t.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,https://app.gofin.io")

	cfg, err := config.Load()
	require.NoError(t, err)

	assert.Equal(t, "testing", cfg.AppEnv)
	assert.True(t, cfg.AppDebug)
	assert.Equal(t, "9090", cfg.HTTPPort)
	assert.Equal(t, "0.0.0.0:9090", cfg.HTTPAddr())
	assert.Equal(t, "db-host", cfg.DBHost)
	assert.Equal(t, 5433, cfg.DBPort)
	assert.Equal(t, "redis-host", cfg.RedisHost)
	assert.Equal(t, "redis-host:6379", cfg.RedisAddr())
	assert.Equal(t, "console", cfg.LogFormat)
	assert.Equal(t, 50, cfg.RateLimitMax)
	assert.Equal(t, "http://localhost:3000,https://app.gofin.io", cfg.CORSAllowedOrigins)
}

func TestLoad_DatabaseDSN(t *testing.T) {
	cfg := &config.Config{
		DBHost:     "localhost",
		DBPort:     5432,
		DBUsername: "gofin",
		DBPassword: "secret",
		DBDatabase: "gofin",
		DBSSLMode:  "prefer",
		DBSchema:   "public",
	}

	dsn := cfg.DatabaseDSN()
	assert.Contains(t, dsn, "host=localhost")
	assert.Contains(t, dsn, "port=5432")
	assert.Contains(t, dsn, "user=gofin")
	assert.Contains(t, dsn, "password=secret")
	assert.Contains(t, dsn, "dbname=gofin")
	assert.Contains(t, dsn, "sslmode=prefer")
	assert.Contains(t, dsn, "search_path=public")
}

func TestConfig_IsLocal(t *testing.T) {
	tests := []struct {
		name string
		env  string
		want bool
	}{
		{"local is local", "local", true},
		{"testing is local", "testing", true},
		{"production is not local", "production", false},
		{"staging is not local", "staging", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{AppEnv: tt.env}
			assert.Equal(t, tt.want, cfg.IsLocal())
		})
	}
}

func TestConfig_IsProduction(t *testing.T) {
	assert.True(t, (&config.Config{AppEnv: "production"}).IsProduction())
	assert.False(t, (&config.Config{AppEnv: "local"}).IsProduction())
	assert.False(t, (&config.Config{AppEnv: "testing"}).IsProduction())
}

func TestConfig_HTTPAddr(t *testing.T) {
	cfg := &config.Config{HTTPHost: "127.0.0.1", HTTPPort: "3000"}
	assert.Equal(t, "127.0.0.1:3000", cfg.HTTPAddr())
}

func TestConfig_KeycloakURLs(t *testing.T) {
	cfg := &config.Config{
		KeycloakURL:     "http://keycloak:8080",
		KeycloakRealm:   "myrealm",
		KeycloakClientID: "myclient",
	}

	assert.Equal(t, "http://keycloak:8080/realms/myrealm", cfg.KeycloakRealmURL())
	assert.Equal(t, "http://keycloak:8080/realms/myrealm/protocol/openid-connect/certs", cfg.KeycloakJWKSURL())
	assert.Equal(t, "http://keycloak:8080/realms/myrealm/protocol/openid-connect/token", cfg.KeycloakTokenURL())
}

func TestConfig_KeycloakURLs_TrailingSlash(t *testing.T) {
	cfg := &config.Config{
		KeycloakURL:   "http://keycloak:8080/",
		KeycloakRealm: "myrealm",
	}

	// Should handle trailing slash
	assert.Equal(t, "http://keycloak:8080/realms/myrealm", cfg.KeycloakRealmURL())
}
