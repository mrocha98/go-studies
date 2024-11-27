package api

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type handleV1CreateUrlBody struct {
	URL string `json:"url"`
}

type handleV1CreateUrlResponse struct {
	Code string `json:"code"`
}

func handleShortenUrl(db map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body handleV1CreateUrlBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			sendJSON(w, makeResponseWithError("Invalid body"), http.StatusBadRequest)
			return
		}

		if _, err := url.Parse(body.URL); err != nil {
			sendJSON(w, makeResponseWithError("Invalid URL"), http.StatusBadRequest)
			return
		}

		code := genCode()
		db[code] = body.URL
		sendJSON(w, makeResponseWithData(handleV1CreateUrlResponse{Code: code}), http.StatusCreated)
	}
}
