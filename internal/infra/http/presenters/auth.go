package presenters

import (
	"health-checker/internal/application/usecases"

	"github.com/gofrs/uuid"
)

type User struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

type SignUpResponse struct {
	User         User   `json:"user"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func SignUpPresenter(output usecases.SignUpOutput) *SignUpResponse {
	return &SignUpResponse{
		User: User{
			ID:    output.User.UserID,
			Name:  output.User.Name,
			Email: output.User.Email,
		},
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
	}
}
