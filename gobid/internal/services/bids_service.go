package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mrocha98/go-studies/gobid/internal/store/pgstore"
)

type BidsService struct {
	pool    *pgxpool.Pool
	queries *pgstore.Queries
}

func NewBidsService(pool *pgxpool.Pool) BidsService {
	return BidsService{
		pool:    pool,
		queries: pgstore.New(pool),
	}
}

var ErrBidIsTooLow = errors.New("the bid value is too low")

func (bs *BidsService) PlaceBid(ctx context.Context, productID, userID uuid.UUID, amount float64) (pgstore.Bid, error) {
	product, err := bs.queries.GetProductById(ctx, productID)
	if err != nil {
		return pgstore.Bid{}, err
	}

	highestBid, err := bs.queries.GetHighestBidByProductId(ctx, productID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return pgstore.Bid{}, err
	}

	if product.BasePrice >= amount || highestBid.Amount >= amount {
		return pgstore.Bid{}, ErrBidIsTooLow
	}

	newBid, err := bs.queries.CreateBid(ctx, pgstore.CreateBidParams{
		ProductID: productID,
		UserID:    userID,
		Amount:    amount,
	})
	if err != nil {
		return pgstore.Bid{}, err
	}

	return newBid, nil
}
