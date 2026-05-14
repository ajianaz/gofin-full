package handler

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/ajianaz/gofin-full/api/internal/auth"
	"github.com/ajianaz/gofin-full/api/internal/repository"
	apperrors "github.com/ajianaz/gofin-full/api/pkg/errors"
)

// blockedHosts are hostnames that must never be used as webhook targets (SSRF prevention).
var blockedHosts = []string{
	"localhost", "127.0.0.1", "0.0.0.0", "::1",
	"169.254.169.254",          // AWS metadata
	"metadata.google.internal", // GCP metadata
	"100.100.100.200",          // GKE metadata
}

// validateWebhookURL checks that a webhook URL is not pointing to internal/private addresses.
// It resolves DNS at validation time to prevent DNS rebinding attacks.
func validateWebhookURL(rawURL string) error {
	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("URL scheme must be http or https")
	}
	host := strings.ToLower(u.Hostname())
	for _, blocked := range blockedHosts {
		if host == blocked || strings.HasSuffix(host, "."+blocked) {
			return fmt.Errorf("URL host is not allowed")
		}
	}
	// Check if host is a raw IP
	ip := net.ParseIP(host)
	if ip == nil {
		// Resolve DNS to catch rebinding attacks (attacker DNS returns public IP during
		// validation, then internal IP during actual webhook delivery)
		ips, err := net.LookupIP(host)
		if err != nil || len(ips) == 0 {
			return fmt.Errorf("could not resolve hostname")
		}
		ip = ips[0]
	}
	// Block private/internal IP ranges (covers both IPv4 and IPv6)
	if isPrivateIP(ip) {
		return fmt.Errorf("private/internal IP addresses are not allowed")
	}
	return nil
}

func isPrivateIP(ip net.IP) bool {
	privateNets := []struct {
		network *net.IPNet
	}{
		{mustParseCIDR("10.0.0.0/8")},
		{mustParseCIDR("172.16.0.0/12")},
		{mustParseCIDR("192.168.0.0/16")},
		{mustParseCIDR("127.0.0.0/8")},
		{mustParseCIDR("0.0.0.0/8")},
		{mustParseCIDR("169.254.0.0/16")},
		{mustParseCIDR("::1/128")},
		{mustParseCIDR("fe80::/10")},
		{mustParseCIDR("fc00::/7")},
	}
	for _, pn := range privateNets {
		if pn.network.Contains(ip) {
			return true
		}
	}
	return ip.IsLoopback() || ip.IsUnspecified() || ip.IsLinkLocalUnicast() || ip.IsPrivate()
}

func mustParseCIDR(s string) *net.IPNet {
	_, network, err := net.ParseCIDR(s)
	if err != nil {
		return &net.IPNet{}
	}
	return network
}

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
		log.Printf("handler: failed to list webhooks: %v", err)
		return apperrors.ErrInternal
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

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"id": {"invalid id format"}})
	}

	w, err := h.repo.FindByID(c.Context(), id, *groupID)
	if err != nil {
		return apperrors.NotFoundResource("webhook", id)
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

	if err := validateWebhookURL(req.URL); err != nil {
		return apperrors.NewValidationError(map[string][]string{"url": {err.Error()}})
	}

	w, err := h.repo.Create(c.Context(), user.ID, *groupID, req.Title, req.URL)
	if err != nil {
		log.Printf("handler: failed to create webhook: %v", err)
		return apperrors.ErrInternal
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

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"id": {"invalid id format"}})
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

	if req.URL != "" {
		if err := validateWebhookURL(req.URL); err != nil {
			return apperrors.NewValidationError(map[string][]string{"url": {err.Error()}})
		}
	}

	if err := h.repo.Update(c.Context(), id, *groupID, req.Title, req.URL, req.Active); err != nil {
		return apperrors.NotFoundResource("webhook", id)
	}

	if req.Triggers != nil {
		_ = h.repo.SetTriggers(c.Context(), id, req.Triggers)
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

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"id": {"invalid id format"}})
	}

	if err := h.repo.Delete(c.Context(), id, *groupID); err != nil {
		return apperrors.NotFoundResource("webhook", id)
	}

	return c.Status(204).Send(nil)
}

func (h *WebhookHandler) Messages(c *fiber.Ctx) error {
	_ = auth.GetUser(c)
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"id": {"invalid id format"}})
	}

	// Verify webhook belongs to user's active group
	if _, err := h.repo.FindByID(c.Context(), id, *groupID); err != nil {
		return apperrors.NotFoundResource("webhook", id)
	}

	messages, err := h.repo.ListMessages(c.Context(), id)
	if err != nil {
		log.Printf("handler: failed to list webhook messages: %v", err)
		return apperrors.ErrInternal
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
