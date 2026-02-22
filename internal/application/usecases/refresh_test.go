package usecases

import (
	"context"
	"errors"
	"health-checker/config"
	entities "health-checker/internal/domain/entity"
	domainerrors "health-checker/internal/domain/errors"
	inmemory "health-checker/internal/infra/persistence/inmemory"
	"health-checker/internal/tests/criptography"
	fakelogger "health-checker/internal/tests/logger"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var refreshTestUserID = uuid.New()

func setupRefreshTest() (*RefreshUseCase, *inmemory.UserRepositoryInMemory, *inmemory.RefreshTokenRepositoryInMemory) {
	userRepo := inmemory.NewUserRepositoryInMemory()
	refreshTokenRepo := inmemory.NewRefreshTokenRepositoryInMemory()

	userRepo.Create(context.Background(), &entities.User{
		ID:        refreshTestUserID,
		Name:      "John Doe",
		Email:     "john.doe@example.com",
		Password:  "password-hashed",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	refreshTokenRepo.Create(context.Background(), &entities.RefreshToken{
		ID:        uuid.New(),
		UserID:    refreshTestUserID,
		TokenHash: "hash",
		ExpiresAt: time.Now().Add(1 * time.Hour),
		Revoked:   false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	uc := NewRefreshUseCase(
		userRepo,
		refreshTokenRepo,
		criptography.NewFakeJWTGenerator(),
		criptography.NewFakeSHA256Hash(),
		config.Config{
			AccessTokenExpiration:  10,
			RefreshTokenExpiration: 20,
		},
		fakelogger.NewFakeLogger(),
	)

	return uc, userRepo, refreshTokenRepo
}

func TestRefreshUseCase_Success(t *testing.T) {
	uc, _, _ := setupRefreshTest()

	authOutput, err := uc.Execute(context.Background(), RefreshCommand{RefreshToken: "any-token"})

	assert.NoError(t, err)
	assert.NotNil(t, authOutput)
	assert.Equal(t, refreshTestUserID, authOutput.User.UserID)
	assert.Equal(t, "John Doe", authOutput.User.Name)
	assert.Equal(t, "john.doe@example.com", authOutput.User.Email)
	assert.Equal(t, "token", authOutput.AccessToken)
	assert.Equal(t, "token", authOutput.RefreshToken)
}

func TestRefreshUseCase_ErrorWhenTokenNotFound(t *testing.T) {
	userRepo := inmemory.NewUserRepositoryInMemory()
	refreshTokenRepo := inmemory.NewRefreshTokenRepositoryInMemory()
	uc := NewRefreshUseCase(
		userRepo,
		refreshTokenRepo,
		criptography.NewFakeJWTGenerator(),
		criptography.NewFakeSHA256Hash(),
		config.Config{AccessTokenExpiration: 10, RefreshTokenExpiration: 20},
		fakelogger.NewFakeLogger(),
	)

	_, err := uc.Execute(context.Background(), RefreshCommand{RefreshToken: "nonexistent"})

	assert.ErrorIs(t, err, domainerrors.ErrRefreshTokenNotFound)
}

func TestRefreshUseCase_ErrorWhenTokenIsRevoked(t *testing.T) {
	userRepo := inmemory.NewUserRepositoryInMemory()
	refreshTokenRepo := inmemory.NewRefreshTokenRepositoryInMemory()
	refreshTokenRepo.Create(context.Background(), &entities.RefreshToken{
		ID:        uuid.New(),
		UserID:    refreshTestUserID,
		TokenHash: "hash",
		ExpiresAt: time.Now().Add(1 * time.Hour),
		Revoked:   true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	uc := NewRefreshUseCase(
		userRepo,
		refreshTokenRepo,
		criptography.NewFakeJWTGenerator(),
		criptography.NewFakeSHA256Hash(),
		config.Config{AccessTokenExpiration: 10, RefreshTokenExpiration: 20},
		fakelogger.NewFakeLogger(),
	)

	_, err := uc.Execute(context.Background(), RefreshCommand{RefreshToken: "any-token"})

	assert.ErrorIs(t, err, domainerrors.ErrRefreshTokenRevoked)
}

func TestRefreshUseCase_ErrorWhenTokenIsExpired(t *testing.T) {
	userRepo := inmemory.NewUserRepositoryInMemory()
	refreshTokenRepo := inmemory.NewRefreshTokenRepositoryInMemory()
	refreshTokenRepo.Create(context.Background(), &entities.RefreshToken{
		ID:        uuid.New(),
		UserID:    refreshTestUserID,
		TokenHash: "hash",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
		Revoked:   false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	uc := NewRefreshUseCase(
		userRepo,
		refreshTokenRepo,
		criptography.NewFakeJWTGenerator(),
		criptography.NewFakeSHA256Hash(),
		config.Config{AccessTokenExpiration: 10, RefreshTokenExpiration: 20},
		fakelogger.NewFakeLogger(),
	)

	_, err := uc.Execute(context.Background(), RefreshCommand{RefreshToken: "any-token"})

	assert.ErrorIs(t, err, domainerrors.ErrRefreshTokenExpired)
}

func TestRefreshUseCase_ErrorWhenUserNotFound(t *testing.T) {
	userRepo := inmemory.NewUserRepositoryInMemory()
	refreshTokenRepo := inmemory.NewRefreshTokenRepositoryInMemory()
	refreshTokenRepo.Create(context.Background(), &entities.RefreshToken{
		ID:        uuid.New(),
		UserID:    uuid.New(),
		TokenHash: "hash",
		ExpiresAt: time.Now().Add(1 * time.Hour),
		Revoked:   false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	uc := NewRefreshUseCase(
		userRepo,
		refreshTokenRepo,
		criptography.NewFakeJWTGenerator(),
		criptography.NewFakeSHA256Hash(),
		config.Config{AccessTokenExpiration: 10, RefreshTokenExpiration: 20},
		fakelogger.NewFakeLogger(),
	)

	_, err := uc.Execute(context.Background(), RefreshCommand{RefreshToken: "any-token"})

	assert.ErrorIs(t, err, domainerrors.ErrUserNotFound)
}

func TestRefreshUseCase_ErrorWhenGenerateAccessTokenFails(t *testing.T) {
	userRepo := inmemory.NewUserRepositoryInMemory()
	refreshTokenRepo := inmemory.NewRefreshTokenRepositoryInMemory()
	tokenGenerator := &criptography.FakeJWTGenerator{ErrOnGenerate: errors.New("generate access token failure")}

	userRepo.Create(context.Background(), &entities.User{
		ID:        refreshTestUserID,
		Name:      "John Doe",
		Email:     "john.doe@example.com",
		Password:  "password-hashed",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	refreshTokenRepo.Create(context.Background(), &entities.RefreshToken{
		ID:        uuid.New(),
		UserID:    refreshTestUserID,
		TokenHash: "hash",
		ExpiresAt: time.Now().Add(1 * time.Hour),
		Revoked:   false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	uc := NewRefreshUseCase(
		userRepo,
		refreshTokenRepo,
		tokenGenerator,
		criptography.NewFakeSHA256Hash(),
		config.Config{AccessTokenExpiration: 10, RefreshTokenExpiration: 20},
		fakelogger.NewFakeLogger(),
	)

	_, err := uc.Execute(context.Background(), RefreshCommand{RefreshToken: "any-token"})

	assert.ErrorIs(t, err, tokenGenerator.ErrOnGenerate)
}

func TestRefreshUseCase_ErrorWhenGenerateRefreshTokenFails(t *testing.T) {
	userRepo := inmemory.NewUserRepositoryInMemory()
	refreshTokenRepo := inmemory.NewRefreshTokenRepositoryInMemory()
	tokenGenerator := &criptography.FakeJWTGenerator{
		ErrOnGenerate: errors.New("generate refresh token failure"),
		FailOnCall:    2,
	}

	userRepo.Create(context.Background(), &entities.User{
		ID:        refreshTestUserID,
		Name:      "John Doe",
		Email:     "john.doe@example.com",
		Password:  "password-hashed",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	refreshTokenRepo.Create(context.Background(), &entities.RefreshToken{
		ID:        uuid.New(),
		UserID:    refreshTestUserID,
		TokenHash: "hash",
		ExpiresAt: time.Now().Add(1 * time.Hour),
		Revoked:   false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	uc := NewRefreshUseCase(
		userRepo,
		refreshTokenRepo,
		tokenGenerator,
		criptography.NewFakeSHA256Hash(),
		config.Config{AccessTokenExpiration: 10, RefreshTokenExpiration: 20},
		fakelogger.NewFakeLogger(),
	)

	_, err := uc.Execute(context.Background(), RefreshCommand{RefreshToken: "any-token"})

	assert.ErrorIs(t, err, tokenGenerator.ErrOnGenerate)
}
