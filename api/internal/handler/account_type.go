package handler

import (
	"log"
"github.com/gofiber/fiber/v2"

	"github.com/ajianaz/gofin-full/api/internal/auth"
	"github.com/ajianaz/gofin-full/api/internal/repository"
	apperrors "github.com/ajianaz/gofin-full/api/pkg/errors")

type AccountTypeHandler struct {
	repo *repository.AccountTypeRepository
}

func NewAccountTypeHandler(repo *repository.AccountTypeRepository) *AccountTypeHandler {
	return &AccountTypeHandler{repo: repo}
}

func (h *AccountTypeHandler) Index(c *fiber.Ctx) error {
	_ = auth.GetUser(c)

	types, err := h.repo.List(c.Context())
	if err != nil {
		log.Printf("handler/Index: failed to list wallet types: %v", err)
		return apperrors.ErrInternal
	}

	var data []fiber.Map
	for _, t := range types {
		data = append(data, fiber.Map{
			"type": "wallet_types",
			"id":   t.ID,
			"attributes": fiber.Map{
				"type": t.Type,
			},
		})
	}
	return c.JSON(fiber.Map{"data": data})
}
