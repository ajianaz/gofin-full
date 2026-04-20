package handler

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/azfirazka/gofin-full/api/internal/auth"
	"github.com/azfirazka/gofin-full/api/internal/domain"
	"github.com/azfirazka/gofin-full/api/internal/repository"
	"github.com/azfirazka/gofin-full/api/internal/service"
	apperrors "github.com/azfirazka/gofin-full/api/pkg/errors"
)

type TransactionHandler struct {
	txService *service.TransactionService
	txRepo    *repository.TransactionRepository
}

func NewTransactionHandler(txService *service.TransactionService, txRepo *repository.TransactionRepository) *TransactionHandler {
	return &TransactionHandler{txService: txService, txRepo: txRepo}
}

func (h *TransactionHandler) Index(c *fiber.Ctx) error {
	_ = auth.GetUser(c)
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	filter := repository.TransactionFilter{
		Page:    c.QueryInt("page", 1),
		PerPage: c.QueryInt("per_page", 50),
	}

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
	if v := c.Query("type"); v != "" {
		filter.Type = v
	}
	if v := c.Query("wallet_id"); v != "" {
		if id, err := strconv.ParseInt(v, 10, 64); err == nil {
			filter.WalletID = &id
		}
	}

	groups, total, err := h.txRepo.ListGroups(c.Context(), *groupID, filter)
	if err != nil {
		return apperrors.NewWithDetail(500, "failed to list transactions", err.Error())
	}

	var data []fiber.Map
	for _, g := range groups {
		data = append(data, fiber.Map{
			"type":       "transactions",
			"id":         g.ID,
			"attributes": fiber.Map{
				"group_title": g.GroupTitle,
				"created_at":  g.CreatedAt,
				"updated_at":  g.UpdatedAt,
			},
		})
	}

	totalPages := int(total) / filter.PerPage
	if int(total)%filter.PerPage > 0 {
		totalPages++
	}

	return c.JSON(fiber.Map{
		"data": data,
		"meta": fiber.Map{
			"pagination": fiber.Map{
				"total":        total,
				"count":        len(data),
				"per_page":     filter.PerPage,
				"current_page": filter.Page,
				"total_pages":  totalPages,
			},
		},
	})
}

func (h *TransactionHandler) Show(c *fiber.Ctx) error {
	_ = auth.GetUser(c)
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"id": {"invalid id"}})
	}

	group, err := h.txRepo.FindGroupByID(c.Context(), int64(id), *groupID)
	if err != nil {
		return apperrors.NotFoundResource("transaction", int64(id))
	}

	return c.JSON(fiber.Map{"data": transactionGroupToMap(group)})
}

func (h *TransactionHandler) Store(c *fiber.Ctx) error {
	user := auth.GetUser(c)
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	var input service.CreateTransactionInput
	if err := c.BodyParser(&input); err != nil {
		return apperrors.NewValidationError(map[string][]string{"body": {"invalid JSON"}})
	}

	fieldErrors := make(map[string][]string)
	if input.Type == "" {
		fieldErrors["type"] = append(fieldErrors["type"], "type is required")
	}
	if input.Amount == "" {
		fieldErrors["amount"] = append(fieldErrors["amount"], "amount is required")
	}
	if input.SourceID == 0 {
		fieldErrors["source_id"] = append(fieldErrors["source_id"], "source_id is required")
	}
	if input.DestinationID == 0 {
		fieldErrors["destination_id"] = append(fieldErrors["destination_id"], "destination_id is required")
	}
	if input.Date.IsZero() {
		input.Date = time.Now().UTC()
	}
	if len(fieldErrors) > 0 {
		return apperrors.NewValidationError(fieldErrors)
	}

	if input.CurrencyID == "" {
		input.CurrencyID = "EUR"
	}

	result, err := h.txService.CreateTransaction(c.Context(), user.ID, *groupID, input)
	if err != nil {
		return apperrors.NewWithDetail(422, "failed to create transaction", err.Error())
	}

	return c.Status(201).JSON(fiber.Map{"data": fiber.Map{
		"type":       "transactions",
		"id":         result.GroupID,
		"attributes": fiber.Map{"journal_id": result.JournalID},
	}})
}

