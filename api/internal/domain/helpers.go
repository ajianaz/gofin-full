package domain

import "fmt"

// Placeholder returns a PostgreSQL positional placeholder ($1, $2, ...).
func Placeholder(n int) string {
	return fmt.Sprintf("$%d", n)
}
