package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/ajianaz/gofin-full/api/internal/auth"
	"github.com/ajianaz/gofin-full/api/internal/repository"
	apperrors "github.com/ajianaz/gofin-full/api/pkg/errors"
)

type AnalyticsHandler struct {
	repo *repository.AnalyticsRepository
}

func NewAnalyticsHandler(repo *repository.AnalyticsRepository) *AnalyticsHandler {
	return &AnalyticsHandler{repo: repo}
}

func (h *AnalyticsHandler) SpendingByCategory(c *fiber.Ctx) error {
	_ = auth.GetUser(c)
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	start, end := parseDateRange(c)

	results, err := h.repo.SpendingByCategory(c.Context(), *groupID, start, end)
	if err != nil {
		return apperrors.NewWithDetail(500, "failed to get spending by category", err.Error())
	}

	var data []fiber.Map
	for _, r := range results {
		data = append(data, fiber.Map{
			"type": "category_spending",
			"attributes": fiber.Map{
				"category_id":   r.CategoryID,
				"category_name": r.CategoryName,
				"total":         r.Total,
				"count":         r.Count,
			},
		})
	}
	return c.JSON(fiber.Map{"data": data})
}

func (h *AnalyticsHandler) SpendingByPeriod(c *fiber.Ctx) error {
	_ = auth.GetUser(c)
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	start, end := parseDateRange(c)

	results, err := h.repo.SpendingByPeriod(c.Context(), *groupID, start, end)
	if err != nil {
		return apperrors.NewWithDetail(500, "failed to get spending by period", err.Error())
	}

	var data []fiber.Map
	for _, r := range results {
		data = append(data, fiber.Map{
			"type": "period_spending",
			"attributes": fiber.Map{
				"period":  r.Period,
				"income":  r.Income,
				"expense": r.Expense,
			},
		})
	}
	return c.JSON(fiber.Map{"data": data})
}

func (h *AnalyticsHandler) NetWorth(c *fiber.Ctx) error {
	_ = auth.GetUser(c)
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	start, end := parseDateRange(c)

	summary, err := h.repo.GetNetWorth(c.Context(), *groupID, start, end)
	if err != nil {
		return apperrors.NewWithDetail(500, "failed to get net worth", err.Error())
	}

	return c.JSON(fiber.Map{"data": fiber.Map{
		"type": "net_worth",
		"attributes": fiber.Map{
			"total_income":      summary.TotalIncome,
			"total_expense":     summary.TotalExpense,
			"net_income":        summary.NetIncome,
			"transaction_count": summary.TransactionCount,
		},
	}})
}

func parseDateRange(c *fiber.Ctx) (time.Time, time.Time) {
	now := time.Now().UTC()
	start := now.AddDate(0, -1, 0)
	end := now

	if s := c.Query("start"); s != "" {
		if t, err := time.Parse(time.RFC3339, s); err == nil {
			start = t
		}
	}
	if e := c.Query("end"); e != "" {
		if t, err := time.Parse(time.RFC3339, e); err == nil {
			end = t
		}
	}

	return start, end
}
