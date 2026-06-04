package handler

import (
	"context"
	"crypto/rand"
	"fmt"
	"os"
	"github.com/rs/zerolog"
	"net/mail"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"

	"github.com/ajianaz/gofin-full/api/internal/auth"
	"github.com/ajianaz/gofin-full/api/internal/config"
	"github.com/ajianaz/gofin-full/api/internal/dto/response"
	"github.com/ajianaz/gofin-full/api/internal/repository"
	apperrors "github.com/ajianaz/gofin-full/api/pkg/errors"
)

// Login attempt lockout defaults (overridden by config).
const (
	loginAttemptsKeyPrefix = "login_attempts:"
	defaultLoginMaxAttempts = 5
	defaultLoginLockoutMinutes = 15
)

// loginMaxAttempts returns the configured max attempts, falling back to default.
func (h *AuthHandler) loginMaxAttempts() int {
	if h.cfg.LoginMaxAttempts > 0 {
		return h.cfg.LoginMaxAttempts
	}
	return defaultLoginMaxAttempts
}

// loginLockoutDuration returns the configured lockout duration, falling back to default.
func (h *AuthHandler) loginLockoutDuration() time.Duration {
	if h.cfg.LoginLockoutMinutes > 0 {
		return time.Duration(h.cfg.LoginLockoutMinutes) * time.Minute
	}
	return time.Duration(defaultLoginLockoutMinutes) * time.Minute
}

// passwordResetKeyPrefix is the Redis key prefix for password reset tokens.
const passwordResetKeyPrefix = "password_reset:"

// emailVerifyKeyPrefix is the Redis key prefix for email verification tokens.
const emailVerifyKeyPrefix = "email_verify:"

// AuthHandler handles authentication endpoints.
type AuthHandler struct {
	log         zerolog.Logger
	jwtMgr      *auth.JWTManager
	provider    auth.AuthProvider
	cfg         *config.Config
	userRepo    *repository.UserRepository
	oauthState  *repository.OAuthStateRepository
	refreshRepo *repository.RefreshTokenRepository
	rdb         redis.Cmdable // optional Redis for login attempt tracking
	mail        MailSender    // optional mail service for password reset
}

// MailSender is an interface for sending emails.
type MailSender interface {
	Configured() bool
	SendEmail(to, subject, body string) error
}

// NewAuthHandler creates a new auth handler.
func NewAuthHandler(log zerolog.Logger, jwtMgr *auth.JWTManager, provider auth.AuthProvider, cfg *config.Config, userRepo *repository.UserRepository, oauthStateRepo *repository.OAuthStateRepository, refreshRepo *repository.RefreshTokenRepository) *AuthHandler {
	return &AuthHandler{log: log, jwtMgr: jwtMgr, provider: provider, cfg: cfg, userRepo: userRepo, oauthState: oauthStateRepo, refreshRepo: refreshRepo}
}

// SetRedis injects an optional Redis client for login attempt tracking.
func (h *AuthHandler) SetRedis(rdb redis.Cmdable) {
	h.rdb = rdb
}

