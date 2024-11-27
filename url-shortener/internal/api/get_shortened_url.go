package api

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mrocha98/go-studies/url-shortener/internal/store"
	"github.com/redis/go-redis/v9"
)

type getShortenedUrlResponse struct {
	URL string `json:"url"`
}

func handleGetShortenedUrl(store store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := chi.URLParam(r, "code")
		url, err := store.GetFullURL(r.Context(), code)
		if err != nil {
			if errors.Is(err, redis.Nil) {
				sendJSON(w, makeResponseWithError("URL not found"), http.StatusNotFound)
				return
			}
			slog.Error("failed to get url in store", slog.Any("error", err))
			sendJSON(w, makeResponseWithError("something went wrong"), http.StatusInternalServerError)
			return
		}

		redirectQueryParam := r.URL.Query().Get("redirect")
		shouldRedirect := redirectQueryParam == "true" || redirectQueryParam == "1"

		if shouldRedirect {
			http.Redirect(w, r, url, http.StatusFound)
			return
		}

		sendJSON(w, makeResponseWithData(getShortenedUrlResponse{URL: url}), http.StatusOK)
	}
}
