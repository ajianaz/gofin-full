package handler

import (
	"github.com/rs/zerolog/log"
	"context"

	"github.com/google/uuid"
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

	// Separate UUIDs from codes to avoid PostgreSQL type mismatch
	var uuids []string
	var codes []string
	for _, id := range unique {
		if _, err := uuid.Parse(id); err == nil {
			uuids = append(uuids, id)
		} else {
			codes = append(codes, id)
		}
	}

	// Build query dynamically based on what we have
	var conditions []string
	var queryArgs []interface{}
	idx := 1

	if len(uuids) > 0 {
		ph := ""
		for i, u := range uuids {
			if i > 0 {
				ph += ","
			}
			ph += domain.Placeholder(idx)
			queryArgs = append(queryArgs, u)
			idx++
		}
		conditions = append(conditions, "id IN ("+ph+")")
	}

	if len(codes) > 0 {
		ph := ""
		for i, c := range codes {
			if i > 0 {
				ph += ","
			}
			ph += domain.Placeholder(idx)
			queryArgs = append(queryArgs, c)
			idx++
		}
		conditions = append(conditions, "code IN ("+ph+")")
	}

	if len(conditions) == 0 {
		return result
	}

	where := conditions[0]
	if len(conditions) > 1 {
		where = "(" + conditions[0] + " OR " + conditions[1] + ")"
	}

	query := `SELECT id, code, symbol, COALESCE(decimal_places, 2)
		  FROM currencies
		  WHERE ` + where + ` AND deleted_at IS NULL`

	rows, err := r.db.Query(ctx, query, queryArgs...)
	if err != nil {
		log.Error().Err(err).Msg("currency: resolve failed")
		return result
	}
	defer rows.Close()

	// Map by both uuid and code so any input format resolves
	for rows.Next() {
		var uuid, code, symbol string
		var dp int
		if err := rows.Scan(&uuid, &code, &symbol, &dp); err != nil {
			log.Error().Err(err).Msg("currency: scan failed")
			continue
		}
		info := CurrencyInfo{Code: code, Symbol: symbol, DecimalPlaces: dp}
		result[uuid] = info
		if code != "" {
			result[code] = info
		}
	}
	if err := rows.Err(); err != nil {
		log.Error().Err(err).Msg("currency: rows iteration failed")
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
