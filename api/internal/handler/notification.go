package handler

import (
	"github.com/rs/zerolog/log"
"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/ajianaz/gofin-full/api/internal/auth"
	"github.com/ajianaz/gofin-full/api/internal/repository"
	apperrors "github.com/ajianaz/gofin-full/api/pkg/errors")

type NotificationHandler struct {
	repo *repository.NotificationRepository
}

func NewNotificationHandler(repo *repository.NotificationRepository) *NotificationHandler {
	return &NotificationHandler{repo: repo}
}

func (h *NotificationHandler) Index(c *fiber.Ctx) error {
	user := auth.GetUser(c)

	page := c.QueryInt("page", 1)
	perPage := c.QueryInt("per_page", 20)
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	notifications, total, err := h.repo.ListPaginated(c.Context(), user.ID, page, perPage)
	if err != nil {
		log.Error().Err(err).Msg("handler/Index: failed to list notifications")
		return apperrors.ErrInternal
	}

	var data []fiber.Map
	for _, n := range notifications {
		data = append(data, fiber.Map{
			"type": "notifications",
			"id":   n.ID,
			"attributes": fiber.Map{
				"channel": n.Channel,
				"type":    n.Type,
				"title":   n.Title,
				"message": n.Message,
				"read":    n.Read,
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

func (h *NotificationHandler) Unread(c *fiber.Ctx) error {
	user := auth.GetUser(c)

	notifications, err := h.repo.ListUnread(c.Context(), user.ID)
	if err != nil {
		log.Error().Err(err).Msg("handler/Index: failed to list unread notifications")
		return apperrors.ErrInternal
	}

	var data []fiber.Map
	for _, n := range notifications {
		data = append(data, fiber.Map{
			"type": "notifications",
			"id":   n.ID,
			"attributes": fiber.Map{
				"channel": n.Channel,
				"type":    n.Type,
				"title":   n.Title,
				"message": n.Message,
				"read":    n.Read,
			},
		})
	}
	return c.JSON(fiber.Map{"data": data})
}

func (h *NotificationHandler) MarkRead(c *fiber.Ctx) error {
	user := auth.GetUser(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"id": {"invalid id format"}})
	}

	if err := h.repo.MarkRead(c.Context(), id, user.ID); err != nil {
		return apperrors.NotFoundResource("notification", id)
	}

	return c.JSON(fiber.Map{"data": fiber.Map{
		"type":       "notifications",
		"id":         id,
		"attributes": fiber.Map{"read": true},
	}})
}

func (h *NotificationHandler) MarkAllRead(c *fiber.Ctx) error {
	user := auth.GetUser(c)

	if err := h.repo.MarkAllRead(c.Context(), user.ID); err != nil {
		log.Error().Err(err).Msg("handler/Index: failed to mark all notifications as read")
		return apperrors.ErrInternal
	}

	return c.JSON(fiber.Map{"data": fiber.Map{
		"type":       "notifications",
		"attributes": fiber.Map{"read": true},
	}})
}

// Stream handles SSE connection for real-time notifications.
// The actual SSE streaming is handled by the SSE middleware; this handler
// validates the user and triggers the stream upgrade.
func (h *NotificationHandler) Stream(c *fiber.Ctx) error {
	user := auth.GetUser(c)
	_ = user.ID // SSE middleware uses user_id param from context
	return c.Next()
}
