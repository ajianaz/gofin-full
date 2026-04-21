package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/ajianaz/gofin-full/api/pkg/pgxuuid"
)

func main() {
	dsn := flag.String("dsn", os.Getenv("DATABASE_URL"), "PostgreSQL connection string")
	direction := flag.String("dir", "up", "Migration direction: up or down")
	migrationsDir := flag.String("path", "migrations/postgres", "Path to migrations directory")
	flag.Parse()

	if *dsn == "" {
		log.Fatal("database DSN is required (use -dsn or DATABASE_URL env)")
	}

	ctx := context.Background()
	poolCfg, err := pgxpool.ParseConfig(*dsn)
	if err != nil {
		log.Fatalf("failed to parse database config: %v", err)
	}
	poolCfg.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		conn.TypeMap().RegisterType(&pgtype.Type{
			Name:  "uuid",
			OID:   pgtype.UUIDOID,
			Codec: pgxuuid.Codec{},
		})
		return nil
	}
	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	ensureMigrationsTable(ctx, pool)

	files, err := os.ReadDir(*migrationsDir)
	if err != nil {
		log.Fatalf("failed to read migrations directory: %v", err)
	}

	var migrations []string
	for _, f := range files {
		name := f.Name()
		if strings.HasSuffix(name, "."+*direction+".sql") {
			migrations = append(migrations, filepath.Join(*migrationsDir, name))
		}
	}
	sort.Strings(migrations)

	for _, m := range migrations {
		base := filepath.Base(m)
		applied := isMigrationApplied(ctx, pool, base)

		if *direction == "up" && applied {
			fmt.Printf("skip   %s (already applied)\n", base)
			continue
		}
		if *direction == "down" && !applied {
			fmt.Printf("skip   %s (not applied)\n", base)
			continue
		}

		raw, err := os.ReadFile(m)
		if err != nil {
			log.Fatalf("failed to read migration %s: %v", base, err)
		}

		// Extract only the relevant section from goose-formatted SQL.
		// Goose files contain both Up and Down sections delimited by markers:
		//   -- +goose Up
		//   -- +goose Down
		sql := extractGooseSection(string(raw), *direction)

		if _, err := pool.Exec(ctx, sql); err != nil {
			log.Fatalf("failed to apply migration %s: %v", base, err)
		}

		if *direction == "up" {
			recordMigration(ctx, pool, base)
			fmt.Printf("apply  %s\n", base)
		} else {
			removeMigration(ctx, pool, base)
			fmt.Printf("revert %s\n", base)
		}
	}

	fmt.Println("done")
}

// extractGooseSection extracts the SQL for the given direction ("up" or "down")
// from a goose-formatted migration file. If no markers are found, returns the
// full content (plain SQL migration).
func extractGooseSection(content, direction string) string {
	upMarker := "-- +goose Up"
	downMarker := "-- +goose Down"

	hasUp := strings.Contains(content, upMarker)
	hasDown := strings.Contains(content, downMarker)

	if !hasUp && !hasDown {
		return content
	}

	var startMarker, endMarker string
	if direction == "up" {
		startMarker = upMarker
		if hasDown {
			endMarker = downMarker
		}
	} else {
		startMarker = downMarker
		endMarker = ""
	}

	startIdx := strings.Index(content, startMarker)
	if startIdx == -1 {
		return ""
	}
	section := content[startIdx+len(startMarker):]

	if endMarker != "" {
		endIdx := strings.Index(section, endMarker)
		if endIdx != -1 {
			section = section[:endIdx]
		}
	}

	// Remove goose StatementBegin/StatementEnd comments and trim whitespace.
	section = strings.ReplaceAll(section, "-- +goose StatementBegin", "")
	section = strings.ReplaceAll(section, "-- +goose StatementEnd", "")
	return strings.TrimSpace(section)
}

func ensureMigrationsTable(ctx context.Context, pool *pgxpool.Pool) {
	_, _ = pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			name      TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`)
}

func isMigrationApplied(ctx context.Context, pool *pgxpool.Pool, name string) bool {
	var exists bool
	_ = pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE name = $1)`, name).Scan(&exists)
	return exists
}

func recordMigration(ctx context.Context, pool *pgxpool.Pool, name string) {
	_, _ = pool.Exec(ctx, `INSERT INTO schema_migrations (name, applied_at) VALUES ($1, $2)`, name, time.Now().UTC())
}

func removeMigration(ctx context.Context, pool *pgxpool.Pool, name string) {
	_, _ = pool.Exec(ctx, `DELETE FROM schema_migrations WHERE name = $1`, name)
}
