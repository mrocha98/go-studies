package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/mrocha98/go-studies/gobid/internal/jsonutils"
	"github.com/mrocha98/go-studies/gobid/internal/usecase/product"
)

func (api *Api) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	data, problems, err := jsonutils.DecodeValidJSON[product.CreateProductReq](r)
	if err != nil {
		jsonutils.EncodeJSON(w, r, http.StatusBadRequest, problems)
	}

	userID, ok := api.Sessions.Get(r.Context(), UserSessionKey).(uuid.UUID)
	if !ok {
		jsonutils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]any{
			"error": "unexpected error",
		})
	}

	id, err := api.ProductService.CreateProduct(
		r.Context(),
		userID,
		data.Name, data.Description, data.BasePrice, data.AuctionEndAt,
	)
	if err != nil {
		jsonutils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]any{
			"error": "failed to create product auction",
		})
	}

	jsonutils.EncodeJSON(w, r, http.StatusCreated, map[string]any{
		"product_id": id,
	})
}
