-- name: CreateOrder :one
INSERT INTO orders (customer_ref)
VALUES ($1)
RETURNING *;

-- name: AddOrderItem :one
INSERT INTO order_items (order_id, product_id, quantity, unit_price)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ListOrderItems :many
SELECT * FROM order_items
WHERE order_id = $1 and is_deleted = false
ORDER BY created_at DESC;

-- name: GetOrder :one
SELECT * FROM orders
WHERE id = $1 and is_deleted = false;

-- name: GetOrdersByCustomerRef :many
SELECT * FROM orders
WHERE customer_ref = $1 and is_deleted = false
ORDER BY created_at DESC;

-- name: GetAllOrders :many
SELECT * FROM orders
WHERE is_deleted = false
ORDER BY created_at DESC;


-- name: UpdateOrderTotalPrice :one
UPDATE orders
SET total_price = $1, created_at = NOW()
WHERE id = $2 and is_deleted = false
RETURNING *;

-- name: DeleteOrder :exec
UPDATE orders
SET is_deleted = true
WHERE id = $1 AND is_deleted = false;

-- name: DeleteOrderItemsByOrderID :exec
UPDATE order_items
SET is_deleted = true
WHERE order_id = $1 AND is_deleted = false;