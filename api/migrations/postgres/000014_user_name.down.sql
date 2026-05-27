-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN IF EXISTS name;
-- +goose StatementEnd
