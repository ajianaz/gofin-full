-- +goose Up
-- +goose StatementBegin

-- User groups (shared workspaces — all financial data scoped to group)
CREATE TABLE user_groups (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title       VARCHAR(255) NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

-- Group-level roles (Tier 2: 22 permissions)
CREATE TABLE user_roles (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title       VARCHAR(255) NOT NULL UNIQUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

-- Seed group-level roles
INSERT INTO user_roles (title, created_at, updated_at) VALUES
    ('read_only', NOW(), NOW()),
    ('manage_transactions', NOW(), NOW()),
    ('manage_meta', NOW(), NOW()),
    ('read_budgets', NOW(), NOW()),
    ('manage_piggy_banks', NOW(), NOW()),
    ('read_subscriptions', NOW(), NOW()),
    ('read_rules', NOW(), NOW()),
    ('read_recurring', NOW(), NOW()),
    ('read_webhooks', NOW(), NOW()),
    ('read_currencies', NOW(), NOW()),
    ('manage_budgets', NOW(), NOW()),
    ('manage_piggy_banks', NOW(), NOW()),
    ('manage_subscriptions', NOW(), NOW()),
    ('manage_rules', NOW(), NOW()),
    ('manage_recurring', NOW(), NOW()),
    ('manage_webhooks', NOW(), NOW()),
    ('manage_currencies', NOW(), NOW()),
    ('view_reports', NOW(), NOW()),
    ('view_memberships', NOW(), NOW()),
    ('full', NOW(), NOW()),
    ('owner', NOW(), NOW())
ON CONFLICT DO NOTHING;

-- Group memberships (user belongs to group with a role)
CREATE TABLE group_memberships (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL,
    user_group_id   UUID NOT NULL REFERENCES user_groups(id) ON DELETE CASCADE,
    user_role_id    UUID NOT NULL REFERENCES user_roles(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, user_group_id)
);

-- Global role assignments (owner, demo)
CREATE TABLE role_user (
    user_id     UUID NOT NULL,
    role_id     UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);

-- Users table (object_guid removed — all IDs are now UUIDs)
CREATE TABLE users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email           VARCHAR(255) NOT NULL UNIQUE,
    password        VARCHAR(255),
    remember_token  VARCHAR(100),
    reset_token     VARCHAR(32),
    blocked         BOOLEAN NOT NULL DEFAULT FALSE,
    blocked_code    VARCHAR(25),
    mfa_secret      VARCHAR(50),
    domain          VARCHAR(255),
    user_group_id   UUID REFERENCES user_groups(id) ON DELETE SET NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

-- Refresh tokens
CREATE TABLE refresh_tokens (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash  VARCHAR(255) NOT NULL UNIQUE,
    expires_at  TIMESTAMPTZ NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_user_group_id ON users(user_group_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_group_memberships_user_id ON group_memberships(user_id);
CREATE INDEX idx_group_memberships_user_group_id ON group_memberships(user_group_id);
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_token_hash ON refresh_tokens(token_hash);
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS role_user;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS group_memberships;
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS user_groups;
-- +goose StatementEnd
