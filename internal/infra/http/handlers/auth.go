package handlers

import (
	"encoding/json"
	"fmt"
	"health-checker/internal/application/usecases"
	"health-checker/internal/infra/http/presenters"
	"health-checker/internal/infra/http/validation"
	"net/http"
)

type AuthHandler struct {
	signUpUseCase usecases.SignUpUseCase
}

func NewAuthHandler(signUpUseCase usecases.SignUpUseCase) *AuthHandler {
	return &AuthHandler{
		signUpUseCase: signUpUseCase,
	}
}

type SignUpRequest struct {
	Name     string `json:"name" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req SignUpRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to decode request: %v", err), http.StatusBadRequest)
		return
	}

	if err := validation.GetValidator().Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	authOutput, err := h.signUpUseCase.Execute(r.Context(), usecases.SignUpCommand{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to sign up: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(presenters.SignUpPresenter(*authOutput))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}
}