// SetMail injects an optional mail service for password reset.
func (h *AuthHandler) SetMail(mail MailSender) {
	h.mail = mail
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

	// Validate required fields before authentication
	if strings.TrimSpace(req.Email) == "" {
		return apperrors.NewValidationError(map[string][]string{
			"email": {"Email is required."},
		})
	}
	if strings.TrimSpace(req.Password) == "" {
		return apperrors.NewValidationError(map[string][]string{
			"password": {"Password is required."},
		})
	}

	clientIP := c.IP()

	// Check if account is temporarily locked due to too many failed attempts
	if h.cfg.LoginRateLimitEnabled {
		if locked, retryMinutes := h.isAccountLocked(c.Context(), req.Email, clientIP); locked {
			return response.SendError(c, 429, fmt.Sprintf("Too many failed login attempts. Try again in %d minutes.", retryMinutes))
		}
	}

	identity, err := h.provider.Authenticate(c.Context(), auth.Credentials{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if h.cfg.LoginRateLimitEnabled {
			h.recordFailedLogin(c.Context(), req.Email, clientIP)
		}
		return apperrors.New(401, "Invalid email or password.")
	}

	if identity.Blocked {
		return apperrors.NewWithDetail(403, "Forbidden", "User account is blocked.")
	}

	// Check email verification if required
	if h.cfg.AuthRequireVerification && !identity.Verified {
		return response.SendError(c, 403, "Please verify your email address.")
	}

	// Successful login — clear failed attempt counter
	if h.cfg.LoginRateLimitEnabled {
		h.clearFailedLogins(c.Context(), req.Email, clientIP)
	}

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

	// Set httpOnly cookies (secure in production, insecure for local HTTP dev)
	secure := !h.cfg.IsLocal()
	auth.SetTokenCookies(c, tokens.AccessToken, tokens.RefreshToken, secure)

	return c.JSON(tokens.PublicResponse())
}

// loginAttemptEntry tracks in-memory login attempt state per key.
type loginAttemptEntry struct {
	timestamps []int64
	mu         sync.Mutex
}

// loginAttemptStore is a global in-memory store for login attempt tracking,
// used as fallback when Redis is unavailable.
var loginAttemptStore sync.Map

// init starts a goroutine to periodically evict stale login attempt entries.
func init() {
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			now := time.Now().UnixMilli()
			windowStart := now - int64(defaultLoginLockoutMinutes)*60*1000
			loginAttemptStore.Range(func(key, val interface{}) bool {
				entry, ok := val.(*loginAttemptEntry)
				if !ok {
					return true
				}
				entry.mu.Lock()
				valid := entry.timestamps[:0]
				for _, ts := range entry.timestamps {
					if ts > windowStart {
						valid = append(valid, ts)
					}
				}
				stale := len(valid) == 0 && len(entry.timestamps) > 0
				entry.timestamps = valid
				entry.mu.Unlock()
				if stale {
					loginAttemptStore.Delete(key)
				}
				return true
			})
		}
	}()
}

// memLoginAttemptCheck performs an in-memory sliding window login attempt check.
// Returns (locked, retryMinutes).
func memLoginAttemptCheck(key string) (bool, int) {
	now := time.Now().UnixMilli()
	windowStart := now - int64(defaultLoginLockoutMinutes)*60*1000

	val, _ := loginAttemptStore.LoadOrStore(key, &loginAttemptEntry{})
	entry, ok := val.(*loginAttemptEntry)
	if !ok {
		fmt.Fprintf(os.Stderr, "warning: unexpected type in loginAttemptStore for key %s: %T\n", key, val)
		return false, 0
	}

	entry.mu.Lock()
	defer entry.mu.Unlock()

	// Remove expired timestamps
	valid := entry.timestamps[:0]
	for _, ts := range entry.timestamps {
		if ts > windowStart {
			valid = append(valid, ts)
		}
	}
	entry.timestamps = valid

	if len(entry.timestamps) >= defaultLoginMaxAttempts {
		// Calculate retry from oldest attempt in window
		oldest := entry.timestamps[0]
		retryMs := (oldest + int64(defaultLoginLockoutMinutes)*60*1000) - now
		retryMinutes := int(retryMs / 60000)
		if retryMinutes < 1 {
			retryMinutes = 1
		}
		return true, retryMinutes
	}

	entry.timestamps = append(entry.timestamps, now)
	return false, 0
}

// memClearLoginAttempts clears the in-memory login attempt counter for a key.
func memClearLoginAttempts(key string) {
	loginAttemptStore.Delete(key)
}

