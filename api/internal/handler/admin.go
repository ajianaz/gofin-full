package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/ajianaz/gofin-full/api/internal/auth"
	"github.com/ajianaz/gofin-full/api/internal/repository"
	apperrors "github.com/ajianaz/gofin-full/api/pkg/errors"
)

type AdminHandler struct {
	userRepo   *repository.UserRepository
	configRepo *repository.ConfigurationRepository
}

func NewAdminHandler(userRepo *repository.UserRepository, configRepo *repository.ConfigurationRepository) *AdminHandler {
	return &AdminHandler{userRepo: userRepo, configRepo: configRepo}
}

// requireAdmin checks that the caller has a global admin or owner role.
func (h *AdminHandler) requireAdmin(c *fiber.Ctx) error {
	claims := auth.GetClaims(c)
	if claims == nil {
		return apperrors.ErrUnauthorized
	}

	// JWT-level role check: the claims should carry admin role info
	// Also verify against DB for defense in depth
	hasRole, err := h.userRepo.HasGlobalRole(c.Context(), claims.UserID, "owner")
	if err == nil && hasRole {
		return nil
	}
	hasRole, err = h.userRepo.HasGlobalRole(c.Context(), claims.UserID, "admin")
	if err == nil && hasRole {
		return nil
	}

	return apperrors.New(403, "Insufficient permissions. Admin access required.")
}

// ListUsers returns all users (admin only).
func (h *AdminHandler) ListUsers(c *fiber.Ctx) error {
	if err := h.requireAdmin(c); err != nil {
		return err
	}

	users, err := h.userRepo.ListAll(c.Context())
	if err != nil {
		return apperrors.NewWithDetail(500, "failed to list users", err.Error())
	}

	var data []fiber.Map
	for _, u := range users {
		role := h.userRepo.GetGlobalRole(c.Context(), u.ID)
		data = append(data, fiber.Map{
			"type": "users",
			"id":   u.ID,
			"attributes": fiber.Map{
				"email":      u.Email,
				"name":       u.Email,
				"role":       role,
				"is_active":  !u.Blocked,
				"created_at": u.CreatedAt.Format("2006-01-02T15:04:05Z"),
			},
		})
	}
	return c.JSON(fiber.Map{"data": data})
}

// CreateUser handles POST /api/v1/admin/users.
func (h *AdminHandler) CreateUser(c *fiber.Ctx) error {
	if err := h.requireAdmin(c); err != nil {
		return err
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

	return c.Status(201).JSON(fiber.Map{
		"data": fiber.Map{
			"type": "users",
			"id":   user.ID,
			"attributes": fiber.Map{
				"email": user.Email,
			},
		},
	})
}

// FeatureFlags returns all system feature flags.
func (h *AdminHandler) FeatureFlags(c *fiber.Ctx) error {
	if err := h.requireAdmin(c); err != nil {
		return err
	}

	flags := map[string]string{
		"two_factor_auth":        "disabled",
		"webhooks":               "enabled",
		"csv_import":             "enabled",
		"budgets":                "enabled",
		"piggy_banks":            "enabled",
		"recurring_transactions": "enabled",
		"rules_engine":           "enabled",
		"export_csv":             "enabled",
		"export_ofx":             "enabled",
		"audit_trail":            "enabled",
		"wallet_sharing":         "enabled",
	}

	// Override from DB configs
	for key := range flags {
		cfg, err := h.configRepo.Get(c.Context(), "feature_"+key)
		if err == nil && cfg.Value != "" {
			flags[key] = cfg.Value
		}
	}

	return c.JSON(fiber.Map{"data": fiber.Map{
		"type":       "feature_flags",
		"attributes": flags,
	}})
}

// SetFeatureFlag updates a system feature flag.
func (h *AdminHandler) SetFeatureFlag(c *fiber.Ctx) error {
	if err := h.requireAdmin(c); err != nil {
		return err
	}

	var req struct {
		Flag  string `json:"flag"`
		Value string `json:"value"`
	}
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewValidationError(map[string][]string{"body": {"invalid JSON"}})
	}
	if req.Flag == "" {
		return apperrors.NewValidationError(map[string][]string{"flag": {"flag is required"}})
	}

	_, err := h.configRepo.Set(c.Context(), "feature_"+req.Flag, req.Value)
	if err != nil {
		return apperrors.NewWithDetail(500, "failed to set feature flag", err.Error())
	}

	return c.JSON(fiber.Map{"data": fiber.Map{
		"type": "feature_flags",
		"attributes": fiber.Map{
			"flag":  req.Flag,
			"value": req.Value,
		},
	}})
}

// isDuplicateKey checks if the error is a PostgreSQL unique constraint violation.
func isDuplicateKey(err error) bool {
	if pgErr, ok := err.(interface{ SQLState() string }); ok {
		return pgErr.SQLState() == "23505"
	}
	return strings.Contains(err.Error(), "duplicate key")
}
