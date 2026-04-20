package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuditRepository struct {
	db *pgxpool.Pool
}

func NewAuditRepository(db *pgxpool.Pool) *AuditRepository {
	return &AuditRepository{db: db}
}

// AuditLog represents a financial mutation log entry.
type AuditLog struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	GroupID     int64     `json:"group_id"`
	Action      string    `json:"action"`
	Entity      string    `json:"entity"`
	EntityID    int64     `json:"entity_id"`
	OldValue    *string   `json:"old_value,omitempty"`
	NewValue    *string   `json:"new_value,omitempty"`
	IPAddress   *string   `json:"ip_address,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

func (r *AuditRepository) Log(ctx context.Context, userID, groupID int64, action, entity string, entityID int64, oldValue, newValue, ipAddress *string) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO audit_logs (user_id, user_group_id, action, entity_type, entity_id, old_value, new_value, ip_address, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		userID, groupID, action, entity, entityID, oldValue, newValue, ipAddress, time.Now().UTC())
	return err
}

func (r *AuditRepository) List(ctx context.Context, groupID int64, entityType string, entityID int64, limit int) ([]AuditLog, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, user_group_id, action, entity_type, entity_id, old_value, new_value, ip_address, created_at
		 FROM audit_logs
		 WHERE user_group_id = $1 AND ($2 = '' OR entity_type = $2) AND ($3 = 0 OR entity_id = $3)
		 ORDER BY created_at DESC LIMIT $4`,
		groupID, entityType, entityID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list audit logs: %w", err)
	}
	defer rows.Close()

	var logs []AuditLog
	for rows.Next() {
		var log AuditLog
		if err := rows.Scan(&log.ID, &log.UserID, &log.GroupID, &log.Action, &log.Entity, &log.EntityID, &log.OldValue, &log.NewValue, &log.IPAddress, &log.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, rows.Err()
}
