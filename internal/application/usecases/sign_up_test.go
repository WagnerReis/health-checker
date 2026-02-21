package usecases

import (
	"context"
	"errors"
	"health-checker/config"
	domainerrors "health-checker/internal/domain/errors"
	inmemory "health-checker/internal/infra/persistence/inmemory"
	"health-checker/internal/tests/criptography"
	fakelogger "health-checker/internal/tests/logger"
	"testing"

	"github.com/stretchr/testify/assert"
)

var validCommand = SignUpCommand{
	Name:     "John Doe",
	Email:    "john.doe@example.com",
	Password: "password",
}

func NewUseCase() *SignUpUseCase {
	repo := inmemory.NewUserRepositoryInMemory()
	hasher := criptography.NewFakeHasher()
	logger := fakelogger.NewFakeLogger()
	tokenGenerator := criptography.NewFakeJWTGenerator()
	config := config.Config{
		AccessTokenExpiration:  10,
		RefreshTokenExpiration: 20,
	}
	return NewSignUpUseCase(
		repo,
		inmemory.NewRefreshTokenRepositoryInMemory(),
		hasher,
		tokenGenerator,
		criptography.NewFakeSHA256Hash(),
		config,
		logger,
	)
}

func TestSignUpUseCase_Success(t *testing.T) {
	repo := inmemory.NewUserRepositoryInMemory()
	hasher := criptography.NewFakeHasher()
	tokenGenerator := criptography.NewFakeJWTGenerator()
	logger := fakelogger.NewFakeLogger()
	config := config.Config{
		AccessTokenExpiration:  10,
		RefreshTokenExpiration: 20,
	}
	uc := NewSignUpUseCase(
		repo,
		inmemory.NewRefreshTokenRepositoryInMemory(),
		hasher,
		tokenGenerator,
		criptography.NewFakeSHA256Hash(),
		config,
		logger,
	)

	authOutput, err := uc.Execute(context.Background(), validCommand)
	assert.NoError(t, err)

	userFound, err := repo.FindByEmail(context.Background(), validCommand.Email)
	assert.NoError(t, err)
	assert.NotNil(t, userFound)
	assert.Equal(t, authOutput.User.Name, userFound.Name)
	assert.Equal(t, authOutput.User.Email, userFound.Email)
	assert.Equal(t, authOutput.AccessToken, "token")
}

func TestSignUpUseCase_ErrorWhenEmailAlreadyExists(t *testing.T) {
	uc := NewUseCase()

	_, err := uc.Execute(context.Background(), validCommand)
	assert.NoError(t, err)

	_, err = uc.Execute(context.Background(), validCommand)

	assert.Error(t, err)
	assert.Equal(t, err, domainerrors.ErrUserEmailAlreadyExists)
}

func TestSignUpUseCase_ErrorWhenHashFails(t *testing.T) {
	hashErr := errors.New("hash entropy failure")
	hasher := &criptography.FakeHasher{ErrOnHash: hashErr}
	repo := inmemory.NewUserRepositoryInMemory()
	tokenGenerator := criptography.NewFakeJWTGenerator()
	config := config.Config{
		AccessTokenExpiration:  10,
		RefreshTokenExpiration: 20,
	}
	uc := NewSignUpUseCase(
		repo,
		inmemory.NewRefreshTokenRepositoryInMemory(),
		hasher,
		tokenGenerator,
		criptography.NewFakeSHA256Hash(),
		config,
		fakelogger.NewFakeLogger(),
	)

	_, err := uc.Execute(context.Background(), validCommand)

	assert.ErrorIs(t, err, hashErr)
}

func TestSignUpUseCase_ErrorWhenCreateUserFails(t *testing.T) {
	createErr := errors.New("duplicate email")
	repo := inmemory.NewUserRepositoryInMemory()
	repo.ErrOnCreate = createErr
	hasher := criptography.NewFakeHasher()
	tokenGenerator := criptography.NewFakeJWTGenerator()
	config := config.Config{
		AccessTokenExpiration:  10,
		RefreshTokenExpiration: 20,
	}
	uc := NewSignUpUseCase(
		repo,
		inmemory.NewRefreshTokenRepositoryInMemory(),
		hasher,
		tokenGenerator,
		criptography.NewFakeSHA256Hash(),
		config,
		fakelogger.NewFakeLogger(),
	)

	_, err := uc.Execute(context.Background(), validCommand)

	assert.ErrorIs(t, err, createErr)
}

func TestSignUpUseCase_ErrorWhenUserEntityIsInvalid(t *testing.T) {
	uc := NewUseCase()

	invalidCommand := SignUpCommand{
		Name:     "",
		Email:    "john.doe@example.com",
		Password: "password",
	}

	_, err := uc.Execute(context.Background(), invalidCommand)

	assert.Error(t, err)
	assert.Equal(t, err, domainerrors.ErrUserNameRequired)
}

func TestSignUpUseCase_ErrorWhenFindUserByEmailFails(t *testing.T) {
	findErr := errors.New("find user by email failure")
	repo := inmemory.NewUserRepositoryInMemory()
	repo.ErrOnFind = findErr
	hasher := criptography.NewFakeHasher()
	tokenGenerator := criptography.NewFakeJWTGenerator()
	config := config.Config{
		AccessTokenExpiration:  10,
		RefreshTokenExpiration: 20,
	}
	uc := NewSignUpUseCase(
		repo,
		inmemory.NewRefreshTokenRepositoryInMemory(),
		hasher,
		tokenGenerator,
		criptography.NewFakeSHA256Hash(),
		config,
		fakelogger.NewFakeLogger(),
	)

	_, err := uc.Execute(context.Background(), validCommand)

	assert.ErrorIs(t, err, findErr)
}

func TestSignUpUseCase_ErrorWhenGenerateAccessTokenFails(t *testing.T) {
	repo := inmemory.NewUserRepositoryInMemory()
	hasher := criptography.NewFakeHasher()
	tokenGenerator := &criptography.FakeJWTGenerator{ErrOnGenerate: errors.New("generate access token failure")}
	config := config.Config{
		AccessTokenExpiration:  10,
		RefreshTokenExpiration: 20,
	}
	uc := NewSignUpUseCase(
		repo,
		inmemory.NewRefreshTokenRepositoryInMemory(),
		hasher,
		tokenGenerator,
		criptography.NewFakeSHA256Hash(),
		config,
		fakelogger.NewFakeLogger(),
	)
	_, err := uc.Execute(context.Background(), validCommand)
	assert.ErrorIs(t, err, tokenGenerator.ErrOnGenerate)
}

func TestSignUpUseCase_ErrorWhenGenerateRefreshTokenFails(t *testing.T) {
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
	uc := NewSignUpUseCase(
		repo,
		inmemory.NewRefreshTokenRepositoryInMemory(),
		hasher,
		tokenGenerator,
		criptography.NewFakeSHA256Hash(),
		config,
		fakelogger.NewFakeLogger(),
	)
	_, err := uc.Execute(context.Background(), validCommand)
	assert.ErrorIs(t, err, tokenGenerator.ErrOnGenerate)
}
