package cryptography

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
)

type JWTTokenGenerator struct {
	secretKey string
}

func NewJWTTokenGenerator(secretKey string) *JWTTokenGenerator {
	return &JWTTokenGenerator{
		secretKey: secretKey,
	}
}

func (e *JWTTokenGenerator) Generate(userID uuid.UUID, email string, expiration int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   userID.String(),
		"email": email,
		"exp":   time.Now().Add(time.Duration(expiration) * time.Minute).Unix(),
	})

	tokenString, err := token.SignedString([]byte(e.secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
