package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azfirazka/gofin-full/api/internal/domain"
)

type CategoryRepository struct {
	db *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(ctx context.Context, userID, groupID int64, name string) (*domain.Category, error) {
	now := time.Now().UTC()
	var c domain.Category
	err := r.db.QueryRow(ctx,
		`INSERT INTO categories (user_id, user_group_id, name, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5) RETURNING id, user_id, user_group_id, name, created_at, updated_at`,
		userID, groupID, name, now, now,
	).Scan(&c.ID, &c.UserID, &c.UserGroupID, &c.Name, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}
	return &c, nil
}

func (r *CategoryRepository) FindByID(ctx context.Context, id, groupID int64) (*domain.Category, error) {
	var c domain.Category
	var deletedAt *time.Time
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, user_group_id, name, created_at, updated_at, deleted_at
		 FROM categories WHERE id = $1 AND user_group_id = $2`, id, groupID,
	).Scan(&c.ID, &c.UserID, &c.UserGroupID, &c.Name, &c.CreatedAt, &c.UpdatedAt, &deletedAt)
	if err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}
	if deletedAt != nil {
		return nil, fmt.Errorf("category not found")
	}
	return &c, nil
}

func (r *CategoryRepository) List(ctx context.Context, groupID int64) ([]domain.Category, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, user_group_id, name, created_at, updated_at
		 FROM categories WHERE user_group_id = $1 AND deleted_at IS NULL ORDER BY name`, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}
	defer rows.Close()

	var categories []domain.Category
	for rows.Next() {
		var c domain.Category
		if err := rows.Scan(&c.ID, &c.UserID, &c.UserGroupID, &c.Name, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, rows.Err()
}

func (r *CategoryRepository) Update(ctx context.Context, id, groupID int64, name string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE categories SET name = $1, updated_at = $2
		 WHERE id = $3 AND user_group_id = $4 AND deleted_at IS NULL`,
		name, time.Now().UTC(), id, groupID)
	return err
}

func (r *CategoryRepository) Delete(ctx context.Context, id, groupID int64) error {
	_, err := r.db.Exec(ctx,
		`UPDATE categories SET deleted_at = $1 WHERE id = $2 AND user_group_id = $3 AND deleted_at IS NULL`,
		time.Now().UTC(), id, groupID)
	return err
}
