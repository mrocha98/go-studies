-- name: CreateUser :one
INSERT INTO users ("user_name", "email", "password_hash", "password_salt", "bio")
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: GetUserById :one
SELECT
	id, user_name, password_hash, password_salt, email, bio, created_at, updated_at
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT
	id, user_name, password_hash, password_salt, email, bio, created_at, updated_at
FROM users
WHERE email = $1;