// isAccountLocked checks if an email+IP combination is locked due to too many failed login attempts.
// Uses email+IP as key to prevent account lockout DoS attacks.
// Falls back to in-memory tracking when Redis is unavailable.
func (h *AuthHandler) isAccountLocked(ctx context.Context, email, clientIP string) (bool, int) {
	key := loginAttemptsKeyPrefix + email + ":" + clientIP

	if h.rdb == nil {
		// In-memory fallback when Redis is unavailable
		locked, retryMinutes := memLoginAttemptCheck(key)
		return locked, retryMinutes
	}

	val, err := h.rdb.Get(ctx, key).Result()
	if err != nil {
		return false, 0
	}

	attempts, err := strconv.Atoi(val)
	if err != nil || attempts < h.loginMaxAttempts() {
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

// recordFailedLogin increments the failed login counter for an email+IP combination.
func (h *AuthHandler) recordFailedLogin(ctx context.Context, email, clientIP string) {
	key := loginAttemptsKeyPrefix + email + ":" + clientIP

	if h.rdb == nil {
		// In-memory fallback — attempt was already recorded by isAccountLocked
		return
	}

	pipe := h.rdb.Pipeline()
	pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, h.loginLockoutDuration())
	if _, err := pipe.Exec(ctx); err != nil {
		h.log.Error().Err(err).Msg("login attempt tracking: redis error")
	}
}

// clearFailedLogins removes the failed login counter for an email+IP after successful login.
func (h *AuthHandler) clearFailedLogins(ctx context.Context, email, clientIP string) {
	key := loginAttemptsKeyPrefix + email + ":" + clientIP

	if h.rdb == nil {
		memClearLoginAttempts(key)
		return
	}

	if err := h.rdb.Del(ctx, key).Err(); err != nil {
		h.log.Error().Err(err).Msg("login attempt tracking: redis error on clear")
	}
}

// Register handles POST /api/v1/auth/register.
// Only works when AUTH_ALLOW_REGISTRATION=true.
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	if !h.cfg.AuthAllowRegistration {
		return response.SendError(c, 403, "Self-registration is disabled. Contact an administrator to create an account.")
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
	if pwErrs := validatePasswordStrength(req.Password); len(pwErrs) > 0 {
		return apperrors.NewValidationError(map[string][]string{
			"password": pwErrs,
		})
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		return apperrors.ErrInternal
	}

	user, err := h.userRepo.Create(c.Context(), req.Email, hash)
	if err != nil {
		if isDuplicateKey(err) {
			return response.SendError(c, 409, "A user with this email already exists.")
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

	// Email verification: if SMTP is configured, send verification email;
	// otherwise auto-verify the user.
	if h.mail != nil && h.mail.Configured() && h.rdb != nil {
		if err := h.sendVerificationEmail(c.Context(), user.Email); err != nil {
			h.log.Error().Err(err).Msg("failed to send verification email")
			// Non-fatal: still return tokens, user can verify later
		}
	} else {
		// No SMTP configured — auto-verify
		_ = h.userRepo.SetVerified(c.Context(), user.ID)
	}

	// Set httpOnly cookies (secure in production, insecure for local HTTP dev)
	secure := !h.cfg.IsLocal()
	auth.SetTokenCookies(c, tokens.AccessToken, tokens.RefreshToken, secure)

	return c.Status(201).JSON(tokens.PublicResponse())
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
// No request body is required — the user is identified from the JWT in httpOnly cookies.
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// Revoke the refresh token from cookie if present
	if refreshToken := c.Cookies(auth.RefreshTokenCookieName); refreshToken != "" && h.refreshRepo != nil {
		tokenHash := auth.HashRefreshToken(refreshToken)
		_ = h.refreshRepo.RevokeByHash(c.Context(), tokenHash)
	}

	// Increment token_version to invalidate all existing JWT tokens for this user
	user := auth.GetUser(c)
	if user != nil {
		_ = h.userRepo.IncrementTokenVersion(c.Context(), user.ID)
		auth.InvalidateTokenCache(c, user.ID)
	}

	// Clear httpOnly cookies
	auth.ClearTokenCookies(c, !h.cfg.IsLocal())

	return c.JSON(fiber.Map{
		"message": "Logged out successfully.",
	})
}

// Refresh handles POST /api/v1/auth/refresh.
// Requires a valid access token (for authentication) and a refresh_token in the body.
// The refresh token is validated cryptographically and rotated (old token is revoked).
func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	// Read refresh token: try request body first, then httpOnly cookie (migration support)
	var refreshToken string

	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.BodyParser(&req); err == nil && req.RefreshToken != "" {
		refreshToken = req.RefreshToken
	}
	if refreshToken == "" {
		refreshToken = c.Cookies(auth.RefreshTokenCookieName)
	}
	if refreshToken == "" {
		return apperrors.NewValidationError(map[string][]string{
			"refresh_token": {"Refresh token is required (in body or httpOnly cookie)."},
		})
	}

	// Validate the refresh token cryptographically (checks signature + expiry)
	claims, err := h.jwtMgr.ValidateRefreshToken(refreshToken)
	if err != nil {
		return apperrors.NewWithDetail(401, "Unauthenticated", "Invalid or expired refresh token.")
	}

	// Verify the token exists in the database (guards against replay of rotated tokens)
	if h.refreshRepo != nil {
		oldHash := auth.HashRefreshToken(refreshToken)
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

	// Set httpOnly cookies (secure in production, insecure for local HTTP dev)
	secure := !h.cfg.IsLocal()
	auth.SetTokenCookies(c, tokens.AccessToken, tokens.RefreshToken, secure)

	return c.JSON(tokens.PublicResponse())
}

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
		return response.SendError(c, 400, "This provider does not support OAuth.")
	}

	// Validate that the requested provider matches the configured provider
	if providerName != h.provider.Name() {
		return response.SendError(c, 400, fmt.Sprintf("Provider '%s' is not configured. Active provider: '%s'.", providerName, h.provider.Name()))
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
		return response.SendError(c, 400, "Current provider does not support OAuth.")
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
		return response.SendError(c, 400, "Missing code or state parameter.")
	}

	// Validate state
	_, redirect, err := h.oauthState.Validate(c.Context(), state)
	if err != nil {
		return response.SendError(c, 400, "Invalid or expired OAuth state.")
	}

	// Authenticate with the provider
	identity, err := h.provider.Authenticate(c.Context(), auth.Credentials{Code: code})
	if err != nil {
		h.log.Error().Err(err).Msg("OAuth callback authentication failed")
		return apperrors.New(401, "Authentication failed.")
	}

	if identity.Blocked {
		return response.SendError(c, 403, "User account is blocked.")
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

	// Set httpOnly cookies (secure in production, insecure for local HTTP dev)
	secure := !h.cfg.IsLocal()
	auth.SetTokenCookies(c, tokens.AccessToken, tokens.RefreshToken, secure)

	// If redirect URL is set, validate against APP_URL allowlist and redirect.
	// Tokens are no longer included in the URL fragment — httpOnly cookies are set above.
	if redirect != "" && isAllowedRedirect(redirect, h.cfg.AppURL) {
		return c.Redirect(redirect)
	}

	return c.JSON(tokens.PublicResponse())
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
	for {
		if _, err := rand.Read(b); err != nil {
			// Fallback: use time-based (not ideal but won't panic)
			for i := range b {
				b[i] = byte(time.Now().UnixNano())
			}
			break
		}
		// Reject bytes that would cause modulo bias
		valid := true
		for i := range b {
			if int(b[i]) >= 256-(256%len(charset)) {
				valid = false
				break
			}
		}
		if valid {
			break
		}
	}
	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}
	return string(b)
}

// isValidEmail validates email format using the standard net/mail parser.
func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// validatePasswordStrength checks password meets minimum security requirements.
// Requires at least 8 characters with a mix of character types.
func validatePasswordStrength(password string) []string {
	var errors []string
	if len(password) < 8 {
		errors = append(errors, "Password must be at least 8 characters.")
	}
	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false
	for _, c := range password {
		switch {
		case 'A' <= c && c <= 'Z':
			hasUpper = true
		case 'a' <= c && c <= 'z':
			hasLower = true
		case '0' <= c && c <= '9':
			hasDigit = true
		default:
			hasSpecial = true
		}
	}
	// Require at least 3 of 4 character types
	types := 0
	if hasUpper {
		types++
	}
	if hasLower {
		types++
	}
	if hasDigit {
		types++
	}
	if hasSpecial {
		types++
	}
	if types < 3 {
		errors = append(errors, "Password must include at least 3 of: uppercase, lowercase, digit, special character.")
	}
	return errors
}

// generateResetToken creates a cryptographically random 32-byte hex token.
func generateResetToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate reset token: %w", err)
	}
	return fmt.Sprintf("%x", b), nil
}

