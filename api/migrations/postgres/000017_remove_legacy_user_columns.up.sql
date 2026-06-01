-- +goose Up
-- +goose StatementBegin

-- Remove legacy unused columns from users table
-- remember_token and reset_token are not referenced in any application code
ALTER TABLE users DROP COLUMN IF EXISTS remember_token;
ALTER TABLE users DROP COLUMN IF EXISTS reset_token;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Restore columns if rollback is needed
ALTER TABLE users ADD COLUMN IF NOT EXISTS remember_token VARCHAR(100);
ALTER TABLE users ADD COLUMN IF NOT EXISTS reset_token VARCHAR(32);

-- +goose StatementEnd
