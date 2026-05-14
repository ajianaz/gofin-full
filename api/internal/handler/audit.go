package handler

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/ajianaz/gofin-full/api/internal/auth"
	"github.com/ajianaz/gofin-full/api/internal/repository"
	apperrors "github.com/ajianaz/gofin-full/api/pkg/errors"
)

type AuditHandler struct {
	repo *repository.AuditRepository
}

func NewAuditHandler(repo *repository.AuditRepository) *AuditHandler {
	return &AuditHandler{repo: repo}
}

func (h *AuditHandler) Index(c *fiber.Ctx) error {
	_ = auth.GetUser(c)
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	entityType := c.Query("entity_type", "")

	var entityID uuid.UUID
	if v := c.Query("entity_id"); v != "" {
		entityID, _ = uuid.Parse(v)
	}

	logs, err := h.repo.List(c.Context(), *groupID, entityType, entityID, 100)
	if err != nil {
		log.Printf("handler: failed to list audit logs: %v", err)
		return apperrors.ErrInternal
	}

	var data []fiber.Map
	for _, log := range logs {
		data = append(data, fiber.Map{
			"type": "audit_logs",
			"id":   log.ID,
			"attributes": fiber.Map{
				"user_id":     log.UserID,
				"user_email":  log.UserEmail,
				"action":      log.Action,
				"entity_type": log.Entity,
				"entity_id":   log.EntityID,
				"old_value":   log.OldValue,
				"new_value":   log.NewValue,
				"ip_address":  log.IPAddress,
				"created_at":  log.CreatedAt.Format("2006-01-02T15:04:05Z"),
			},
		})
	}
	return c.JSON(fiber.Map{"data": data})
}