// ForgotPassword handles POST /api/v1/auth/forgot-password.
// Always returns 200 to prevent email enumeration.
func (h *AuthHandler) ForgotPassword(c *fiber.Ctx) error {
	// Check if SMTP is configured
	if h.mail == nil || !h.mail.Configured() {
		return response.SendError(c, 503, "Password reset is not configured.")
	}

	// Check if Redis is available
	if h.rdb == nil {
		return response.SendError(c, 503, "Password reset is not available.")
	}

	var req struct {
		Email string `json:"email"`
	}
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewValidationError(map[string][]string{
			"body": {"Invalid request body."},
		})
	}

	if strings.TrimSpace(req.Email) == "" {
		return apperrors.NewValidationError(map[string][]string{
			"email": {"Email is required."},
		})
	}

	if !isValidEmail(req.Email) {
		return apperrors.NewValidationError(map[string][]string{
			"email": {"Invalid email format."},
		})
	}

	// Check if user exists — but don't reveal the result
	user, err := h.userRepo.FindByEmail(c.Context(), req.Email)
	if err != nil || user == nil {
		// User not found — return success anyway to prevent enumeration
		return c.JSON(fiber.Map{
			"message": "If an account with this email exists, a password reset link has been sent.",
		})
	}

	// Generate reset token
	token, err := generateResetToken()
	if err != nil {
		h.log.Error().Err(err).Msg("failed to generate reset token")
		return apperrors.ErrInternal
	}

	// Store token in Redis with 1 hour TTL
	key := passwordResetKeyPrefix + token
	if err := h.rdb.Set(c.Context(), key, req.Email, 1*time.Hour).Err(); err != nil {
		h.log.Error().Err(err).Msg("failed to store reset token in redis")
		return apperrors.ErrInternal
	}

	// Build reset link
	appURL := h.cfg.AppURL
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", strings.TrimSuffix(appURL, "/"), token)

	// Send email
	subject := "Password Reset Request"
	body := fmt.Sprintf(
		"Hello,\n\nYou have requested a password reset for your GoFin account.\n\n"+
			"Click the link below to reset your password:\n%s\n\n"+
			"This link will expire in 1 hour.\n\n"+
			"If you did not request this, you can safely ignore this email.\n",
		resetLink,
	)

	if err := h.mail.SendEmail(req.Email, subject, body); err != nil {
		h.log.Error().Err(err).Msg("failed to send reset email")
		// Still return success to prevent enumeration
	}

	return c.JSON(fiber.Map{
		"message": "If an account with this email exists, a password reset link has been sent.",
	})
}

