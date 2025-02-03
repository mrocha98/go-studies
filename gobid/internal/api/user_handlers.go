package api

import (
	"errors"
	"net/http"

	"github.com/mrocha98/go-studies/gobid/internal/jsonutils"
	"github.com/mrocha98/go-studies/gobid/internal/services"
	"github.com/mrocha98/go-studies/gobid/internal/usecase/user"
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
	panic("TODO - NOT IMPLEMENTED")
}

func (api *Api) handleLogoutUser(w http.ResponseWriter, r *http.Request) {
	panic("TODO - NOT IMPLEMENTED")
}
