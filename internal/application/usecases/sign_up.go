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

type SignUpOutput struct {
	User         User
	AccessToken  string
	RefreshToken string
}

type SignUpUseCase struct {
	userRepository repository.UserRepository
	hasher         cryptography.Hasher
	tokenGenerator cryptography.TokenGenerator
	config         config.Config
	logger         application.Logger
}

func NewSignUpUseCase(userRepository repository.UserRepository, hasher cryptography.Hasher, tokenGenerator cryptography.TokenGenerator, config config.Config, logger application.Logger) *SignUpUseCase {
	return &SignUpUseCase{
		userRepository: userRepository,
		hasher:         hasher,
		tokenGenerator: tokenGenerator,
		config:         config,
		logger:         logger,
	}
}

func (u *SignUpUseCase) Execute(ctx context.Context, cmd SignUpCommand) (*SignUpOutput, error) {
	id := valueobject.NewID(uuid.Nil).Value()

	user, err := u.userRepository.FindByEmail(ctx, cmd.Email)
	if err != nil && !errors.Is(err, domainerrors.ErrUserNotFound) && !errors.Is(err, domainerrors.ErrUserNotFound) {
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

	user.RefreshToken = &refreshToken
	err = u.userRepository.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return &SignUpOutput{
		User: User{
			UserID: user.ID,
			Name:   user.Name,
			Email:  user.Email,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
