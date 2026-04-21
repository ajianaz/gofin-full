package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/ajianaz/gofin-full/api/internal/auth"
	"github.com/ajianaz/gofin-full/api/internal/repository"
	apperrors "github.com/ajianaz/gofin-full/api/pkg/errors"
)

type WalletMemberHandler struct {
	memberRepo *repository.WalletMemberRepository
}

func NewWalletMemberHandler(memberRepo *repository.WalletMemberRepository) *WalletMemberHandler {
	return &WalletMemberHandler{memberRepo: memberRepo}
}

func (h *WalletMemberHandler) Index(c *fiber.Ctx) error {
	_ = auth.GetUser(c)

	walletID, err := uuid.Parse(c.Params("wallet_id"))
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"wallet_id": {"invalid wallet id format"}})
	}

	members, err := h.memberRepo.ListByWallet(c.Context(), walletID)
	if err != nil {
		return apperrors.NewWithDetail(500, "failed to list wallet members", err.Error())
	}

	var data []fiber.Map
	for _, m := range members {
		data = append(data, fiber.Map{
			"type": "wallet_members",
			"id":   m.ID,
			"attributes": fiber.Map{
				"wallet_id": m.WalletID,
				"user_id":   m.UserID,
				"role":      m.Role,
			},
		})
	}
	return c.JSON(fiber.Map{"data": data})
}

func (h *WalletMemberHandler) Store(c *fiber.Ctx) error {
	user := auth.GetUser(c)

	walletID, err := uuid.Parse(c.Params("wallet_id"))
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"wallet_id": {"invalid wallet id format"}})
	}

	// Only owner can add members
	isOwner, err := h.memberRepo.IsWalletOwner(c.Context(), walletID, user.ID)
	if err != nil || !isOwner {
		return apperrors.New(403, "only wallet owner can add members")
	}

	var req struct {
		UserID uuid.UUID `json:"user_id"`
		Role   string    `json:"role"`
	}
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewValidationError(map[string][]string{"body": {"invalid JSON"}})
	}
	if req.UserID == uuid.Nil {
		return apperrors.NewValidationError(map[string][]string{"user_id": {"user_id is required"}})
	}
	if req.Role == "" {
		req.Role = "viewer"
	}
	if req.Role != "editor" && req.Role != "viewer" {
		return apperrors.NewValidationError(map[string][]string{"role": {"role must be editor or viewer"}})
	}

	m, err := h.memberRepo.AddMember(c.Context(), walletID, req.UserID, req.Role)
	if err != nil {
		return apperrors.NewWithDetail(500, "failed to add wallet member", err.Error())
	}

	return c.Status(201).JSON(fiber.Map{"data": fiber.Map{
		"type": "wallet_members",
		"id":   m.ID,
		"attributes": fiber.Map{
			"wallet_id": m.WalletID,
			"user_id":   m.UserID,
			"role":      m.Role,
		},
	}})
}

func (h *WalletMemberHandler) Update(c *fiber.Ctx) error {
	user := auth.GetUser(c)

	walletID, err := uuid.Parse(c.Params("wallet_id"))
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"wallet_id": {"invalid wallet id format"}})
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"id": {"invalid id format"}})
	}

	isOwner, err := h.memberRepo.IsWalletOwner(c.Context(), walletID, user.ID)
	if err != nil || !isOwner {
		return apperrors.New(403, "only wallet owner can update members")
	}

	var req struct {
		Role string `json:"role"`
	}
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewValidationError(map[string][]string{"body": {"invalid JSON"}})
	}
	if req.Role != "editor" && req.Role != "viewer" {
		return apperrors.NewValidationError(map[string][]string{"role": {"role must be editor or viewer"}})
	}

	if err := h.memberRepo.UpdateRole(c.Context(), id, walletID, req.Role); err != nil {
		return apperrors.NotFoundResource("wallet_member", id)
	}

	return c.JSON(fiber.Map{"data": fiber.Map{
		"type": "wallet_members", "id": id,
		"attributes": fiber.Map{"role": req.Role},
	}})
}

func (h *WalletMemberHandler) Delete(c *fiber.Ctx) error {
	user := auth.GetUser(c)

	walletID, err := uuid.Parse(c.Params("wallet_id"))
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"wallet_id": {"invalid wallet id format"}})
	}

	userID, err := uuid.Parse(c.Params("user_id"))
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"user_id": {"invalid user_id format"}})
	}

	isOwner, err := h.memberRepo.IsWalletOwner(c.Context(), walletID, user.ID)
	if err != nil || !isOwner {
		return apperrors.New(403, "only wallet owner can remove members")
	}

	if err := h.memberRepo.RemoveMember(c.Context(), walletID, userID); err != nil {
		return apperrors.NewWithDetail(500, "failed to remove wallet member", err.Error())
	}

	return c.Status(204).Send(nil)
}
