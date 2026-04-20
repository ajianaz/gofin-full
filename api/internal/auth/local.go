package auth

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"github.com/azfirazka/gofin-full/api/internal/config"
)

// localProvider implements email + password authentication.
type localProvider struct {
	db *pgxpool.Pool
}

func newLocalProvider(cfg *config.Config) *localProvider {
	return &localProvider{}
}

// SetDB sets the database pool (called after DB init).
func (p *localProvider) SetDB(db *pgxpool.Pool) { p.db = db }

func (p *localProvider) Name() string { return "local" }

func (p *localProvider) AuthURL(_ string) string { return "" }

func (p *localProvider) Authenticate(ctx context.Context, creds Credentials) (*UserIdentity, error) {
	if p.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	if creds.Email == "" || creds.Password == "" {
		return nil, fmt.Errorf("email and password required")
	}

	var id int64
	var hashedPassword string
	var blocked bool
	var deletedAt *string

	err := p.db.QueryRow(ctx,
		`SELECT id, password, blocked, deleted_at::text
			 FROM users WHERE email = $1`, creds.Email,
	).Scan(&id, &hashedPassword, &blocked, &deletedAt)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}
	if deletedAt != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(creds.Password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return &UserIdentity{
		ID:       id,
		Email:    creds.Email,
		Blocked:  blocked,
		DemoUser: false,
	}, nil
}

// HashPassword hashes a plain-text password using bcrypt.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(bytes), nil
}

// CheckPassword compares a plain-text password against a bcrypt hash.
func CheckPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
