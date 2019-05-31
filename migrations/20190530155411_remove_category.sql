-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE log_items
    ADD COLUMN category;

UPDATE log_items
set category = (
    select categories.name
    from categories
    where categories.id = log_items.category_id
);

CREATE TEMPORARY TABLE log_items_backup
(
    id,
    created_at,
    name,
    amount,
    message_id,
    telegram_user_id,
    category
);
INSERT INTO log_items_backup
SELECT id, created_at, name, amount, message_id, telegram_user_id, category
FROM log_items;
DROP TABLE log_items;
CREATE TABLE log_items
(
    id,
    created_at,
    name,
    amount,
    message_id,
    telegram_user_id,
    category
);
INSERT INTO log_items
SELECT id, created_at, name, amount, message_id, telegram_user_id, category
FROM log_items_backup;
DROP TABLE log_items_backup;


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
