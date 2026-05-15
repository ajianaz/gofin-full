package handler

import (
"bytes"
	"time"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/ajianaz/gofin-full/api/internal/auth"
	"github.com/ajianaz/gofin-full/api/internal/repository"
	"github.com/ajianaz/gofin-full/api/internal/service"
	apperrors "github.com/ajianaz/gofin-full/api/pkg/errors")

type ExportHandler struct {
	exportService *service.ExportService
}

func NewExportHandler(exportService *service.ExportService) *ExportHandler {
	return &ExportHandler{exportService: exportService}
}

func (h *ExportHandler) CSV(c *fiber.Ctx) error {
	_ = auth.GetUser(c)
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	filter := parseExportFilter(c)

	buf := &bytes.Buffer{}
	if err := h.exportService.ExportTransactionsCSV(c.Context(), *groupID, buf, filter); err != nil {
		log.Printf("handler/CSV: failed to export CSV: %v", err)
		return apperrors.ErrInternal
	}

	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", "attachment; filename=transactions.csv")
	return c.Send(buf.Bytes())
}

func (h *ExportHandler) OFX(c *fiber.Ctx) error {
	_ = auth.GetUser(c)
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	filter := parseExportFilter(c)

	buf := &bytes.Buffer{}
	if err := h.exportService.ExportTransactionsOFX(c.Context(), *groupID, buf, filter); err != nil {
		log.Printf("handler/CSV: failed to export OFX: %v", err)
		return apperrors.ErrInternal
	}

	c.Set("Content-Type", "application/x-ofx")
	c.Set("Content-Disposition", "attachment; filename=transactions.ofx")
	return c.Send(buf.Bytes())
}

func (h *ExportHandler) Reconcile(c *fiber.Ctx) error {
	_ = auth.GetUser(c)
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	var req []service.CSVRow
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewValidationError(map[string][]string{"body": {"invalid JSON"}})
	}

	result, err := h.exportService.Reconcile(c.Context(), *groupID, req)
	if err != nil {
		log.Printf("handler/CSV: failed to reconcile: %v", err)
		return apperrors.ErrInternal
	}

	return c.JSON(fiber.Map{"data": fiber.Map{
		"type": "reconciliation",
		"attributes": fiber.Map{
			"matched":       result.Matched,
			"unmatched":     result.Unmatched,
			"total_checked": result.TotalChecked,
		},
	}})
}

// parseExportFilter builds a TransactionFilter from export query parameters.
func parseExportFilter(c *fiber.Ctx) repository.TransactionFilter {
	var filter repository.TransactionFilter
	if v := c.Query("start"); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			filter.DateFrom = &t
		}
	}
	if v := c.Query("end"); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			filter.DateTo = &t
		}
	}
	if v := c.Query("wallet_id"); v != "" {
		if id, err := uuid.Parse(v); err == nil {
			filter.WalletID = &id
		}
	}
	return filter
}
