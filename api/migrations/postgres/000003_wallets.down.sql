-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS wallet_members;
DROP TABLE IF EXISTS wallets;
-- +goose StatementEnd
