package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ajianaz/gofin-full/api/internal/domain"
)

// UserGroupRepository handles user group data access.
type UserGroupRepository struct {
	db *pgxpool.Pool
}

// NewUserGroupRepository creates a new user group repository.
func NewUserGroupRepository(db *pgxpool.Pool) *UserGroupRepository {
	return &UserGroupRepository{db: db}
}

// Create creates a new user group.
func (r *UserGroupRepository) Create(ctx context.Context, title string) (*domain.UserGroup, error) {
	now := time.Now().UTC()
	var id uuid.UUID
	err := r.db.QueryRow(ctx,
		`INSERT INTO user_groups (title, created_at, updated_at) VALUES ($1, $2, $3) RETURNING id`,
		title, now, now,
	).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("failed to create group: %w", err)
	}

	return &domain.UserGroup{
		ID:        id,
		Title:     title,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// FindByID finds a group by ID (soft-delete aware).
func (r *UserGroupRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.UserGroup, error) {
	var g domain.UserGroup
	var deletedAt *time.Time

	err := r.db.QueryRow(ctx,
		`SELECT id, title, created_at, updated_at, deleted_at
		 FROM user_groups WHERE id = $1`, id,
	).Scan(&g.ID, &g.Title, &g.CreatedAt, &g.UpdatedAt, &deletedAt)
	if err != nil {
		return nil, fmt.Errorf("group not found: %w", err)
	}
	if deletedAt != nil {
		return nil, fmt.Errorf("group not found")
	}

	return &g, nil
}

// List returns all groups a user is a member of.
func (r *UserGroupRepository) List(ctx context.Context, userID uuid.UUID) ([]domain.UserGroup, error) {
	rows, err := r.db.Query(ctx,
		`SELECT ug.id, ug.title, ug.created_at, ug.updated_at
		 FROM user_groups ug
		 JOIN group_memberships gm ON gm.user_group_id = ug.id
		 WHERE gm.user_id = $1 AND ug.deleted_at IS NULL
		 ORDER BY ug.title`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list groups: %w", err)
	}
	defer rows.Close()

	var groups []domain.UserGroup
	for rows.Next() {
		var g domain.UserGroup
		if err := rows.Scan(&g.ID, &g.Title, &g.CreatedAt, &g.UpdatedAt); err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}

	return groups, rows.Err()
}

// Update updates a group's title.
func (r *UserGroupRepository) Update(ctx context.Context, id uuid.UUID, title string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE user_groups SET title = $1, updated_at = $2 WHERE id = $3 AND deleted_at IS NULL`,
		title, time.Now().UTC(), id,
	)
	return err
}

// Delete soft-deletes a group.
func (r *UserGroupRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`UPDATE user_groups SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`,
		time.Now().UTC(), id,
	)
	return err
}

// IsMember checks if a user is a member of a group.
func (r *UserGroupRepository) IsMember(ctx context.Context, userID, groupID uuid.UUID) (bool, error) {
	var count int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM group_memberships WHERE user_id = $1 AND user_group_id = $2`,
		userID, groupID,
	).Scan(&count)
	return count > 0, err
}

// AddMember adds a user to a group with the owner role.
func (r *UserGroupRepository) AddMember(ctx context.Context, userID, groupID uuid.UUID) error {
	now := time.Now().UTC()
	// Resolve the owner role ID
	var ownerRoleID uuid.UUID
	err := r.db.QueryRow(ctx, `SELECT id FROM user_roles WHERE title = 'owner' LIMIT 1`).Scan(&ownerRoleID)
	if err != nil {
		return fmt.Errorf("owner role not found: %w", err)
	}

	_, err = r.db.Exec(ctx,
		`INSERT INTO group_memberships (user_id, user_group_id, user_role_id, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5)
		 ON CONFLICT (user_id, user_group_id) DO NOTHING`,
		userID, groupID, ownerRoleID, now, now,
	)
	return err
}
