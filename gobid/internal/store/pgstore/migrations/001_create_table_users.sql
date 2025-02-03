CREATE TABLE IF NOT EXISTS users (
	id uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
	user_name VARCHAR(50) UNIQUE NOT NULL,
	email TEXT UNIQUE NOT NULL,
	password_hash BYTEA NOT NULL,
	password_salt BYTEA NOT NULL,
	bio TEXT NOT NULL DEFAULT '',
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
---- create above / drop below ----
DROP TABLE IF EXISTS users;
