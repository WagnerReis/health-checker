package helpers

import (
	"encoding/json"
	"fmt"
	"health-checker/internal/infra/http/validation"
	"net/http"
)

func DecodeAndValidateRequest[T any](w http.ResponseWriter, r *http.Request) (*T, error) {
	var req T
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, fmt.Errorf("failed to decode request: %w", err)
	}

	if err := validation.GetValidator().Struct(req); err != nil {
		return nil, fmt.Errorf("failed to validate request: %w", err)
	}

	return &req, nil
}
