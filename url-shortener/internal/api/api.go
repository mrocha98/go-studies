package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type response struct {
	Date  time.Time `json:"date"`
	Data  any       `json:"data,omitempty"`
	Error string    `json:"error,omitempty"`
}

func makeResponseWithData(data any) response {
	return response{
		Data: data,
		Date: time.Now(),
	}
}

func makeResponseWithError(error string) response {
	return response{
		Error: error,
		Date:  time.Now(),
	}
}

func NewHandler(db map[string]string) http.Handler {
	r := chi.NewMux()

	r.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
	)

	r.Route("/api", func(r chi.Router) {
		r.Route("/url", func(r chi.Router) {
			r.Post("/shorten", handleShortenUrl(db))
			r.Get("/{code}", handleGetShortenedUrl(db))
		})
	})

	return r
}
