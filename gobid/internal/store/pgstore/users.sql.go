// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: users.sql

package pgstore

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users ("user_name", "email", "password_hash", "password_salt", "bio")
VALUES ($1, $2, $3, $4, $5)
RETURNING id
`

type CreateUserParams struct {
	UserName     string `db:"user_name" json:"userName"`
	Email        string `db:"email" json:"email"`
	PasswordHash []byte `db:"password_hash" json:"passwordHash"`
	PasswordSalt []byte `db:"password_salt" json:"passwordSalt"`
	Bio          string `db:"bio" json:"bio"`
}

// CreateUser
//
//	INSERT INTO users ("user_name", "email", "password_hash", "password_salt", "bio")
//	VALUES ($1, $2, $3, $4, $5)
//	RETURNING id
func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.UserName,
		arg.Email,
		arg.PasswordHash,
		arg.PasswordSalt,
		arg.Bio,
	)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const getUserById = `-- name: GetUserById :one
SELECT
	id, user_name, password_hash, password_salt, email, bio, created_at, updated_at
FROM users
WHERE id = $1
`

type GetUserByIdRow struct {
	ID           uuid.UUID          `db:"id" json:"id"`
	UserName     string             `db:"user_name" json:"userName"`
	PasswordHash []byte             `db:"password_hash" json:"passwordHash"`
	PasswordSalt []byte             `db:"password_salt" json:"passwordSalt"`
	Email        string             `db:"email" json:"email"`
	Bio          string             `db:"bio" json:"bio"`
	CreatedAt    pgtype.Timestamptz `db:"created_at" json:"createdAt"`
	UpdatedAt    pgtype.Timestamptz `db:"updated_at" json:"updatedAt"`
}

// GetUserById
//
//	SELECT
//		id, user_name, password_hash, password_salt, email, bio, created_at, updated_at
//	FROM users
//	WHERE id = $1
func (q *Queries) GetUserById(ctx context.Context, id uuid.UUID) (GetUserByIdRow, error) {
	row := q.db.QueryRow(ctx, getUserById, id)
	var i GetUserByIdRow
	err := row.Scan(
		&i.ID,
		&i.UserName,
		&i.PasswordHash,
		&i.PasswordSalt,
		&i.Email,
		&i.Bio,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
