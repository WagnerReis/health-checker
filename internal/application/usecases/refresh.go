package usecases

import (
	"context"
	"errors"
	"health-checker/config"
	"health-checker/internal/application/cryptography"
	application "health-checker/internal/application/logger"
	entities "health-checker/internal/domain/entity"
	domainerrors "health-checker/internal/domain/errors"
	"health-checker/internal/domain/repository"
	"time"

	"github.com/google/uuid"
)

type RefreshCommand struct {
	RefreshToken string
}

type RefreshUseCase struct {
	userRepository         repository.UserRepository
	refreshTokenRepository repository.RefreshTokenRepository
	tokenGenerator         cryptography.TokenGenerator
	sha256Hash             cryptography.SHA256Hash
	config                 config.Config
	logger                 application.Logger
}

func NewRefreshUseCase(
	userRepository repository.UserRepository,
	refreshTokenRepository repository.RefreshTokenRepository,
	tokenGenerator cryptography.TokenGenerator,
	sha256Hash cryptography.SHA256Hash,
	config config.Config,
	logger application.Logger,
) *RefreshUseCase {
	return &RefreshUseCase{
		userRepository:         userRepository,
		refreshTokenRepository: refreshTokenRepository,
		tokenGenerator:         tokenGenerator,
		sha256Hash:             sha256Hash,
		config:                 config,
		logger:                 logger,
	}
}

func (u *RefreshUseCase) Execute(ctx context.Context, cmd RefreshCommand) (*AuthOutput, error) {
	refreshTokenHash := u.sha256Hash.Hash([]byte(cmd.RefreshToken))
	refreshToken, err := u.refreshTokenRepository.FindByTokenHash(ctx, refreshTokenHash)
	if err != nil {
		if errors.Is(err, domainerrors.ErrRefreshTokenNotFound) {
			return nil, domainerrors.ErrRefreshTokenNotFound
		}
		return nil, err
	}
	if refreshToken.Revoked {
		return nil, domainerrors.ErrRefreshTokenRevoked
	}
	if refreshToken.ExpiresAt.Before(time.Now()) {
		return nil, domainerrors.ErrRefreshTokenExpired
	}

	user, err := u.userRepository.FindByID(ctx, refreshToken.UserID)
	if err != nil {
		if errors.Is(err, domainerrors.ErrUserNotFound) {
			return nil, domainerrors.ErrUserNotFound
		}
		return nil, err
	}

	err = u.refreshTokenRepository.Revoke(ctx, refreshTokenHash)
	if err != nil {
		return nil, err
	}

	accessToken, err := u.tokenGenerator.Generate(user.ID, user.Email, u.config.AccessTokenSecret, u.config.AccessTokenExpiration)
	if err != nil {
		return nil, err
	}

	refreshTokenString, err := u.tokenGenerator.Generate(user.ID, user.Email, u.config.RefreshTokenSecret, u.config.RefreshTokenExpiration)
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(time.Duration(u.config.RefreshTokenExpiration) * time.Second)
	newRefreshTokenHash := u.sha256Hash.Hash([]byte(refreshTokenString))
	err = u.refreshTokenRepository.Create(ctx, entities.NewRefreshToken(uuid.Nil, user.ID, newRefreshTokenHash, expiresAt))
	if err != nil {
		return nil, err
	}

	u.logger.Info("User refreshed token successfully", application.Field{Key: "user_id", Value: user.ID.String()})

	return &AuthOutput{
		User: User{
			UserID: user.ID,
			Name:   user.Name,
			Email:  user.Email,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshTokenString,
	}, nil
}
