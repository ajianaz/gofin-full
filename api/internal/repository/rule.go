package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ajianaz/gofin-full/api/internal/domain"
)

type RuleRepository struct {
	db *pgxpool.Pool
}

func NewRuleRepository(db *pgxpool.Pool) *RuleRepository {
	return &RuleRepository{db: db}
}

func (r *RuleRepository) Create(ctx context.Context, userID, groupID uuid.UUID, title string, priority int, ruleGroupID *uuid.UUID) (*domain.Rule, error) {
	now := time.Now().UTC()
	var rule domain.Rule
	err := r.db.QueryRow(ctx,
		`INSERT INTO rules (user_id, user_group_id, rule_group_id, title, priority, active, strict, stop_processing, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,TRUE,FALSE,FALSE,$6,$7)
		 RETURNING id, user_id, user_group_id, rule_group_id, title, priority, active, strict, stop_processing, created_at, updated_at`,
		userID, groupID, ruleGroupID, title, priority, now, now,
	).Scan(&rule.ID, &rule.UserID, &rule.UserGroupID, &rule.RuleGroupID, &rule.Title, &rule.Priority, &rule.Active, &rule.Strict, &rule.StopProcessing, &rule.CreatedAt, &rule.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create rule: %w", err)
	}
	return &rule, nil
}

func (r *RuleRepository) FindByID(ctx context.Context, id, groupID uuid.UUID) (*domain.Rule, error) {
	var rule domain.Rule
	var deletedAt *time.Time
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, user_group_id, rule_group_id, title, priority, active, strict, stop_processing, created_at, updated_at, deleted_at
		 FROM rules WHERE id = $1 AND user_group_id = $2`, id, groupID,
	).Scan(&rule.ID, &rule.UserID, &rule.UserGroupID, &rule.RuleGroupID, &rule.Title, &rule.Priority, &rule.Active, &rule.Strict, &rule.StopProcessing, &rule.CreatedAt, &rule.UpdatedAt, &deletedAt)
	if err != nil {
		return nil, fmt.Errorf("rule not found: %w", err)
	}
	if deletedAt != nil {
		return nil, fmt.Errorf("rule not found")
	}

	triggers, _ := r.findTriggers(ctx, rule.ID)
	actions, _ := r.findActions(ctx, rule.ID)
	rule.Triggers = triggers
	rule.Actions = actions
	return &rule, nil
}

func (r *RuleRepository) List(ctx context.Context, groupID uuid.UUID) ([]domain.Rule, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, user_group_id, rule_group_id, title, priority, active, strict, stop_processing, created_at, updated_at
		 FROM rules WHERE user_group_id = $1 AND deleted_at IS NULL ORDER BY priority, title`, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to list rules: %w", err)
	}
	defer rows.Close()

	var rules []domain.Rule
	for rows.Next() {
		var rule domain.Rule
		if err := rows.Scan(&rule.ID, &rule.UserID, &rule.UserGroupID, &rule.RuleGroupID, &rule.Title, &rule.Priority, &rule.Active, &rule.Strict, &rule.StopProcessing, &rule.CreatedAt, &rule.UpdatedAt); err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	return rules, rows.Err()
}

func (r *RuleRepository) Update(ctx context.Context, id, groupID uuid.UUID, title string, active *bool, strict *bool, stopProcessing *bool) error {
	_, err := r.db.Exec(ctx,
		`UPDATE rules SET
		  title = COALESCE(NULLIF($1, ''), title),
		  active = COALESCE($2, active),
		  strict = COALESCE($3, strict),
		  stop_processing = COALESCE($4, stop_processing),
		  updated_at = $5
		 WHERE id = $6 AND user_group_id = $7 AND deleted_at IS NULL`,
		title, active, strict, stopProcessing, time.Now().UTC(), id, groupID)
	return err
}

func (r *RuleRepository) Delete(ctx context.Context, id, groupID uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`UPDATE rules SET deleted_at = $1 WHERE id = $2 AND user_group_id = $3 AND deleted_at IS NULL`,
		time.Now().UTC(), id, groupID)
	return err
}

func (r *RuleRepository) SetTriggers(ctx context.Context, ruleID uuid.UUID, triggers []domain.RuleTrigger) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `DELETE FROM rule_triggers WHERE rule_id = $1`, ruleID)
	if err != nil {
		return err
	}
	for _, t := range triggers {
		_, err = tx.Exec(ctx,
			`INSERT INTO rule_triggers (rule_id, trigger_type, trigger_value, stop_processing, created_at, updated_at)
			 VALUES ($1,$2,$3,$4,$5,$6)`,
			ruleID, t.TriggerType, t.TriggerValue, t.StopProcessing, time.Now().UTC(), time.Now().UTC())
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (r *RuleRepository) SetActions(ctx context.Context, ruleID uuid.UUID, actions []domain.RuleAction) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `DELETE FROM rule_actions WHERE rule_id = $1`, ruleID)
	if err != nil {
		return err
	}
	for _, a := range actions {
		_, err = tx.Exec(ctx,
			`INSERT INTO rule_actions (rule_id, action_type, action_value, "order", stop_processing, created_at, updated_at)
			 VALUES ($1,$2,$3,$4,$5,$6,$7)`,
			ruleID, a.ActionType, a.ActionValue, a.Order, a.StopProcessing, time.Now().UTC(), time.Now().UTC())
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (r *RuleRepository) findTriggers(ctx context.Context, ruleID uuid.UUID) ([]domain.RuleTrigger, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, rule_id, trigger_type, trigger_value, stop_processing, created_at, updated_at
		 FROM rule_triggers WHERE rule_id = $1`, ruleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var triggers []domain.RuleTrigger
	for rows.Next() {
		var t domain.RuleTrigger
		if err := rows.Scan(&t.ID, &t.RuleID, &t.TriggerType, &t.TriggerValue, &t.StopProcessing, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		triggers = append(triggers, t)
	}
	return triggers, rows.Err()
}

func (r *RuleRepository) findActions(ctx context.Context, ruleID uuid.UUID) ([]domain.RuleAction, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, rule_id, action_type, action_value, "order", stop_processing, created_at, updated_at
		 FROM rule_actions WHERE rule_id = $1 ORDER BY "order"`, ruleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var actions []domain.RuleAction
	for rows.Next() {
		var a domain.RuleAction
		if err := rows.Scan(&a.ID, &a.RuleID, &a.ActionType, &a.ActionValue, &a.Order, &a.StopProcessing, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		actions = append(actions, a)
	}
	return actions, rows.Err()
}
