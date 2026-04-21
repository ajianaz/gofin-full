package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ajianaz/gofin-full/api/internal/domain"
)

type WebhookRepository struct {
	db *pgxpool.Pool
}

func NewWebhookRepository(db *pgxpool.Pool) *WebhookRepository {
	return &WebhookRepository{db: db}
}

func (r *WebhookRepository) Create(ctx context.Context, userID, groupID uuid.UUID, title, url string) (*domain.Webhook, error) {
	now := time.Now().UTC()
	var w domain.Webhook
	err := r.db.QueryRow(ctx,
		`INSERT INTO webhooks (user_id, user_group_id, title, url, active, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,TRUE,$5,$6)
		 RETURNING id, user_id, user_group_id, title, url, active, created_at, updated_at`,
		userID, groupID, title, url, now, now,
	).Scan(&w.ID, &w.UserID, &w.UserGroupID, &w.Title, &w.URL, &w.Active, &w.CreatedAt, &w.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create webhook: %w", err)
	}
	return &w, nil
}

func (r *WebhookRepository) FindByID(ctx context.Context, id, groupID uuid.UUID) (*domain.Webhook, error) {
	var w domain.Webhook
	var deletedAt *time.Time
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, user_group_id, title, url, active, created_at, updated_at, deleted_at
		 FROM webhooks WHERE id = $1 AND user_group_id = $2`, id, groupID,
	).Scan(&w.ID, &w.UserID, &w.UserGroupID, &w.Title, &w.URL, &w.Active, &w.CreatedAt, &w.UpdatedAt, &deletedAt)
	if err != nil {
		return nil, fmt.Errorf("webhook not found: %w", err)
	}
	if deletedAt != nil {
		return nil, fmt.Errorf("webhook not found")
	}
	return &w, nil
}

func (r *WebhookRepository) List(ctx context.Context, groupID uuid.UUID) ([]domain.Webhook, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, user_group_id, title, url, active, created_at, updated_at
		 FROM webhooks WHERE user_group_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC`, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to list webhooks: %w", err)
	}
	defer rows.Close()

	var webhooks []domain.Webhook
	for rows.Next() {
		var w domain.Webhook
		if err := rows.Scan(&w.ID, &w.UserID, &w.UserGroupID, &w.Title, &w.URL, &w.Active, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, err
		}
		webhooks = append(webhooks, w)
	}
	return webhooks, rows.Err()
}

func (r *WebhookRepository) Update(ctx context.Context, id, groupID uuid.UUID, title, url string, active *bool) error {
	_, err := r.db.Exec(ctx,
		`UPDATE webhooks SET title = COALESCE(NULLIF($1, ''), title), url = COALESCE(NULLIF($2, ''), url), active = COALESCE($3, active), updated_at = $4
		 WHERE id = $5 AND user_group_id = $6 AND deleted_at IS NULL`,
		title, url, active, time.Now().UTC(), id, groupID)
	return err
}

func (r *WebhookRepository) Delete(ctx context.Context, id, groupID uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`UPDATE webhooks SET deleted_at = $1 WHERE id = $2 AND user_group_id = $3 AND deleted_at IS NULL`,
		time.Now().UTC(), id, groupID)
	return err
}

// Trigger operations

func (r *WebhookRepository) SetTriggers(ctx context.Context, webhookID uuid.UUID, triggers []string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `DELETE FROM webhook_triggers WHERE webhook_id = $1`, webhookID); err != nil {
		return err
	}

	now := time.Now().UTC()
	for _, t := range triggers {
		if _, err := tx.Exec(ctx,
			`INSERT INTO webhook_triggers (webhook_id, trigger, created_at, updated_at) VALUES ($1,$2,$3,$4)`,
			webhookID, t, now, now); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (r *WebhookRepository) ListTriggers(ctx context.Context, webhookID uuid.UUID) ([]domain.WebhookTrigger, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, webhook_id, trigger, created_at, updated_at FROM webhook_triggers WHERE webhook_id = $1`, webhookID)
	if err != nil {
		return nil, fmt.Errorf("failed to list triggers: %w", err)
	}
	defer rows.Close()

	var triggers []domain.WebhookTrigger
	for rows.Next() {
		var t domain.WebhookTrigger
		if err := rows.Scan(&t.ID, &t.WebhookID, &t.Trigger, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		triggers = append(triggers, t)
	}
	return triggers, rows.Err()
}

// Message operations

func (r *WebhookRepository) CreateMessage(ctx context.Context, webhookID uuid.UUID, message string) (*domain.WebhookMessage, error) {
	now := time.Now().UTC()
	var m domain.WebhookMessage
	err := r.db.QueryRow(ctx,
		`INSERT INTO webhook_messages (webhook_id, message, created_at, updated_at)
		 VALUES ($1,$2,$3,$4)
		 RETURNING id, webhook_id, message, created_at, updated_at`,
		webhookID, message, now, now,
	).Scan(&m.ID, &m.WebhookID, &m.Message, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create webhook message: %w", err)
	}
	return &m, nil
}

func (r *WebhookRepository) ListMessages(ctx context.Context, webhookID uuid.UUID) ([]domain.WebhookMessage, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, webhook_id, message, created_at, updated_at
		 FROM webhook_messages WHERE webhook_id = $1 ORDER BY created_at DESC`, webhookID)
	if err != nil {
		return nil, fmt.Errorf("failed to list webhook messages: %w", err)
	}
	defer rows.Close()

	var messages []domain.WebhookMessage
	for rows.Next() {
		var m domain.WebhookMessage
		if err := rows.Scan(&m.ID, &m.WebhookID, &m.Message, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, rows.Err()
}
