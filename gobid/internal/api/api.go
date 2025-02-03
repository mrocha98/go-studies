package api

import (
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/mrocha98/go-studies/gobid/internal/services"
)

type Api struct {
	Router      chi.Mux
	UserService services.UserService
	Sessions    *scs.SessionManager
}
