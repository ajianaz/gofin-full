package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type AnalyticsRepository struct {
	db *pgxpool.Pool
}

func NewAnalyticsRepository(db *pgxpool.Pool) *AnalyticsRepository {
	return &AnalyticsRepository{db: db}
}

// CategorySpending represents spending totals grouped by category.
type CategorySpending struct {
	CategoryID   uuid.UUID `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Total        string  `json:"total"`
	Count        int     `json:"count"`
}

// SpendingByCategory returns spending totals grouped by category for a period.
func (r *AnalyticsRepository) SpendingByCategory(ctx context.Context, groupID uuid.UUID, start, end time.Time) ([]CategorySpending, error) {
	rows, err := r.db.Query(ctx,
		`SELECT COALESCE(c.id, 0), COALESCE(c.name, 'Uncategorized'),
		        COUNT(DISTINCT tj.id), COALESCE(SUM(ABS(t.amount)), 0)
		 FROM transaction_journals tj
		 JOIN transactions t ON t.transaction_journal_id = tj.id
		 LEFT JOIN category_transaction ct ON ct.transaction_journal_id = tj.id
		 LEFT JOIN categories c ON c.id = ct.category_id
		 WHERE tj.user_group_id = $1 AND tj.date >= $2 AND tj.date <= $3
		 GROUP BY c.id, c.name
		 ORDER BY SUM(ABS(t.amount)) DESC`,
		groupID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get spending by category: %w", err)
	}
	defer rows.Close()

	var results []CategorySpending
	for rows.Next() {
		var cs CategorySpending
		var total decimal.Decimal
		if err := rows.Scan(&cs.CategoryID, &cs.CategoryName, &cs.Count, &total); err != nil {
			return nil, err
		}
		cs.Total = total.StringFixed(2)
		results = append(results, cs)
	}
	return results, rows.Err()
}

// PeriodSpending represents spending/income totals grouped by time period.
type PeriodSpending struct {
	Period  string `json:"period"`
	Income  string `json:"income"`
	Expense string `json:"expense"`
}

// SpendingByPeriod returns spending/income grouped by month.
func (r *AnalyticsRepository) SpendingByPeriod(ctx context.Context, groupID uuid.UUID, start, end time.Time) ([]PeriodSpending, error) {
	rows, err := r.db.Query(ctx,
		`SELECT TO_CHAR(tj.date, 'YYYY-MM') as period,
		        COALESCE(SUM(CASE WHEN tt.type = 'deposit' THEN ABS(t.amount) ELSE 0 END), 0),
		        COALESCE(SUM(CASE WHEN tt.type = 'withdrawal' THEN ABS(t.amount) ELSE 0 END), 0)
		 FROM transaction_journals tj
		 JOIN transactions t ON t.transaction_journal_id = tj.id
		 JOIN transaction_types tt ON tt.id = tj.transaction_type_id
		 WHERE tj.user_group_id = $1 AND tj.date >= $2 AND tj.date <= $3
		 GROUP BY TO_CHAR(tj.date, 'YYYY-MM')
		 ORDER BY period`,
		groupID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get spending by period: %w", err)
	}
	defer rows.Close()

	var results []PeriodSpending
	for rows.Next() {
		var ps PeriodSpending
		var income, expense decimal.Decimal
		if err := rows.Scan(&ps.Period, &income, &expense); err != nil {
			return nil, err
		}
		ps.Income = income.StringFixed(2)
		ps.Expense = expense.StringFixed(2)
		results = append(results, ps)
	}
	return results, rows.Err()
}

// NetWorthSummary returns income/expense/total summary for a period.
type NetWorthSummary struct {
	TotalIncome  string `json:"total_income"`
	TotalExpense string `json:"total_expense"`
	NetIncome    string `json:"net_income"`
	TransactionCount int `json:"transaction_count"`
}

// GetNetWorth returns a financial summary for a group over a period.
func (r *AnalyticsRepository) GetNetWorth(ctx context.Context, groupID uuid.UUID, start, end time.Time) (*NetWorthSummary, error) {
	var income, expense decimal.Decimal
	var count int
	err := r.db.QueryRow(ctx,
		`SELECT
		        COALESCE(SUM(CASE WHEN tt.type = 'deposit' THEN ABS(t.amount) ELSE 0 END), 0),
		        COALESCE(SUM(CASE WHEN tt.type = 'withdrawal' THEN ABS(t.amount) ELSE 0 END), 0),
		        COUNT(DISTINCT tj.id)
		 FROM transaction_journals tj
		 JOIN transactions t ON t.transaction_journal_id = tj.id
		 JOIN transaction_types tt ON tt.id = tj.transaction_type_id
		 WHERE tj.user_group_id = $1 AND tj.date >= $2 AND tj.date <= $3`,
		groupID, start, end,
	).Scan(&income, &expense, &count)
	if err != nil {
		return nil, fmt.Errorf("failed to get net worth: %w", err)
	}

	return &NetWorthSummary{
		TotalIncome:      income.StringFixed(2),
		TotalExpense:     expense.StringFixed(2),
		NetIncome:       income.Sub(expense).StringFixed(2),
		TransactionCount: count,
	}, nil
}
