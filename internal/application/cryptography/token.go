package cryptography

import "github.com/google/uuid"

type TokenGenerator interface {
	Generate(userID uuid.UUID, email string, secretKey string, expiration int) (string, error)
}
