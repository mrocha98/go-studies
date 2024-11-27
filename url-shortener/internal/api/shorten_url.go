package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/mrocha98/go-studies/url-shortener/internal/store"
)

type handleV1CreateUrlBody struct {
	URL string `json:"url"`
}

type handleV1CreateUrlResponse struct {
	Code string `json:"code"`
}

func handleShortenUrl(store store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body handleV1CreateUrlBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			sendJSON(w, makeResponseWithError("Invalid body"), http.StatusBadRequest)
			return
		}

		if _, err := url.Parse(body.URL); err != nil {
			sendJSON(w, makeResponseWithError("Invalid URL"), http.StatusBadRequest)
			return
		}

		code, err := store.SaveShortenedURL(r.Context(), body.URL)
		if err != nil {
			slog.Error("failed to save shortened url", slog.Any("error", err))
			sendJSON(w, makeResponseWithError("something went wrong"), http.StatusInternalServerError)
			return
		}
		sendJSON(w, makeResponseWithData(handleV1CreateUrlResponse{Code: code}), http.StatusCreated)
	}
}
