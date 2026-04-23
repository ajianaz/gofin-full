package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config uses flat keys that match env var names directly.
// This ensures viper's AutomaticEnv works correctly.
type Config struct {
	// App
	AppEnv      string `mapstructure:"APP_ENV"`
	AppDebug    bool   `mapstructure:"APP_DEBUG"`
	AppURL      string `mapstructure:"APP_URL"`
	AppTimezone string `mapstructure:"TZ"`

	// HTTP
	HTTPPort string `mapstructure:"HTTP_PORT"`
	HTTPHost string `mapstructure:"HTTP_HOST"`

	// Database
	DBHost            string        `mapstructure:"DB_HOST"`
	DBPort            int           `mapstructure:"DB_PORT"`
	DBDatabase        string        `mapstructure:"DB_DATABASE"`
	DBUsername        string        `mapstructure:"DB_USERNAME"`
	DBPassword        string        `mapstructure:"DB_PASSWORD"`
	DBSSLMode         string        `mapstructure:"DB_SSL_MODE"`
	DBSchema          string        `mapstructure:"DB_SCHEMA"`
	DBMaxOpenConns    int           `mapstructure:"DB_MAX_OPEN_CONNS"`
	DBMaxIdleConns    int           `mapstructure:"DB_MAX_IDLE_CONNS"`
	DBConnMaxLifetime int           `mapstructure:"DB_CONN_MAX_LIFETIME"`

	// Redis
	RedisHost    string `mapstructure:"REDIS_HOST"`
	RedisPort    int    `mapstructure:"REDIS_PORT"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDB      int    `mapstructure:"REDIS_DB"`
	RedisCacheDB  int    `mapstructure:"REDIS_CACHE_DB"`

	// Keycloak
	KeycloakURL          string `mapstructure:"KEYCLOAK_URL"`
	KeycloakRealm        string `mapstructure:"KEYCLOAK_REALM"`
	KeycloakClientID     string `mapstructure:"KEYCLOAK_CLIENT_ID"`
	KeycloakClientSecret string `mapstructure:"KEYCLOAK_CLIENT_SECRET"`

	// Log
	LogLevel  string `mapstructure:"LOG_LEVEL"`
	LogFormat string `mapstructure:"LOG_FORMAT"`

	// Auth
	AuthProvider         string `mapstructure:"AUTH_PROVIDER"`
	AuthJWTSecret        string `mapstructure:"AUTH_JWT_SECRET"`
	AuthJWTExpiry        int    `mapstructure:"AUTH_JWT_EXPIRY_MINUTES"`
	AuthRefreshExpiry    int    `mapstructure:"AUTH_REFRESH_EXPIRY_DAYS"`
	AuthAllowRegistration bool  `mapstructure:"AUTH_ALLOW_REGISTRATION"`
	StaticCronToken      string `mapstructure:"STATIC_CRON_TOKEN"`

	// OAuth - Google
	GoogleClientID     string `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `mapstructure:"GOOGLE_CLIENT_SECRET"`

	// OAuth - GitHub
	GitHubClientID     string `mapstructure:"GITHUB_CLIENT_ID"`
	GitHubClientSecret string `mapstructure:"GITHUB_CLIENT_SECRET"`

	// Security
	RateLimitMax             int    `mapstructure:"RATE_LIMIT_MAX"`
	RateLimitWindowSeconds   int    `mapstructure:"RATE_LIMIT_WINDOW_SECONDS"`
	MaxRequestBodyBytes      int64  `mapstructure:"MAX_REQUEST_BODY_BYTES"`
	CORSAllowedOrigins       string `mapstructure:"CORS_ALLOWED_ORIGINS"`
	Allow2FABypass           bool   `mapstructure:"ALLOW_2FA_BYPASS"`
	DisablePrometheus        bool   `mapstructure:"DISABLE_PROMETHEUS"`

	// Features
	FeatureExport           bool `mapstructure:"FEATURE_EXPORT"`
	FeatureWebhooks         bool `mapstructure:"FEATURE_WEBHOOKS"`
	FeatureHandleDebts      bool `mapstructure:"FEATURE_HANDLE_DEBTS"`
	FeatureExpressionEngine bool `mapstructure:"FEATURE_EXPRESSION_ENGINE"`
	FeatureRunningBalance   bool `mapstructure:"FEATURE_RUNNING_BALANCE"`

	// Business
	BusinessMaxUploadSize      int64 `mapstructure:"MAX_UPLOAD_SIZE"`
	BusinessAllowWebhooks      bool  `mapstructure:"ALLOW_WEBHOOKS"`
	BusinessWebhookMaxAttempts int   `mapstructure:"WEBHOOK_MAX_ATTEMPTS"`
	BusinessEnableExternalRates bool  `mapstructure:"ENABLE_EXTERNAL_RATES"`
	BusinessEnableExchangeRates bool  `mapstructure:"ENABLE_EXCHANGE_RATES"`
}

