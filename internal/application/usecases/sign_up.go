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
	valueobject "health-checker/internal/shared/value-object"
	"time"

	"github.com/google/uuid"
)

type SignUpCommand struct {
	Name     string
	Email    string
	Password string
}

type User struct {
	UserID uuid.UUID
	Name   string
	Email  string
}

type AuthOutput struct {
	User         User
	AccessToken  string
	RefreshToken string
}

type SignUpUseCase struct {
	userRepository         repository.UserRepository
	refreshTokenRepository repository.RefreshTokenRepository
	hasher                 cryptography.Hasher
	tokenGenerator         cryptography.TokenGenerator
	sha256Hash             cryptography.SHA256Hash
	config                 config.Config
	logger                 application.Logger
}

func NewSignUpUseCase(
	userRepository repository.UserRepository,
	refreshTokenRepository repository.RefreshTokenRepository,
	hasher cryptography.Hasher,
	tokenGenerator cryptography.TokenGenerator,
	sha256Hash cryptography.SHA256Hash,
	config config.Config,
	logger application.Logger) *SignUpUseCase {
	return &SignUpUseCase{
		userRepository:         userRepository,
		refreshTokenRepository: refreshTokenRepository,
		hasher:                 hasher,
		tokenGenerator:         tokenGenerator,
		sha256Hash:             sha256Hash,
		config:                 config,
		logger:                 logger,
	}
}

func (u *SignUpUseCase) Execute(ctx context.Context, cmd SignUpCommand) (*AuthOutput, error) {
	id := valueobject.NewID(uuid.Nil).Value()

	user, err := u.userRepository.FindByEmail(ctx, cmd.Email)
	if err != nil && !errors.Is(err, domainerrors.ErrUserNotFound) {
		return nil, err
	}
	if user != nil {
		return nil, domainerrors.ErrUserEmailAlreadyExists
	}

	hashedPassword, err := u.hasher.Hash(cmd.Password)
	if err != nil {
		return nil, err
	}

	user, err = entities.NewUser(id, cmd.Name, cmd.Email, *hashedPassword, nil)
	if err != nil {
		return nil, err
	}

	err = u.userRepository.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	u.logger.Info("User created successfully", application.Field{Key: "user_id", Value: user.ID.String()})

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
