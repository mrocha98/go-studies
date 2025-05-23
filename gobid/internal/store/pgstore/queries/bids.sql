-- name: CreateBid :one
INSERT INTO bids (
	product_id, user_id, amount
) VALUES ( $1, $2, $3 )
RETURNING *;

-- name: GetBidsByProductId :many
SELECT * FROM bids
WHERE product_id = $1
ORDER BY amount DESC;

-- name: GetHighestBidByProductId :one
SELECT * FROM bids
WHERE product_id = $1
ORDER BY amount DESC
LIMIT 1;
