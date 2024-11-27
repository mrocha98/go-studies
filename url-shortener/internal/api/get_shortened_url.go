package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type getShortenedUrlResponse struct {
	URL string `json:"url"`
}

func handleGetShortenedUrl(db map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := chi.URLParam(r, "code")
		url, ok := db[code]
		if !ok {
			sendJSON(w, makeResponseWithError("URL not found"), http.StatusNotFound)
			return
		}

		sendJSON(w, makeResponseWithData(getShortenedUrlResponse{URL: url}), http.StatusOK)
	}
}
