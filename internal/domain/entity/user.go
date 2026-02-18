package entities

import (
	"time"

	valueobject "health-checker/internal/shared/value-object"

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

func NewUser(id uuid.UUID, name, email, password string, refreshToken *string) *User {
	if id == uuid.Nil {
		id = valueobject.NewID(uuid.Nil).Value()
	}

	return &User{
		ID:           id,
		Name:         name,
		Email:        email,
		Password:     password,
		RefreshToken: refreshToken,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}
