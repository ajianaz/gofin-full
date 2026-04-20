package testhelpers

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"github.com/azfirazka/gofin-full/api/internal/auth"
)

// SeedData holds the IDs and tokens created by SeedTestData.
type SeedData struct {
	OwnerUserID    int64
	OwnerEmail     string
	OwnerToken     string
	ReadOnlyUserID int64
	ReadOnlyEmail  string
	ReadOnlyToken  string
	FullUserID     int64
	FullEmail      string
	FullToken      string
	TxUserID       int64
	TxEmail        string
	TxUserToken    string
	GroupID        int64
	WalletID       int64
}

// SetupTestDB connects to the test database and runs all pending up-migrations.
// It mirrors cmd/migrate/main.go logic (extractGooseSection, schema_migrations table).
func SetupTestDB(cfg *TestConfig) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping: %w", err)
	}

	ensureMigrationsTable(ctx, pool)

	// Resolve migrations dir relative to project root (where go test runs from).
	migrationsDir := filepath.Join("..", "..", "migrations", "postgres")
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		// Fallback: try from project root
		migrationsDir = "migrations/postgres"
	}
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("read migrations dir: %w", err)
	}

	var migrations []string
	for _, f := range files {
		name := f.Name()
		if strings.HasSuffix(name, ".up.sql") {
			migrations = append(migrations, filepath.Join(migrationsDir, name))
		}
	}
	sort.Strings(migrations)

	for _, m := range migrations {
		base := filepath.Base(m)
		if isMigrationApplied(ctx, pool, base) {
			continue
		}

		raw, err := os.ReadFile(m)
		if err != nil {
			pool.Close()
			return nil, fmt.Errorf("read migration %s: %w", base, err)
		}

		sql := extractGooseSection(string(raw), "up")
		if _, err := pool.Exec(ctx, sql); err != nil {
			pool.Close()
			return nil, fmt.Errorf("apply migration %s: %w", base, err)
		}

		recordMigration(ctx, pool, base)
	}

	return pool, nil
}

// TruncateAllTables removes all rows from data tables, preserving reference data
// (roles, user_roles) that are seeded by migrations.
func TruncateAllTables(db *pgxpool.Pool) {
	ctx := context.Background()
	tables := []string{
		"wallet_members",
		"wallets",
		"oauth_states",
		"api_keys",
		"refresh_tokens",
		"role_user",
		"group_memberships",
		"users",
		"user_groups",
	}
	for _, t := range tables {
		_, _ = db.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE", t))
	}
}

