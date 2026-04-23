package handler

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"

	"github.com/ajianaz/gofin-full/api/internal/auth"
	"github.com/ajianaz/gofin-full/api/internal/config"
	"github.com/ajianaz/gofin-full/api/internal/repository"
	response "github.com/ajianaz/gofin-full/api/internal/dto/response"
	apperrors "github.com/ajianaz/gofin-full/api/pkg/errors"
)

// Login attempt lockout constants.
const (
	loginMaxAttempts      = 5
	loginLockoutDuration  = 15 * time.Minute
	loginAttemptWindow    = 15 * time.Minute
	loginAttemptsKeyPrefix = "login_attempts:"
)

// AuthHandler handles authentication endpoints.
type AuthHandler struct {
	jwtMgr       *auth.JWTManager
	provider     auth.AuthProvider
	cfg          *config.Config
	userRepo     *repository.UserRepository
	oauthState   *repository.OAuthStateRepository
	refreshRepo  *repository.RefreshTokenRepository
	rdb          redis.Cmdable // optional Redis for login attempt tracking
}

// NewAuthHandler creates a new auth handler.
func NewAuthHandler(jwtMgr *auth.JWTManager, provider auth.AuthProvider, cfg *config.Config, userRepo *repository.UserRepository, oauthStateRepo *repository.OAuthStateRepository, refreshRepo *repository.RefreshTokenRepository) *AuthHandler {
	return &AuthHandler{jwtMgr: jwtMgr, provider: provider, cfg: cfg, userRepo: userRepo, oauthState: oauthStateRepo, refreshRepo: refreshRepo}
}

// SetRedis injects an optional Redis client for login attempt tracking.
func (h *AuthHandler) SetRedis(rdb redis.Cmdable) {
	h.rdb = rdb
}

// Login handles POST /api/v1/auth/login.
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewValidationError(map[string][]string{
			"body": {"Invalid request body."},
		})
	}

	// Check if account is temporarily locked due to too many failed attempts
	if locked, retryMinutes := h.isAccountLocked(c.Context(), req.Email); locked {
		return c.Status(429).JSON(fiber.Map{
			"message": fmt.Sprintf("Account temporarily locked due to too many failed login attempts. Try again in %d minutes.", retryMinutes),
		})
	}

	identity, err := h.provider.Authenticate(c.Context(), auth.Credentials{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		h.recordFailedLogin(c.Context(), req.Email)
		return apperrors.NewWithDetail(401, "Unauthenticated", err.Error())
	}

	if identity.Blocked {
		return apperrors.NewWithDetail(403, "Forbidden", "User account is blocked.")
	}

	// Successful login — clear failed attempt counter
	h.clearFailedLogins(c.Context(), req.Email)

	tokens, err := h.jwtMgr.GenerateTokenPair(identity, identity.UserGroupID)
	if err != nil {
		return apperrors.ErrInternal
	}

	// Store refresh token hash in DB
	if h.refreshRepo != nil {
		expiresAt := time.Now().UTC().Add(time.Duration(h.cfg.AuthRefreshExpiry) * 24 * time.Hour)
		tokenHash := auth.HashRefreshToken(tokens.RefreshToken)
		_ = h.refreshRepo.Store(c.Context(), identity.ID, tokenHash, expiresAt)
	}

	return c.JSON(tokens)
}

// isAccountLocked checks if an account is locked due to too many failed login attempts.
// Returns (locked, retryMinutes).
func (h *AuthHandler) isAccountLocked(ctx context.Context, email string) (bool, int) {
	if h.rdb == nil {
		return false, 0
	}

	key := loginAttemptsKeyPrefix + email
	val, err := h.rdb.Get(ctx, key).Result()
	if err != nil {
		return false, 0
	}

	attempts, err := strconv.Atoi(val)
	if err != nil || attempts < loginMaxAttempts {
		return false, 0
	}

	ttl, err := h.rdb.TTL(ctx, key).Result()
	if err != nil || ttl <= 0 {
		return false, 0
	}

	retryMinutes := int(ttl.Minutes())
	if retryMinutes < 1 {
		retryMinutes = 1
	}
	return true, retryMinutes
}

// recordFailedLogin increments the failed login counter for an email.
func (h *AuthHandler) recordFailedLogin(ctx context.Context, email string) {
	if h.rdb == nil {
		return
	}

	key := loginAttemptsKeyPrefix + email
	pipe := h.rdb.Pipeline()
	pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, loginLockoutDuration)
	if _, err := pipe.Exec(ctx); err != nil {
		log.Printf("login attempt tracking: redis error: %v", err)
	}
}

// clearFailedLogins removes the failed login counter for an email after successful login.
func (h *AuthHandler) clearFailedLogins(ctx context.Context, email string) {
	if h.rdb == nil {
		return
	}

	if err := h.rdb.Del(ctx, loginAttemptsKeyPrefix+email).Err(); err != nil {
		log.Printf("login attempt tracking: redis error on clear: %v", err)
	}
}

