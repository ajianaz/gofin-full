package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ajianaz/gofin-full/api/internal/domain"
)

type NoteRepository struct {
	db *pgxpool.Pool
}

func NewNoteRepository(db *pgxpool.Pool) *NoteRepository {
	return &NoteRepository{db: db}
}

func (r *NoteRepository) Create(ctx context.Context, userID, groupID uuid.UUID, noteableType string, noteableID uuid.UUID, note string) (*domain.Note, error) {
	now := time.Now().UTC()
	var n domain.Note
	err := r.db.QueryRow(ctx,
		`INSERT INTO notes (user_id, user_group_id, noteable_type, noteable_id, note, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7)
		 RETURNING id, user_id, user_group_id, noteable_type, noteable_id, note, created_at, updated_at`,
		userID, groupID, noteableType, noteableID, note, now, now,
	).Scan(&n.ID, &n.UserID, &n.GroupID, &n.NoteableType, &n.NoteableID, &n.Note, &n.CreatedAt, &n.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create note: %w", err)
	}
	return &n, nil
}

func (r *NoteRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Note, error) {
	var n domain.Note
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, user_group_id, noteable_type, noteable_id, note, created_at, updated_at
		 FROM notes WHERE id = $1`,
		id,
	).Scan(&n.ID, &n.UserID, &n.GroupID, &n.NoteableType, &n.NoteableID, &n.Note, &n.CreatedAt, &n.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("note not found: %w", err)
	}
	return &n, nil
}

func (r *NoteRepository) ListByEntity(ctx context.Context, noteableType string, noteableID uuid.UUID) ([]domain.Note, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, user_group_id, noteable_type, noteable_id, note, created_at, updated_at
		 FROM notes WHERE noteable_type = $1 AND noteable_id = $2 ORDER BY created_at DESC`,
		noteableType, noteableID)
	if err != nil {
		return nil, fmt.Errorf("failed to list notes: %w", err)
	}
	defer rows.Close()

	var notes []domain.Note
	for rows.Next() {
		var n domain.Note
		if err := rows.Scan(&n.ID, &n.UserID, &n.GroupID, &n.NoteableType, &n.NoteableID, &n.Note, &n.CreatedAt, &n.UpdatedAt); err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}
	return notes, rows.Err()
}

func (r *NoteRepository) Update(ctx context.Context, id uuid.UUID, note string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE notes SET note = $1, updated_at = $2 WHERE id = $3`,
		note, time.Now().UTC(), id)
	return err
}

func (r *NoteRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM notes WHERE id = $1`, id)
	return err
}
