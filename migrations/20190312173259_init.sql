-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
  id integer,
  first_name varchar(255),
  last_name varchar(255),
  username varchar(255),
  language_code varchar(255),
  created_at bigint,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS log_items (
  id varchar(255),
  created_at bigint,
  name varchar(255),
  amount real,
  message_id bigint,
  telegram_user_id bigint,
  category_id bigint DEFAULT 9999 ,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS categories (
  id integer primary key autoincrement,
  name varchar(255),
  telegram_user_id bigint
);

CREATE INDEX IF NOT EXISTS idx_log_items_message_id ON log_items(message_id);
CREATE INDEX IF NOT EXISTS idx_log_items_category_id ON log_items(category_id);
CREATE INDEX IF NOT EXISTS idx_categories_name ON categories(name);
CREATE INDEX IF NOT EXISTS idx_categories_telegram_user_id ON categories(telegram_user_id);

INSERT INTO categories(id, name)
SELECT 9999, 'інші'
WHERE NOT EXISTS(SELECT 1 FROM categories WHERE id = 9999);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE categories;

DROP TABLE log_items;

DROP TABLE users;
-- +goose StatementEnd
