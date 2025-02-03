// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package pgstore

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	//CreateUser
	//
	//  INSERT INTO users ("user_name", "email", "password_hash", "password_salt", "bio")
	//  VALUES ($1, $2, $3, $4, $5)
	//  RETURNING id
	CreateUser(ctx context.Context, arg CreateUserParams) (uuid.UUID, error)
	//GetUserById
	//
	//  SELECT
	//  	id, user_name, password_hash, password_salt, email, bio, created_at, updated_at
	//  FROM users
	//  WHERE id = $1
	GetUserById(ctx context.Context, id uuid.UUID) (GetUserByIdRow, error)
}

var _ Querier = (*Queries)(nil)
