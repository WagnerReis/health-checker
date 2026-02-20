package handlers

import (
	"fmt"
	"health-checker/internal/application/usecases"
	"health-checker/internal/infra/http/helpers"
	"health-checker/internal/infra/http/presenters"
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
	req, err := helpers.DecodeAndValidateRequest[SignUpRequest](w, r)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	authOutput, err := h.signUpUseCase.Execute(r.Context(), usecases.SignUpCommand{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to sign up: %v", err))
		return
	}

	helpers.WriteJSONResponse(w, http.StatusCreated, presenters.SignUpPresenter(*authOutput))
}
