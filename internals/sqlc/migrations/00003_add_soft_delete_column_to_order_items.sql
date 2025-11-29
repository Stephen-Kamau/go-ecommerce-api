-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE order_items
ADD COLUMN is_deleted BOOLEAN NOT NULL DEFAULT FALSE
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE order_items
DROP COLUMN is_deleted;
-- +goose StatementEnd
