-- name: CreateProduct :one
INSERT INTO products (
	seller_id,
	name,
	description,
	base_price,
	auction_end_at
) VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: GetProductById :one
SELECT * FROM products
WHERE id = $1;
