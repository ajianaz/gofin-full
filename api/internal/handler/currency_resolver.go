package handler

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ajianaz/gofin-full/api/internal/domain"
)

// CurrencyInfo holds resolved currency display data.
type CurrencyInfo struct {
	Code          string `json:"currency_code"`
	Symbol        string `json:"currency_symbol"`
	DecimalPlaces int    `json:"currency_decimal_places"`
}

// CurrencyResolver resolves currency IDs to display info.
type CurrencyResolver struct {
	db *pgxpool.Pool
}

// NewCurrencyResolver creates a new resolver backed by the DB pool.
func NewCurrencyResolver(db *pgxpool.Pool) *CurrencyResolver {
	return &CurrencyResolver{db: db}
}

// ResolveMany resolves multiple currency IDs in a single query.
// Returns a map[currencyID]CurrencyInfo.
func (r *CurrencyResolver) ResolveMany(ctx context.Context, ids []string) map[string]CurrencyInfo {
	result := make(map[string]CurrencyInfo, len(ids))
	if len(ids) == 0 {
		return result
	}

	// Deduplicate
	seen := make(map[string]bool, len(ids))
	var unique []string
	for _, id := range ids {
		if id != "" && !seen[id] {
			seen[id] = true
			unique = append(unique, id)
		}
	}
	if len(unique) == 0 {
		return result
	}

	// Build query with positional args
	placeholders := ""
	args := []interface{}{}
	for i, id := range unique {
		if i > 0 {
			placeholders += ","
		}
		placeholders += domain.Placeholder(i + 1)
		args = append(args, id)
	}

	query := `SELECT id, code, symbol, COALESCE(decimal_places, 2)
			  FROM currencies WHERE id IN (` + placeholders + `)
			  AND deleted_at IS NULL`

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		log.Printf("currency: resolve failed: %v", err)
		return result
	}
	defer rows.Close()

	for rows.Next() {
		var id, code, symbol string
		var dp int
		if err := rows.Scan(&id, &code, &symbol, &dp); err != nil {
			log.Printf("currency: scan failed: %v", err)
			continue
		}
		result[id] = CurrencyInfo{Code: code, Symbol: symbol, DecimalPlaces: dp}
	}
	return result
}

// ResolveSingle resolves a single currency ID.
func (r *CurrencyResolver) ResolveSingle(ctx context.Context, id string) CurrencyInfo {
	if id == "" {
		return CurrencyInfo{}
	}
	m := r.ResolveMany(ctx, []string{id})
	return m[id]
}
