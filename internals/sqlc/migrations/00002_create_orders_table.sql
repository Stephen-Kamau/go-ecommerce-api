-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,
    customer_ref TEXT NOT NULL DEFAULT "CustomerRefX", -- random string
    total_price INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS order_items (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES products(id),
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    unit_price INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- add indexes for performance
CREATE INDEX IF NOT EXISTS idx_orders_customer_ref ON orders(customer_ref);
CREATE INDEX  IF NOT EXISTS idx_order_items_order_id ON order_items(order_id);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE  IF EXISTS order_items;
DROP TABLE IF EXISTS orders;

-- drop indexes
DROP INDEX IF EXISTS idx_orders_customer_ref;
DROP INDEX IF EXISTS idx_order_items_order_id;

-- +goose StatementEnd
