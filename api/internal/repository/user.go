package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ajianaz/gofin-full/api/internal/auth"
	"github.com/ajianaz/gofin-full/api/internal/domain"
)

// UserRepository handles user data access.
type UserRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository creates a new user repository.
func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// Create inserts a new user and creates their default group.
func (r *UserRepository) Create(ctx context.Context, email, passwordHash string) (*domain.User, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	now := time.Now().UTC()

	// Create default group for user
	var groupID int64
	err = tx.QueryRow(ctx,
		`INSERT INTO user_groups (title, created_at, updated_at) VALUES ($1, $2, $3) RETURNING id`,
		email, now, now,
	).Scan(&groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to create user group: %w", err)
	}

	// Get the 'owner' group role
	var ownerRoleID int64
	err = tx.QueryRow(ctx,
		`SELECT id FROM user_roles WHERE title = 'owner' LIMIT 1`,
	).Scan(&ownerRoleID)
	if err != nil {
		return nil, fmt.Errorf("failed to find owner role: %w", err)
	}

	// Create user
	var userID int64
	err = tx.QueryRow(ctx,
		`INSERT INTO users (email, password, user_group_id, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		email, passwordHash, groupID, now, now,
	).Scan(&userID)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Add user to group as owner
	_, err = tx.Exec(ctx,
		`INSERT INTO group_memberships (user_id, user_group_id, user_role_id, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		userID, groupID, ownerRoleID, now, now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create group membership: %w", err)
	}

	// Get the global 'owner' role
	var globalOwnerRoleID int64
	err = tx.QueryRow(ctx,
		`SELECT id FROM roles WHERE title = 'owner' LIMIT 1`,
	).Scan(&globalOwnerRoleID)
	if err != nil {
		return nil, fmt.Errorf("failed to find global owner role: %w", err)
	}

	// Only assign global owner role to the very first user
	var userCount int
	err = tx.QueryRow(ctx,
		`SELECT COUNT(*) FROM users WHERE deleted_at IS NULL`,
	).Scan(&userCount)
	if err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}
	if userCount == 1 {
		_, err = tx.Exec(ctx,
			`INSERT INTO role_user (user_id, role_id) VALUES ($1, $2)`,
			userID, globalOwnerRoleID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to assign owner role: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit: %w", err)
	}

	return &domain.User{
		ID:          userID,
		Email:       email,
		UserGroupID: &groupID,
	}, nil
}

// FindByID finds a user by ID (soft-delete aware).
func (r *UserRepository) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	var u domain.User
	var deletedAt *time.Time

	err := r.db.QueryRow(ctx,
		`SELECT id, email, password, blocked, user_group_id, created_at, updated_at, deleted_at
		 FROM users WHERE id = $1`, id,
	).Scan(&u.ID, &u.Email, &u.Password, &u.Blocked, &u.UserGroupID, &u.CreatedAt, &u.UpdatedAt, &deletedAt)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	if deletedAt != nil {
		return nil, fmt.Errorf("user not found")
	}

	return &u, nil
}

// FindByEmail finds a user by email (soft-delete aware).
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var u domain.User
	var deletedAt *time.Time

	err := r.db.QueryRow(ctx,
		`SELECT id, email, password, blocked, user_group_id, created_at, updated_at, deleted_at
		 FROM users WHERE email = $1`, email,
	).Scan(&u.ID, &u.Email, &u.Password, &u.Blocked, &u.UserGroupID, &u.CreatedAt, &u.UpdatedAt, &deletedAt)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	if deletedAt != nil {
		return nil, fmt.Errorf("user not found")
	}

	return &u, nil
}

// Update updates user fields.
func (r *UserRepository) Update(ctx context.Context, id int64, email, password string) error {
	if password != "" {
		_, err := r.db.Exec(ctx,
			`UPDATE users SET email = $1, password = $2, updated_at = $3 WHERE id = $4 AND deleted_at IS NULL`,
			email, password, time.Now().UTC(), id,
		)
		return err
	}

	_, err := r.db.Exec(ctx,
		`UPDATE users SET email = $1, updated_at = $2 WHERE id = $3 AND deleted_at IS NULL`,
		email, time.Now().UTC(), id,
	)
	return err
}

// Exists checks if a non-deleted user exists.
func (r *UserRepository) Exists(ctx context.Context) (bool, error) {
	var count int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM users WHERE deleted_at IS NULL`,
	).Scan(&count)
	return count > 0, err
}

// HasGlobalRole checks if a user has a specific global role.
func (r *UserRepository) HasGlobalRole(ctx context.Context, userID int64, roleTitle string) (bool, error) {
	var count int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM role_user ru
		 JOIN roles r ON r.id = ru.role_id
		 WHERE ru.user_id = $1 AND r.title = $2`,
		userID, roleTitle,
	).Scan(&count)
	return count > 0, err
}

// GetUserRoleInGroup returns the user's role in a specific group.
func (r *UserRepository) GetUserRoleInGroup(ctx context.Context, userID, groupID int64) (auth.GroupRole, error) {
	var roleTitle string
	err := r.db.QueryRow(ctx,
		`SELECT ur.title FROM group_memberships gm
		 JOIN user_roles ur ON ur.id = gm.user_role_id
		 WHERE gm.user_id = $1 AND gm.user_group_id = $2`,
		userID, groupID,
	).Scan(&roleTitle)
	if err != nil {
		return "", fmt.Errorf("no membership found")
	}
	return auth.GroupRole(roleTitle), nil
}

// SetActiveGroup updates the user's active group.
func (r *UserRepository) SetActiveGroup(ctx context.Context, userID, groupID int64) error {
	// Verify membership
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM group_memberships WHERE user_id = $1 AND user_group_id = $2)`,
		userID, groupID,
	).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("user is not a member of this group")
	}

	_, err = r.db.Exec(ctx,
		`UPDATE users SET user_group_id = $1, updated_at = $2 WHERE id = $3 AND deleted_at IS NULL`,
		groupID, time.Now().UTC(), userID,
	)
	return err
}

// ListAll returns all users (admin use).
func (r *UserRepository) ListAll(ctx context.Context) ([]domain.User, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, email, blocked, user_group_id, created_at, updated_at
			 FROM users WHERE deleted_at IS NULL ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Email, &u.Blocked, &u.UserGroupID, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}
