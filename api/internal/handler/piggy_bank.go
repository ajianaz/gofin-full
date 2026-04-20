package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"

	"github.com/ajianaz/gofin-full/api/internal/auth"
	"github.com/ajianaz/gofin-full/api/internal/domain"
	"github.com/ajianaz/gofin-full/api/internal/repository"
	apperrors "github.com/ajianaz/gofin-full/api/pkg/errors"
)

type PiggyBankHandler struct {
	repo *repository.PiggyBankRepository
}

func NewPiggyBankHandler(repo *repository.PiggyBankRepository) *PiggyBankHandler {
	return &PiggyBankHandler{repo: repo}
}

// requireGroupID extracts the active group ID or returns an error.
func requireGroupID(c *fiber.Ctx) (int64, error) {
	gid := auth.GetActiveGroupID(c)
	if gid == nil || *gid == 0 {
		return 0, apperrors.NewWithDetail(400, "Bad Request", "No active group selected.")
	}
	return *gid, nil
}

func (h *PiggyBankHandler) Index(c *fiber.Ctx) error {
	groupID, err := requireGroupID(c)
	if err != nil {
		return err
	}

	accountID, err := c.ParamsInt("wallet_id")
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"wallet_id": {"invalid wallet_id"}})
	}

	pbs, err := h.repo.List(c.Context(), int64(accountID), groupID)
	if err != nil {
		return apperrors.NewWithDetail(500, "failed to list piggy banks", err.Error())
	}

	var data []fiber.Map
	for _, pb := range pbs {
		data = append(data, fiber.Map{
			"type":       "piggy_banks",
			"id":         pb.ID,
			"attributes": fiber.Map{
				"wallet_id": pb.AccountID, "name": pb.Name,
				"target_amount": pb.TargetAmount.StringFixed(2),
				"start_date":    pb.StartDate, "target_date": pb.TargetDate,
				"order":         pb.Order,
			},
		})
	}
	return c.JSON(fiber.Map{"data": data})
}

func (h *PiggyBankHandler) Show(c *fiber.Ctx) error {
	groupID, err := requireGroupID(c)
	if err != nil {
		return err
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"id": {"invalid id"}})
	}

	pb, err := h.repo.FindByID(c.Context(), int64(id), groupID)
	if err != nil {
		return apperrors.NotFoundResource("piggy_bank", int64(id))
	}

	return c.JSON(fiber.Map{"data": fiber.Map{
		"type": "piggy_banks", "id": pb.ID,
		"attributes": fiber.Map{
			"account_id":      pb.AccountID, "name": pb.Name,
			"target_amount":   pb.TargetAmount.StringFixed(2),
			"current_amount":  pb.CurrentAmount.StringFixed(2),
			"left_to_target":  pb.LeftToTarget.StringFixed(2),
			"percentage":      pb.Percentage,
			"start_date":      pb.StartDate, "target_date": pb.TargetDate,
			"order":           pb.Order, "notes": pb.Notes,
		},
	}})
}

func (h *PiggyBankHandler) Store(c *fiber.Ctx) error {
	groupID, err := requireGroupID(c)
	if err != nil {
		return err
	}

	var req struct {
		WalletID    int64   `json:"wallet_id"`
		Name         string  `json:"name"`
		TargetAmount string  `json:"target_amount"`
		Order        int     `json:"order"`
		Notes        *string `json:"notes"`
	}
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewValidationError(map[string][]string{"body": {"invalid JSON"}})
	}

	fieldErrors := make(map[string][]string)
	if req.WalletID == 0 {
		fieldErrors["wallet_id"] = append(fieldErrors["wallet_id"], "wallet_id is required")
	}
	if req.Name == "" {
		fieldErrors["name"] = append(fieldErrors["name"], "name is required")
	}
	if len(fieldErrors) > 0 {
		return apperrors.NewValidationError(fieldErrors)
	}

	targetAmt, _ := decimal.NewFromString(req.TargetAmount)
	pb := &domain.PiggyBank{
		AccountID:    req.WalletID,
		Name:         req.Name,
		TargetAmount: targetAmt,
		Order:        req.Order,
		Notes:        req.Notes,
	}

	pb, err = h.repo.Create(c.Context(), pb, groupID)
	if err != nil {
		return apperrors.NewWithDetail(500, "failed to create piggy bank", err.Error())
	}

	return c.Status(201).JSON(fiber.Map{"data": fiber.Map{
		"type": "piggy_banks", "id": pb.ID,
		"attributes": fiber.Map{
			"wallet_id": pb.AccountID, "name": pb.Name,
			"target_amount": pb.TargetAmount.StringFixed(2),
			"order": pb.Order,
		},
	}})
}

func (h *PiggyBankHandler) Update(c *fiber.Ctx) error {
	groupID, err := requireGroupID(c)
	if err != nil {
		return err
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"id": {"invalid id"}})
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

	if err := h.repo.Update(c.Context(), int64(id), groupID, req.Name, targetAmount, nil, nil, req.Notes); err != nil {
		return apperrors.NotFoundResource("piggy_bank", int64(id))
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

	id, err := c.ParamsInt("id")
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"id": {"invalid id"}})
	}

	if err := h.repo.Delete(c.Context(), int64(id), groupID); err != nil {
		return apperrors.NotFoundResource("piggy_bank", int64(id))
	}

	return c.Status(204).Send(nil)
}

// AddMoney handles POST /wallets/:wallet_id/piggy_banks/:id/add-money
func (h *PiggyBankHandler) AddMoney(c *fiber.Ctx) error {
	groupID, err := requireGroupID(c)
	if err != nil {
		return err
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"id": {"invalid id"}})
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
	evt, err := h.repo.AddMoney(c.Context(), int64(id), groupID, amount)
	if err != nil {
		return apperrors.NewWithDetail(422, "failed to add money", err.Error())
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

	id, err := c.ParamsInt("id")
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"id": {"invalid id"}})
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
	evt, err := h.repo.RemoveMoney(c.Context(), int64(id), groupID, amount)
	if err != nil {
		return apperrors.NewWithDetail(422, "failed to remove money", err.Error())
	}

	return c.JSON(fiber.Map{"data": fiber.Map{
		"type": "piggy_bank_events", "id": evt.ID,
		"attributes": fiber.Map{"piggy_bank_id": evt.PiggyBankID, "amount": evt.Amount.StringFixed(2)},
	}})
}
