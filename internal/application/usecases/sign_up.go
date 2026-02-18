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

func (u *SignUpUseCase) Execute(ctx context.Context, cmd SignUpCommand) error {
	id := valueobject.NewID(uuid.Nil).Value()

	user, err := u.userRepository.FindByEmail(ctx, cmd.Email)
	if err != nil && !errors.Is(err, domainerrors.ErrUserNotFound) {
		u.logger.Error("Failed to find user by email", application.Field{Key: "error", Value: err.Error()})
		return err
	}
	if user != nil {
		return errors.New("email already exists")
	}

	hashedPassword, err := u.hasher.Hash(cmd.Password)
	if err != nil {
		u.logger.Error("Failed to hash password", application.Field{Key: "error", Value: err.Error()})
		return err
	}

	user = entities.NewUser(id, cmd.Name, cmd.Email, *hashedPassword, nil)
	err = u.userRepository.Create(ctx, user)
	if err != nil {
		u.logger.Error("Failed to create user", application.Field{Key: "error", Value: err.Error()})
		return err
	}
	u.logger.Info("User created successfully", application.Field{Key: "user_id", Value: user.ID.String()})
	return nil
}
