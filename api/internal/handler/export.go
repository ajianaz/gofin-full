package handler

import (
	"bytes"
	"github.com/gofiber/fiber/v2"

	"github.com/ajianaz/gofin-full/api/internal/auth"
	"github.com/ajianaz/gofin-full/api/internal/service"
	apperrors "github.com/ajianaz/gofin-full/api/pkg/errors"
)

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

	buf := &bytes.Buffer{}
	if err := h.exportService.ExportTransactionsCSV(c.Context(), *groupID, buf); err != nil {
		return apperrors.NewWithDetail(500, "failed to export CSV", err.Error())
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

	buf := &bytes.Buffer{}
	if err := h.exportService.ExportTransactionsOFX(c.Context(), *groupID, buf); err != nil {
		return apperrors.NewWithDetail(500, "failed to export OFX", err.Error())
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
		return apperrors.NewWithDetail(500, "failed to reconcile", err.Error())
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
