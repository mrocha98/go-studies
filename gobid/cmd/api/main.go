package main

import (
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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

func makeSessions(pool *pgxpool.Pool) *scs.SessionManager {
	s := scs.New()
	s.Store = pgxstore.New(pool)
	s.Lifetime = 24 * time.Hour
	s.Cookie.HttpOnly = true
	s.Cookie.SameSite = http.SameSiteLaxMode

	return s
}

func makeApi(pool *pgxpool.Pool) api.Api {
	api := api.Api{
		Router:      *chi.NewMux(),
		Env:         envutils.NewOSEnv(),
		UserService: services.NewUserService(pool),
		Sessions:    makeSessions(pool),
	}

	api.BindRoutes()

	return api
}

func makeServerAddress() string {
	env := envutils.NewOSEnv()
	return strings.Join([]string{env.APIHost(), env.APIPort()}, ":")
}

func init() {
	gob.Register(uuid.UUID{})
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
