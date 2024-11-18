package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mrocha98/go-studies/omdb/crawler"
)

type response struct {
	Date  time.Time `json:"date"`
	Data  any       `json:"data,omitempty"`
	Error string    `json:"error,omitempty"`
}

func (response) WithData(data any) response {
	return response{
		Data: data,
		Date: time.Now(),
	}
}

func (response) WithError(error string) response {
	return response{
		Error: error,
		Date:  time.Now(),
	}
}

func NewHandler(apiKey string) http.Handler {
	r := chi.NewMux()

	r.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
	)

	r.Route("/v1", func(r chi.Router) {
		r.Get("/movies", handleV1SearchMovies(apiKey))
	})

	return r
}

func handleV1SearchMovies(apiKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		search := r.URL.Query().Get("s")
		result, err := crawler.SearchMovie(apiKey, search)
		if err != nil {
			sendJSON(w, response{}.WithError("Something went wrong"), http.StatusBadGateway)
			return
		}

		sendJSON(w, response{}.WithData(result), http.StatusOK)
	}
}

func sendJSON(w http.ResponseWriter, response response, status int) {
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(response)
	if err != nil {
		slog.Error("Failed to marshal json", "error", err)
		sendJSON(w, response.WithError("something went wrong"), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	if _, err := w.Write(data); err != nil {
		slog.Error("Failed to write response to client", "error", err)
		return
	}
}
