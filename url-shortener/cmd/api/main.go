package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/mrocha98/go-studies/url-shortener/internal/api"
	"github.com/mrocha98/go-studies/url-shortener/internal/store"
	"github.com/redis/go-redis/v9"
)

func main() {
	if err := run(); err != nil {
		slog.Error("Failed to execute code", "error", err)
		os.Exit(1)
	}
}

func run() error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "cache:6379",
		Password: "9a1c6fbde8614645b543ef703153f295",
		DB:       0,
	})

	store := store.NewStore(rdb)

	handler := api.NewHandler(store)

	s := http.Server{
		Addr:         "0.0.0.0:8008",
		Handler:      handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  time.Minute,
	}

	slog.Info("Listening on 8008")
	if err := s.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}

	return nil
}
