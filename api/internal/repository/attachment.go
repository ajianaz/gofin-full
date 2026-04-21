package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ajianaz/gofin-full/api/internal/domain"
)

type AttachmentRepository struct {
	db *pgxpool.Pool
}

func NewAttachmentRepository(db *pgxpool.Pool) *AttachmentRepository {
	return &AttachmentRepository{db: db}
}

func (r *AttachmentRepository) Create(ctx context.Context, userID uuid.UUID, attachableType string, attachableID uuid.UUID, filename, mimeType string, size int64) (*domain.Attachment, error) {
	now := time.Now().UTC()
	var a domain.Attachment
	err := r.db.QueryRow(ctx,
		`INSERT INTO attachments (user_id, attachable_type, attachable_id, filename, mime_type, size, uploaded, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,FALSE,$7,$8)
		 RETURNING id, user_id, attachable_type, attachable_id, filename, mime_type, size, uploaded, created_at, updated_at`,
		userID, attachableType, attachableID, filename, mimeType, size, now, now,
	).Scan(&a.ID, &a.UserID, &a.AttachableType, &a.AttachableID, &a.Filename, &a.MimeType, &a.Size, &a.Uploaded, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create attachment: %w", err)
	}
	return &a, nil
}

func (r *AttachmentRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Attachment, error) {
	var a domain.Attachment
	var deletedAt *time.Time
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, attachable_type, attachable_id, filename, mime_type, size, uploaded, created_at, updated_at, deleted_at
		 FROM attachments WHERE id = $1`, id,
	).Scan(&a.ID, &a.UserID, &a.AttachableType, &a.AttachableID, &a.Filename, &a.MimeType, &a.Size, &a.Uploaded, &a.CreatedAt, &a.UpdatedAt, &deletedAt)
	if err != nil {
		return nil, fmt.Errorf("attachment not found: %w", err)
	}
	if deletedAt != nil {
		return nil, fmt.Errorf("attachment not found")
	}
	return &a, nil
}

func (r *AttachmentRepository) ListByEntity(ctx context.Context, attachableType string, attachableID uuid.UUID) ([]domain.Attachment, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, attachable_type, attachable_id, filename, mime_type, size, uploaded, created_at, updated_at
		 FROM attachments WHERE attachable_type = $1 AND attachable_id = $2 AND deleted_at IS NULL ORDER BY created_at DESC`,
		attachableType, attachableID)
	if err != nil {
		return nil, fmt.Errorf("failed to list attachments: %w", err)
	}
	defer rows.Close()

	var attachments []domain.Attachment
	for rows.Next() {
		var a domain.Attachment
		if err := rows.Scan(&a.ID, &a.UserID, &a.AttachableType, &a.AttachableID, &a.Filename, &a.MimeType, &a.Size, &a.Uploaded, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		attachments = append(attachments, a)
	}
	return attachments, rows.Err()
}

func (r *AttachmentRepository) MarkUploaded(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`UPDATE attachments SET uploaded = TRUE, updated_at = $1 WHERE id = $2 AND deleted_at IS NULL`,
		time.Now().UTC(), id)
	return err
}

func (r *AttachmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`UPDATE attachments SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`,
		time.Now().UTC(), id)
	return err
}
