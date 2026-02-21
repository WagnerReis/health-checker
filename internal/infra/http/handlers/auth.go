package handlers

import (
	"errors"
	"fmt"
	"health-checker/internal/application/usecases"
	domainerrors "health-checker/internal/domain/errors"
	"health-checker/internal/infra/http/helpers"
	"health-checker/internal/infra/http/presenters"
	"net/http"

	"github.com/google/uuid"
)

type AuthHandler struct {
	signUpUseCase usecases.SignUpUseCase
	loginUseCase  usecases.LoginUseCase
	logoutUseCase usecases.LogoutUseCase
}

func NewAuthHandler(
	signUpUseCase usecases.SignUpUseCase,
	loginUseCase usecases.LoginUseCase,
	logoutUseCase usecases.LogoutUseCase,
) *AuthHandler {
	return &AuthHandler{
		signUpUseCase: signUpUseCase,
		loginUseCase:  loginUseCase,
		logoutUseCase: logoutUseCase,
	}
}

type SignUpRequest struct {
	Name     string `json:"name" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
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
		if errors.Is(err, domainerrors.ErrUserEmailAlreadyExists) {
			helpers.WriteError(w, http.StatusConflict, domainerrors.ErrUserEmailAlreadyExists.Error())
			return
		}
		helpers.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to sign up: %v", err))
		return
	}

	helpers.WriteJSONResponse(w, http.StatusCreated, presenters.AuthPresenter(*authOutput))
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	req, err := helpers.DecodeAndValidateRequest[LoginRequest](w, r)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	authOutput, err := h.loginUseCase.Execute(r.Context(), usecases.LoginCommand{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, domainerrors.ErrUserInvalidCredentials) {
			helpers.WriteError(w, http.StatusUnauthorized, domainerrors.ErrUserInvalidCredentials.Error())
			return
		}
		helpers.WriteError(w, http.StatusInternalServerError, domainerrors.ErrUserInvalidCredentials.Error())
		return
	}

	helpers.WriteJSONResponse(w, http.StatusOK, presenters.AuthPresenter(*authOutput))
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	req, err := helpers.DecodeAndValidateRequest[LogoutRequest](w, r)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	userID, ok := r.Context().Value("userID").(uuid.UUID)
	if !ok {
		helpers.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	err = h.logoutUseCase.Execute(r.Context(), usecases.LogoutCommand{
		UserID:           userID,
		RefreshTokenHash: req.RefreshToken,
	})
	if err != nil {
		if errors.Is(err, domainerrors.ErrRefreshTokenNotFound) {
			helpers.WriteError(w, http.StatusUnauthorized, domainerrors.ErrRefreshTokenNotFound.Error())
			return
		}
		if errors.Is(err, domainerrors.ErrUserNotFound) {
			helpers.WriteError(w, http.StatusNotFound, domainerrors.ErrUserNotFound.Error())
			return
		}
		helpers.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to logout: %v", err))
		return
	}

	helpers.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Logged out successfully"})
}
