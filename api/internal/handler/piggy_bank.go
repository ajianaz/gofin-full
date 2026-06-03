package handler

import (
	"github.com/rs/zerolog/log"
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/ajianaz/gofin-full/api/internal/auth"
	"github.com/ajianaz/gofin-full/api/internal/domain"
	"github.com/ajianaz/gofin-full/api/internal/repository"
	apperrors "github.com/ajianaz/gofin-full/api/pkg/errors"
)

type PiggyBankHandler struct {
	repo       *repository.PiggyBankRepository
	walletRepo *repository.WalletRepository
	curry      *CurrencyResolver
}

func NewPiggyBankHandler(repo *repository.PiggyBankRepository, walletRepo *repository.WalletRepository, curry *CurrencyResolver) *PiggyBankHandler {
	return &PiggyBankHandler{repo: repo, walletRepo: walletRepo, curry: curry}
}

// resolvePiggyBankCurrency resolves the currency info for a piggy bank via its wallet.
func (h *PiggyBankHandler) resolvePiggyBankCurrency(ctx context.Context, accountID uuid.UUID) CurrencyInfo {
	if accountID == uuid.Nil {
		return CurrencyInfo{DecimalPlaces: 2}
	}
	w, err := h.walletRepo.FindByID(ctx, accountID, uuid.Nil)
	if err != nil || w == nil || w.CurrencyID == nil || *w.CurrencyID == "" {
		return CurrencyInfo{DecimalPlaces: 2}
	}
	return h.curry.ResolveSingle(ctx, *w.CurrencyID)
}

// requireGroupID extracts the active group ID or returns an error.
func requireGroupID(c *fiber.Ctx) (uuid.UUID, error) {
	gid := auth.GetActiveGroupID(c)
	if gid == nil || *gid == uuid.Nil {
		return uuid.Nil, apperrors.NewWithDetail(400, "Bad Request", "No active group selected.")
	}
	return *gid, nil
}

func (h *PiggyBankHandler) Index(c *fiber.Ctx) error {
	groupID, err := requireGroupID(c)
	if err != nil {
		return err
	}

	accountID, err := uuid.Parse(c.Params("wallet_id"))
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"wallet_id": {"invalid wallet_id"}})
	}

	page := c.QueryInt("page", 1)
	perPage := c.QueryInt("per_page", 20)
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	pbs, total, err := h.repo.ListPaginated(c.Context(), accountID, groupID, page, perPage)
	if err != nil {
		log.Error().Err(err).Msg("handler/Index: failed to list piggy banks")
		return apperrors.ErrInternal
	}

	ci := h.resolvePiggyBankCurrency(c.Context(), accountID)
	dp := int32(ci.DecimalPlaces)

	var data []fiber.Map
	for _, pb := range pbs {
		data = append(data, fiber.Map{
			"type": "piggy_banks",
			"id":   pb.ID,
			"attributes": fiber.Map{
				"wallet_id": pb.AccountID, "name": pb.Name,
				"target_amount": pb.TargetAmount.StringFixed(dp),
				"start_date":    pb.StartDate, "target_date": pb.TargetDate,
				"order": pb.Order,
			},
		})
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

func (h *PiggyBankHandler) Show(c *fiber.Ctx) error {
	groupID, err := requireGroupID(c)
	if err != nil {
		return err
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"id": {"invalid id format"}})
	}

	pb, err := h.repo.FindByID(c.Context(), id, groupID)
	if err != nil {
		return apperrors.NotFoundResource("piggy_bank", id)
	}

	ci := h.resolvePiggyBankCurrency(c.Context(), pb.AccountID)
	dp := int32(ci.DecimalPlaces)
	attrs := fiber.Map{
		"account_id": pb.AccountID, "name": pb.Name,
		"target_amount":  pb.TargetAmount.StringFixed(dp),
		"current_amount": pb.CurrentAmount.StringFixed(dp),
		"left_to_target": pb.LeftToTarget.StringFixed(dp),
		"percentage":     pb.Percentage,
		"start_date":     pb.StartDate, "target_date": pb.TargetDate,
		"order": pb.Order, "notes": pb.Notes,
	}
	if ci.Code != "" {
		attrs["currency_code"] = ci.Code
		attrs["currency_symbol"] = ci.Symbol
		attrs["currency_decimal_places"] = ci.DecimalPlaces
	}

	return c.JSON(fiber.Map{"data": fiber.Map{
		"type": "piggy_banks", "id": pb.ID,
		"attributes": attrs,
	}})
}

