package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

type handleV1CreateUrlBody struct {
	URL string `json:"url"`
}

type handleV1CreateUrlResponse struct {
	Code string `json:"code"`
}

func NewHandler(db map[string]string) http.Handler {
	r := chi.NewMux()

	r.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
	)

	r.Get("/{code}", handleRootGetUrlByCode(db))
	r.Route("/v1", func(r chi.Router) {
		r.Route("/urls", func(r chi.Router) {
			r.Post("/", handleV1CreateUrl(db))
		})
	})

	return r
}

func handleV1CreateUrl(db map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body handleV1CreateUrlBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			sendJSON(w, response{}.WithError("Invalid body"), http.StatusBadRequest)
			return
		}

		if _, err := url.Parse(body.URL); err != nil {
			sendJSON(w, response{}.WithError("Invalid URL"), http.StatusBadRequest)
			return
		}

		code := GenerateRandomString()
		db[code] = body.URL
		sendJSON(w, response{}.WithData(handleV1CreateUrlResponse{Code: code}), http.StatusCreated)
	}
}

func handleRootGetUrlByCode(db map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := chi.URLParam(r, "code")
		url, ok := db[code]
		if !ok {
			sendJSON(w, response{}.WithError("URL not found"), http.StatusNotFound)
			return
		}

		http.Redirect(w, r, url, http.StatusPermanentRedirect)
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
