CREATE TABLE IF NOT EXISTS products (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  seller_id UUID NOT NULL REFERENCES users(id),
	name TEXT NOT NULL,
	description TEXT NOT NULL DEFAULT '',
	base_price FLOAT NOT NULL,
	auction_end_at TIMESTAMPTZ NOT NULL,
	is_sold BOOLEAN NOT NULL DEFAULT false,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
---- create above / drop below ----
DROP TABLE IF EXISTS products;
