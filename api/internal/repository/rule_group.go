package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azfirazka/gofin-full/api/internal/domain"
)

type RuleGroupRepository struct {
	db *pgxpool.Pool
}

func NewRuleGroupRepository(db *pgxpool.Pool) *RuleGroupRepository {
	return &RuleGroupRepository{db: db}
}

func (r *RuleGroupRepository) Create(ctx context.Context, userID, groupID int64, title string, order int) (*domain.RuleGroup, error) {
	now := time.Now().UTC()
	var rg domain.RuleGroup
	err := r.db.QueryRow(ctx,
		`INSERT INTO rule_groups (user_id, user_group_id, title, active, "order", created_at, updated_at)
		 VALUES ($1,$2,$3,TRUE,$4,$5,$6)
		 RETURNING id, user_id, user_group_id, title, active, "order", created_at, updated_at`,
		userID, groupID, title, order, now, now,
	).Scan(&rg.ID, &rg.UserID, &rg.UserGroupID, &rg.Title, &rg.Active, &rg.Order, &rg.CreatedAt, &rg.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create rule group: %w", err)
	}
	return &rg, nil
}

func (r *RuleGroupRepository) List(ctx context.Context, groupID int64) ([]domain.RuleGroup, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, user_group_id, title, active, "order", created_at, updated_at
		 FROM rule_groups WHERE user_group_id = $1 AND deleted_at IS NULL ORDER BY "order", title`, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to list rule groups: %w", err)
	}
	defer rows.Close()

	var groups []domain.RuleGroup
	for rows.Next() {
		var rg domain.RuleGroup
		if err := rows.Scan(&rg.ID, &rg.UserID, &rg.UserGroupID, &rg.Title, &rg.Active, &rg.Order, &rg.CreatedAt, &rg.UpdatedAt); err != nil {
			return nil, err
		}
		groups = append(groups, rg)
	}
	return groups, rows.Err()
}

func (r *RuleGroupRepository) FindByID(ctx context.Context, id, groupID int64) (*domain.RuleGroup, error) {
	var rg domain.RuleGroup
	var deletedAt *time.Time
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, user_group_id, title, active, "order", created_at, updated_at, deleted_at
		 FROM rule_groups WHERE id = $1 AND user_group_id = $2`, id, groupID,
	).Scan(&rg.ID, &rg.UserID, &rg.UserGroupID, &rg.Title, &rg.Active, &rg.Order, &rg.CreatedAt, &rg.UpdatedAt, &deletedAt)
	if err != nil {
		return nil, fmt.Errorf("rule group not found: %w", err)
	}
	if deletedAt != nil {
		return nil, fmt.Errorf("rule group not found")
	}
	return &rg, nil
}

func (r *RuleGroupRepository) Update(ctx context.Context, id, groupID int64, title string, active *bool) error {
	_, err := r.db.Exec(ctx,
		`UPDATE rule_groups SET title = COALESCE(NULLIF($1, ''), title), active = COALESCE($2, active), updated_at = $3
		 WHERE id = $4 AND user_group_id = $5 AND deleted_at IS NULL`,
		title, active, time.Now().UTC(), id, groupID)
	return err
}

func (r *RuleGroupRepository) Delete(ctx context.Context, id, groupID int64) error {
	_, err := r.db.Exec(ctx,
		`UPDATE rule_groups SET deleted_at = $1 WHERE id = $2 AND user_group_id = $3 AND deleted_at IS NULL`,
		time.Now().UTC(), id, groupID)
	return err
}
