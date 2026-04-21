-- Currencies (reference data, not group-scoped)
CREATE TABLE IF NOT EXISTS currencies (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code            VARCHAR(3) NOT NULL UNIQUE,
    name            VARCHAR(255) NOT NULL,
    symbol          VARCHAR(10) NOT NULL DEFAULT '',
    decimal_places  INT NOT NULL DEFAULT 2,
    enabled         BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_currencies_code ON currencies(code) WHERE deleted_at IS NULL;

-- Account types (reference data)
CREATE TABLE IF NOT EXISTS account_types (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type       VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Bills
CREATE TABLE IF NOT EXISTS bills (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id),
    user_group_id   UUID NOT NULL REFERENCES user_groups(id),
    name            VARCHAR(255) NOT NULL,
    amount_min      DECIMAL(32,16) NOT NULL DEFAULT 0,
    amount_max      DECIMAL(32,16) NOT NULL DEFAULT 0,
    date            DATE NOT NULL,
    end_date        DATE,
    repeat_freq     VARCHAR(255) NOT NULL DEFAULT 'monthly',
    skip            INT NOT NULL DEFAULT 0,
    active          BOOLEAN NOT NULL DEFAULT TRUE,
    "order"         INT NOT NULL DEFAULT 0,
    notes           TEXT,
    currency_id     VARCHAR(255) NOT NULL DEFAULT '',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_bills_user_group ON bills(user_group_id) WHERE deleted_at IS NULL;

-- Exchange rates
CREATE TABLE IF NOT EXISTS exchange_rates (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id          UUID NOT NULL REFERENCES users(id),
    user_group_id     UUID NOT NULL REFERENCES user_groups(id),
    from_currency_id VARCHAR(3) NOT NULL,
    to_currency_id   VARCHAR(3) NOT NULL,
    rate             DECIMAL(32,16) NOT NULL DEFAULT 0,
    date             DATE NOT NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at       TIMESTAMPTZ
);

CREATE INDEX idx_exchange_rates_group ON exchange_rates(user_group_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_exchange_rates_date ON exchange_rates(user_group_id, date) WHERE deleted_at IS NULL;

-- Webhooks
CREATE TABLE IF NOT EXISTS webhooks (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id),
    user_group_id   UUID NOT NULL REFERENCES user_groups(id),
    title           VARCHAR(255) NOT NULL,
    url             VARCHAR(255) NOT NULL,
    active          BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_webhooks_user_group ON webhooks(user_group_id) WHERE deleted_at IS NULL;

-- Webhook triggers
CREATE TABLE IF NOT EXISTS webhook_triggers (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    webhook_id  UUID NOT NULL REFERENCES webhooks(id) ON DELETE CASCADE,
    trigger     VARCHAR(255) NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_webhook_triggers_webhook ON webhook_triggers(webhook_id);

-- Webhook messages (outgoing)
CREATE TABLE IF NOT EXISTS webhook_messages (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    webhook_id      UUID NOT NULL REFERENCES webhooks(id) ON DELETE CASCADE,
    message         TEXT NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_webhook_messages_webhook ON webhook_messages(webhook_id);

-- Webhook deliveries (delivery attempts)
CREATE TABLE IF NOT EXISTS webhook_deliveries (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    webhook_message_id UUID NOT NULL REFERENCES webhook_messages(id) ON DELETE CASCADE,
    response_code   INT,
    response_body    TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Attachments (polymorphic)
CREATE TABLE IF NOT EXISTS attachments (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id),
    attachable_type VARCHAR(255) NOT NULL,
    attachable_id   UUID NOT NULL,
    filename        VARCHAR(255) NOT NULL,
    mime_type       VARCHAR(255) NOT NULL DEFAULT '',
    size            BIGINT NOT NULL DEFAULT 0,
    uploaded        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_attachments_user ON attachments(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_attachments_entity ON attachments(attachable_type, attachable_id) WHERE deleted_at IS NULL;

-- Notifications
CREATE TABLE IF NOT EXISTS notifications (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id),
    channel     VARCHAR(255) NOT NULL DEFAULT 'email',
    type        VARCHAR(255) NOT NULL,
    title       VARCHAR(255) NOT NULL,
    message     TEXT NOT NULL DEFAULT '',
    read        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_notifications_user ON notifications(user_id);
CREATE INDEX idx_notifications_read ON notifications(user_id) WHERE read = FALSE;

-- User preferences
CREATE TABLE IF NOT EXISTS preferences (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id),
    name        VARCHAR(255) NOT NULL,
    data        TEXT NOT NULL DEFAULT '{}',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, name)
);

-- System configurations (key-value, not user-scoped)
CREATE TABLE IF NOT EXISTS configurations (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(255) NOT NULL UNIQUE,
    value       TEXT NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Object groups (user-scoped collections)
CREATE TABLE IF NOT EXISTS object_groups (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id),
    user_group_id   UUID NOT NULL REFERENCES user_groups(id),
    title           VARCHAR(255) NOT NULL,
    "order"         INT NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_object_groups_user_group ON object_groups(user_group_id);

-- Notes (polymorphic)
CREATE TABLE IF NOT EXISTS notes (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    noteable_type   VARCHAR(255) NOT NULL,
    noteable_id     UUID NOT NULL,
    note            TEXT NOT NULL DEFAULT '',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_notes_entity ON notes(noteable_type, noteable_id);

-- Locations (polymorphic)
CREATE TABLE IF NOT EXISTS locations (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    locatable_type VARCHAR(255) NOT NULL,
    locatable_id   UUID NOT NULL,
    latitude        DOUBLE PRECISION,
    longitude       DOUBLE PRECISION,
    zoom_level      INT NOT NULL DEFAULT 10,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_locations_entity ON locations(locatable_type, locatable_id);

-- Seed default currencies
INSERT INTO currencies (code, name, symbol, decimal_places) VALUES
    ('EUR', 'Euro', '€', 2),
    ('USD', 'US Dollar', '$', 2),
    ('GBP', 'British Pound', '£', 2),
    ('JPY', 'Japanese Yen', '¥', 0),
    ('IDR', 'Indonesian Rupiah', 'Rp', 2),
    ('CHF', 'Swiss Franc', 'CHF', 2),
    ('CAD', 'Canadian Dollar', 'CA$', 2),
    ('AUD', 'Australian Dollar', 'A$', 2)
ON CONFLICT (code) DO NOTHING;

-- Seed default account types
INSERT INTO account_types (type) VALUES
    ('asset'),
    ('expense'),
    ('revenue'),
    ('liability'),
    ('opening-balance'),
    ('reconciliation'),
    ('import'),
    ('cash'),
    ('credit-card'),
    ('beneficiary'),
    ('initial-balance'),
    ('liability-credit'),
    ('shared-asset'),
    ('saving-asset'),
    ('cc-asset'),
    ('cash-wallet-asset')
ON CONFLICT (type) DO NOTHING;
