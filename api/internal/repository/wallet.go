package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"

	"github.com/ajianaz/gofin-full/api/internal/domain"
)

// WalletRepository handles wallet data access.
type WalletRepository struct {
	db *pgxpool.Pool
}

// NewWalletRepository creates a new wallet repository.
func NewWalletRepository(db *pgxpool.Pool) *WalletRepository {
	return &WalletRepository{db: db}
}

// Create inserts a new wallet.
func (r *WalletRepository) Create(ctx context.Context, w *domain.Wallet) (*domain.Wallet, error) {
	now := time.Now().UTC()
	var id uuid.UUID

	err := r.db.QueryRow(ctx,
		`INSERT INTO wallets (user_id, user_group_id, name, account_type, iban, bic,
		  currency_id, active, virtual_balance, include_net_worth,
		  latitude, longitude, liability_type, liability_direction,
		  interest_rate, interest_period, current_debt,
		  credit_card_type, monthly_payment_date, monthly_payment_amount, notes,
		  created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23)
		 RETURNING id`,
		w.UserID, w.UserGroupID, w.Name, w.AccountType,
		w.IBAN, w.BIC, w.CurrencyID, w.Active, w.VirtualBalance, w.IncludeNetWorth,
		w.Latitude, w.Longitude, w.LiabilityType, w.LiabilityDirection,
		w.InterestRate, w.InterestPeriod, w.CurrentDebt,
		w.CreditCardType, w.MonthlyPaymentDate, w.MonthlyPaymentAmount, w.Notes,
		now, now,
	).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("failed to create wallet: %w", err)
	}

	w.ID = id
	w.CreatedAt = now
	w.UpdatedAt = now
	return w, nil
}

// FindByID finds a wallet by ID within a group (soft-delete aware).
func (r *WalletRepository) FindByID(ctx context.Context, id, groupID uuid.UUID) (*domain.Wallet, error) {
	var w domain.Wallet
	var deletedAt *time.Time
	var iban, bic, currencyID, liabilityType, liabilityDirection, interestPeriod, creditCardType sqlString
	var notes sqlString
	var lat, long *float64
	var interestRate, currentDebt, monthlyPaymentAmt *decimal.Decimal
	var monthlyPaymentDate *time.Time
	var active, includeNetWorth bool

	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, user_group_id, name, account_type,
		  COALESCE(iban, ''), COALESCE(bic, ''), COALESCE(currency_id, ''),
		  active, virtual_balance, include_net_worth,
		  latitude, longitude,
		  COALESCE(liability_type::text, ''), COALESCE(liability_direction::text, ''),
		  interest_rate, COALESCE(interest_period::text, ''), current_debt,
		  COALESCE(credit_card_type::text, ''), monthly_payment_date, monthly_payment_amount,
		  COALESCE(notes, ''), created_at, updated_at, deleted_at
		 FROM wallets WHERE id = $1 AND user_group_id = $2`, id, groupID,
	).Scan(&w.ID, &w.UserID, &w.UserGroupID, &w.Name, &w.AccountType,
		&iban, &bic, &currencyID,
		&active, &w.VirtualBalance, &includeNetWorth,
		&lat, &long,
		&liabilityType, &liabilityDirection,
		&interestRate, &interestPeriod, &currentDebt,
		&creditCardType, &monthlyPaymentDate, &monthlyPaymentAmt,
		&notes, &w.CreatedAt, &w.UpdatedAt, &deletedAt)
	if err != nil {
		return nil, fmt.Errorf("wallet not found: %w", err)
	}
	if deletedAt != nil {
		return nil, fmt.Errorf("wallet not found")
	}

	w.Active = active
	w.IncludeNetWorth = includeNetWorth
	w.Latitude = lat
	w.Longitude = long
	if iban.Valid && iban.s != "" { w.IBAN = &iban.s }
	if bic.Valid && bic.s != "" { w.BIC = &bic.s }
	if currencyID.Valid && currencyID.s != "" { w.CurrencyID = &currencyID.s }
	if notes.Valid && notes.s != "" { w.Notes = &notes.s }
	if interestPeriod.Valid && interestPeriod.s != "" { w.InterestPeriod = &interestPeriod.s }
	if liabilityType.Valid && liabilityType.s != "" { w.LiabilityType = &liabilityType.s }
	if liabilityDirection.Valid && liabilityDirection.s != "" { w.LiabilityDirection = &liabilityDirection.s }
	if creditCardType.Valid && creditCardType.s != "" { w.CreditCardType = &creditCardType.s }
	if interestRate != nil { w.InterestRate = interestRate }
	if currentDebt != nil { w.CurrentDebt = currentDebt }
	if monthlyPaymentAmt != nil { w.MonthlyPaymentAmount = monthlyPaymentAmt }
	w.MonthlyPaymentDate = monthlyPaymentDate

	return &w, nil
}

// List returns all wallets in a group, optionally filtered by type and active status.
func (r *WalletRepository) List(ctx context.Context, groupID uuid.UUID, walletType string, activeOnly bool) ([]domain.Wallet, error) {
	query := `SELECT id, user_id, user_group_id, name, account_type,
		  COALESCE(iban, ''), COALESCE(bic, ''), COALESCE(currency_id, ''),
		  active, virtual_balance, include_net_worth,
		  latitude, longitude,
		  COALESCE(liability_type::text, ''), COALESCE(liability_direction::text, ''),
		  interest_rate, COALESCE(interest_period::text, ''), current_debt,
		  COALESCE(credit_card_type::text, ''), monthly_payment_date, monthly_payment_amount,
		  COALESCE(notes, ''), created_at, updated_at
		 FROM wallets WHERE user_group_id = $1 AND deleted_at IS NULL`
	args := []interface{}{groupID}
	argN := 2

	if walletType != "" {
		query += fmt.Sprintf(" AND account_type = $%d", argN)
		args = append(args, walletType)
		argN++
	}
	if activeOnly {
		query += " AND active = true"
	}
	query += " ORDER BY name"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list wallets: %w", err)
	}
	defer rows.Close()

	return scanWallets(rows)
}

// Update updates wallet fields.
func (r *WalletRepository) Update(ctx context.Context, id, groupID uuid.UUID, name string, active, includeNetWorth *bool, notes *string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE wallets SET name = COALESCE(NULLIF($1, ''), name),
		  active = COALESCE($2, active),
		  include_net_worth = COALESCE($3, include_net_worth),
		  notes = COALESCE($4, notes),
		  updated_at = $5
		 WHERE id = $6 AND user_group_id = $7 AND deleted_at IS NULL`,
		name, active, includeNetWorth, notes, time.Now().UTC(), id, groupID,
	)
	return err
}

