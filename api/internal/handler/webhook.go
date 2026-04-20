package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/azfirazka/gofin-full/api/internal/auth"
	"github.com/azfirazka/gofin-full/api/internal/repository"
	apperrors "github.com/azfirazka/gofin-full/api/pkg/errors"
)

type WebhookHandler struct {
	repo *repository.WebhookRepository
}

func NewWebhookHandler(repo *repository.WebhookRepository) *WebhookHandler {
	return &WebhookHandler{repo: repo}
}

func (h *WebhookHandler) Index(c *fiber.Ctx) error {
	_ = auth.GetUser(c)
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	webhooks, err := h.repo.List(c.Context(), *groupID)
	if err != nil {
		return apperrors.NewWithDetail(500, "failed to list webhooks", err.Error())
	}

	var data []fiber.Map
	for _, w := range webhooks {
		data = append(data, fiber.Map{
			"type": "webhooks",
			"id":   w.ID,
			"attributes": fiber.Map{
				"title":  w.Title,
				"url":    w.URL,
				"active": w.Active,
			},
		})
	}
	return c.JSON(fiber.Map{"data": data})
}

func (h *WebhookHandler) Show(c *fiber.Ctx) error {
	_ = auth.GetUser(c)
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"id": {"invalid id"}})
	}

	w, err := h.repo.FindByID(c.Context(), int64(id), *groupID)
	if err != nil {
		return apperrors.NotFoundResource("webhook", int64(id))
	}

	triggers, _ := h.repo.ListTriggers(c.Context(), w.ID)
	var triggerList []string
	for _, t := range triggers {
		triggerList = append(triggerList, t.Trigger)
	}

	return c.JSON(fiber.Map{"data": fiber.Map{
		"type": "webhooks",
		"id":   w.ID,
		"attributes": fiber.Map{
			"title":    w.Title,
			"url":      w.URL,
			"active":   w.Active,
			"triggers": triggerList,
		},
	}})
}

func (h *WebhookHandler) Store(c *fiber.Ctx) error {
	user := auth.GetUser(c)
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	var req struct {
		Title    string   `json:"title"`
		URL      string   `json:"url"`
		Triggers []string `json:"triggers"`
	}
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewValidationError(map[string][]string{"body": {"invalid JSON"}})
	}
	if req.Title == "" || req.URL == "" {
		return apperrors.NewValidationError(map[string][]string{"title": {"title and url are required"}})
	}

	w, err := h.repo.Create(c.Context(), user.ID, *groupID, req.Title, req.URL)
	if err != nil {
		return apperrors.NewWithDetail(500, "failed to create webhook", err.Error())
	}

	if len(req.Triggers) > 0 {
		_ = h.repo.SetTriggers(c.Context(), w.ID, req.Triggers)
	}

	return c.Status(201).JSON(fiber.Map{"data": fiber.Map{
		"type": "webhooks",
		"id":   w.ID,
		"attributes": fiber.Map{
			"title":  w.Title,
			"url":    w.URL,
			"active": w.Active,
		},
	}})
}

func (h *WebhookHandler) Update(c *fiber.Ctx) error {
	_ = auth.GetUser(c)
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"id": {"invalid id"}})
	}

	var req struct {
		Title    string   `json:"title"`
		URL      string   `json:"url"`
		Active   *bool    `json:"active"`
		Triggers []string `json:"triggers"`
	}
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewValidationError(map[string][]string{"body": {"invalid JSON"}})
	}

	if err := h.repo.Update(c.Context(), int64(id), *groupID, req.Title, req.URL, req.Active); err != nil {
		return apperrors.NotFoundResource("webhook", int64(id))
	}

	if req.Triggers != nil {
		_ = h.repo.SetTriggers(c.Context(), int64(id), req.Triggers)
	}

	return c.JSON(fiber.Map{"data": fiber.Map{
		"type": "webhooks", "id": id,
		"attributes": fiber.Map{"title": req.Title},
	}})
}

func (h *WebhookHandler) Delete(c *fiber.Ctx) error {
	_ = auth.GetUser(c)
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"id": {"invalid id"}})
	}

	if err := h.repo.Delete(c.Context(), int64(id), *groupID); err != nil {
		return apperrors.NotFoundResource("webhook", int64(id))
	}

	return c.Status(204).Send(nil)
}

func (h *WebhookHandler) Messages(c *fiber.Ctx) error {
	_ = auth.GetUser(c)

	id, err := c.ParamsInt("id")
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"id": {"invalid id"}})
	}

	messages, err := h.repo.ListMessages(c.Context(), int64(id))
	if err != nil {
		return apperrors.NewWithDetail(500, "failed to list webhook messages", err.Error())
	}

	var data []fiber.Map
	for _, m := range messages {
		data = append(data, fiber.Map{
			"type": "webhook_messages",
			"id":   m.ID,
			"attributes": fiber.Map{
				"message": m.Message,
			},
		})
	}
	return c.JSON(fiber.Map{"data": data})
}
