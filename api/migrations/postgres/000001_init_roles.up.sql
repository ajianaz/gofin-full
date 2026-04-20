-- +goose Up
-- +goose StatementBegin

-- Global roles table (Tier 1: system-level roles)
CREATE TABLE IF NOT EXISTS roles (
    id          BIGSERIAL PRIMARY KEY,
    title       VARCHAR(255) NOT NULL UNIQUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Seed 22 user roles from Firefly III
INSERT INTO roles (title, created_at, updated_at) VALUES
    ('owner', NOW(), NOW()),
    ('full', NOW(), NOW()),
    ('manage_transactions', NOW(), NOW()),
    ('manage_meta', NOW(), NOW()),
    ('read_budgets', NOW(), NOW()),
    ('manage_budgets', NOW(), NOW()),
    ('read_piggy_banks', NOW(), NOW()),
    ('manage_piggy_banks', NOW(), NOW()),
    ('read_subscriptions', NOW(), NOW()),
    ('manage_subscriptions', NOW(), NOW()),
    ('read_rules', NOW(), NOW()),
    ('manage_rules', NOW(), NOW()),
    ('read_recurring', NOW(), NOW()),
    ('manage_recurring', NOW(), NOW()),
    ('read_webhooks', NOW(), NOW()),
    ('manage_webhooks', NOW(), NOW()),
    ('read_currencies', NOW(), NOW()),
    ('manage_currencies', NOW(), NOW()),
    ('view_reports', NOW(), NOW()),
    ('view_memberships', NOW(), NOW()),
    ('read_only', NOW(), NOW()),
    ('demo', NOW(), NOW())
ON CONFLICT DO NOTHING;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS roles;
-- +goose StatementEnd