// SeedTestData creates a full set of test fixtures: owner user with group and wallet,
// plus additional users with different group roles.
func SeedTestData(db *pgxpool.Pool, jwtMgr *auth.JWTManager) (*SeedData, error) {
	ctx := context.Background()

	// ---------- resolve user_role IDs ----------
	var ownerRoleID, readOnlyRoleID, fullRoleID, txRoleID int64
	if err := db.QueryRow(ctx, `SELECT id FROM user_roles WHERE title = 'owner'`).Scan(&ownerRoleID); err != nil {
		return nil, fmt.Errorf("find owner role: %w", err)
	}
	if err := db.QueryRow(ctx, `SELECT id FROM user_roles WHERE title = 'read_only'`).Scan(&readOnlyRoleID); err != nil {
		return nil, fmt.Errorf("find read_only role: %w", err)
	}
	if err := db.QueryRow(ctx, `SELECT id FROM user_roles WHERE title = 'full'`).Scan(&fullRoleID); err != nil {
		return nil, fmt.Errorf("find full role: %w", err)
	}
	if err := db.QueryRow(ctx, `SELECT id FROM user_roles WHERE title = 'manage_transactions'`).Scan(&txRoleID); err != nil {
		return nil, fmt.Errorf("find manage_transactions role: %w", err)
	}

	// ---------- resolve global role IDs ----------
	var globalOwnerRoleID int64
	if err := db.QueryRow(ctx, `SELECT id FROM roles WHERE title = 'owner' LIMIT 1`).Scan(&globalOwnerRoleID); err != nil {
		return nil, fmt.Errorf("find global owner role: %w", err)
	}

	// ---------- create user group ----------
	var groupID int64
	err := db.QueryRow(ctx,
		`INSERT INTO user_groups (title, created_at, updated_at) VALUES ($1, NOW(), NOW()) RETURNING id`,
		"Test Group",
	).Scan(&groupID)
	if err != nil {
		return nil, fmt.Errorf("insert group: %w", err)
	}

	// ---------- insert users ----------
	type userSpec struct {
		email  string
		roleID int64
		userID *int64
		token  *string
	}

	specs := []userSpec{
		{email: "test@gofin.io", roleID: ownerRoleID},
		{email: "readonly_user@gofin.io", roleID: readOnlyRoleID},
		{email: "full_user@gofin.io", roleID: fullRoleID},
		{email: "tx_user@gofin.io", roleID: txRoleID},
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	for i := range specs {
		var uid int64
		err := db.QueryRow(ctx,
			`INSERT INTO users (email, password, user_group_id, created_at, updated_at)
			 VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id`,
			specs[i].email, string(hashedPassword), groupID,
		).Scan(&uid)
		if err != nil {
			return nil, fmt.Errorf("insert user %s: %w", specs[i].email, err)
		}

		// group membership
		_, err = db.Exec(ctx,
			`INSERT INTO group_memberships (user_id, user_group_id, user_role_id, created_at, updated_at)
			 VALUES ($1, $2, $3, NOW(), NOW())`,
			uid, groupID, specs[i].roleID,
		)
		if err != nil {
			return nil, fmt.Errorf("insert membership for %s: %w", specs[i].email, err)
		}

		// Assign global owner role to the first user (owner)
		if i == 0 {
			_, _ = db.Exec(ctx,
				`INSERT INTO role_user (user_id, role_id) VALUES ($1, $2)`,
				uid, globalOwnerRoleID,
			)
		}

		specs[i].userID = &uid

		// generate JWT token
		identity := &auth.UserIdentity{ID: uid, Email: specs[i].email}
		pair, err := jwtMgr.GenerateTokenPair(identity, &groupID)
		if err != nil {
			return nil, fmt.Errorf("generate token for %s: %w", specs[i].email, err)
		}
		specs[i].token = &pair.AccessToken
	}

	// ---------- create wallet for owner ----------
	var walletID int64
	err = db.QueryRow(ctx,
		`INSERT INTO wallets (user_id, user_group_id, name, account_type, active, virtual_balance, include_net_worth, created_at, updated_at)
		 VALUES ($1, $2, 'Test Wallet', 'asset', TRUE, 0, TRUE, NOW(), NOW()) RETURNING id`,
		*specs[0].userID, groupID,
	).Scan(&walletID)
	if err != nil {
		return nil, fmt.Errorf("insert wallet: %w", err)
	}

	return &SeedData{
		OwnerUserID:    *specs[0].userID,
		OwnerEmail:     specs[0].email,
		OwnerToken:     *specs[0].token,
		ReadOnlyUserID: *specs[1].userID,
		ReadOnlyEmail:  specs[1].email,
		ReadOnlyToken:  *specs[1].token,
		FullUserID:     *specs[2].userID,
		FullEmail:      specs[2].email,
		FullToken:      *specs[2].token,
		TxUserID:       *specs[3].userID,
		TxEmail:        specs[3].email,
		TxUserToken:    *specs[3].token,
		GroupID:        groupID,
		WalletID:       walletID,
	}, nil
}

// ---------- internal migration helpers ----------

// ensureMigrationsTable creates the schema_migrations tracking table if absent.
func ensureMigrationsTable(ctx context.Context, pool *pgxpool.Pool) {
	_, _ = pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			name       TEXT PRIMARY KEY,
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

	section = strings.ReplaceAll(section, "-- +goose StatementBegin", "")
	section = strings.ReplaceAll(section, "-- +goose StatementEnd", "")
	return strings.TrimSpace(section)
}
