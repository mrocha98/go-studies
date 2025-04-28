package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mrocha98/go-studies/gobid/internal/store/pgstore"
)

type ProductService struct {
	pool    *pgxpool.Pool
	queries *pgstore.Queries
}

func NewProductService(pool *pgxpool.Pool) ProductService {
	return ProductService{
		pool:    pool,
		queries: pgstore.New(pool),
	}
}

func (ps *ProductService) CreateProduct(
	ctx context.Context,
	sellerId uuid.UUID,
	name string,
	description string,
	basePrice float64,
	auctionEnd time.Time,
) (uuid.UUID, error) {
	return ps.queries.CreateProduct(ctx, pgstore.CreateProductParams{
		SellerID:     sellerId,
		Name:         name,
		Description:  description,
		BasePrice:    basePrice,
		AuctionEndAt: auctionEnd,
	})
}
