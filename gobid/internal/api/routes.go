package api

import "github.com/go-chi/chi/v5"

func (api *Api) BindRoutes() {
	api.Router.Route("/api", func(r chi.Router) {
		r.Get("/health", api.handleHealth)
		r.Route("/v1", func(r chi.Router) {
			r.Route("/users", func(r chi.Router) {
				r.Post("/sign-up", api.handleSignUpUser)
				r.Post("/login", api.handleLoginUser)
				r.Post("/logout", api.handleLogoutUser)
			})
		})
	})
}
