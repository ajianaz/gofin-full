-- +goose Up
-- +goose StatementBegin

-- Categories
CREATE TABLE IF NOT EXISTS categories (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id),
    user_group_id   UUID NOT NULL REFERENCES user_groups(id),
    name            VARCHAR(255) NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_categories_user_group ON categories(user_group_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_categories_user ON categories(user_id) WHERE deleted_at IS NULL;

-- Tags
CREATE TABLE IF NOT EXISTS tags (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id),
    user_group_id   UUID NOT NULL REFERENCES user_groups(id),
    tag             VARCHAR(255) NOT NULL,
    date            DATE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_tags_user_group ON tags(user_group_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_tags_user ON tags(user_id) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX idx_tags_unique ON tags(user_id, user_group_id, tag) WHERE deleted_at IS NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS categories;
-- +goose StatementEnd
