package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
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
	ID        int64     `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	GroupID   uuid.UUID `json:"group_id"`
	Action    string    `json:"action"`
	Entity    string    `json:"entity"`
	EntityID  uuid.UUID `json:"entity_id"`
	OldValue  *string   `json:"old_value,omitempty"`
	NewValue  *string   `json:"new_value,omitempty"`
	IPAddress *string   `json:"ip_address,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UserEmail string    `json:"user_email,omitempty"`
}

func (r *AuditRepository) Log(ctx context.Context, userID, groupID uuid.UUID, action, entity string, entityID uuid.UUID, oldValue, newValue, ipAddress *string) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO audit_logs (user_id, user_group_id, action, entity_type, entity_id, old_value, new_value, ip_address, created_at)
			 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		userID, groupID, action, entity, entityID, oldValue, newValue, ipAddress, time.Now().UTC())
	return err
}

func (r *AuditRepository) List(ctx context.Context, groupID uuid.UUID, entityType string, entityID uuid.UUID, limit int) ([]AuditLog, error) {
	rows, err := r.db.Query(ctx,
		`SELECT a.id, a.user_id, a.user_group_id, a.action, a.entity_type, a.entity_id,
		        a.old_value, a.new_value, a.ip_address, a.created_at,
		        COALESCE(u.email, '') as user_email
			 FROM audit_logs a
			 LEFT JOIN users u ON u.id = a.user_id
			 WHERE a.user_group_id = $1 AND ($2 = '' OR a.entity_type = $2) AND ($3 = uuid_nil() OR a.entity_id = $3)
			 ORDER BY a.created_at DESC LIMIT $4`,
		groupID, entityType, entityID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list audit logs: %w", err)
	}
	defer rows.Close()

	var logs []AuditLog
	for rows.Next() {
		var log AuditLog
		if err := rows.Scan(&log.ID, &log.UserID, &log.GroupID, &log.Action, &log.Entity, &log.EntityID, &log.OldValue, &log.NewValue, &log.IPAddress, &log.CreatedAt, &log.UserEmail); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, rows.Err()
}