// StoreSplit handles POST /transactions/split for multi-journal transactions.
func (h *TransactionHandler) StoreSplit(c *fiber.Ctx) error {
	user := auth.GetUser(c)
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	var req struct {
		Type     string                      `json:"type"`
		Date     time.Time                   `json:"date"`
		Title    string                      `json:"group_title"`
		Journals []service.SplitJournalInput `json:"journals"`
	}
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewValidationError(map[string][]string{"body": {"invalid JSON"}})
	}

	fieldErrors := make(map[string][]string)
	if req.Type == "" {
		fieldErrors["type"] = append(fieldErrors["type"], "type is required")
	}
	if len(req.Journals) < 2 {
		fieldErrors["journals"] = append(fieldErrors["journals"], "at least 2 journals required")
	}
	if req.Date.IsZero() {
		req.Date = time.Now().UTC()
	}
	if len(fieldErrors) > 0 {
		return apperrors.NewValidationError(fieldErrors)
	}

	result, err := h.txService.CreateSplitTransaction(c.Context(), user.ID, *groupID, req.Type, req.Date, req.Title, req.Journals)
	if err != nil {
		return apperrors.NewWithDetail(422, "failed to create split transaction", err.Error())
	}

	return c.Status(201).JSON(fiber.Map{"data": fiber.Map{
		"type":       "transactions",
		"id":         result.GroupID,
		"attributes": fiber.Map{"journal_id": result.JournalID},
	}})
}

func (h *TransactionHandler) Update(c *fiber.Ctx) error {
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
		Description string     `json:"description"`
		Date        *time.Time `json:"date"`
		Notes       *string    `json:"notes"`
		CategoryIDs []int64    `json:"category_ids"`
		TagIDs      []int64    `json:"tag_ids"`
	}
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewValidationError(map[string][]string{"body": {"invalid JSON"}})
	}

	if err := h.txRepo.UpdateJournal(c.Context(), int64(id), *groupID, req.Description, req.Date, req.Notes); err != nil {
		return apperrors.NotFoundResource("transaction", int64(id))
	}

	if req.CategoryIDs != nil {
		h.txRepo.SetJournalCategories(c.Context(), int64(id), req.CategoryIDs)
	}
	if req.TagIDs != nil {
		h.txRepo.SetJournalTags(c.Context(), int64(id), req.TagIDs)
	}

	return c.JSON(fiber.Map{"data": fiber.Map{
		"type": "transactions", "id": id,
		"attributes": fiber.Map{
			"description": req.Description,
			"notes":       req.Notes,
		},
	}})
}

func (h *TransactionHandler) Delete(c *fiber.Ctx) error {
	user := auth.GetUser(c)
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"id": {"invalid id"}})
	}

	if err := h.txService.DeleteTransaction(c.Context(), int64(id), user.ID, *groupID); err != nil {
		return apperrors.NotFoundResource("transaction", int64(id))
	}

	return c.Status(204).Send(nil)
}

func transactionGroupToMap(g *domain.TransactionGroup) fiber.Map {
	var journals []fiber.Map
	for _, j := range g.Journals {
		journal := fiber.Map{
			"transaction_journal_id": j.ID,
			"type":                  string(j.Type),
			"date":                  j.Date,
			"description":           j.Description,
			"currency_id":           j.CurrencyID,
			"reconciled":            j.Reconciled,
			"notes":                 j.Notes,
			"created_at":            j.CreatedAt,
			"updated_at":            j.UpdatedAt,
		}

		if len(j.SourceTransactions) > 0 {
			st := j.SourceTransactions[0]
			journal["source_id"] = st.AccountID
			journal["amount"] = st.Amount.StringFixed(2)
		}
		if len(j.DestinationTransactions) > 0 {
			dt := j.DestinationTransactions[0]
			journal["destination_id"] = dt.AccountID
		}

		journals = append(journals, journal)
	}

	return fiber.Map{
		"type":       "transactions",
		"id":         g.ID,
		"attributes": fiber.Map{
			"group_title":  g.GroupTitle,
			"transactions": journals,
			"created_at":   g.CreatedAt,
			"updated_at":   g.UpdatedAt,
		},
	}
}
