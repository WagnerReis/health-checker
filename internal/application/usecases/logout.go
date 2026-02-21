package usecases

import (
	"context"
	"errors"
	"health-checker/internal/application/cryptography"
	application "health-checker/internal/application/logger"
	domainerrors "health-checker/internal/domain/errors"
	"health-checker/internal/domain/repository"

	"github.com/google/uuid"
)

type LogoutCommand struct {
	UserID           uuid.UUID
	RefreshTokenHash string
}

type LogoutUseCase struct {
	userRepository         repository.UserRepository
	refreshTokenRepository repository.RefreshTokenRepository
	sha256Hash             cryptography.SHA256Hash
	logger                 application.Logger
}

func NewLogoutUseCase(
	userRepository repository.UserRepository,
	refreshTokenRepository repository.RefreshTokenRepository,
	sha256Hash cryptography.SHA256Hash,
	logger application.Logger,
) *LogoutUseCase {
	return &LogoutUseCase{
		userRepository:         userRepository,
		refreshTokenRepository: refreshTokenRepository,
		sha256Hash:             sha256Hash,
		logger:                 logger,
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

	refreshTokenHash := u.sha256Hash.Hash([]byte(cmd.RefreshTokenHash))
	err = u.refreshTokenRepository.Revoke(ctx, refreshTokenHash)
	if err != nil {
		if errors.Is(err, domainerrors.ErrRefreshTokenNotFound) {
			return domainerrors.ErrRefreshTokenNotFound
		}
		return err
	}

	u.logger.Info("User logged out successfully", application.Field{Key: "user_id", Value: user.ID.String()})
	return nil
}