func (h *PiggyBankHandler) Store(c *fiber.Ctx) error {
	groupID, err := requireGroupID(c)
	if err != nil {
		return err
	}

	var req struct {
		WalletID     uuid.UUID `json:"wallet_id"`
		Name         string    `json:"name"`
		TargetAmount string    `json:"target_amount"`
		Order        int       `json:"order"`
		Notes        *string   `json:"notes"`
	}
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewValidationError(map[string][]string{"body": {"invalid JSON"}})
	}

	fieldErrors := make(map[string][]string)
	if req.WalletID == uuid.Nil {
		fieldErrors["wallet_id"] = append(fieldErrors["wallet_id"], "wallet_id is required")
	}
	if strings.TrimSpace(req.Name) == "" {
		fieldErrors["name"] = append(fieldErrors["name"], "name is required")
	}
	if req.TargetAmount != "" {
		if ta, err := decimal.NewFromString(req.TargetAmount); err != nil {
			fieldErrors["target_amount"] = append(fieldErrors["target_amount"], "must be a valid number")
		} else if ta.IsNegative() {
			fieldErrors["target_amount"] = append(fieldErrors["target_amount"], "must be zero or positive")
		}
	}
	if len(fieldErrors) > 0 {
		return apperrors.NewValidationError(fieldErrors)
	}

	targetAmt, _ := decimal.NewFromString(req.TargetAmount)
	pb := &domain.PiggyBank{
		AccountID:    req.WalletID,
		Name:         sanitizeStr(req.Name),
		TargetAmount: targetAmt,
		Order:        req.Order,
		Notes:        sanitizePtr(req.Notes),
	}

	pb, err = h.repo.Create(c.Context(), pb, groupID)
	if err != nil {
		log.Error().Err(err).Msg("handler/Index: failed to create piggy bank")
		return apperrors.ErrInternal
	}

	ci := h.resolvePiggyBankCurrency(c.Context(), pb.AccountID)
	dp := int32(ci.DecimalPlaces)

	return c.Status(201).JSON(fiber.Map{"data": fiber.Map{
		"type": "piggy_banks", "id": pb.ID,
		"attributes": fiber.Map{
			"wallet_id": pb.AccountID, "name": pb.Name,
			"target_amount": pb.TargetAmount.StringFixed(dp),
			"order":         pb.Order,
		},
	}})
}

func (h *PiggyBankHandler) Update(c *fiber.Ctx) error {
	groupID, err := requireGroupID(c)
	if err != nil {
		return err
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"id": {"invalid id format"}})
	}

	var req struct {
		Name         string  `json:"name"`
		TargetAmount *string `json:"target_amount"`
		Notes        *string `json:"notes"`
	}
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewValidationError(map[string][]string{"body": {"invalid JSON"}})
	}

	var targetAmount *decimal.Decimal
	if req.TargetAmount != nil {
		amt, _ := decimal.NewFromString(*req.TargetAmount)
		targetAmount = &amt
	}

	// Verify piggy bank exists before update
	if _, err := h.repo.FindByID(c.Context(), id, groupID); err != nil {
		return apperrors.NotFoundResource("piggy_bank", id)
	}
	if err := h.repo.Update(c.Context(), id, groupID, sanitizeStr(req.Name), targetAmount, nil, nil, sanitizePtr(req.Notes)); err != nil {
		return apperrors.NotFoundResource("piggy_bank", id)
	}

	return c.JSON(fiber.Map{"data": fiber.Map{
		"type": "piggy_banks", "id": id,
		"attributes": fiber.Map{"name": req.Name, "notes": req.Notes},
	}})
}

func (h *PiggyBankHandler) Delete(c *fiber.Ctx) error {
	groupID, err := requireGroupID(c)
	if err != nil {
		return err
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"id": {"invalid id format"}})
	}

	// Verify piggy bank exists before delete
	if _, err := h.repo.FindByID(c.Context(), id, groupID); err != nil {
		return apperrors.NotFoundResource("piggy_bank", id)
	}
	if err := h.repo.Delete(c.Context(), id, groupID); err != nil {
		return apperrors.NotFoundResource("piggy_bank", id)
	}

	return c.Status(204).Send(nil)
}

// AddMoney handles POST /wallets/:wallet_id/piggy_banks/:id/add-money
func (h *PiggyBankHandler) AddMoney(c *fiber.Ctx) error {
	groupID, err := requireGroupID(c)
	if err != nil {
		return err
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"id": {"invalid id format"}})
	}

	var req struct {
		Amount string `json:"amount"`
	}
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewValidationError(map[string][]string{"body": {"invalid JSON"}})
	}
	if req.Amount == "" {
		return apperrors.NewValidationError(map[string][]string{"amount": {"amount is required"}})
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"amount": {"invalid amount"}})
	}
	evt, err := h.repo.AddMoney(c.Context(), id, groupID, amount)
	if err != nil {
		log.Error().Err(err).Msg("piggy bank add money failed")
		return apperrors.New(422, "Could not add money to piggy bank. Check your wallet balance.")
	}

	return c.JSON(fiber.Map{"data": fiber.Map{
		"type": "piggy_bank_events", "id": evt.ID,
		"attributes": fiber.Map{"piggy_bank_id": evt.PiggyBankID, "amount": evt.Amount.StringFixed(2)},
	}})
}

// RemoveMoney handles POST /wallets/:wallet_id/piggy_banks/:id/remove-money
func (h *PiggyBankHandler) RemoveMoney(c *fiber.Ctx) error {
	groupID, err := requireGroupID(c)
	if err != nil {
		return err
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"id": {"invalid id format"}})
	}

	var req struct {
		Amount string `json:"amount"`
	}
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewValidationError(map[string][]string{"body": {"invalid JSON"}})
	}
	if req.Amount == "" {
		return apperrors.NewValidationError(map[string][]string{"amount": {"amount is required"}})
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"amount": {"invalid amount"}})
	}
	evt, err := h.repo.RemoveMoney(c.Context(), id, groupID, amount)
	if err != nil {
		log.Error().Err(err).Msg("piggy bank remove money failed")
		return apperrors.New(422, "Could not remove money from piggy bank. Insufficient piggy bank balance.")
	}

	return c.JSON(fiber.Map{"data": fiber.Map{
		"type": "piggy_bank_events", "id": evt.ID,
		"attributes": fiber.Map{"piggy_bank_id": evt.PiggyBankID, "amount": evt.Amount.StringFixed(2)},
	}})
}
