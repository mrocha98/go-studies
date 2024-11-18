package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/mrocha98/go-studies/omdb/api"
)

func main() {
	loadEnvs()
	if err := run(); err != nil {
		slog.Error("Failed to execute code", "error", err)
		os.Exit(1)
	}
	slog.Info("Finished execution")
}

func loadEnvs() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
}

func run() error {
	apiKey := os.Getenv("OMDB_API_KEY")
	handler := api.NewHandler(apiKey)

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