// Delete soft-deletes a wallet.
func (r *WalletRepository) Delete(ctx context.Context, id, groupID uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`UPDATE wallets SET deleted_at = $1 WHERE id = $2 AND user_group_id = $3 AND deleted_at IS NULL`,
		time.Now().UTC(), id, groupID,
	)
	return err
}

// ValidateType checks if the wallet type is valid for the given operation.
func (r *WalletRepository) ValidateType(walletType string, asSource, asDestination bool) error {
	wt := domain.WalletType(walletType)
	if asSource && !domain.IsSourceValid(wt) {
		return fmt.Errorf("wallet type '%s' is not valid as a transaction source", walletType)
	}
	if asDestination && !domain.IsDestinationValid(wt) {
		return fmt.Errorf("wallet type '%s' is not valid as a transaction destination", walletType)
	}
	return nil
}

// sqlString is a nullable string scanner.
type sqlString struct {
	s     string
	Valid bool
}

func (s *sqlString) Scan(value interface{}) error {
	if value == nil {
		s.Valid = false
		return nil
	}
	s.Valid = true
	switch v := value.(type) {
	case string:
		s.s = v
	case []byte:
		s.s = string(v)
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
	return nil
}

func scanWallets(rows interface{ Next() bool; Scan(...interface{}) error; Err() error }) ([]domain.Wallet, error) {
	var wallets []domain.Wallet
	for rows.Next() {
		var w domain.Wallet
		var iban, bic, currencyID, liabilityType, liabilityDirection, interestPeriod, creditCardType sqlString
		var notes sqlString
		var lat, long *float64
		var interestRate, currentDebt, monthlyPaymentAmt *decimal.Decimal
		var monthlyPaymentDate *time.Time

		if err := rows.Scan(&w.ID, &w.UserID, &w.UserGroupID, &w.Name, &w.AccountType,
			&iban, &bic, &currencyID,
			&w.Active, &w.VirtualBalance, &w.IncludeNetWorth,
			&lat, &long,
			&liabilityType, &liabilityDirection,
			&interestRate, &interestPeriod, &currentDebt,
			&creditCardType, &monthlyPaymentDate, &monthlyPaymentAmt,
			&notes, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, err
		}

		w.Latitude = lat
		w.Longitude = long
		if iban.Valid && iban.s != "" { w.IBAN = &iban.s }
		if bic.Valid && bic.s != "" { w.BIC = &bic.s }
		if currencyID.Valid && currencyID.s != "" { w.CurrencyID = &currencyID.s }
		if notes.Valid && notes.s != "" { w.Notes = &notes.s }
		if interestPeriod.Valid && interestPeriod.s != "" { w.InterestPeriod = &interestPeriod.s }
		if liabilityType.Valid && liabilityType.s != "" { w.LiabilityType = &liabilityType.s }
		if liabilityDirection.Valid && liabilityDirection.s != "" { w.LiabilityDirection = &liabilityDirection.s }
		if creditCardType.Valid && creditCardType.s != "" { w.CreditCardType = &creditCardType.s }
		if interestRate != nil { w.InterestRate = interestRate }
		if currentDebt != nil { w.CurrentDebt = currentDebt }
		if monthlyPaymentAmt != nil { w.MonthlyPaymentAmount = monthlyPaymentAmt }
		w.MonthlyPaymentDate = monthlyPaymentDate

		wallets = append(wallets, w)
	}

	return wallets, rows.Err()
}

// init ensures strings is used.
var _ = strings.TrimSpace
