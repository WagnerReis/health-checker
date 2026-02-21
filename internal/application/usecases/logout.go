package usecases

import (
	"context"
	"errors"
	application "health-checker/internal/application/logger"
	domainerrors "health-checker/internal/domain/errors"
	"health-checker/internal/domain/repository"

	"github.com/google/uuid"
)

type LogoutCommand struct {
	UserID uuid.UUID
}

type LogoutUseCase struct {
	userRepository repository.UserRepository
	logger         application.Logger
}

func NewLogoutUseCase(userRepository repository.UserRepository, logger application.Logger) *LogoutUseCase {
	return &LogoutUseCase{
		userRepository: userRepository,
		logger:         logger,
	}
}

func (u *LogoutUseCase) Execute(ctx context.Context, cmd LogoutCommand) error {
	user, err := u.userRepository.FindByID(ctx, cmd.UserID)
	if err != nil {
		if errors.Is(err, domainerrors.ErrUserNotFound) {
			return domainerrors.ErrUserNotFound
		}
		return err
	}

	user.RefreshToken = nil
	err = u.userRepository.Update(ctx, user)
	if err != nil {
		return err
	}

	u.logger.Info("User logged out successfully", application.Field{Key: "user_id", Value: user.ID.String()})
	return nil
}
