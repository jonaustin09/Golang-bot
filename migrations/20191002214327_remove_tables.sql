-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
DROP TABLE users;
DROP TABLE input_logs;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
