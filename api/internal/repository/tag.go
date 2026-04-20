package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ajianaz/gofin-full/api/internal/domain"
)

type TagRepository struct {
	db *pgxpool.Pool
}

func NewTagRepository(db *pgxpool.Pool) *TagRepository {
	return &TagRepository{db: db}
}

func (r *TagRepository) Create(ctx context.Context, userID, groupID int64, tag string, date *time.Time) (*domain.Tag, error) {
	now := time.Now().UTC()
	var t domain.Tag
	err := r.db.QueryRow(ctx,
		`INSERT INTO tags (user_id, user_group_id, tag, date, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6) RETURNING id, user_id, user_group_id, tag, date, created_at, updated_at`,
		userID, groupID, tag, date, now, now,
	).Scan(&t.ID, &t.UserID, &t.UserGroupID, &t.Tag, &t.Date, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create tag: %w", err)
	}
	return &t, nil
}

func (r *TagRepository) FindByID(ctx context.Context, id, groupID int64) (*domain.Tag, error) {
	var t domain.Tag
	var deletedAt *time.Time
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, user_group_id, tag, date, created_at, updated_at, deleted_at
		 FROM tags WHERE id = $1 AND user_group_id = $2`, id, groupID,
	).Scan(&t.ID, &t.UserID, &t.UserGroupID, &t.Tag, &t.Date, &t.CreatedAt, &t.UpdatedAt, &deletedAt)
	if err != nil {
		return nil, fmt.Errorf("tag not found: %w", err)
	}
	if deletedAt != nil {
		return nil, fmt.Errorf("tag not found")
	}
	return &t, nil
}

func (r *TagRepository) List(ctx context.Context, groupID int64) ([]domain.Tag, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, user_group_id, tag, date, created_at, updated_at
		 FROM tags WHERE user_group_id = $1 AND deleted_at IS NULL ORDER BY tag`, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}
	defer rows.Close()

	var tags []domain.Tag
	for rows.Next() {
		var t domain.Tag
		if err := rows.Scan(&t.ID, &t.UserID, &t.UserGroupID, &t.Tag, &t.Date, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	return tags, rows.Err()
}

func (r *TagRepository) Update(ctx context.Context, id, groupID int64, tag string, date *time.Time) error {
	_, err := r.db.Exec(ctx,
		`UPDATE tags SET tag = $1, date = $2, updated_at = $3
		 WHERE id = $4 AND user_group_id = $5 AND deleted_at IS NULL`,
		tag, date, time.Now().UTC(), id, groupID)
	return err
}

func (r *TagRepository) Delete(ctx context.Context, id, groupID int64) error {
	_, err := r.db.Exec(ctx,
		`UPDATE tags SET deleted_at = $1 WHERE id = $2 AND user_group_id = $3 AND deleted_at IS NULL`,
		time.Now().UTC(), id, groupID)
	return err
}

// FindOrCreate returns an existing tag or creates one.
func (r *TagRepository) FindOrCreate(ctx context.Context, userID, groupID int64, tag string) (*domain.Tag, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, user_group_id, tag, date, created_at, updated_at
		 FROM tags WHERE user_id = $1 AND user_group_id = $2 AND tag = $3 AND deleted_at IS NULL LIMIT 1`,
		userID, groupID, tag)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var t domain.Tag
		if err := rows.Scan(&t.ID, &t.UserID, &t.UserGroupID, &t.Tag, &t.Date, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		return &t, nil
	}
	return r.Create(ctx, userID, groupID, tag, nil)
}
