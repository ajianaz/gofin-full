-- Transaction types (seed data)
CREATE TABLE IF NOT EXISTS transaction_types (
    id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO transaction_types (type) VALUES
    ('withdrawal'),
    ('deposit'),
    ('transfer'),
    ('opening-balance'),
    ('reconciliation'),
    ('liability-credit'),
    ('invalid')
ON CONFLICT (type) DO NOTHING;

-- Link types for journal-to-journal links
CREATE TABLE IF NOT EXISTS link_types (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name       VARCHAR(255) NOT NULL UNIQUE,
    inward     VARCHAR(255) NOT NULL DEFAULT '',
    outward    VARCHAR(255) NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO link_types (name, inward, outward) VALUES
    ('refund', 'is refunded by', 'refunds')
ON CONFLICT (name) DO NOTHING;

-- Transaction groups (top-level container)
CREATE TABLE IF NOT EXISTS transaction_groups (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID NOT NULL REFERENCES users(id),
    user_group_id UUID NOT NULL REFERENCES user_groups(id),
    group_title  VARCHAR(255) NOT NULL DEFAULT '',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ
);

CREATE INDEX idx_transaction_groups_user_group ON transaction_groups(user_group_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_transaction_groups_user ON transaction_groups(user_id) WHERE deleted_at IS NULL;

-- Transaction journals
CREATE TABLE IF NOT EXISTS transaction_journals (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_group_id    UUID NOT NULL REFERENCES transaction_groups(id),
    user_id                 UUID NOT NULL REFERENCES users(id),
    user_group_id           UUID NOT NULL REFERENCES user_groups(id),
    transaction_type_id     UUID NOT NULL REFERENCES transaction_types(id),
    date                    DATE NOT NULL,
    "order"                 INT NOT NULL DEFAULT 0,
    description             VARCHAR(65536) NOT NULL DEFAULT '',
    transaction_currency_id VARCHAR(255) NOT NULL DEFAULT '',
    foreign_currency_id     VARCHAR(255),
    budget_id               UUID,
    bill_id                 UUID,
    piggy_bank_id           UUID,
    reconciled              BOOLEAN NOT NULL DEFAULT FALSE,
    notes                   TEXT,
    interest_date           DATE,
    book_date               DATE,
    process_date            DATE,
    due_date                DATE,
    payment_date            DATE,
    invoice_date            DATE,
    external_id             VARCHAR(255),
    external_url            VARCHAR(255),
    internal_reference      VARCHAR(255),
    recurrence_id           UUID,
    recurrence_total        INT,
    recurrence_count        INT,
    import_hash_v2          VARCHAR(255),
    original_source         VARCHAR(255),
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at              TIMESTAMPTZ
);

CREATE INDEX idx_transaction_journals_group ON transaction_journals(transaction_group_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_transaction_journals_user_group ON transaction_journals(user_group_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_transaction_journals_date ON transaction_journals(user_group_id, date) WHERE deleted_at IS NULL;
CREATE INDEX idx_transaction_journals_type ON transaction_journals(user_group_id, transaction_type_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_transaction_journals_desc ON transaction_journals(user_group_id, description) WHERE deleted_at IS NULL;

-- SEPA fields on journals
ALTER TABLE transaction_journals ADD COLUMN IF NOT EXISTS sepa_cc        VARCHAR(255);
ALTER TABLE transaction_journals ADD COLUMN IF NOT EXISTS sepa_ct_op     VARCHAR(255);
ALTER TABLE transaction_journals ADD COLUMN IF NOT EXISTS sepa_ct_id     VARCHAR(255);
ALTER TABLE transaction_journals ADD COLUMN IF NOT EXISTS sepa_db        VARCHAR(255);
ALTER TABLE transaction_journals ADD COLUMN IF NOT EXISTS sepa_country   VARCHAR(255);
ALTER TABLE transaction_journals ADD COLUMN IF NOT EXISTS sepa_ep        VARCHAR(255);
ALTER TABLE transaction_journals ADD COLUMN IF NOT EXISTS sepa_ci        VARCHAR(255);
ALTER TABLE transaction_journals ADD COLUMN IF NOT EXISTS sepa_batch_id  VARCHAR(255);

-- Transactions (actual monetary movements — debit/credit per journal)
CREATE TABLE IF NOT EXISTS transactions (
    id                       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_journal_id   UUID NOT NULL REFERENCES transaction_journals(id),
    account_id               UUID NOT NULL REFERENCES wallets(id),
    amount                   DECIMAL(32,16) NOT NULL DEFAULT 0,
    native_amount            DECIMAL(32,16) NOT NULL DEFAULT 0,
    foreign_amount           DECIMAL(32,16),
    native_foreign_amount    DECIMAL(32,16),
    foreign_currency_id      VARCHAR(255),
    reconciled               BOOLEAN NOT NULL DEFAULT FALSE,
    created_at               TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at               TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_transactions_journal ON transactions(transaction_journal_id);
CREATE INDEX idx_transactions_account ON transactions(account_id);

-- Journal-to-category pivot
CREATE TABLE IF NOT EXISTS category_transaction (
    category_id             UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    transaction_journal_id  UUID NOT NULL REFERENCES transaction_journals(id) ON DELETE CASCADE,
    PRIMARY KEY (category_id, transaction_journal_id)
);

-- Journal-to-tag pivot
CREATE TABLE IF NOT EXISTS journal_tag (
    tag_id                  UUID NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    transaction_journal_id  UUID NOT NULL REFERENCES transaction_journals(id) ON DELETE CASCADE,
    PRIMARY KEY (tag_id, transaction_journal_id)
);

-- Transaction journal metadata (key-value)
CREATE TABLE IF NOT EXISTS transaction_journal_meta (
    id                     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_journal_id UUID NOT NULL REFERENCES transaction_journals(id) ON DELETE CASCADE,
    name                   VARCHAR(255) NOT NULL,
    value                  TEXT NOT NULL DEFAULT '',
    created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tj_meta_journal ON transaction_journal_meta(transaction_journal_id);

-- Transaction journal links
CREATE TABLE IF NOT EXISTS transaction_journal_links (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    link_type_id  UUID NOT NULL REFERENCES link_types(id),
    source_id     UUID NOT NULL REFERENCES transaction_journals(id),
    destination_id UUID NOT NULL REFERENCES transaction_journals(id),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tj_links_source ON transaction_journal_links(source_id);
CREATE INDEX idx_tj_links_dest ON transaction_journal_links(destination_id);
