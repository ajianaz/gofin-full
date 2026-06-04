package handler

import (
	"strings"

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
	repo  *repository.WalletRepository
	curry *CurrencyResolver
}

// NewWalletHandler creates a new wallet handler.
func NewWalletHandler(repo *repository.WalletRepository, curry *CurrencyResolver) *WalletHandler {
	return &WalletHandler{repo: repo, curry: curry}
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

	page := c.QueryInt("page", 1)
	perPage := c.QueryInt("per_page", 20)
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	} else if perPage > 100 {
		perPage = 100
	}

	wallets, total, err := h.repo.ListPaginated(c.Context(), *groupID, walletType, activeOnly, page, perPage)
	if err != nil {
		return apperrors.ErrInternal
	}

	// Collect currency IDs for batch resolve
	cIDs := make([]string, 0, len(wallets))
	for _, w := range wallets {
		if w.CurrencyID != nil && *w.CurrencyID != "" {
			cIDs = append(cIDs, *w.CurrencyID)
		}
	}
	cMap := h.curry.ResolveMany(c.Context(), cIDs)

	var data []fiber.Map
	for _, w := range wallets {
		data = append(data, walletToMap(&w, cMap))
	}

	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}
	return c.JSON(fiber.Map{
		"data": data,
		"meta": fiber.Map{
			"pagination": fiber.Map{
				"total":        total,
				"count":        len(data),
				"per_page":     perPage,
				"current_page": page,
				"total_pages":  totalPages,
			},
		},
	})
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

	cMap := h.curry.ResolveMany(c.Context(), currencyIDsFromWallet(wallet))

	return c.JSON(fiber.Map{"data": walletToMap(wallet, cMap)})
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
		Name           string `json:"name"`
		WalletType     string `json:"wallet_type"`
		CurrencyID     string `json:"currency_id"`
		Active         *bool  `json:"active"`
		IBAN           string `json:"iban"`
		BIC            string `json:"bic"`
		Notes          string `json:"notes"`
		OpeningBalance string `json:"opening_balance"`
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

	// Validate name
	if strings.TrimSpace(req.Name) == "" {
		return apperrors.NewValidationError(map[string][]string{
			"name": {"Name is required."},
		})
	}
	if len(req.Name) > 255 {
		return apperrors.NewValidationError(map[string][]string{
			"name": {"Name must be 255 characters or less."},
		})
	}

	active := true
	if req.Active != nil {
		active = *req.Active
	}

	// Parse opening balance (default 0) — must be non-negative
	virtualBalance := decimal.Zero
	if req.OpeningBalance != "" {
		ob, err := decimal.NewFromString(req.OpeningBalance)
		if err != nil {
			return apperrors.NewValidationError(map[string][]string{
				"opening_balance": {"Must be a valid number."},
			})
		}
		if ob.IsNegative() {
			return apperrors.NewValidationError(map[string][]string{
				"opening_balance": {"Must be zero or positive."},
			})
		}
		virtualBalance = ob
	}

	wallet := &domain.Wallet{
		UserID:          user.ID,
		UserGroupID:     *groupID,
		Name:            sanitizeStr(req.Name),
		AccountType:     req.WalletType,
		Active:          active,
		VirtualBalance:  virtualBalance,
		IncludeNetWorth: true,
	}
	if req.IBAN != "" {
		iban := sanitizeStr(req.IBAN)
		wallet.IBAN = &iban
	}
	if req.BIC != "" {
		bic := sanitizeStr(req.BIC)
		wallet.BIC = &bic
	}
	if req.CurrencyID != "" {
		info := h.curry.ResolveSingle(c.Context(), req.CurrencyID)
		if info.Code == "" {
			return apperrors.NewValidationError(map[string][]string{"currency_id": {"currency not found"}})
		}
		wallet.CurrencyID = &req.CurrencyID
	}
	if req.Notes != "" {
		notes := sanitizeStr(req.Notes)
		wallet.Notes = &notes
	}

	created, err := h.repo.Create(c.Context(), wallet)
	if err != nil {
		return apperrors.ErrInternal
	}

	cMap := h.curry.ResolveMany(c.Context(), currencyIDsFromWallet(created))

	return c.Status(201).JSON(fiber.Map{"data": walletToMap(created, cMap)})
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
		CurrencyID      *string `json:"currency_id"`
		Notes           *string `json:"notes"`
	}
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewValidationError(map[string][]string{
			"body": {"Invalid request body."},
		})
	}

	if strings.TrimSpace(req.Name) == "" {
		return apperrors.NewValidationError(map[string][]string{
			"name": {"name is required"},
		})
	}
	if len(req.Name) > 255 {
		return apperrors.NewValidationError(map[string][]string{
			"name": {"Name must be 255 characters or less."},
		})
	}

	// Verify wallet exists before update
	if _, err := h.repo.FindByID(c.Context(), id, *groupID); err != nil {
		return apperrors.NotFoundResource("wallet", id)
	}

	sanitizedName := sanitizeStr(req.Name)
	if err := h.repo.Update(c.Context(), id, *groupID, sanitizedName, req.Active, req.IncludeNetWorth, req.CurrencyID, sanitizePtr(req.Notes)); err != nil {
		return apperrors.ErrInternal
	}

	updated, err := h.repo.FindByID(c.Context(), id, *groupID)
	if err != nil {
		return apperrors.ErrInternal
	}

	cMap := h.curry.ResolveMany(c.Context(), currencyIDsFromWallet(updated))

	return c.JSON(fiber.Map{"data": walletToMap(updated, cMap)})
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

	// Verify wallet exists before delete
	if _, err := h.repo.FindByID(c.Context(), id, *groupID); err != nil {
		return apperrors.NotFoundResource("wallet", id)
	}

	if err := h.repo.Delete(c.Context(), id, *groupID); err != nil {
		return apperrors.ErrInternal
	}

	return c.Status(204).Send(nil)
}

