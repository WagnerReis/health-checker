package usecases

import (
	"context"
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

var logoutTestUserID = uuid.New()

func setupLogoutTest() (*LogoutUseCase, *inmemory.UserRepositoryInMemory, *inmemory.RefreshTokenRepositoryInMemory) {
	userRepo := inmemory.NewUserRepositoryInMemory()
	refreshTokenRepo := inmemory.NewRefreshTokenRepositoryInMemory()

	userRepo.Create(context.Background(), &entities.User{
		ID:        logoutTestUserID,
		Name:      "John Doe",
		Email:     "john.doe@example.com",
		Password:  "password-hashed",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	refreshTokenRepo.Create(context.Background(), &entities.RefreshToken{
		ID:        uuid.New(),
		UserID:    logoutTestUserID,
		TokenHash: "hash",
		ExpiresAt: time.Now().Add(1 * time.Hour),
		Revoked:   false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	uc := NewLogoutUseCase(
		userRepo,
		refreshTokenRepo,
		criptography.NewFakeSHA256Hash(),
		fakelogger.NewFakeLogger(),
	)

	return uc, userRepo, refreshTokenRepo
}

func TestLogoutUseCase_Success(t *testing.T) {
	uc, _, refreshTokenRepo := setupLogoutTest()

	err := uc.Execute(context.Background(), LogoutCommand{
		UserID:           logoutTestUserID,
		RefreshTokenHash: "any-token",
	})

	assert.NoError(t, err)

	token, _ := refreshTokenRepo.FindByTokenHash(context.Background(), "hash")
	assert.True(t, token.Revoked)
}

func TestLogoutUseCase_ErrorWhenUserNotFound(t *testing.T) {
	uc, _, _ := setupLogoutTest()

	err := uc.Execute(context.Background(), LogoutCommand{
		UserID:           uuid.New(),
		RefreshTokenHash: "any-token",
	})

	assert.ErrorIs(t, err, domainerrors.ErrUserNotFound)
}

func TestLogoutUseCase_ErrorWhenRefreshTokenNotFound(t *testing.T) {
	userRepo := inmemory.NewUserRepositoryInMemory()
	userRepo.Create(context.Background(), &entities.User{
		ID:        logoutTestUserID,
		Name:      "John Doe",
		Email:     "john.doe@example.com",
		Password:  "password-hashed",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	uc := NewLogoutUseCase(
		userRepo,
		inmemory.NewRefreshTokenRepositoryInMemory(),
		criptography.NewFakeSHA256Hash(),
		fakelogger.NewFakeLogger(),
	)

	err := uc.Execute(context.Background(), LogoutCommand{
		UserID:           logoutTestUserID,
		RefreshTokenHash: "nonexistent",
	})

	assert.ErrorIs(t, err, domainerrors.ErrRefreshTokenNotFound)
}