// Convenience accessors that group related fields.
func (c Config) IsLocal() bool      { return c.AppEnv == "local" || c.AppEnv == "testing" }
func (c Config) IsProduction() bool { return c.AppEnv == "production" }
func (c Config) HTTPAddr() string    { return fmt.Sprintf("%s:%s", c.HTTPHost, c.HTTPPort) }
func (c Config) RedisAddr() string   { return fmt.Sprintf("%s:%d", c.RedisHost, c.RedisPort) }
func (c Config) DatabaseDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s search_path=%s",
		c.DBHost, c.DBPort, c.DBUsername, c.DBPassword, c.DBDatabase, c.DBSSLMode, c.DBSchema,
	)
}
func (c Config) KeycloakRealmURL() string {
	return strings.TrimSuffix(c.KeycloakURL, "/") + "/realms/" + c.KeycloakRealm
}
func (c Config) KeycloakJWKSURL() string {
	return c.KeycloakRealmURL() + "/protocol/openid-connect/certs"
}
func (c Config) KeycloakTokenURL() string {
	return c.KeycloakRealmURL() + "/protocol/openid-connect/token"
}

func Load() (*Config, error) {
	v := viper.New()

	// Don't try to read a .env file in tests or if it doesn't exist
	v.SetConfigFile(".env")
	v.SetConfigType("env")
	v.AutomaticEnv()

	setDefaults(v)

	// Read config file (ignore error if not found)
	_ = v.ReadInConfig()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := validate(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func validate(cfg *Config) error {
	if cfg.AppEnv == "development" || cfg.AppEnv == "local" || cfg.AppEnv == "testing" {
		return nil
	}
	if cfg.AuthJWTSecret == "change-me-in-production-32chars!" {
		return fmt.Errorf("AUTH_JWT_SECRET must be changed from the default value in %s environment", cfg.AppEnv)
	}
	if len(cfg.AuthJWTSecret) < 32 {
		return fmt.Errorf("AUTH_JWT_SECRET must be at least 32 characters, got %d", len(cfg.AuthJWTSecret))
	}
	return nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("APP_ENV", "production")
	v.SetDefault("APP_DEBUG", false)
	v.SetDefault("APP_URL", "http://localhost")
	v.SetDefault("TZ", "UTC")

	v.SetDefault("HTTP_PORT", "8080")
	v.SetDefault("HTTP_HOST", "0.0.0.0")

	v.SetDefault("DB_HOST", "localhost")
	v.SetDefault("DB_PORT", 5432)
	v.SetDefault("DB_DATABASE", "gofin")
	v.SetDefault("DB_USERNAME", "gofin")
	v.SetDefault("DB_PASSWORD", "")
	v.SetDefault("DB_SSL_MODE", "prefer")
	v.SetDefault("DB_SCHEMA", "public")
	v.SetDefault("DB_MAX_OPEN_CONNS", 5)
	v.SetDefault("DB_MAX_IDLE_CONNS", 2)
	v.SetDefault("DB_CONN_MAX_LIFETIME", 300)

	v.SetDefault("REDIS_HOST", "localhost")
	v.SetDefault("REDIS_PORT", 6379)
	v.SetDefault("REDIS_PASSWORD", "")
	v.SetDefault("REDIS_DB", 0)
	v.SetDefault("REDIS_CACHE_DB", 1)

	v.SetDefault("KEYCLOAK_URL", "http://localhost:8088")
	v.SetDefault("KEYCLOAK_REALM", "gofin")
	v.SetDefault("KEYCLOAK_CLIENT_ID", "gofin-api")

	v.SetDefault("LOG_LEVEL", "info")
	v.SetDefault("LOG_FORMAT", "json")

	v.SetDefault("AUTH_PROVIDER", "local")
	v.SetDefault("AUTH_JWT_SECRET", "change-me-in-production-32chars!")
	v.SetDefault("AUTH_JWT_EXPIRY_MINUTES", 60)
	v.SetDefault("AUTH_REFRESH_EXPIRY_DAYS", 30)
	v.SetDefault("AUTH_ALLOW_REGISTRATION", false)
	v.SetDefault("STATIC_CRON_TOKEN", "PLEASE_REPLACE_WITH_32_CHAR_CODE")
	v.SetDefault("GOOGLE_CLIENT_ID", "")
	v.SetDefault("GOOGLE_CLIENT_SECRET", "")
	v.SetDefault("GITHUB_CLIENT_ID", "")
	v.SetDefault("GITHUB_CLIENT_SECRET", "")
	v.SetDefault("KEYCLOAK_CLIENT_SECRET", "")

	v.SetDefault("RATE_LIMIT_MAX", 100)
	v.SetDefault("RATE_LIMIT_WINDOW_SECONDS", 60)
	v.SetDefault("MAX_REQUEST_BODY_BYTES", 10485760)
	v.SetDefault("CORS_ALLOWED_ORIGINS", "http://localhost:5173")
	v.SetDefault("ALLOW_2FA_BYPASS", false)
	v.SetDefault("DISABLE_PROMETHEUS", false)

	v.SetDefault("FEATURE_EXPORT", true)
	v.SetDefault("FEATURE_WEBHOOKS", true)
	v.SetDefault("FEATURE_HANDLE_DEBTS", true)
	v.SetDefault("FEATURE_EXPRESSION_ENGINE", true)
	v.SetDefault("FEATURE_RUNNING_BALANCE", true)

	v.SetDefault("MAX_UPLOAD_SIZE", 1073741824)
	v.SetDefault("ALLOW_WEBHOOKS", false)
	v.SetDefault("WEBHOOK_MAX_ATTEMPTS", 3)
	v.SetDefault("ENABLE_EXTERNAL_RATES", false)
	v.SetDefault("ENABLE_EXCHANGE_RATES", false)
}
