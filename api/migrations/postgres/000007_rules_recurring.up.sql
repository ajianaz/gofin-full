-- Rule groups
CREATE TABLE IF NOT EXISTS rule_groups (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID NOT NULL REFERENCES users(id),
    user_group_id UUID NOT NULL REFERENCES user_groups(id),
    title        VARCHAR(255) NOT NULL,
    active       BOOLEAN NOT NULL DEFAULT TRUE,
    "order"      INT NOT NULL DEFAULT 0,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ
);

CREATE INDEX idx_rule_groups_user_group ON rule_groups(user_group_id) WHERE deleted_at IS NULL;

-- Rules
CREATE TABLE IF NOT EXISTS rules (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id          UUID NOT NULL REFERENCES users(id),
    user_group_id    UUID NOT NULL REFERENCES user_groups(id),
    rule_group_id    UUID REFERENCES rule_groups(id),
    title            VARCHAR(255) NOT NULL,
    description      TEXT,
    priority         INT NOT NULL DEFAULT 0,
    active           BOOLEAN NOT NULL DEFAULT TRUE,
    strict           BOOLEAN NOT NULL DEFAULT FALSE,
    stop_processing  BOOLEAN NOT NULL DEFAULT FALSE,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at       TIMESTAMPTZ
);

CREATE INDEX idx_rules_user_group ON rules(user_group_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_rules_rule_group ON rules(rule_group_id) WHERE deleted_at IS NULL;

-- Rule triggers
CREATE TABLE IF NOT EXISTS rule_triggers (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rule_id          UUID NOT NULL REFERENCES rules(id) ON DELETE CASCADE,
    trigger_type     VARCHAR(255) NOT NULL,
    trigger_value    TEXT NOT NULL DEFAULT '',
    stop_processing  BOOLEAN NOT NULL DEFAULT FALSE,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_rule_triggers_rule ON rule_triggers(rule_id);

-- Rule actions
CREATE TABLE IF NOT EXISTS rule_actions (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rule_id          UUID NOT NULL REFERENCES rules(id) ON DELETE CASCADE,
    action_type      VARCHAR(255) NOT NULL,
    action_value     TEXT NOT NULL DEFAULT '',
    "order"          INT NOT NULL DEFAULT 0,
    stop_processing  BOOLEAN NOT NULL DEFAULT FALSE,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_rule_actions_rule ON rule_actions(rule_id);

-- Recurrences (recurring transaction schedules)
CREATE TABLE IF NOT EXISTS recurrences (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id          UUID NOT NULL REFERENCES users(id),
    user_group_id    UUID NOT NULL REFERENCES user_groups(id),
    title            VARCHAR(255) NOT NULL,
    description      TEXT,
    first_date       DATE NOT NULL,
    latest_date      DATE,
    repeat_until     DATE,
    repeat_freq      VARCHAR(255) NOT NULL DEFAULT 'monthly',
    skip             INT NOT NULL DEFAULT 0,
    active           BOOLEAN NOT NULL DEFAULT TRUE,
    apply_rules      BOOLEAN NOT NULL DEFAULT FALSE,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at       TIMESTAMPTZ
);

CREATE INDEX idx_recurrences_user_group ON recurrences(user_group_id) WHERE deleted_at IS NULL;

-- Recurring transactions (templates within a recurrence)
CREATE TABLE IF NOT EXISTS recurring_transactions (
    id                        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recurrence_id             UUID NOT NULL REFERENCES recurrences(id) ON DELETE CASCADE,
    type                      VARCHAR(255) NOT NULL DEFAULT 'withdrawal',
    description               TEXT NOT NULL DEFAULT '',
    amount                    DECIMAL(32,16) NOT NULL DEFAULT 0,
    transaction_currency_id   VARCHAR(255) NOT NULL DEFAULT '',
    source_id                 UUID,
    destination_id            UUID,
    budget_id                 UUID,
    category_id               UUID,
    piggy_bank_id             UUID,
    "order"                   INT NOT NULL DEFAULT 0,
    created_at                TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at                TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_recurring_tx_recurrence ON recurring_transactions(recurrence_id);

-- Recurring repetitions (generated transaction references)
CREATE TABLE IF NOT EXISTS recurring_repetitions (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recurrence_id    UUID NOT NULL REFERENCES recurrences(id) ON DELETE CASCADE,
    relevant_date    DATE NOT NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_recurring_reps_recurrence ON recurring_repetitions(recurrence_id);

-- Recurrence meta (key-value)
CREATE TABLE IF NOT EXISTS recurrence_meta (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recurrence_id    UUID NOT NULL REFERENCES recurrences(id) ON DELETE CASCADE,
    name             VARCHAR(255) NOT NULL,
    value            TEXT NOT NULL DEFAULT '',
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_recurrence_meta_rec ON recurrence_meta(recurrence_id);

-- Recurring transaction meta (key-value)
CREATE TABLE IF NOT EXISTS recurring_transaction_meta (
    id                        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recurring_transaction_id  UUID NOT NULL REFERENCES recurring_transactions(id) ON DELETE CASCADE,
    name                      VARCHAR(255) NOT NULL,
    value                     TEXT NOT NULL DEFAULT '',
    created_at                TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at                TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_rec_tx_meta_rec ON recurring_transaction_meta(recurring_transaction_id);
