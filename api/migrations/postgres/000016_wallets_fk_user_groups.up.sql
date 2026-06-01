-- +goose Up
-- +goose StatementBegin

-- Add foreign key constraint from wallets.user_group_id to user_groups.id
ALTER TABLE wallets
    ADD CONSTRAINT fk_wallets_user_group
    FOREIGN KEY (user_group_id) REFERENCES user_groups(id) ON DELETE CASCADE;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE wallets DROP CONSTRAINT IF EXISTS fk_wallets_user_group;

-- +goose StatementEnd
