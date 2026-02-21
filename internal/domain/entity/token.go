package entities

import (
	valueobjects "health-checker/internal/shared/value-object"
	"time"

	uuid "github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	TokenHash string    `json:"token_hash"`
	ExpiresAt time.Time `json:"expires_at"`
	Revoked   bool      `json:"revoked"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewRefreshToken(id uuid.UUID, userID uuid.UUID, tokenHash string, expiresAt time.Time) *RefreshToken {
	if id == uuid.Nil {
		id = valueobjects.NewID(uuid.Nil).Value()
	}

	refreshToken := &RefreshToken{
		ID:        id,
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
		Revoked:   false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return refreshToken
}
