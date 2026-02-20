package cryptography

import "github.com/gofrs/uuid"

type TokenGenerator interface {
	Generate(userID uuid.UUID, email string, secretKey string, expiration int) (string, error)
}
