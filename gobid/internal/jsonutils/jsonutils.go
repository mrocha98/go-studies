package jsonutils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mrocha98/go-studies/gobid/internal/validator"
)

func EncodeJSON[T any](w http.ResponseWriter, r *http.Request, statusCode int, data T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON %w", err)
	}
	return nil
}

func DecodeJSON[T validator.Validator](r *http.Request) (T, error) {
	var data T

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return data, fmt.Errorf("decode json failed: %w", err)
	}

	return data, nil
}

func DecodeValidJSON[T validator.Validator](r *http.Request) (T, validator.Evaluator, error) {
	data, err := DecodeJSON[T](r)

	if err != nil {
		return data, nil, err
	}

	if problems := data.Valid(r.Context()); len(problems) > 0 {
		return data, problems, fmt.Errorf("invalid %T: %d problems", data, len(problems))
	}

	return data, nil, nil
}
