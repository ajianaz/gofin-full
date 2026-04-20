package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azfirazka/gofin-full/api/internal/domain"
)

type ObjectGroupRepository struct {
	db *pgxpool.Pool
}

func NewObjectGroupRepository(db *pgxpool.Pool) *ObjectGroupRepository {
	return &ObjectGroupRepository{db: db}
}

func (r *ObjectGroupRepository) Create(ctx context.Context, userID, groupID int64, title string, order int) (*domain.ObjectGroup, error) {
	now := time.Now().UTC()
	var og domain.ObjectGroup
	err := r.db.QueryRow(ctx,
		`INSERT INTO object_groups (user_id, user_group_id, title, "order", created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6)
		 RETURNING id, user_id, user_group_id, title, "order", created_at, updated_at`,
		userID, groupID, title, order, now, now,
	).Scan(&og.ID, &og.UserID, &og.UserGroupID, &og.Title, &og.Order, &og.CreatedAt, &og.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create object group: %w", err)
	}
	return &og, nil
}

func (r *ObjectGroupRepository) List(ctx context.Context, groupID int64) ([]domain.ObjectGroup, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, user_group_id, title, "order", created_at, updated_at
		 FROM object_groups WHERE user_group_id = $1 ORDER BY "order", title`, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to list object groups: %w", err)
	}
	defer rows.Close()

	var groups []domain.ObjectGroup
	for rows.Next() {
		var og domain.ObjectGroup
		if err := rows.Scan(&og.ID, &og.UserID, &og.UserGroupID, &og.Title, &og.Order, &og.CreatedAt, &og.UpdatedAt); err != nil {
			return nil, err
		}
		groups = append(groups, og)
	}
	return groups, rows.Err()
}

func (r *ObjectGroupRepository) FindByID(ctx context.Context, id, groupID int64) (*domain.ObjectGroup, error) {
	var og domain.ObjectGroup
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, user_group_id, title, "order", created_at, updated_at
		 FROM object_groups WHERE id = $1 AND user_group_id = $2`, id, groupID,
	).Scan(&og.ID, &og.UserID, &og.UserGroupID, &og.Title, &og.Order, &og.CreatedAt, &og.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("object group not found: %w", err)
	}
	return &og, nil
}

func (r *ObjectGroupRepository) Update(ctx context.Context, id, groupID int64, title string, order *int) error {
	_, err := r.db.Exec(ctx,
		`UPDATE object_groups SET title = COALESCE(NULLIF($1, ''), title), "order" = COALESCE($2, "order"), updated_at = $3
		 WHERE id = $4 AND user_group_id = $5`,
		title, order, time.Now().UTC(), id, groupID)
	return err
}

func (r *ObjectGroupRepository) Delete(ctx context.Context, id, groupID int64) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM object_groups WHERE id = $1 AND user_group_id = $2`, id, groupID)
	return err
}
