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

var validLoginCommand = LoginCommand{
	Email:    "john.doe@example.com",
	Password: "password",
}

func setup() (*LoginUseCase, *inmemory.UserRepositoryInMemory) {
	repo := inmemory.NewUserRepositoryInMemory()
	hasher := criptography.NewFakeHasher()
	tokenGenerator := criptography.NewFakeJWTGenerator()
	config := config.Config{
		AccessTokenExpiration:  10,
		RefreshTokenExpiration: 20,
	}
	return NewLoginUseCase(repo, hasher, tokenGenerator, config, fakelogger.NewFakeLogger()), repo
}

func TestLoginUseCase_Success(t *testing.T) {
	uc, repo := setup()
	repo.Create(context.Background(), &entities.User{
		ID:           uuid.New(),
		Name:         "John Doe",
		Email:        "john.doe@example.com",
		Password:     "password-hashed",
		RefreshToken: nil,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	})

	authOutput, err := uc.Execute(context.Background(), validLoginCommand)
	assert.NoError(t, err)

	userFound, err := repo.FindByEmail(context.Background(), validCommand.Email)
	assert.NoError(t, err)
	assert.NotNil(t, userFound)
	assert.Equal(t, authOutput.User.Name, userFound.Name)
	assert.Equal(t, authOutput.User.Email, userFound.Email)
	assert.Equal(t, authOutput.RefreshToken, *userFound.RefreshToken)
	assert.Equal(t, authOutput.AccessToken, "token")
}

func TestLoginUseCase_ErrorWhenUserNotFound(t *testing.T) {
	uc, _ := setup()

	_, err := uc.Execute(context.Background(), validLoginCommand)
	assert.ErrorIs(t, err, domainerrors.ErrUserNotFound)
}

func TestLoginUseCase_ErrorWhenFindUserByEmailFails(t *testing.T) {
	uc, repo := setup()
	repo.ErrOnFind = errors.New("find user by email failure")
	_, err := uc.Execute(context.Background(), validLoginCommand)
	assert.ErrorIs(t, err, repo.ErrOnFind)
}

func TestLoginUseCase_ErrorWhenUpdateUserFails(t *testing.T) {
	uc, repo := setup()
	repo.ErrOnUpdate = errors.New("update user failure")
	repo.Create(context.Background(), &entities.User{
		ID:           uuid.New(),
		Name:         "John Doe",
		Email:        "john.doe@example.com",
		Password:     "password-hashed",
		RefreshToken: nil,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	})
	_, err := uc.Execute(context.Background(), validLoginCommand)
	assert.ErrorIs(t, err, repo.ErrOnUpdate)
}

func TestLoginUseCase_ErrorWhenInvalidCredentials(t *testing.T) {
	uc, repo := setup()
	repo.Create(context.Background(), &entities.User{
		ID:           uuid.New(),
		Name:         "John Doe",
		Email:        "john.doe@example.com",
		Password:     "password-hashed",
		RefreshToken: nil,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	})

	_, err := uc.Execute(context.Background(), LoginCommand{Email: "john.doe@example.com", Password: "wrong-password"})
	assert.ErrorIs(t, err, domainerrors.ErrUserInvalidCredentials)
}

func TestLoginUseCase_ErrorWhenGenerateAccessTokenFails(t *testing.T) {
	repo := inmemory.NewUserRepositoryInMemory()
	hasher := criptography.NewFakeHasher()
	tokenGenerator := &criptography.FakeJWTGenerator{ErrOnGenerate: errors.New("generate access token failure")}
	config := config.Config{
		AccessTokenExpiration:  10,
		RefreshTokenExpiration: 20,
	}
	uc := NewLoginUseCase(repo, hasher, tokenGenerator, config, fakelogger.NewFakeLogger())
	repo.Create(context.Background(), &entities.User{
		ID:           uuid.New(),
		Name:         "John Doe",
		Email:        "john.doe@example.com",
		Password:     "password-hashed",
		RefreshToken: nil,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	})

	_, err := uc.Execute(context.Background(), validLoginCommand)
	assert.ErrorIs(t, err, tokenGenerator.ErrOnGenerate)
}

func TestLoginUseCase_ErrorWhenGenerateRefreshTokenFails(t *testing.T) {
	repo := inmemory.NewUserRepositoryInMemory()
	hasher := criptography.NewFakeHasher()
	tokenGenerator := &criptography.FakeJWTGenerator{
		ErrOnGenerate: errors.New("generate refresh token failure"),
		FailOnCall:    2,
	}
	config := config.Config{
		AccessTokenExpiration:  10,
		RefreshTokenExpiration: 20,
	}
	uc := NewLoginUseCase(repo, hasher, tokenGenerator, config, fakelogger.NewFakeLogger())
	repo.Create(context.Background(), &entities.User{
		ID:           uuid.New(),
		Name:         "John Doe",
		Email:        "john.doe@example.com",
		Password:     "password-hashed",
		RefreshToken: nil,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	})

	_, err := uc.Execute(context.Background(), validLoginCommand)
	assert.ErrorIs(t, err, tokenGenerator.ErrOnGenerate)
}
