package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ajianaz/gofin-full/api/internal/domain"
)

type NotificationRepository struct {
	db *pgxpool.Pool
}

func NewNotificationRepository(db *pgxpool.Pool) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(ctx context.Context, userID int64, channel, notifType, title, message string) (*domain.Notification, error) {
	now := time.Now().UTC()
	var n domain.Notification
	err := r.db.QueryRow(ctx,
		`INSERT INTO notifications (user_id, channel, type, title, message, "read", created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,FALSE,$6,$7)
		 RETURNING id, user_id, channel, type, title, message, "read", created_at, updated_at`,
		userID, channel, notifType, title, message, now, now,
	).Scan(&n.ID, &n.UserID, &n.Channel, &n.Type, &n.Title, &n.Message, &n.Read, &n.CreatedAt, &n.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}
	return &n, nil
}

func (r *NotificationRepository) List(ctx context.Context, userID int64) ([]domain.Notification, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, channel, type, title, message, "read", created_at, updated_at
		 FROM notifications WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list notifications: %w", err)
	}
	defer rows.Close()

	var notifications []domain.Notification
	for rows.Next() {
		var n domain.Notification
		if err := rows.Scan(&n.ID, &n.UserID, &n.Channel, &n.Type, &n.Title, &n.Message, &n.Read, &n.CreatedAt, &n.UpdatedAt); err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}
	return notifications, rows.Err()
}

func (r *NotificationRepository) ListUnread(ctx context.Context, userID int64) ([]domain.Notification, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, channel, type, title, message, "read", created_at, updated_at
		 FROM notifications WHERE user_id = $1 AND "read" = FALSE ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list unread notifications: %w", err)
	}
	defer rows.Close()

	var notifications []domain.Notification
	for rows.Next() {
		var n domain.Notification
		if err := rows.Scan(&n.ID, &n.UserID, &n.Channel, &n.Type, &n.Title, &n.Message, &n.Read, &n.CreatedAt, &n.UpdatedAt); err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}
	return notifications, rows.Err()
}

func (r *NotificationRepository) MarkRead(ctx context.Context, id, userID int64) error {
	_, err := r.db.Exec(ctx,
		`UPDATE notifications SET "read" = TRUE, updated_at = $1 WHERE id = $2 AND user_id = $3`,
		time.Now().UTC(), id, userID)
	return err
}

func (r *NotificationRepository) MarkAllRead(ctx context.Context, userID int64) error {
	_, err := r.db.Exec(ctx,
		`UPDATE notifications SET "read" = TRUE, updated_at = $1 WHERE user_id = $2 AND "read" = FALSE`,
		time.Now().UTC(), userID)
	return err
}
