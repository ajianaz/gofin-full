package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/ajianaz/gofin-full/api/internal/auth"
	"github.com/ajianaz/gofin-full/api/internal/domain"
	"github.com/ajianaz/gofin-full/api/internal/repository"
	apperrors "github.com/ajianaz/gofin-full/api/pkg/errors"
)

// WalletHandler handles wallet endpoints.
type WalletHandler struct {
	repo *repository.WalletRepository
}

// NewWalletHandler creates a new wallet handler.
func NewWalletHandler(repo *repository.WalletRepository) *WalletHandler {
	return &WalletHandler{repo: repo}
}

// Index handles GET /api/v1/wallets.
func (h *WalletHandler) Index(c *fiber.Ctx) error {
	user := auth.GetUser(c)
	if user == nil {
		return apperrors.ErrUnauthorized
	}

	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.NewWithDetail(400, "Bad Request", "No active group. Switch to a group first.")
	}

	walletType := c.Query("type")
	activeOnly := c.QueryBool("active", true)

	wallets, err := h.repo.List(c.Context(), *groupID, walletType, activeOnly)
	if err != nil {
		return apperrors.ErrInternal
	}

	var data []fiber.Map
	for _, w := range wallets {
		data = append(data, walletToMap(&w))
	}

	return c.JSON(fiber.Map{"data": data})
}

// Show handles GET /api/v1/wallets/:id.
func (h *WalletHandler) Show(c *fiber.Ctx) error {
	user := auth.GetUser(c)
	if user == nil {
		return apperrors.ErrUnauthorized
	}

	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.NewWithDetail(400, "Bad Request", "No active group.")
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperrors.ErrBadRequest
	}

	wallet, err := h.repo.FindByID(c.Context(), id, *groupID)
	if err != nil {
		return apperrors.NotFoundResource("wallet", id)
	}

	return c.JSON(fiber.Map{"data": walletToMap(wallet)})
}

// Store handles POST /api/v1/wallets.
func (h *WalletHandler) Store(c *fiber.Ctx) error {
	user := auth.GetUser(c)
	if user == nil {
		return apperrors.ErrUnauthorized
	}

	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.NewWithDetail(400, "Bad Request", "No active group.")
	}

	var req struct {
		Name        string   `json:"name"`
		WalletType string   `json:"wallet_type"`
		CurrencyID  string   `json:"currency_id"`
		Active      *bool    `json:"active"`
		IBAN        string   `json:"iban"`
		BIC         string   `json:"bic"`
		Notes       string   `json:"notes"`
	}
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewValidationError(map[string][]string{
			"body": {"Invalid request body."},
		})
	}

	// Validate wallet type
	if req.WalletType == "" {
		req.WalletType = "asset"
	}
	if !isValidWalletType(domain.WalletType(req.WalletType)) {
		return apperrors.NewValidationError(map[string][]string{
				 "wallet_type": {"Invalid wallet type."},
		})
	}

	active := true
	if req.Active != nil {
		active = *req.Active
	}

	wallet := &domain.Wallet{
		UserID:          user.ID,
		UserGroupID:     *groupID,
		Name:            req.Name,
		AccountType:     req.WalletType,
		Active:          active,
		VirtualBalance:  decimal.Zero,
		IncludeNetWorth: true,
	}
	if req.IBAN != "" { wallet.IBAN = &req.IBAN }
	if req.BIC != "" { wallet.BIC = &req.BIC }
	if req.CurrencyID != "" { wallet.CurrencyID = &req.CurrencyID }
	if req.Notes != "" { wallet.Notes = &req.Notes }

	created, err := h.repo.Create(c.Context(), wallet)
	if err != nil {
		return apperrors.ErrInternal
	}

	return c.Status(201).JSON(fiber.Map{"data": walletToMap(created)})
}

// Update handles PUT /api/v1/wallets/:id.
func (h *WalletHandler) Update(c *fiber.Ctx) error {
	user := auth.GetUser(c)
	if user == nil {
		return apperrors.ErrUnauthorized
	}

	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.NewWithDetail(400, "Bad Request", "No active group.")
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperrors.ErrBadRequest
	}

	var req struct {
		Name            string  `json:"name"`
		Active          *bool   `json:"active"`
		IncludeNetWorth *bool   `json:"include_net_worth"`
		Notes           *string `json:"notes"`
	}
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewValidationError(map[string][]string{
			"body": {"Invalid request body."},
		})
	}

	if err := h.repo.Update(c.Context(), id, *groupID, req.Name, req.Active, req.IncludeNetWorth, req.Notes); err != nil {
		return apperrors.ErrInternal
	}

	updated, err := h.repo.FindByID(c.Context(), id, *groupID)
	if err != nil {
		return apperrors.ErrInternal
	}

	return c.JSON(fiber.Map{"data": walletToMap(updated)})
}

// Delete handles DELETE /api/v1/wallets/:id.
func (h *WalletHandler) Delete(c *fiber.Ctx) error {
	user := auth.GetUser(c)
	if user == nil {
		return apperrors.ErrUnauthorized
	}

	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.NewWithDetail(400, "Bad Request", "No active group.")
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperrors.ErrBadRequest
	}

	if err := h.repo.Delete(c.Context(), id, *groupID); err != nil {
		return apperrors.ErrInternal
	}

	return c.Status(204).Send(nil)
}

func walletToMap(w *domain.Wallet) fiber.Map {
	m := fiber.Map{
		"type":       "wallets",
		"id":         w.ID,
		"attributes": fiber.Map{
			"name":              w.Name,
				 "wallet_type":      w.AccountType,
			"active":            w.Active,
			"virtual_balance":   w.VirtualBalance.StringFixed(2),
			"include_net_worth": w.IncludeNetWorth,
			"created_at":        w.CreatedAt,
			"updated_at":        w.UpdatedAt,
		},
	}
	if w.IBAN != nil { m["attributes"].(fiber.Map)["iban"] = *w.IBAN }
	if w.BIC != nil { m["attributes"].(fiber.Map)["bic"] = *w.BIC }
	if w.CurrencyID != nil { m["attributes"].(fiber.Map)["currency_id"] = *w.CurrencyID }
	if w.Notes != nil { m["attributes"].(fiber.Map)["notes"] = *w.Notes }
	if w.Latitude != nil { m["attributes"].(fiber.Map)["latitude"] = *w.Latitude }
	if w.Longitude != nil { m["attributes"].(fiber.Map)["longitude"] = *w.Longitude }
	return m
}

func isValidWalletType(wt domain.WalletType) bool {
	switch wt {
	case domain.WalletTypeAsset, domain.WalletTypeDefault, domain.WalletTypeCash,
		domain.WalletTypeDebt, domain.WalletTypeInitialBalance, domain.WalletTypeLoan,
		domain.WalletTypeMortgage, domain.WalletTypeReconciliation, domain.WalletTypeExpense,
		domain.WalletTypeRevenue, domain.WalletTypeBeneficiary, domain.WalletTypeCreditCard,
		domain.WalletTypeImport, domain.WalletTypeLiabilityCredit:
		return true
	}
	return false
}
