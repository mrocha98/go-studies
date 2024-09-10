package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/mrocha98/go-studies/ama/backend/api"
	"github.com/mrocha98/go-studies/ama/backend/internal/store/pgstore"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	ctx := context.Background()

	pool, err := pgxpool.New(
		ctx,
		fmt.Sprintf(
			"user=%s password=%s host=%s port=%s dbname=%s",
			os.Getenv("AMA_BACKEND_DATABASE_USER"),
			os.Getenv("AMA_BACKEND_DATABASE_PASSWORD"),
			os.Getenv("AMA_BACKEND_DATABASE_HOST"),
			os.Getenv("AMA_BACKEND_DATABASE_PORT"),
			os.Getenv("AMA_BACKEND_DATABASE_NAME"),
		))
	if err != nil {
		panic(err)
	}

	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		panic(err)
	}

	handler := api.NewHandler(pgstore.New(pool))
	go func() {
		addr := strings.Join(
			[]string{
				os.Getenv("AMA_BACKEND_API_HOST"),
				os.Getenv("AMA_BACKEND_API_PORT"),
			},
			":",
		)
		fmt.Printf("Listening on %s\n", addr)
		if err := http.ListenAndServe(addr, handler); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				panic(err)
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
