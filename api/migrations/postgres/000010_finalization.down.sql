-- 000010_finalization.down.sql
-- Drop all indexes and constraints from finalization migration

-- Constraints
ALTER TABLE wallet_members DROP CONSTRAINT IF EXISTS fk_wallet_members_wallet;
ALTER TABLE wallet_members DROP CONSTRAINT IF EXISTS fk_wallet_members_user;

-- Indexes
DROP INDEX IF EXISTS idx_users_email_unique;
DROP INDEX IF EXISTS idx_users_group_id;
DROP INDEX IF EXISTS idx_user_groups_title_user;
DROP INDEX IF EXISTS idx_wallets_group_id;
DROP INDEX IF EXISTS idx_wallets_account_type;
DROP INDEX IF EXISTS idx_wallets_currency;
DROP INDEX IF EXISTS idx_categories_group_id;
DROP INDEX IF EXISTS idx_tags_group_id;
DROP INDEX IF EXISTS idx_transaction_groups_group_id;
DROP INDEX IF EXISTS idx_transaction_groups_created;
DROP INDEX IF EXISTS idx_tx_journals_group_id;
DROP INDEX IF EXISTS idx_tx_journals_source;
DROP INDEX IF EXISTS idx_tx_journals_destination;
DROP INDEX IF EXISTS idx_tx_journals_category;
DROP INDEX IF EXISTS idx_tx_journals_date;
DROP INDEX IF EXISTS idx_transactions_journal_id;
DROP INDEX IF EXISTS idx_budgets_group_id;
DROP INDEX IF EXISTS idx_budgets_period;
DROP INDEX IF EXISTS idx_piggy_banks_wallet_id;
DROP INDEX IF EXISTS idx_bills_group_id;
DROP INDEX IF EXISTS idx_bills_next_date;
DROP INDEX IF EXISTS idx_notifications_user_id;
DROP INDEX IF EXISTS idx_notifications_unread;
DROP INDEX IF EXISTS idx_audit_logs_user_id;
DROP INDEX IF EXISTS idx_audit_logs_entity;
DROP INDEX IF EXISTS idx_exchange_rates_pair;
DROP INDEX IF EXISTS idx_webhooks_group_id;
DROP INDEX IF EXISTS idx_attachments_entity;
DROP INDEX IF EXISTS idx_notes_entity;
DROP INDEX IF EXISTS idx_locations_entity;
DROP INDEX IF EXISTS idx_wallet_members_wallet_id;
DROP INDEX IF EXISTS idx_wallet_members_unique;
DROP INDEX IF EXISTS idx_preferences_unique;
DROP INDEX IF EXISTS idx_configurations_unique;
DROP INDEX IF EXISTS idx_object_groups_group_id;
DROP INDEX IF EXISTS idx_rules_group_id;
DROP INDEX IF EXISTS idx_rules_active;
DROP INDEX IF EXISTS idx_recurrences_rule_id;
DROP INDEX IF EXISTS idx_recurrences_next_date;
