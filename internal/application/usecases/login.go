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

type LoginCommand struct {
	Email    string
	Password string
}

type LoginUseCase struct {
	userRepository         repository.UserRepository
	refreshTokenRepository repository.RefreshTokenRepository
	hasher                 cryptography.Hasher
	tokenGenerator         cryptography.TokenGenerator
	sha256Hash             cryptography.SHA256Hash
	config                 config.Config
	logger                 application.Logger
}

func NewLoginUseCase(
	userRepository repository.UserRepository,
	refreshTokenRepository repository.RefreshTokenRepository,
	hasher cryptography.Hasher,
	tokenGenerator cryptography.TokenGenerator,
	sha256Hash cryptography.SHA256Hash,
	config config.Config,
	logger application.Logger,
) *LoginUseCase {
	return &LoginUseCase{
		userRepository:         userRepository,
		refreshTokenRepository: refreshTokenRepository,
		hasher:                 hasher,
		tokenGenerator:         tokenGenerator,
		sha256Hash:             sha256Hash,
		config:                 config,
		logger:                 logger,
	}
}

func (u *LoginUseCase) Execute(ctx context.Context, cmd LoginCommand) (*AuthOutput, error) {
	user, err := u.userRepository.FindByEmail(ctx, cmd.Email)
	if err != nil {
		if errors.Is(err, domainerrors.ErrUserNotFound) {
			return nil, domainerrors.ErrUserNotFound
		}
		return nil, err
	}

	if !u.hasher.Compare(cmd.Password, user.Password) {
		return nil, domainerrors.ErrUserInvalidCredentials
	}

	accessToken, err := u.tokenGenerator.Generate(user.ID, user.Email, u.config.AccessTokenSecret, u.config.AccessTokenExpiration)
	if err != nil {
		return nil, err
	}
	refreshToken, err := u.tokenGenerator.Generate(user.ID, user.Email, u.config.RefreshTokenSecret, u.config.RefreshTokenExpiration)
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(time.Duration(u.config.RefreshTokenExpiration) * time.Second)
	refreshTokenHash := u.sha256Hash.Hash([]byte(refreshToken))
	err = u.refreshTokenRepository.Create(ctx, entities.NewRefreshToken(uuid.Nil, user.ID, refreshTokenHash, expiresAt))
	if err != nil {
		return nil, err
	}

	u.logger.Info("User logged in successfully", application.Field{Key: "user_id", Value: user.ID.String()})

	return &AuthOutput{
		User: User{
			UserID: user.ID,
			Name:   user.Name,
			Email:  user.Email,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
