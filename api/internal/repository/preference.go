package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ajianaz/gofin-full/api/internal/domain"
)

type PreferenceRepository struct {
	db *pgxpool.Pool
}

func NewPreferenceRepository(db *pgxpool.Pool) *PreferenceRepository {
	return &PreferenceRepository{db: db}
}

func (r *PreferenceRepository) Set(ctx context.Context, userID int64, name, data string) (*domain.Preference, error) {
	now := time.Now().UTC()
	var p domain.Preference
	err := r.db.QueryRow(ctx,
		`INSERT INTO preferences (user_id, name, data, created_at, updated_at) VALUES ($1,$2,$3,$4,$5)
		 ON CONFLICT (user_id, name) DO UPDATE SET data = $3, updated_at = $5
		 RETURNING id, user_id, name, data, created_at, updated_at`,
		userID, name, data, now, now,
	).Scan(&p.ID, &p.UserID, &p.Name, &p.Data, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to set preference: %w", err)
	}
	return &p, nil
}

func (r *PreferenceRepository) Get(ctx context.Context, userID int64, name string) (*domain.Preference, error) {
	var p domain.Preference
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, name, data, created_at, updated_at FROM preferences WHERE user_id = $1 AND name = $2`,
		userID, name,
	).Scan(&p.ID, &p.UserID, &p.Name, &p.Data, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("preference not found: %w", err)
	}
	return &p, nil
}

func (r *PreferenceRepository) List(ctx context.Context, userID int64) ([]domain.Preference, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, name, data, created_at, updated_at FROM preferences WHERE user_id = $1 ORDER BY name`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list preferences: %w", err)
	}
	defer rows.Close()

	var prefs []domain.Preference
	for rows.Next() {
		var p domain.Preference
		if err := rows.Scan(&p.ID, &p.UserID, &p.Name, &p.Data, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		prefs = append(prefs, p)
	}
	return prefs, rows.Err()
}

func (r *PreferenceRepository) Delete(ctx context.Context, userID int64, name string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM preferences WHERE user_id = $1 AND name = $2`, userID, name)
	return err
}

// Configuration (system-level, not user-scoped)

type ConfigurationRepository struct {
	db *pgxpool.Pool
}

func NewConfigurationRepository(db *pgxpool.Pool) *ConfigurationRepository {
	return &ConfigurationRepository{db: db}
}

func (r *ConfigurationRepository) Set(ctx context.Context, name, value string) (*domain.Configuration, error) {
	now := time.Now().UTC()
	var c domain.Configuration
	err := r.db.QueryRow(ctx,
		`INSERT INTO configurations (name, value, created_at, updated_at) VALUES ($1,$2,$3,$4)
		 ON CONFLICT (name) DO UPDATE SET value = $2, updated_at = $4
		 RETURNING id, name, value, created_at, updated_at`,
		name, value, now, now,
	).Scan(&c.ID, &c.Name, &c.Value, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to set configuration: %w", err)
	}
	return &c, nil
}

func (r *ConfigurationRepository) Get(ctx context.Context, name string) (*domain.Configuration, error) {
	var c domain.Configuration
	err := r.db.QueryRow(ctx,
		`SELECT id, name, value, created_at, updated_at FROM configurations WHERE name = $1`, name,
	).Scan(&c.ID, &c.Name, &c.Value, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("configuration not found: %w", err)
	}
	return &c, nil
}

func (r *ConfigurationRepository) List(ctx context.Context) ([]domain.Configuration, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, name, value, created_at, updated_at FROM configurations ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("failed to list configurations: %w", err)
	}
	defer rows.Close()

	var configs []domain.Configuration
	for rows.Next() {
		var c domain.Configuration
		if err := rows.Scan(&c.ID, &c.Name, &c.Value, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		configs = append(configs, c)
	}
	return configs, rows.Err()
}
