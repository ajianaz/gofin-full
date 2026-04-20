-- +goose Up
-- +goose StatementBegin

-- Wallets (accounts) table
CREATE TABLE wallets (
    id                      BIGSERIAL PRIMARY KEY,
    user_id                 BIGINT NOT NULL,
    user_group_id           BIGINT NOT NULL,
    name                    VARCHAR(255) NOT NULL,
    account_type            VARCHAR(255) NOT NULL DEFAULT 'asset',
    iban                    VARCHAR(255),
    bic                     VARCHAR(255),
    currency_id             VARCHAR(3),
    active                  BOOLEAN NOT NULL DEFAULT TRUE,
    virtual_balance         DECIMAL(32,4) NOT NULL DEFAULT 0,
    include_net_worth       BOOLEAN NOT NULL DEFAULT TRUE,
    latitude                DOUBLE PRECISION,
    longitude               DOUBLE PRECISION,
    liability_type          VARCHAR(255),
    liability_direction     VARCHAR(255),
    interest_rate           DECIMAL(32,4),
    interest_period         VARCHAR(255),
    current_debt            DECIMAL(32,4) DEFAULT 0,
    credit_card_type        VARCHAR(255),
    monthly_payment_date    DATE,
    monthly_payment_amount  DECIMAL(32,4),
    notes                   TEXT,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at              TIMESTAMPTZ
);

-- Wallet members (sharing)
CREATE TABLE wallet_members (
    id          BIGSERIAL PRIMARY KEY,
    wallet_id   BIGINT NOT NULL,
    user_id     BIGINT NOT NULL,
    role        VARCHAR(255) NOT NULL DEFAULT 'viewer',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(wallet_id, user_id)
);

-- Indexes
CREATE INDEX idx_wallets_user_group_id ON wallets(user_group_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_wallets_user_id ON wallets(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_wallets_account_type ON wallets(account_type) WHERE deleted_at IS NULL;
CREATE INDEX idx_wallets_active ON wallets(active) WHERE deleted_at IS NULL;
CREATE INDEX idx_wallet_members_wallet_id ON wallet_members(wallet_id);
CREATE INDEX idx_wallet_members_user_id ON wallet_members(user_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS wallet_members;
DROP TABLE IF EXISTS wallets;
-- +goose StatementEnd
