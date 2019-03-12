-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS input_logs (
  id varchar(255),
  created_at bigint,
  text varchar(255),
  telegram_user_id bigint,
  PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE input_logs;
-- +goose StatementEnd
