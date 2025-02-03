package api

import (
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/mrocha98/go-studies/gobid/internal/jsonutils"
)

func (api *Api) handleGetCSRFToken(w http.ResponseWriter, r *http.Request) {
	token := csrf.Token(r)
	jsonutils.EncodeJSON(w, r, http.StatusOK, map[string]any{
		"csrf_token": token,
	})
}

func (api *Api) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !api.Sessions.Exists(r.Context(), UserSessionKey) {
			jsonutils.EncodeJSON(w, r, http.StatusUnauthorized, map[string]any{
				"error": "authentication is required to access this resource",
			})
			return
		}
		next.ServeHTTP(w, r)
	})
}
