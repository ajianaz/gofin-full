-- +goose Up
-- +goose StatementBegin

-- Global roles table (Tier 1: system-level roles)
CREATE TABLE IF NOT EXISTS roles (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title       VARCHAR(255) NOT NULL UNIQUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Seed 22 user roles from Firefly III (explicit UUIDs for application code references)
INSERT INTO roles (id, title, created_at, updated_at) VALUES
    ('a0000000-0000-0000-0000-000000000001', 'owner', NOW(), NOW()),
    ('a0000000-0000-0000-0000-000000000002', 'full', NOW(), NOW()),
    ('a0000000-0000-0000-0000-000000000003', 'manage_transactions', NOW(), NOW()),
    ('a0000000-0000-0000-0000-000000000004', 'manage_meta', NOW(), NOW()),
    ('a0000000-0000-0000-0000-000000000005', 'read_budgets', NOW(), NOW()),
    ('a0000000-0000-0000-0000-000000000006', 'manage_budgets', NOW(), NOW()),
    ('a0000000-0000-0000-0000-000000000007', 'read_piggy_banks', NOW(), NOW()),
    ('a0000000-0000-0000-0000-000000000008', 'manage_piggy_banks', NOW(), NOW()),
    ('a0000000-0000-0000-0000-000000000009', 'read_subscriptions', NOW(), NOW()),
    ('a0000000-0000-0000-0000-00000000000a', 'manage_subscriptions', NOW(), NOW()),
    ('a0000000-0000-0000-0000-00000000000b', 'read_rules', NOW(), NOW()),
    ('a0000000-0000-0000-0000-00000000000c', 'manage_rules', NOW(), NOW()),
    ('a0000000-0000-0000-0000-00000000000d', 'read_recurring', NOW(), NOW()),
    ('a0000000-0000-0000-0000-00000000000e', 'manage_recurring', NOW(), NOW()),
    ('a0000000-0000-0000-0000-00000000000f', 'read_webhooks', NOW(), NOW()),
    ('a0000000-0000-0000-0000-000000000010', 'manage_webhooks', NOW(), NOW()),
    ('a0000000-0000-0000-0000-000000000011', 'read_currencies', NOW(), NOW()),
    ('a0000000-0000-0000-0000-000000000012', 'manage_currencies', NOW(), NOW()),
    ('a0000000-0000-0000-0000-000000000013', 'view_reports', NOW(), NOW()),
    ('a0000000-0000-0000-0000-000000000014', 'view_memberships', NOW(), NOW()),
    ('a0000000-0000-0000-0000-000000000015', 'read_only', NOW(), NOW()),
    ('a0000000-0000-0000-0000-000000000016', 'demo', NOW(), NOW())
ON CONFLICT DO NOTHING;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS roles;
-- +goose StatementEnd
