package handler

import (
	"fmt"
	"log"
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
		log.Printf("handler/SpendingByCategory: failed to get spending by category: %v", err)
		return apperrors.ErrInternal
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
		log.Printf("handler/SpendingByCategory: failed to get spending by period: %v", err)
		return apperrors.ErrInternal
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
		log.Printf("handler/SpendingByCategory: failed to get net worth: %v", err)
		return apperrors.ErrInternal
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
		if t, err := parseFlexDate(s); err == nil {
			start = t
		}
	}
	if e := c.Query("end"); e != "" {
		if t, err := parseFlexDate(e); err == nil {
			end = t
		}
	}

	return start, end
}

// parseFlexDate tries multiple date layouts to support both RFC3339 and plain YYYY-MM-DD.
func parseFlexDate(s string) (time.Time, error) {
	for _, layout := range []string{
		time.RFC3339,
		"2006-01-02",
		"2006-01-02T15:04:05Z",
	} {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unrecognized date format: %s", s)
}
