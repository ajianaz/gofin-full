package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/ajianaz/gofin-full/api/internal/auth"
	"github.com/ajianaz/gofin-full/api/internal/domain"
	"github.com/ajianaz/gofin-full/api/internal/repository"
	apperrors "github.com/ajianaz/gofin-full/api/pkg/errors"
)

var validTriggerTypes = []string{
	"title_contains", "description_contains", "amount_is",
	"amount_less_than", "amount_greater_than", "category_is", "account_is",
}

var validActionTypes = []string{
	"set_category", "set_account", "add_tag", "remove_tag", "set_description",
}

type RuleGroupHandler struct {
	repo *repository.RuleGroupRepository
}

func NewRuleGroupHandler(repo *repository.RuleGroupRepository) *RuleGroupHandler {
	return &RuleGroupHandler{repo: repo}
}

func (h *RuleGroupHandler) Index(c *fiber.Ctx) error {
	_ = auth.GetUser(c)
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	page := c.QueryInt("page", 1)
	perPage := c.QueryInt("per_page", 20)
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	} else if perPage > 100 {
		perPage = 100
	}

	groups, total, err := h.repo.ListPaginated(c.Context(), *groupID, page, perPage)
	if err != nil {
		log.Error().Err(err).Msg("handler/Index: failed to list rule groups")
		return apperrors.ErrInternal
	}

	var data []fiber.Map
	for _, g := range groups {
		data = append(data, fiber.Map{
			"type": "rule_groups", "id": g.ID,
			"attributes": fiber.Map{"title": g.Title, "active": g.Active, "order": g.Order},
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

func (h *RuleGroupHandler) Store(c *fiber.Ctx) error {
	user := auth.GetUser(c)
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	var req struct {
		Title string `json:"title"`
		Order int    `json:"order"`
	}
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewValidationError(map[string][]string{"body": {"invalid JSON"}})
	}
	if req.Title == "" {
		return apperrors.NewValidationError(map[string][]string{"title": {"title is required"}})
	}

	g, err := h.repo.Create(c.Context(), user.ID, *groupID, sanitizeStr(req.Title), req.Order)
	if err != nil {
		log.Error().Err(err).Msg("handler/Index: failed to create rule group")
		return apperrors.ErrInternal
	}

	return c.Status(201).JSON(fiber.Map{"data": fiber.Map{
		"type": "rule_groups", "id": g.ID,
		"attributes": fiber.Map{"title": g.Title, "active": g.Active, "order": g.Order},
	}})
}

func (h *RuleGroupHandler) Show(c *fiber.Ctx) error {
	_ = auth.GetUser(c)
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"id": {"invalid id format"}})
	}

	g, err := h.repo.FindByID(c.Context(), id, *groupID)
	if err != nil {
		return apperrors.NotFoundResource("rule_group", id)
	}

	return c.JSON(fiber.Map{"data": fiber.Map{
		"type": "rule_groups", "id": g.ID,
		"attributes": fiber.Map{"title": g.Title, "active": g.Active, "order": g.Order},
	}})
}

func (h *RuleGroupHandler) Update(c *fiber.Ctx) error {
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
		Title  string `json:"title"`
		Active *bool  `json:"active"`
	}
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewValidationError(map[string][]string{"body": {"invalid JSON"}})
	}

	if err := h.repo.Update(c.Context(), id, *groupID, sanitizeStr(req.Title), req.Active); err != nil {
		return apperrors.NotFoundResource("rule_group", id)
	}

	return c.JSON(fiber.Map{"data": fiber.Map{
		"type": "rule_groups", "id": id,
		"attributes": fiber.Map{"title": req.Title},
	}})
}

func (h *RuleGroupHandler) Delete(c *fiber.Ctx) error {
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
		return apperrors.NotFoundResource("rule_group", id)
	}

	return c.Status(204).Send(nil)
}

// RuleHandler
type RuleHandler struct {
	repo *repository.RuleRepository
}

func NewRuleHandler(repo *repository.RuleRepository) *RuleHandler {
	return &RuleHandler{repo: repo}
}

