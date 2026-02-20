package usecases

import (
	"context"
	"errors"
	"health-checker/config"
	inmemory "health-checker/internal/infra/persistence/inmemory/repository"
	valueobject "health-checker/internal/shared/value-object"
	"health-checker/internal/tests/criptography"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

var tokenCmd = GenerateTokensCommand{
	UserID: valueobject.NewID(uuid.Nil).Value(),
	Email:  "john.doe@example.com",
}

var tokenCfg = &config.Config{
	AccessTokenExpiration:  10,
	RefreshTokenExpiration: 20,
}

func TestGenerateTokensUseCase_Success(t *testing.T) {
	repo := inmemory.NewUserRepositoryInMemory()
	tokenGenerator := criptography.NewFakeJWTGenerator()
	uc := NewGenerateTokensUseCase(repo, tokenGenerator, tokenCfg)

	authOutput, err := uc.Execute(context.Background(), tokenCmd)

	assert.NoError(t, err)
	assert.NotNil(t, authOutput)
	assert.Equal(t, "token", authOutput.AccessToken)
	assert.Equal(t, "token", authOutput.RefreshToken)
}

func TestGenerateTokensUseCase_ErrorWhenAccessTokenFails(t *testing.T) {
	tokenErr := errors.New("signing key expired")
	repo := inmemory.NewUserRepositoryInMemory()
	tokenGenerator := &criptography.FakeJWTGenerator{ErrOnGenerate: tokenErr}
	uc := NewGenerateTokensUseCase(repo, tokenGenerator, tokenCfg)

	output, err := uc.Execute(context.Background(), tokenCmd)

	assert.Nil(t, output)
	assert.ErrorIs(t, err, tokenErr)
}

func TestGenerateTokensUseCase_ErrorWhenRefreshTokenFails(t *testing.T) {
	tokenErr := errors.New("signing key expired")
	repo := inmemory.NewUserRepositoryInMemory()
	tokenGenerator := &criptography.FakeJWTGenerator{ErrOnGenerate: tokenErr, FailOnCall: 2}
	uc := NewGenerateTokensUseCase(repo, tokenGenerator, tokenCfg)

	output, err := uc.Execute(context.Background(), tokenCmd)

	assert.Nil(t, output)
	assert.ErrorIs(t, err, tokenErr)
}
