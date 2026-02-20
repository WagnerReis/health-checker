package criptography

import "github.com/gofrs/uuid"

type FakeJWTGenerator struct {
	ErrOnGenerate error
	FailOnCall    int
	callCount     int
}

func NewFakeJWTGenerator() *FakeJWTGenerator {
	return &FakeJWTGenerator{}
}

func (g *FakeJWTGenerator) Generate(userID uuid.UUID, email string, secretKey string, expiration int) (string, error) {
	g.callCount++
	if g.ErrOnGenerate != nil {
		if g.FailOnCall == 0 || g.callCount == g.FailOnCall {
			return "", g.ErrOnGenerate
		}
	}
	return "token", nil
}
