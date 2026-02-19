package usecases

import (
	"context"
	"errors"
	"health-checker/internal/application/hasher"
	application "health-checker/internal/application/logger"
	entities "health-checker/internal/domain/entity"
	domainerrors "health-checker/internal/domain/errors"
	"health-checker/internal/domain/repository"
	valueobject "health-checker/internal/shared/value-object"

	"github.com/gofrs/uuid"
)

type SignUpCommand struct {
	Name     string
	Email    string
	Password string
}

type SignUpOutput struct {
	UserID uuid.UUID
	Name   string
	Email  string
}

type SignUpUseCase struct {
	userRepository repository.UserRepository
	hasher         hasher.Hasher
	logger         application.Logger
}

func NewSignUpUseCase(userRepository repository.UserRepository, hasher hasher.Hasher, logger application.Logger) *SignUpUseCase {
	return &SignUpUseCase{
		userRepository: userRepository,
		hasher:         hasher,
		logger:         logger,
	}
}

func (u *SignUpUseCase) Execute(ctx context.Context, cmd SignUpCommand) (*SignUpOutput, error) {
	id := valueobject.NewID(uuid.Nil).Value()

	user, err := u.userRepository.FindByEmail(ctx, cmd.Email)
	if err != nil && !errors.Is(err, domainerrors.ErrUserNotFound) {
		u.logger.Error("Failed to find user by email", application.Field{Key: "error", Value: err.Error()})
		return nil, err
	}
	if user != nil {
		return nil, errors.New("email already exists")
	}

	hashedPassword, err := u.hasher.Hash(cmd.Password)
	if err != nil {
		u.logger.Error("Failed to hash password", application.Field{Key: "error", Value: err.Error()})
		return nil, err
	}

	user, err = entities.NewUser(id, cmd.Name, cmd.Email, *hashedPassword, nil)
	if err != nil {
		u.logger.Error("Failed to create user entity", application.Field{Key: "error", Value: err.Error()})
		return nil, err
	}

	err = u.userRepository.Create(ctx, user)
	if err != nil {
		u.logger.Error("Failed to create user", application.Field{Key: "error", Value: err.Error()})
		return nil, err
	}
	u.logger.Info("User created successfully", application.Field{Key: "user_id", Value: user.ID.String()})
	return &SignUpOutput{
		UserID: user.ID,
		Name:   user.Name,
		Email:  user.Email,
	}, nil
}
