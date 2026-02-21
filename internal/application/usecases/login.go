package usecases

import (
	"context"
	"errors"
	"health-checker/config"
	"health-checker/internal/application/cryptography"
	application "health-checker/internal/application/logger"
	domainerrors "health-checker/internal/domain/errors"
	"health-checker/internal/domain/repository"
)

type LoginCommand struct {
	Email    string
	Password string
}

type LoginUseCase struct {
	userRepository repository.UserRepository
	hasher         cryptography.Hasher
	tokenGenerator cryptography.TokenGenerator
	config         config.Config
	logger         application.Logger
}

func NewLoginUseCase(
	userRepository repository.UserRepository,
	hasher cryptography.Hasher,
	tokenGenerator cryptography.TokenGenerator,
	config config.Config,
	logger application.Logger,
) *LoginUseCase {
	return &LoginUseCase{
		userRepository: userRepository,
		hasher:         hasher,
		tokenGenerator: tokenGenerator,
		config:         config,
		logger:         logger,
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

	user.RefreshToken = &refreshToken
	err = u.userRepository.Update(ctx, user)
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
