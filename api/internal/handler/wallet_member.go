package handler

import (
	"github.com/gofiber/fiber/v2"

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

	walletID, err := c.ParamsInt("wallet_id")
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"wallet_id": {"invalid wallet id"}})
	}

	members, err := h.memberRepo.ListByWallet(c.Context(), int64(walletID))
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

	walletID, err := c.ParamsInt("wallet_id")
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"wallet_id": {"invalid wallet id"}})
	}

	// Only owner can add members
	isOwner, err := h.memberRepo.IsWalletOwner(c.Context(), int64(walletID), user.ID)
	if err != nil || !isOwner {
		return apperrors.New(403, "only wallet owner can add members")
	}

	var req struct {
		UserID int64  `json:"user_id"`
		Role   string `json:"role"`
	}
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewValidationError(map[string][]string{"body": {"invalid JSON"}})
	}
	if req.UserID == 0 {
		return apperrors.NewValidationError(map[string][]string{"user_id": {"user_id is required"}})
	}
	if req.Role == "" {
		req.Role = "viewer"
	}
	if req.Role != "editor" && req.Role != "viewer" {
		return apperrors.NewValidationError(map[string][]string{"role": {"role must be editor or viewer"}})
	}

	m, err := h.memberRepo.AddMember(c.Context(), int64(walletID), req.UserID, req.Role)
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

	walletID, err := c.ParamsInt("wallet_id")
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"wallet_id": {"invalid wallet id"}})
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"id": {"invalid id"}})
	}

	isOwner, err := h.memberRepo.IsWalletOwner(c.Context(), int64(walletID), user.ID)
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

	if err := h.memberRepo.UpdateRole(c.Context(), int64(id), int64(walletID), req.Role); err != nil {
		return apperrors.NotFoundResource("wallet_member", int64(id))
	}

	return c.JSON(fiber.Map{"data": fiber.Map{
		"type": "wallet_members", "id": id,
		"attributes": fiber.Map{"role": req.Role},
	}})
}

func (h *WalletMemberHandler) Delete(c *fiber.Ctx) error {
	user := auth.GetUser(c)

	walletID, err := c.ParamsInt("wallet_id")
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"wallet_id": {"invalid wallet id"}})
	}

	userID, err := c.ParamsInt("user_id")
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"user_id": {"invalid user_id"}})
	}

	isOwner, err := h.memberRepo.IsWalletOwner(c.Context(), int64(walletID), user.ID)
	if err != nil || !isOwner {
		return apperrors.New(403, "only wallet owner can remove members")
	}

	if err := h.memberRepo.RemoveMember(c.Context(), int64(walletID), int64(userID)); err != nil {
		return apperrors.NewWithDetail(500, "failed to remove wallet member", err.Error())
	}

	return c.Status(204).Send(nil)
}
