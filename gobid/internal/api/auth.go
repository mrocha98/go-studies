package api

import (
	"net/http"

	"github.com/mrocha98/go-studies/gobid/internal/jsonutils"
)

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