// ResetPassword handles POST /api/v1/auth/reset-password.
func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	// Check if Redis is available
	if h.rdb == nil {
		return response.SendError(c, 503, "Password reset is not available.")
	}

	var req struct {
		Token          string `json:"token"`
		Password       string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewValidationError(map[string][]string{
			"body": {"Invalid request body."},
		})
	}

	if strings.TrimSpace(req.Token) == "" {
		return apperrors.NewValidationError(map[string][]string{
			"token": {"Reset token is required."},
		})
	}

	if strings.TrimSpace(req.Password) == "" {
		return apperrors.NewValidationError(map[string][]string{
			"password": {"Password is required."},
		})
	}

	if len(req.Password) < 8 {
		return apperrors.NewValidationError(map[string][]string{
			"password": {"Password must be at least 8 characters."},
		})
	}

	if req.Password != req.ConfirmPassword {
		return apperrors.NewValidationError(map[string][]string{
			"confirm_password": {"Passwords do not match."},
		})
	}

	// Look up token in Redis
	key := passwordResetKeyPrefix + req.Token
	email, err := h.rdb.Get(c.Context(), key).Result()
	if err != nil {
		return response.SendError(c, 400, "Invalid or expired reset link.")
	}

	// Find user by email
	user, err := h.userRepo.FindByEmail(c.Context(), email)
	if err != nil || user == nil {
		return response.SendError(c, 400, "Invalid or expired reset link.")
	}

	// Hash new password
	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		return apperrors.ErrInternal
	}

	// Update user password in DB
	if err := h.userRepo.UpdatePassword(c.Context(), user.ID, hash); err != nil {
		h.log.Error().Err(err).Msg("failed to update password")
		return apperrors.ErrInternal
	}

	// Delete token from Redis (single use)
	_ = h.rdb.Del(c.Context(), key)

	// Invalidate all existing JWT tokens
	_ = h.userRepo.IncrementTokenVersion(c.Context(), user.ID)

	return c.JSON(fiber.Map{
		"message": "Password has been reset successfully.",
	})
}

