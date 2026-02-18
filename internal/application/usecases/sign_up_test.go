package usecases

import (
	"context"
	"errors"
	inmemory "health-checker/internal/infra/persistence/inmemory/repository"
	fakehasher "health-checker/internal/tests/hasher"
	fakelogger "health-checker/internal/tests/logger"
	"testing"

	"github.com/stretchr/testify/assert"
)

var validCommand = SignUpCommand{
	Name:     "John Doe",
	Email:    "john.doe@example.com",
	Password: "password",
}

func TestSignUpUseCase_Success(t *testing.T) {
	repo := inmemory.NewUserRepositoryInMemory()
	hasher := fakehasher.NewFakeHasher()
	logger := fakelogger.NewFakeLogger()
	uc := NewSignUpUseCase(repo, hasher, logger)

	err := uc.Execute(context.Background(), validCommand)
	assert.NoError(t, err)

	user, err := repo.FindByEmail(context.Background(), validCommand.Email)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, validCommand.Name, user.Name)
	assert.Equal(t, validCommand.Email, user.Email)
	assert.True(t, hasher.Compare(validCommand.Password, user.Password))
}

func TestSignUpUseCase_ErrorWhenEmailAlreadyExists(t *testing.T) {
	repo := inmemory.NewUserRepositoryInMemory()
	hasher := fakehasher.NewFakeHasher()
	uc := NewSignUpUseCase(repo, hasher, fakelogger.NewFakeLogger())

	_ = uc.Execute(context.Background(), validCommand)

	err := uc.Execute(context.Background(), validCommand)

	assert.Error(t, err)
	assert.Equal(t, "email already exists", err.Error())
}

func TestSignUpUseCase_ErrorWhenHashFails(t *testing.T) {
	hashErr := errors.New("hash entropy failure")
	hasher := &fakehasher.FakeHasher{ErrOnHash: hashErr}
	repo := inmemory.NewUserRepositoryInMemory()
	uc := NewSignUpUseCase(repo, hasher, fakelogger.NewFakeLogger())

	err := uc.Execute(context.Background(), validCommand)

	assert.ErrorIs(t, err, hashErr)
}

func TestSignUpUseCase_ErrorWhenCreateUserFails(t *testing.T) {
	createErr := errors.New("duplicate email")
	repo := inmemory.NewUserRepositoryInMemory()
	repo.ErrOnCreate = createErr
	hasher := fakehasher.NewFakeHasher()
	uc := NewSignUpUseCase(repo, hasher, fakelogger.NewFakeLogger())

	err := uc.Execute(context.Background(), validCommand)

	assert.ErrorIs(t, err, createErr)
}

func TestSignUpUseCase_ErrorWhenUserEntityIsInvalid(t *testing.T) {
	repo := inmemory.NewUserRepositoryInMemory()
	hasher := fakehasher.NewFakeHasher()
	uc := NewSignUpUseCase(repo, hasher, fakelogger.NewFakeLogger())

	invalidCommand := SignUpCommand{
		Name:     "",
		Email:    "john.doe@example.com",
		Password: "password",
	}

	err := uc.Execute(context.Background(), invalidCommand)

	assert.Error(t, err)
	assert.EqualError(t, err, "name is required")
}

func TestSignUpUseCase_ErrorWhenFindUserByEmailFails(t *testing.T) {
	findErr := errors.New("find user by email failure")
	repo := inmemory.NewUserRepositoryInMemory()
	repo.ErrOnFind = findErr
	hasher := fakehasher.NewFakeHasher()
	uc := NewSignUpUseCase(repo, hasher, fakelogger.NewFakeLogger())

	err := uc.Execute(context.Background(), validCommand)

	assert.ErrorIs(t, err, findErr)
}
