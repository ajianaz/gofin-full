package request

import "fmt"

// PageFilter holds pagination parameters for list endpoints.
type PageFilter struct {
	Page    int
	PerPage int
}

// Defaults applies default values (page=1, perPage=20) if zero/invalid.
func (f *PageFilter) Defaults() {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.PerPage < 1 || f.PerPage > 100 {
		f.PerPage = 20
	}
}

// Offset returns the SQL OFFSET value.
func (f *PageFilter) Offset() int {
	return (f.Page - 1) * f.PerPage
}

// TotalPages computes total pages from a count.
func (f *PageFilter) TotalPages(total int) int {
	pages := total / f.PerPage
	if total%f.PerPage > 0 {
		pages++
	}
	return pages
}

// Meta returns the standard pagination metadata for API responses.
func (f *PageFilter) Meta(total, count int) map[string]interface{} {
	return map[string]interface{}{
		"pagination": map[string]interface{}{
			"total":        int64(total),
			"count":        count,
			"per_page":     f.PerPage,
			"current_page": f.Page,
			"total_pages":  f.TotalPages(total),
		},
	}
}

// CountQuery wraps a base WHERE clause with COUNT(*) for pagination.
func (f *PageFilter) CountQuery(countSQL string, countArgs ...interface{}) string {
	return countSQL
}

// ApplyLimitOffset appends LIMIT/OFFSET to a query with numbered args.
// argN is the next available argument number. Returns the appended clause and the new argN.
func (f *PageFilter) ApplyLimitOffset(argN int) (clause string, newArgN int) {
	clause = fmt.Sprintf(" LIMIT $%d OFFSET $%d", argN, argN+1)
	return clause, argN + 2
}