// VerifyEmail handles POST /api/v1/auth/verify-email.
func (h *AuthHandler) VerifyEmail(c *fiber.Ctx) error {
	if h.rdb == nil {
		return response.SendError(c, 503, "Email verification is not available.")
	}

	var req struct {
		Token string `json:"token"`
	}
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewValidationError(map[string][]string{
			"body": {"Invalid request body."},
		})
	}

	if strings.TrimSpace(req.Token) == "" {
		return apperrors.NewValidationError(map[string][]string{
			"token": {"Verification token is required."},
		})
	}

	// Look up token in Redis
	key := emailVerifyKeyPrefix + req.Token
	email, err := h.rdb.Get(c.Context(), key).Result()
	if err != nil {
		return response.SendError(c, 400, "Invalid or expired verification link.")
	}

	// Find user by email
	user, err := h.userRepo.FindByEmail(c.Context(), email)
	if err != nil || user == nil {
		return response.SendError(c, 400, "Invalid or expired verification link.")
	}

	// Set verified = true
	if err := h.userRepo.SetVerified(c.Context(), user.ID); err != nil {
		h.log.Error().Err(err).Msg("failed to set user verified")
		return apperrors.ErrInternal
	}

	// Delete token from Redis (single use)
	_ = h.rdb.Del(c.Context(), key)

	return c.JSON(fiber.Map{
		"message": "Email verified successfully. You can now log in.",
	})
}

// ResendVerification handles POST /api/v1/auth/resend-verification.
// Requires authentication.
func (h *AuthHandler) ResendVerification(c *fiber.Ctx) error {
	user := auth.GetUser(c)
	if user == nil {
		return apperrors.ErrUnauthorized
	}

	// Check if already verified
	dbUser, err := h.userRepo.FindByEmail(c.Context(), user.Email)
	if err != nil || dbUser == nil {
		return apperrors.ErrInternal
	}
	if dbUser.Verified {
		return response.SendError(c, 400, "Email is already verified.")
	}

	if h.mail == nil || !h.mail.Configured() || h.rdb == nil {
		return response.SendError(c, 503, "Email verification is not configured.")
	}

	if err := h.sendVerificationEmail(c.Context(), user.Email); err != nil {
		h.log.Error().Err(err).Msg("failed to resend verification email")
		return apperrors.ErrInternal
	}

	return c.JSON(fiber.Map{
		"message": "Verification email has been sent.",
	})
}

// sendVerificationEmail generates a token, stores it in Redis, and sends the verification email.
func (h *AuthHandler) sendVerificationEmail(ctx context.Context, email string) error {
	token, err := generateResetToken()
	if err != nil {
		return fmt.Errorf("failed to generate verification token: %w", err)
	}

	// Store token in Redis with 24h TTL
	key := emailVerifyKeyPrefix + token
	if err := h.rdb.Set(ctx, key, email, 24*time.Hour).Err(); err != nil {
		return fmt.Errorf("failed to store verification token: %w", err)
	}

	// Build verification link
	appURL := h.cfg.AppURL
	verifyLink := fmt.Sprintf("%s/verify-email?token=%s", strings.TrimSuffix(appURL, "/"), token)

	// Send email
	subject := "Verify Your Email Address"
	body := fmt.Sprintf(
		"Hello,\n\nPlease verify your email address for your GoFin account.\n\n"+
			"Click the link below to verify:\n%s\n\n"+
			"This link will expire in 24 hours.\n\n"+
			"If you did not create an account, you can safely ignore this email.\n",
		verifyLink,
	)

	return h.mail.SendEmail(email, subject, body)
}
