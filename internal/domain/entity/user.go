package entities

import (
	"time"

	domainerrors "health-checker/internal/domain/errors"
	valueobjects "health-checker/internal/shared/value-object"

	uuid "github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	RefreshToken *string   `json:"refresh_token"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func NewUser(id uuid.UUID, name, emailStr, password string, refreshToken *string) (*User, error) {
	if id == uuid.Nil {
		id = valueobjects.NewID(uuid.Nil).Value()
	}

	email, err := valueobjects.NewEmail(emailStr)
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:           id,
		Name:         name,
		Email:        email.String(),
		Password:     password,
		RefreshToken: refreshToken,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = user.validate()
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *User) validate() error {
	if u.Name == "" {
		return domainerrors.ErrUserNameRequired
	}
	if u.Password == "" {
		return domainerrors.ErrUserPasswordRequired
	}
	if len(u.Password) < 8 {
		return domainerrors.ErrUserPasswordTooShort
	}
	return nil
}
