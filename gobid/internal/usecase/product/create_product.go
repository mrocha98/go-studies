package product

import (
	"context"
	"time"

	"github.com/mrocha98/go-studies/gobid/internal/validator"
)

type CreateProductReq struct {
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	BasePrice    float64   `json:"basePrice"`
	AuctionEndAt time.Time `json:"auctionEndAt"`
}

const minAudictionDuration = 2 * time.Hour

func (req CreateProductReq) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator

	eval.
		CheckIsNotBlank("name", req.Name).
		CheckIsNotBlank("description", req.Description).
		CheckMinChars("description", req.Description, 10).
		CheckMaxChars("description", req.Description, 255).
		CheckIsGreaterThanOrEqualFloat("basePrice", req.BasePrice, 0).
		CheckField(
			time.Until(req.AuctionEndAt) >= minAudictionDuration,
			"basePrice",
			"must be at least two hours duration",
		)

	return eval
}
