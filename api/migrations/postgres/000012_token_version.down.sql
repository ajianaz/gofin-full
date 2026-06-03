-- +goose Down
-- +goose StatementBegin

ALTER TABLE users DROP COLUMN IF EXISTS token_version;

-- +goose StatementEnd
