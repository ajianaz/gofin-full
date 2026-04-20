-- 000010_finalization.up.sql
-- Final production indexes, constraints, and optimizations

-- Users
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_unique ON users (email) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_users_group_id ON users (user_group_id) WHERE deleted_at IS NULL;

-- Wallets
CREATE INDEX IF NOT EXISTS idx_wallets_group_id ON wallets (user_group_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_wallets_account_type ON wallets (account_type) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_wallets_currency ON wallets (currency_id) WHERE deleted_at IS NULL;

-- Categories
CREATE INDEX IF NOT EXISTS idx_categories_group_id ON categories (user_group_id) WHERE deleted_at IS NULL;

-- Tags
CREATE INDEX IF NOT EXISTS idx_tags_group_id ON tags (user_group_id) WHERE deleted_at IS NULL;

-- Transaction groups
CREATE INDEX IF NOT EXISTS idx_transaction_groups_group_id ON transaction_groups (user_group_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_transaction_groups_created ON transaction_groups (created_at DESC) WHERE deleted_at IS NULL;

-- Transaction journals
CREATE INDEX IF NOT EXISTS idx_tx_journals_group_id ON transaction_journals (transaction_group_id);
CREATE INDEX IF NOT EXISTS idx_tx_journals_date ON transaction_journals (date DESC);

-- Transactions
CREATE INDEX IF NOT EXISTS idx_transactions_journal_id ON transactions (transaction_journal_id);

-- Budgets
CREATE INDEX IF NOT EXISTS idx_budgets_group_id ON budgets (user_group_id) WHERE deleted_at IS NULL;

-- Piggy banks
CREATE INDEX IF NOT EXISTS idx_piggy_banks_account_id ON piggy_banks (account_id) WHERE deleted_at IS NULL;

-- Bills
CREATE INDEX IF NOT EXISTS idx_bills_group_id ON bills (user_group_id) WHERE deleted_at IS NULL;

-- Notifications
CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications (user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_notifications_unread ON notifications (user_id) WHERE "read" = FALSE;

-- Audit logs
CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON audit_logs (user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_logs_entity ON audit_logs (entity_type, entity_id);

-- Exchange rates
CREATE INDEX IF NOT EXISTS idx_exchange_rates_pair ON exchange_rates (from_currency_id, to_currency_id, date DESC);

-- Webhooks
CREATE INDEX IF NOT EXISTS idx_webhooks_group_id ON webhooks (user_group_id) WHERE deleted_at IS NULL;

-- Attachments
CREATE INDEX IF NOT EXISTS idx_attachments_entity ON attachments (attachable_type, attachable_id) WHERE deleted_at IS NULL;

-- Notes
CREATE INDEX IF NOT EXISTS idx_notes_entity ON notes (noteable_type, noteable_id);

-- Locations
CREATE INDEX IF NOT EXISTS idx_locations_entity ON locations (locatable_type, locatable_id);

-- Wallet members
CREATE INDEX IF NOT EXISTS idx_wallet_members_wallet_id ON wallet_members (wallet_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_wallet_members_unique ON wallet_members (wallet_id, user_id);

-- Preferences
CREATE UNIQUE INDEX IF NOT EXISTS idx_preferences_unique ON preferences (user_id, name);

-- Configurations
CREATE UNIQUE INDEX IF NOT EXISTS idx_configurations_unique ON configurations (name);

-- Object groups
CREATE INDEX IF NOT EXISTS idx_object_groups_group_id ON object_groups (user_group_id);

-- Rules
CREATE INDEX IF NOT EXISTS idx_rules_group_id ON rules (rule_group_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_rules_active ON rules (active) WHERE deleted_at IS NULL;

-- Recurrences
CREATE INDEX IF NOT EXISTS idx_recurrences_user_group ON recurrences (user_group_id) WHERE deleted_at IS NULL;

-- FK constraint: ensure wallet members reference valid wallets
ALTER TABLE wallet_members ADD CONSTRAINT fk_wallet_members_wallet
    FOREIGN KEY (wallet_id) REFERENCES wallets(id) ON DELETE CASCADE;

-- FK constraint: ensure wallet members reference valid users
ALTER TABLE wallet_members ADD CONSTRAINT fk_wallet_members_user
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
