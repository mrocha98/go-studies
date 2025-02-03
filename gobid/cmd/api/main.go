package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
	"github.com/mrocha98/go-studies/gobid/internal/api"
	"github.com/mrocha98/go-studies/gobid/internal/envutils"
	"github.com/mrocha98/go-studies/gobid/internal/services"
)

func makeDatabaseConnectionPool(ctx context.Context) (*pgxpool.Pool, error) {
	env := envutils.NewOSEnv()
	pool, err := pgxpool.New(ctx,
		fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s",
			env.DBUser(),
			env.DBPassword(),
			env.DBHost(),
			env.DBPort(),
			env.DBName(),
		))

	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return pool, nil
}

func makeApi(pool *pgxpool.Pool) api.Api {
	api := api.Api{
		Router:      *chi.NewMux(),
		UserService: services.NewUserService(pool),
	}

	api.BindRoutes()

	return api
}

func makeServerAddress() string {
	env := envutils.NewOSEnv()
	return strings.Join([]string{env.APIHost(), env.APIPort()}, ":")
}

func main() {
	ctx := context.Background()
	pool, err := makeDatabaseConnectionPool(ctx)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	api := makeApi(pool)
	go func() {
		addr := makeServerAddress()
		fmt.Printf("Listening on %s\n", addr)
		if err := http.ListenAndServe(addr, &api.Router); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				panic(err)
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
