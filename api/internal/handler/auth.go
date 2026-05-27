package handler

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
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
func NewAuthHandler(jwtMgr *auth.JWTManager, provider auth.AuthProvider, cfg *config.Config, userRepo *repository.UserRepository, oauthStateRepo *repository.OAuthStateRepository, refreshRepo *repository.RefreshTokenRepository) *AuthHandler {
	return &AuthHandler{jwtMgr: jwtMgr, provider: provider, cfg: cfg, userRepo: userRepo, oauthState: oauthStateRepo, refreshRepo: refreshRepo}
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
			return c.Status(429).JSON(fiber.Map{
				"message": fmt.Sprintf("Too many failed login attempts. Try again in %d minutes.", retryMinutes),
			})
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
		return c.Status(403).JSON(fiber.Map{
			"message":  "Please verify your email address.",
			"verified": false,
		})
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

	return c.JSON(tokens)
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
				entry := val.(*loginAttemptEntry)
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
	entry := val.(*loginAttemptEntry)

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
		log.Printf("login attempt tracking: redis error: %v", err)
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

	// Email verification: if SMTP is configured, send verification email;
	// otherwise auto-verify the user.
	if h.mail != nil && h.mail.Configured() && h.rdb != nil {
		if err := h.sendVerificationEmail(c.Context(), user.Email); err != nil {
			log.Printf("failed to send verification email: %v", err)
			// Non-fatal: still return tokens, user can verify later
		}
	} else {
		// No SMTP configured — auto-verify
		_ = h.userRepo.SetVerified(c.Context(), user.ID)
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
// Requires a valid access token (for authentication) and a refresh_token in the body.
// The refresh token is validated cryptographically and rotated (old token is revoked).
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
		log.Printf("OAuth callback authentication failed: %v", err)
		return apperrors.New(401, "Authentication failed.")
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
		return c.Status(503).JSON(fiber.Map{
			"message": "Password reset is not configured.",
		})
	}

	// Check if Redis is available
	if h.rdb == nil {
		return c.Status(503).JSON(fiber.Map{
			"message": "Password reset is not available.",
		})
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
		log.Printf("failed to generate reset token: %v", err)
		return apperrors.ErrInternal
	}

	// Store token in Redis with 1 hour TTL
	key := passwordResetKeyPrefix + token
	if err := h.rdb.Set(c.Context(), key, req.Email, 1*time.Hour).Err(); err != nil {
		log.Printf("failed to store reset token in redis: %v", err)
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
		log.Printf("failed to send reset email: %v", err)
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
		return c.Status(503).JSON(fiber.Map{
			"message": "Password reset is not available.",
		})
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
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid or expired reset link.",
		})
	}

	// Find user by email
	user, err := h.userRepo.FindByEmail(c.Context(), email)
	if err != nil || user == nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid or expired reset link.",
		})
	}

	// Hash new password
	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		return apperrors.ErrInternal
	}

	// Update user password in DB
	if err := h.userRepo.UpdatePassword(c.Context(), user.ID, hash); err != nil {
		log.Printf("failed to update password: %v", err)
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
		return c.Status(503).JSON(fiber.Map{
			"message": "Email verification is not available.",
		})
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
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid or expired verification link.",
		})
	}

	// Find user by email
	user, err := h.userRepo.FindByEmail(c.Context(), email)
	if err != nil || user == nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid or expired verification link.",
		})
	}

	// Set verified = true
	if err := h.userRepo.SetVerified(c.Context(), user.ID); err != nil {
		log.Printf("failed to set user verified: %v", err)
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
		return c.Status(400).JSON(fiber.Map{
			"message": "Email is already verified.",
		})
	}

	if h.mail == nil || !h.mail.Configured() || h.rdb == nil {
		return c.Status(503).JSON(fiber.Map{
			"message": "Email verification is not configured.",
		})
	}

	if err := h.sendVerificationEmail(c.Context(), user.Email); err != nil {
		log.Printf("failed to resend verification email: %v", err)
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
