// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package pgstore

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	//CreateBid
	//
	//  INSERT INTO bids (
	//  	product_id, user_id, amount
	//  ) VALUES ( $1, $2, $3 )
	//  RETURNING id, product_id, user_id, amount, created_at
	CreateBid(ctx context.Context, arg CreateBidParams) (Bid, error)
	//CreateProduct
	//
	//  INSERT INTO products (
	//  	seller_id,
	//  	name,
	//  	description,
	//  	base_price,
	//  	auction_end_at
	//  ) VALUES ($1, $2, $3, $4, $5)
	//  RETURNING id
	CreateProduct(ctx context.Context, arg CreateProductParams) (uuid.UUID, error)
	//CreateUser
	//
	//  INSERT INTO users ("user_name", "email", "password_hash", "password_salt", "bio")
	//  VALUES ($1, $2, $3, $4, $5)
	//  RETURNING id
	CreateUser(ctx context.Context, arg CreateUserParams) (uuid.UUID, error)
	//GetBidsByProductId
	//
	//  SELECT id, product_id, user_id, amount, created_at FROM bids
	//  WHERE product_id = $1
	//  ORDER BY amount DESC
	GetBidsByProductId(ctx context.Context, productID uuid.UUID) ([]Bid, error)
	//GetHighestBidByProductId
	//
	//  SELECT id, product_id, user_id, amount, created_at FROM bids
	//  WHERE product_id = $1
	//  ORDER BY amount DESC
	//  LIMIT 1
	GetHighestBidByProductId(ctx context.Context, productID uuid.UUID) (Bid, error)
	//GetProductById
	//
	//  SELECT id, seller_id, name, description, base_price, auction_end_at, is_sold, created_at, updated_at FROM products
	//  WHERE id = $1
	GetProductById(ctx context.Context, id uuid.UUID) (Product, error)
	//GetUserByEmail
	//
	//  SELECT
	//  	id, user_name, password_hash, password_salt, email, bio, created_at, updated_at
	//  FROM users
	//  WHERE email = $1
	GetUserByEmail(ctx context.Context, email string) (GetUserByEmailRow, error)
	//GetUserById
	//
	//  SELECT
	//  	id, user_name, password_hash, password_salt, email, bio, created_at, updated_at
	//  FROM users
	//  WHERE id = $1
	GetUserById(ctx context.Context, id uuid.UUID) (GetUserByIdRow, error)
}

var _ Querier = (*Queries)(nil)
