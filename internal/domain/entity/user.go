package entities

import (
	"errors"
	"time"

	valueobjects "health-checker/internal/shared/value-object"

	uuid "github.com/gofrs/uuid"
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
		return errors.New("name is required")
	}
	if u.Password == "" {
		return errors.New("password is required")
	}
	if len(u.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	return nil
}
