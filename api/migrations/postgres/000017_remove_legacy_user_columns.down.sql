-- +goose Down
-- +goose StatementBegin

-- Restore columns if rollback is needed
ALTER TABLE users ADD COLUMN IF NOT EXISTS remember_token VARCHAR(100);
ALTER TABLE users ADD COLUMN IF NOT EXISTS reset_token VARCHAR(32);

-- +goose StatementEnd
