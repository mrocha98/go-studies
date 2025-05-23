CREATE TABLE IF NOT EXISTS bids (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	product_id UUID NOT NULL REFERENCES products (id),
	user_id UUID NOT NULL REFERENCES users (id),
	amount FLOAT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
---- create above / drop below ----
DROP TABLE IF EXISTS bids;