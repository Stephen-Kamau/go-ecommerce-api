-- name: CreateProduct :one
INSERT INTO products (name, description, price, stock)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ListProducts :many
SELECT * FROM products ORDER BY id;

-- name: FindProductByID :one
SELECT * FROM products WHERE id = $1;


-- name: UpdateProductStock :one
UPDATE products
SET stock = stock - $1
WHERE id = $2 AND stock >= $1
RETURNING *;


-- name: DeleteProduct :exec
DELETE FROM products WHERE id = $1;

-- name: SearchProductsByName :many
SELECT * FROM products
WHERE name ILIKE '%' || $1 || '%'
ORDER BY id;


-- name: UpdateProductDetails :one
UPDATE products
SET name = $1, description = $2, price = $3, updated_at = NOW()
WHERE id = $4
RETURNING *;


-- name: GetProductsByIDs :many
SELECT * FROM products
WHERE id = ANY($1)
ORDER BY id;


-- name: GetProductByName :one
SELECT id, name, description, price, stock, created_at, updated_at
FROM products
WHERE name = $1;


-- name: ProductExists :one
SELECT EXISTS(
    SELECT 1 FROM products WHERE name = $1
);
