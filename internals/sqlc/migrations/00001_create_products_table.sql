-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE  IF NOT EXISTS products (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL DEFAULT 'UNNAMED PRODUCT',
    description TEXT NOT NULL DEFAULT '',
    price INTEGER NOT NULL CHECK (price > 0),
    stock INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- indexes for performance
CREATE INDEX IF NOT EXISTS idx_products_name ON products(name);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP TABLE IF EXISTS products;
DROP INDEX IF EXISTS idx_products_name;
-- +goose StatementEnd
