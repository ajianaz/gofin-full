package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ajianaz/gofin-full/api/internal/domain"
)

type WalletMemberRepository struct {
	db *pgxpool.Pool
}

func NewWalletMemberRepository(db *pgxpool.Pool) *WalletMemberRepository {
	return &WalletMemberRepository{db: db}
}

func (r *WalletMemberRepository) AddMember(ctx context.Context, walletID, userID int64, role string) (*domain.WalletMember, error) {
	now := time.Now().UTC()
	var m domain.WalletMember
	err := r.db.QueryRow(ctx,
		`INSERT INTO wallet_members (wallet_id, user_id, role, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5)
		 ON CONFLICT (wallet_id, user_id) DO UPDATE SET role = $3, updated_at = $5
		 RETURNING id, wallet_id, user_id, role, created_at, updated_at`,
		walletID, userID, role, now, now,
	).Scan(&m.ID, &m.WalletID, &m.UserID, &m.Role, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to add wallet member: %w", err)
	}
	return &m, nil
}

func (r *WalletMemberRepository) ListByWallet(ctx context.Context, walletID int64) ([]domain.WalletMember, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, wallet_id, user_id, role, created_at, updated_at
		 FROM wallet_members WHERE wallet_id = $1 ORDER BY created_at`, walletID)
	if err != nil {
		return nil, fmt.Errorf("failed to list wallet members: %w", err)
	}
	defer rows.Close()

	var members []domain.WalletMember
	for rows.Next() {
		var m domain.WalletMember
		if err := rows.Scan(&m.ID, &m.WalletID, &m.UserID, &m.Role, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, rows.Err()
}

func (r *WalletMemberRepository) FindByWalletAndUser(ctx context.Context, walletID, userID int64) (*domain.WalletMember, error) {
	var m domain.WalletMember
	err := r.db.QueryRow(ctx,
		`SELECT id, wallet_id, user_id, role, created_at, updated_at
		 FROM wallet_members WHERE wallet_id = $1 AND user_id = $2`,
		walletID, userID,
	).Scan(&m.ID, &m.WalletID, &m.UserID, &m.Role, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("wallet member not found: %w", err)
	}
	return &m, nil
}

func (r *WalletMemberRepository) UpdateRole(ctx context.Context, id, walletID int64, role string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE wallet_members SET role = $1, updated_at = $2 WHERE id = $3 AND wallet_id = $4`,
		role, time.Now().UTC(), id, walletID)
	return err
}

func (r *WalletMemberRepository) RemoveMember(ctx context.Context, walletID, userID int64) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM wallet_members WHERE wallet_id = $1 AND user_id = $2`,
		walletID, userID)
	return err
}

// GetWalletRole returns the user's role for a wallet. Returns empty string if not a member.
func (r *WalletMemberRepository) GetWalletRole(ctx context.Context, walletID, userID int64) (string, error) {
	var role string
	err := r.db.QueryRow(ctx,
		`SELECT role FROM wallet_members WHERE wallet_id = $1 AND user_id = $2`,
		walletID, userID,
	).Scan(&role)
	if err != nil {
		return "", nil // not a member
	}
	return role, nil
}

// IsWalletOwner checks if the user is the wallet owner (wallet.user_id matches).
func (r *WalletMemberRepository) IsWalletOwner(ctx context.Context, walletID, userID int64) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM wallets WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL)`,
		walletID, userID,
	).Scan(&exists)
	return exists, err
}
