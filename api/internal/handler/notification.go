package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/azfirazka/gofin-full/api/internal/auth"
	"github.com/azfirazka/gofin-full/api/internal/repository"
	apperrors "github.com/azfirazka/gofin-full/api/pkg/errors"
)

type NotificationHandler struct {
	repo *repository.NotificationRepository
}

func NewNotificationHandler(repo *repository.NotificationRepository) *NotificationHandler {
	return &NotificationHandler{repo: repo}
}

func (h *NotificationHandler) Index(c *fiber.Ctx) error {
	user := auth.GetUser(c)

	notifications, err := h.repo.List(c.Context(), user.ID)
	if err != nil {
		return apperrors.NewWithDetail(500, "failed to list notifications", err.Error())
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

func (h *NotificationHandler) Unread(c *fiber.Ctx) error {
	user := auth.GetUser(c)

	notifications, err := h.repo.ListUnread(c.Context(), user.ID)
	if err != nil {
		return apperrors.NewWithDetail(500, "failed to list unread notifications", err.Error())
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

	id, err := c.ParamsInt("id")
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"id": {"invalid id"}})
	}

	if err := h.repo.MarkRead(c.Context(), int64(id), user.ID); err != nil {
		return apperrors.NotFoundResource("notification", int64(id))
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
		return apperrors.NewWithDetail(500, "failed to mark all notifications as read", err.Error())
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
