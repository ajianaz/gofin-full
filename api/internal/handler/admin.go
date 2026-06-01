package handler

import (
	"strings"

	"github.com/rs/zerolog"
	"github.com/gofiber/fiber/v2"

	"github.com/ajianaz/gofin-full/api/internal/auth"
	"github.com/ajianaz/gofin-full/api/internal/config"
	"github.com/ajianaz/gofin-full/api/internal/dto/response"
	"github.com/ajianaz/gofin-full/api/internal/repository"
	"github.com/ajianaz/gofin-full/api/internal/validation"
	apperrors "github.com/ajianaz/gofin-full/api/pkg/errors"
)

type AdminHandler struct {
	log        zerolog.Logger
	cfg        *config.Config
	userRepo   *repository.UserRepository
	configRepo *repository.ConfigurationRepository
}

func NewAdminHandler(log zerolog.Logger, cfg *config.Config, userRepo *repository.UserRepository, configRepo *repository.ConfigurationRepository) *AdminHandler {
	return &AdminHandler{log: log, cfg: cfg, userRepo: userRepo, configRepo: configRepo}
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

// ListUsers returns all users (admin only) with pagination.
func (h *AdminHandler) ListUsers(c *fiber.Ctx) error {
	if err := h.requireAdmin(c); err != nil {
		return err
	}

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("per_page", 20)
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	users, total, err := h.userRepo.ListAllWithRolesPaginated(c.Context(), limit, offset)
	if err != nil {
		h.log.Error().Err(err).Msg("handler/ListUsers: failed to list users")
		return apperrors.ErrInternal
	}

	totalPages := 0
	if total > 0 {
		totalPages = (total + limit - 1) / limit
	}

	data := make([]fiber.Map, 0, len(users))
	for _, uwr := range users {
		u := uwr.User
		data = append(data, fiber.Map{
			"type": "users",
			"id":   u.ID,
			"attributes": fiber.Map{
				"email":      u.Email,
				"name":       u.Email,
				"role":       uwr.Role,
				"is_active":  !u.Blocked,
				"created_at": u.CreatedAt.Format("2006-01-02T15:04:05Z"),
			},
		})
	}

	return c.JSON(fiber.Map{
		"data": data,
		"meta": fiber.Map{
			"current_page": page,
			"per_page":     limit,
			"total":        total,
			"total_pages":  totalPages,
		},
	})
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
		return validation.FieldErrors{"body": {"Invalid request body."}}.ToAppError()
	}

	// Validate input using validation framework
	errs := make(validation.FieldErrors)
	validation.Required("email", req.Email, errs)
	validation.Email("email", req.Email, errs)
	validation.Required("password", req.Password, errs)
	validation.MinLength("password", req.Password, 8, errs)
	if errs.Has() {
		return errs.ToAppError()
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

	// boolStr converts a bool to "enabled" or "disabled".
	boolStr := func(b bool) string {
		if b {
			return "enabled"
		}
		return "disabled"
	}

	flags := map[string]string{
		"webhooks":               boolStr(h.cfg.FeatureWebhooks),
		"csv_import":             "enabled",
		"budgets":                "enabled",
		"piggy_banks":            "enabled",
		"recurring_transactions": "enabled",
		"rules_engine":           "enabled",
		"export_csv":             boolStr(h.cfg.FeatureExport),
		"export_ofx":             boolStr(h.cfg.FeatureExport),
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
		h.log.Error().Err(err).Msg("handler/requireAdmin: failed to set feature flag")
		return apperrors.ErrInternal
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
