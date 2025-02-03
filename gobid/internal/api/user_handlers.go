package api

import (
	"errors"
	"net/http"

	"github.com/mrocha98/go-studies/gobid/internal/jsonutils"
	"github.com/mrocha98/go-studies/gobid/internal/services"
	"github.com/mrocha98/go-studies/gobid/internal/usecase/user"
)

const (
	UserSessionKey = "AuthenticatedUserId"
)

func (api *Api) handleSignUpUser(w http.ResponseWriter, r *http.Request) {
	data, problems, err := jsonutils.DecodeValidJSON[user.CreateUserReq](r)

	if err != nil {
		_ = jsonutils.EncodeJSON(w, r, http.StatusBadRequest, map[string]any{"errors": problems})
		return
	}

	id, err := api.UserService.CreateUser(
		r.Context(),
		data.UserName, data.Email, data.Password, data.Bio,
	)
	if err != nil {
		var errorMessage string
		var statusCode int
		if errors.Is(err, services.ErrDuplicatedUserNameOrEmail) {
			errorMessage = "email of username already in use"
			statusCode = http.StatusUnprocessableEntity
		} else {
			errorMessage = "unexpected error"
			statusCode = http.StatusInternalServerError
		}
		_ = jsonutils.EncodeJSON(w, r, statusCode, map[string]any{"error": errorMessage})
		return
	}

	jsonutils.EncodeJSON(w, r, http.StatusCreated, map[string]any{"id": id})
}

func (api *Api) handleLoginUser(w http.ResponseWriter, r *http.Request) {
	data, problems, err := jsonutils.DecodeValidJSON[user.LoginUserReq](r)
	if err != nil {
		jsonutils.EncodeJSON(w, r, http.StatusBadRequest, map[string]any{"errors": problems})
	}

	id, err := api.UserService.AuthenticateUser(r.Context(), data.Email, data.Password)
	if err != nil {
		var errorMessage string
		var statusCode int
		if errors.Is(err, services.ErrInvalidCredentials) {
			errorMessage = "invalid email or password"
			statusCode = http.StatusUnprocessableEntity
		} else {
			errorMessage = "unexpected error"
			statusCode = http.StatusInternalServerError
		}
		_ = jsonutils.EncodeJSON(w, r, statusCode, map[string]any{"error": errorMessage})
		return
	}

	if err := api.Sessions.RenewToken(r.Context()); err != nil {
		_ = jsonutils.EncodeJSON(
			w, r, http.StatusInternalServerError, map[string]any{"error": "unexpected error"},
		)
		return
	}
	api.Sessions.Put(r.Context(), UserSessionKey, id)
	w.WriteHeader(http.StatusNoContent)
}

func (api *Api) handleLogoutUser(w http.ResponseWriter, r *http.Request) {
	if err := api.Sessions.RenewToken(r.Context()); err != nil {
		_ = jsonutils.EncodeJSON(
			w, r, http.StatusInternalServerError, map[string]any{"error": "unexpected error"},
		)
		return
	}
	api.Sessions.Remove(r.Context(), UserSessionKey)
	w.WriteHeader(http.StatusNoContent)
}