func walletToMap(w *domain.Wallet, cMap map[string]CurrencyInfo) fiber.Map {
	dp := 2
	if w.CurrencyID != nil && *w.CurrencyID != "" {
		if ci, ok := cMap[*w.CurrencyID]; ok {
			dp = ci.DecimalPlaces
		}
	}
	m := fiber.Map{
		"type": "wallets",
		"id":   w.ID,
		"attributes": fiber.Map{
			"name":              w.Name,
			"wallet_type":       w.AccountType,
			"active":            w.Active,
			"virtual_balance":   w.VirtualBalance.StringFixed(int32(dp)),
			"include_net_worth": w.IncludeNetWorth,
			"created_at":        w.CreatedAt,
			"updated_at":        w.UpdatedAt,
		},
	}
	if w.IBAN != nil {
		m["attributes"].(fiber.Map)["iban"] = *w.IBAN
	}
	if w.BIC != nil {
		m["attributes"].(fiber.Map)["bic"] = *w.BIC
	}
	if w.CurrencyID != nil {
		m["attributes"].(fiber.Map)["currency_id"] = *w.CurrencyID
	}
	if w.Notes != nil {
		m["attributes"].(fiber.Map)["notes"] = *w.Notes
	}
	if w.Latitude != nil {
		m["attributes"].(fiber.Map)["latitude"] = *w.Latitude
	}
	if w.Longitude != nil {
		m["attributes"].(fiber.Map)["longitude"] = *w.Longitude
	}
	// Resolve currency info
	if w.CurrencyID != nil && *w.CurrencyID != "" {
		if ci, ok := cMap[*w.CurrencyID]; ok {
			m["attributes"].(fiber.Map)["currency_code"] = ci.Code
			m["attributes"].(fiber.Map)["currency_symbol"] = ci.Symbol
			m["attributes"].(fiber.Map)["currency_decimal_places"] = ci.DecimalPlaces
		}
	}
	return m
}

func currencyIDsFromWallet(w *domain.Wallet) []string {
	if w.CurrencyID != nil && *w.CurrencyID != "" {
		return []string{*w.CurrencyID}
	}
	return nil
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