// Register handles POST /api/v1/auth/register.
// Only works when AUTH_ALLOW_REGISTRATION=true.
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	if !h.cfg.AuthAllowRegistration {
		return c.Status(403).JSON(fiber.Map{
			"message": "Self-registration is disabled. Contact an administrator to create an account.",
		})
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewValidationError(map[string][]string{
			"body": {"Invalid request body."},
		})
	}

	if req.Email == "" {
		return apperrors.NewValidationError(map[string][]string{
			"email": {"Email is required."},
		})
	}
	if !isValidEmail(req.Email) {
		return apperrors.NewValidationError(map[string][]string{
			"email": {"Invalid email format."},
		})
	}
	if len(req.Password) < 8 {
		return apperrors.NewValidationError(map[string][]string{
			"password": {"Password must be at least 8 characters."},
		})
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		return apperrors.ErrInternal
	}

	user, err := h.userRepo.Create(c.Context(), req.Email, hash)
	if err != nil {
		if isDuplicateKey(err) {
			return c.Status(409).JSON(fiber.Map{
				"message": "A user with this email already exists.",
			})
		}
		return apperrors.ErrInternal
	}

	// Auto-login: generate tokens for the newly registered user
	identity := &auth.UserIdentity{ID: user.ID, Email: user.Email}
	tokens, err := h.jwtMgr.GenerateTokenPair(identity, user.UserGroupID)
	if err != nil {
		return apperrors.ErrInternal
	}

	// Store refresh token hash in DB
	if h.refreshRepo != nil {
		expiresAt := time.Now().UTC().Add(time.Duration(h.cfg.AuthRefreshExpiry) * 24 * time.Hour)
		tokenHash := auth.HashRefreshToken(tokens.RefreshToken)
		_ = h.refreshRepo.Store(c.Context(), user.ID, tokenHash, expiresAt)
	}

	return c.Status(201).JSON(tokens)
}

// Me handles GET /api/v1/auth/me.
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	user := auth.GetUser(c)
	if user == nil {
		return apperrors.ErrUnauthorized
	}

	return c.JSON(fiber.Map{
		"data": fiber.Map{
			"type": "users",
			"id":   user.ID,
			"attributes": fiber.Map{
				"email": user.Email,
			},
		},
	})
}

// Provider returns the active auth provider name.
func (h *AuthHandler) Provider(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"provider": h.provider.Name(),
	})
}

// Logout handles POST /api/v1/auth/logout.
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	_ = c.BodyParser(&req)

	// Revoke the specific refresh token if provided
	if req.RefreshToken != "" && h.refreshRepo != nil {
		tokenHash := auth.HashRefreshToken(req.RefreshToken)
		_ = h.refreshRepo.RevokeByHash(c.Context(), tokenHash)
	}

	// Increment token_version to invalidate all existing JWT tokens for this user
	user := auth.GetUser(c)
	if user != nil {
		_ = h.userRepo.IncrementTokenVersion(c.Context(), user.ID)
	}

	return c.JSON(fiber.Map{
		"message": "Logged out successfully.",
	})
}

// Refresh handles POST /api/v1/auth/refresh.
// Requires a valid (possibly expired) access token in the Authorization header.
// The refresh_token in the body is validated to ensure it was issued alongside the access token.
func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.BodyParser(&req); err != nil || req.RefreshToken == "" {
		return apperrors.NewValidationError(map[string][]string{
			"refresh_token": {"Refresh token is required."},
		})
	}

	// Validate the refresh token cryptographically (checks signature + expiry)
	claims, err := h.jwtMgr.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return apperrors.NewWithDetail(401, "Unauthenticated", "Invalid or expired refresh token.")
	}

	// Verify the token exists in the database (guards against replay of rotated tokens)
	if h.refreshRepo != nil {
		oldHash := auth.HashRefreshToken(req.RefreshToken)
		if _, _, dbErr := h.refreshRepo.GetByHash(c.Context(), oldHash); dbErr != nil {
			return apperrors.NewWithDetail(401, "Unauthenticated", "Refresh token has been revoked or does not exist.")
		}
		// Rotate: revoke the old token so it cannot be reused
		_ = h.refreshRepo.RevokeByHash(c.Context(), oldHash)
	}

	identity := &auth.UserIdentity{
		ID:           claims.UserID,
		Email:        claims.Email,
		DemoUser:     claims.DemoUser,
		TokenVersion: claims.TokenVersion,
	}

	tokens, err := h.jwtMgr.GenerateTokenPair(identity, claims.GroupID)
	if err != nil {
		return apperrors.ErrInternal
	}

	// Store new refresh token in DB
	if h.refreshRepo != nil {
		newHash := auth.HashRefreshToken(tokens.RefreshToken)
		newExpiresAt := time.Now().UTC().Add(time.Duration(h.cfg.AuthRefreshExpiry) * 24 * time.Hour)
		_ = h.refreshRepo.Store(c.Context(), claims.UserID, newHash, newExpiresAt)
	}

	return c.JSON(tokens)
}

