package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
)

func (api *Api) BindRoutes() {
	api.Router.Use(
		middleware.RequestID,
		middleware.Logger,
		middleware.Recoverer,
		api.Sessions.LoadAndSave,
	)

	if mode := api.Env.Mode(); mode == "production" {
		csrfMiddleware := csrf.Protect(
			[]byte(api.Env.CSRFKey()),
			csrf.Secure(true),
		)
		api.Router.Use(csrfMiddleware)
	}

	api.Router.Route("/api", func(r chi.Router) {
		r.Get("/health", api.handleHealth)
		r.Route("/v1", func(r chi.Router) {
			r.Get("/csrftoken", api.handleGetCSRFToken)
			r.Route("/users", func(r chi.Router) {
				r.Post("/sign-up", api.handleSignUpUser)
				r.Post("/login", api.handleLoginUser)
				r.With(api.AuthMiddleware).Post("/logout", api.handleLogoutUser)
			})
			r.Route("/products", func(r chi.Router) {
				r.Group(func(r chi.Router) {
					r.Use(api.AuthMiddleware)
					r.Post("/", api.handleCreateProduct)
				})
			})
		})
	})
}
