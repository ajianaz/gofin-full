package handler

import (
	"github.com/gofiber/fiber/v2"

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
	entityIDStr := c.Query("entity_id", "")

	// Parse entity_id as UUID for the query; pass zero if empty
	var entityID int64
	if entityIDStr != "" {
		// The audit repository uses int64 for entity IDs
		// For now, we just pass 0 if not parseable as int64
		entityID = 0
	}

	logs, err := h.repo.List(c.Context(), 0, entityType, entityID, 100)
	if err != nil {
		return apperrors.NewWithDetail(500, "failed to list audit logs", err.Error())
	}

	var data []fiber.Map
	for _, log := range logs {
		data = append(data, fiber.Map{
			"type": "audit_logs",
			"id":   log.ID,
			"attributes": fiber.Map{
				"user_id":     log.UserID,
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