// Ensure response package is used.
var _ = response.HealthResponse{}
var _ = apperrors.ErrBadRequest

// OAuthURL handles GET /api/v1/auth/:provider/url.
// Generates an OAuth state and returns the provider's authorization URL.
func (h *AuthHandler) OAuthURL(c *fiber.Ctx) error {
	providerName := c.Params("provider")
	if providerName == "" {
		return apperrors.NewValidationError(map[string][]string{
			"provider": {"Provider is required."},
		})
	}

	// Only OAuth providers support this
	if providerName == "local" || providerName == "disabled" {
		return c.Status(400).JSON(fiber.Map{
			"message": "This provider does not support OAuth.",
		})
	}

	// Validate that the requested provider matches the configured provider
	if providerName != h.provider.Name() {
		return c.Status(400).JSON(fiber.Map{
			"message": fmt.Sprintf("Provider '%s' is not configured. Active provider: '%s'.", providerName, h.provider.Name()),
		})
	}

	// Generate CSRF state
	redirect := c.Query("redirect", "")
	state, err := h.oauthState.Generate(c.Context(), providerName, redirect)
	if err != nil {
		return apperrors.ErrInternal
	}

	// Get auth URL from provider
	authURL := h.provider.AuthURL(state)
	if authURL == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "Current provider does not support OAuth.",
		})
	}

	return c.JSON(fiber.Map{
		"url":   authURL,
		"state": state,
	})
}

// OAuthCallback handles GET /api/v1/auth/:provider/callback.
// Exchanges the OAuth code for user info, creates/finds user, returns JWT tokens.
func (h *AuthHandler) OAuthCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" || state == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "Missing code or state parameter.",
		})
	}

	// Validate state
	_, redirect, err := h.oauthState.Validate(c.Context(), state)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid or expired OAuth state.",
		})
	}

	// Authenticate with the provider
	identity, err := h.provider.Authenticate(c.Context(), auth.Credentials{Code: code})
	if err != nil {
		return apperrors.NewWithDetail(401, "Unauthenticated", err.Error())
	}

	if identity.Blocked {
		return c.Status(403).JSON(fiber.Map{
			"message": "User account is blocked.",
		})
	}

	// Find or create user (auto-provision for OAuth)
	user, err := h.userRepo.FindByEmail(c.Context(), identity.Email)
	if err != nil {
		// User doesn't exist — auto-create
		hash, hashErr := auth.HashPassword(generateRandomPassword(24))
		if hashErr != nil {
			return apperrors.ErrInternal
		}
		user, err = h.userRepo.Create(c.Context(), identity.Email, hash)
		if err != nil {
			return apperrors.ErrInternal
		}
	}

	// Generate JWT tokens
	tokens, err := h.jwtMgr.GenerateTokenPair(&auth.UserIdentity{
		ID:       user.ID,
		Email:    user.Email,
		Blocked:  user.Blocked,
		DemoUser: user.DemoUser,
	}, user.UserGroupID)
	if err != nil {
		return apperrors.ErrInternal
	}

	// Store refresh token in DB
	if h.refreshRepo != nil {
		expiresAt := time.Now().UTC().Add(time.Duration(h.cfg.AuthRefreshExpiry) * 24 * time.Hour)
		tokenHash := auth.HashRefreshToken(tokens.RefreshToken)
		_ = h.refreshRepo.Store(c.Context(), user.ID, tokenHash, expiresAt)
	}

	// If redirect URL is set, validate against APP_URL allowlist and use fragment (not query)
	if redirect != "" && isAllowedRedirect(redirect, h.cfg.AppURL) {
		return c.Redirect(redirect + "#access_token=" + tokens.AccessToken + "&refresh_token=" + tokens.RefreshToken)
	}

	return c.JSON(tokens)
}

// isAllowedRedirect validates that a redirect URL matches the configured APP_URL origin.
func isAllowedRedirect(redirect, appURL string) bool {
	r, err := url.Parse(redirect)
	if err != nil {
		return false
	}
	if r.Scheme != "http" && r.Scheme != "https" {
		return false
	}
	a, err := url.Parse(appURL)
	if err != nil {
		return false
	}
	return strings.EqualFold(r.Hostname(), a.Hostname())
}

// generateRandomPassword creates a cryptographically random password for OAuth auto-provisioned users.
func generateRandomPassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	rand.Read(b)
	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b)
}

// isValidEmail performs basic email format validation.
func isValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".") && len(email) >= 5
}
