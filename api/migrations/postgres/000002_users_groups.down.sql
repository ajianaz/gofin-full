-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS role_user;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS group_memberships;
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS user_groups;
-- +goose StatementEnd
