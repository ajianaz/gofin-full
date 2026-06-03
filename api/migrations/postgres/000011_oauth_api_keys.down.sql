-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_api_keys_key_hash;
DROP INDEX IF EXISTS idx_api_keys_user_id;
DROP INDEX IF EXISTS idx_oauth_states_state;
DROP TABLE IF EXISTS api_keys;
DROP TABLE IF EXISTS oauth_states;

-- +goose StatementEnd
