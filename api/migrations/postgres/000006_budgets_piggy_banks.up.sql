-- Budgets
CREATE TABLE IF NOT EXISTS budgets (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID NOT NULL REFERENCES users(id),
    user_group_id UUID NOT NULL REFERENCES user_groups(id),
    name         VARCHAR(255) NOT NULL,
    active       BOOLEAN NOT NULL DEFAULT TRUE,
    "order"      INT NOT NULL DEFAULT 0,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ
);

CREATE INDEX idx_budgets_user_group ON budgets(user_group_id) WHERE deleted_at IS NULL;

-- Budget limits (time-bounded spending limits per budget)
CREATE TABLE IF NOT EXISTS budget_limits (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    budget_id  UUID NOT NULL REFERENCES budgets(id) ON DELETE CASCADE,
    start      DATE NOT NULL,
    "end"      DATE NOT NULL,
    amount     DECIMAL(32,16) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_budget_limits_budget ON budget_limits(budget_id);

-- Auto-budget configuration
CREATE TABLE IF NOT EXISTS auto_budgets (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    budget_id       UUID NOT NULL REFERENCES budgets(id) ON DELETE CASCADE,
    auto_budget_type VARCHAR(255) NOT NULL DEFAULT 'none',
    period          VARCHAR(255) NOT NULL DEFAULT 'monthly',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_auto_budgets_budget ON auto_budgets(budget_id);

-- Piggy banks (savings goals linked to a wallet)
CREATE TABLE IF NOT EXISTS piggy_banks (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id    UUID NOT NULL REFERENCES wallets(id),
    name          VARCHAR(255) NOT NULL,
    target_amount DECIMAL(32,16) NOT NULL DEFAULT 0,
    start_date    DATE,
    target_date   DATE,
    "order"       INT NOT NULL DEFAULT 0,
    notes         TEXT,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ
);

CREATE INDEX idx_piggy_banks_account ON piggy_banks(account_id) WHERE deleted_at IS NULL;

-- Piggy bank events (money added/removed from piggy bank)
CREATE TABLE IF NOT EXISTS piggy_bank_events (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    piggy_bank_id UUID NOT NULL REFERENCES piggy_banks(id) ON DELETE CASCADE,
    amount        DECIMAL(32,16) NOT NULL DEFAULT 0,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_piggy_bank_events_pb ON piggy_bank_events(piggy_bank_id);
