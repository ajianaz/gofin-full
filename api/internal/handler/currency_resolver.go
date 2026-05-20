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

// ResolveMany resolves multiple currency identifiers in a single query.
// Accepts both UUIDs and currency codes (e.g. "IDR", "USD").
// Returns a map[inputIdentifier]CurrencyInfo so callers can look up by whatever they passed.
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

	// Build query with positional args — match by id OR code
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
		  FROM currencies
		  WHERE (id IN (` + placeholders + `) OR code IN (` + placeholders + `))
		  AND deleted_at IS NULL`

	// Duplicate args for both IN clauses (avoid append mutating args slice)
	allArgs := make([]interface{}, len(args)*2)
	copy(allArgs, args)
	copy(allArgs[len(args):], args)

	rows, err := r.db.Query(ctx, query, allArgs...)
	if err != nil {
		log.Printf("currency: resolve failed: %v", err)
		return result
	}
	defer rows.Close()

	// Map by both uuid and code so any input format resolves
	for rows.Next() {
		var uuid, code, symbol string
		var dp int
		if err := rows.Scan(&uuid, &code, &symbol, &dp); err != nil {
			log.Printf("currency: scan failed: %v", err)
			continue
		}
		info := CurrencyInfo{Code: code, Symbol: symbol, DecimalPlaces: dp}
		result[uuid] = info
		if code != "" {
			result[code] = info
		}
	}
	return result
}

// ResolveSingle resolves a single currency identifier (UUID or code).
func (r *CurrencyResolver) ResolveSingle(ctx context.Context, id string) CurrencyInfo {
	if id == "" {
		return CurrencyInfo{}
	}
	m := r.ResolveMany(ctx, []string{id})
	return m[id]
}
