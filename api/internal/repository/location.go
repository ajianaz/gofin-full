package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ajianaz/gofin-full/api/internal/domain"
)

type LocationRepository struct {
	db *pgxpool.Pool
}

func NewLocationRepository(db *pgxpool.Pool) *LocationRepository {
	return &LocationRepository{db: db}
}

func (r *LocationRepository) Set(ctx context.Context, userID, groupID uuid.UUID, locatableType string, locatableID uuid.UUID, latitude, longitude *float64, zoomLevel int) (*domain.Location, error) {
	now := time.Now().UTC()
	var loc domain.Location
	err := r.db.QueryRow(ctx,
		`INSERT INTO locations (user_id, user_group_id, locatable_type, locatable_id, latitude, longitude, zoom_level, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		 RETURNING id, user_id, user_group_id, locatable_type, locatable_id, latitude, longitude, zoom_level, created_at, updated_at`,
		userID, groupID, locatableType, locatableID, latitude, longitude, zoomLevel, now, now,
	).Scan(&loc.ID, &loc.UserID, &loc.GroupID, &loc.LocatableType, &loc.LocatableID, &loc.Latitude, &loc.Longitude, &loc.ZoomLevel, &loc.CreatedAt, &loc.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to set location: %w", err)
	}
	return &loc, nil
}

func (r *LocationRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Location, error) {
	var loc domain.Location
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, user_group_id, locatable_type, locatable_id, latitude, longitude, zoom_level, created_at, updated_at
		 FROM locations WHERE id = $1`,
		id,
	).Scan(&loc.ID, &loc.UserID, &loc.GroupID, &loc.LocatableType, &loc.LocatableID, &loc.Latitude, &loc.Longitude, &loc.ZoomLevel, &loc.CreatedAt, &loc.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("location not found: %w", err)
	}
	return &loc, nil
}

func (r *LocationRepository) GetByEntity(ctx context.Context, locatableType string, locatableID uuid.UUID) (*domain.Location, error) {
	var loc domain.Location
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, user_group_id, locatable_type, locatable_id, latitude, longitude, zoom_level, created_at, updated_at
		 FROM locations WHERE locatable_type = $1 AND locatable_id = $2`,
		locatableType, locatableID,
	).Scan(&loc.ID, &loc.UserID, &loc.GroupID, &loc.LocatableType, &loc.LocatableID, &loc.Latitude, &loc.Longitude, &loc.ZoomLevel, &loc.CreatedAt, &loc.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("location not found: %w", err)
	}
	return &loc, nil
}

func (r *LocationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM locations WHERE id = $1`, id)
	return err
}
