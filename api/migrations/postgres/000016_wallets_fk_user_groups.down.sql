-- +goose Down
-- +goose StatementBegin

ALTER TABLE wallets DROP CONSTRAINT IF EXISTS fk_wallets_user_group;

-- +goose StatementEnd