func (h *RuleHandler) Index(c *fiber.Ctx) error {
	_ = auth.GetUser(c)
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	page := c.QueryInt("page", 1)
	perPage := c.QueryInt("per_page", 20)
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	} else if perPage > 100 {
		perPage = 100
	}

	rules, total, err := h.repo.ListPaginated(c.Context(), *groupID, page, perPage)
	if err != nil {
		log.Error().Err(err).Msg("handler/Index: failed to list rules")
		return apperrors.ErrInternal
	}

	var data []fiber.Map
	for _, r := range rules {
		data = append(data, fiber.Map{
			"type": "rules", "id": r.ID,
			"attributes": fiber.Map{
				"title": r.Title, "priority": r.Priority, "active": r.Active,
				"strict": r.Strict, "stop_processing": r.StopProcessing,
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

func (h *RuleHandler) Show(c *fiber.Ctx) error {
	_ = auth.GetUser(c)
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"id": {"invalid id format"}})
	}

	rule, err := h.repo.FindByID(c.Context(), id, *groupID)
	if err != nil {
		return apperrors.NotFoundResource("rule", id)
	}

	return c.JSON(fiber.Map{"data": ruleToMap(rule)})
}

func (h *RuleHandler) Store(c *fiber.Ctx) error {
	user := auth.GetUser(c)
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	var req struct {
		Title      string               `json:"title"`
		Priority   int                  `json:"priority"`
		RuleGroupID *uuid.UUID          `json:"rule_group_id"`
		Triggers   []domain.RuleTrigger `json:"triggers"`
		Actions    []domain.RuleAction  `json:"actions"`
	}
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewValidationError(map[string][]string{"body": {"invalid JSON"}})
	}
	if req.Title == "" {
		return apperrors.NewValidationError(map[string][]string{"title": {"title is required"}})
	}

	// Validate triggers and actions
	fieldErrors := make(map[string][]string)
	for i, t := range req.Triggers {
		valid := false
		for _, vt := range validTriggerTypes {
			if t.TriggerType == vt {
				valid = true
				break
			}
		}
		if !valid {
			fieldErrors[fmt.Sprintf("triggers.%d.trigger_type", i)] = append(fieldErrors[fmt.Sprintf("triggers.%d.trigger_type", i)], "must be one of: "+fmt.Sprintf("%v", validTriggerTypes))
		}
	}
	for i, a := range req.Actions {
		valid := false
		for _, va := range validActionTypes {
			if a.ActionType == va {
				valid = true
				break
			}
		}
		if !valid {
			fieldErrors[fmt.Sprintf("actions.%d.action_type", i)] = append(fieldErrors[fmt.Sprintf("actions.%d.action_type", i)], "must be one of: "+fmt.Sprintf("%v", validActionTypes))
		}
	}
	if len(fieldErrors) > 0 {
		return apperrors.NewValidationError(fieldErrors)
	}

	rule, err := h.repo.Create(c.Context(), user.ID, *groupID, sanitizeStr(req.Title), req.Priority, req.RuleGroupID)
	if err != nil {
		log.Error().Err(err).Msg("handler/Index: failed to create rule")
		return apperrors.ErrInternal
	}

	if len(req.Triggers) > 0 {
		h.repo.SetTriggers(c.Context(), rule.ID, req.Triggers)
	}
	if len(req.Actions) > 0 {
		h.repo.SetActions(c.Context(), rule.ID, req.Actions)
	}

	return c.Status(201).JSON(fiber.Map{"data": fiber.Map{
		"type": "rules", "id": rule.ID,
		"attributes": fiber.Map{"title": rule.Title, "priority": rule.Priority, "active": rule.Active},
	}})
}

func (h *RuleHandler) Update(c *fiber.Ctx) error {
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
		Title           string                `json:"title"`
		Active          *bool                 `json:"active"`
		Strict          *bool                 `json:"strict"`
		StopProcessing  *bool                 `json:"stop_processing"`
		Triggers        []domain.RuleTrigger  `json:"triggers"`
		Actions         []domain.RuleAction   `json:"actions"`
	}
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewValidationError(map[string][]string{"body": {"invalid JSON"}})
	}

	// Validate triggers and actions
	fieldErrors := make(map[string][]string)
	for i, t := range req.Triggers {
		valid := false
		for _, vt := range validTriggerTypes {
			if t.TriggerType == vt {
				valid = true
				break
			}
		}
		if !valid {
			fieldErrors[fmt.Sprintf("triggers.%d.trigger_type", i)] = append(fieldErrors[fmt.Sprintf("triggers.%d.trigger_type", i)], "must be one of: "+fmt.Sprintf("%v", validTriggerTypes))
		}
	}
	for i, a := range req.Actions {
		valid := false
		for _, va := range validActionTypes {
			if a.ActionType == va {
				valid = true
				break
			}
		}
		if !valid {
			fieldErrors[fmt.Sprintf("actions.%d.action_type", i)] = append(fieldErrors[fmt.Sprintf("actions.%d.action_type", i)], "must be one of: "+fmt.Sprintf("%v", validActionTypes))
		}
	}
	if len(fieldErrors) > 0 {
		return apperrors.NewValidationError(fieldErrors)
	}

	if err := h.repo.Update(c.Context(), id, *groupID, sanitizeStr(req.Title), req.Active, req.Strict, req.StopProcessing); err != nil {
		return apperrors.NotFoundResource("rule", id)
	}

	if req.Triggers != nil {
		h.repo.SetTriggers(c.Context(), id, req.Triggers)
	}
	if req.Actions != nil {
		h.repo.SetActions(c.Context(), id, req.Actions)
	}

	return c.JSON(fiber.Map{"data": fiber.Map{
		"type": "rules", "id": id,
		"attributes": fiber.Map{"title": req.Title},
	}})
}

func (h *RuleHandler) Delete(c *fiber.Ctx) error {
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
		return apperrors.NotFoundResource("rule", id)
	}

	return c.Status(204).Send(nil)
}

func ruleToMap(r *domain.Rule) fiber.Map {
	var triggers []fiber.Map
	for _, t := range r.Triggers {
		triggers = append(triggers, fiber.Map{
			"id": t.ID, "trigger_type": t.TriggerType,
			"trigger_value": t.TriggerValue, "stop_processing": t.StopProcessing,
		})
	}
	var actions []fiber.Map
	for _, a := range r.Actions {
		actions = append(actions, fiber.Map{
			"id": a.ID, "action_type": a.ActionType,
			"action_value": a.ActionValue, "order": a.Order,
		})
	}
	return fiber.Map{
		"type": "rules", "id": r.ID,
		"attributes": fiber.Map{
			"title": r.Title, "priority": r.Priority, "active": r.Active,
			"strict": r.Strict, "stop_processing": r.StopProcessing,
			"triggers": triggers, "actions": actions,
		},
	}
}
