package usecases

import (
	"context"
	"health-checker/config"
	"health-checker/internal/application/cryptography"
	"health-checker/internal/domain/repository"

	"github.com/gofrs/uuid"
)

type GenerateTokensCommand struct {
	UserID uuid.UUID
	Email  string
}

type GenerateTokensOutput struct {
	AccessToken  string
	RefreshToken string
}

type GenerateTokensUseCase struct {
	userRepository repository.UserRepository
	tokenGenerator cryptography.TokenGenerator
	config         *config.Config
}

func NewGenerateTokensUseCase(userRepository repository.UserRepository, tokenGenerator cryptography.TokenGenerator, config *config.Config) *GenerateTokensUseCase {
	return &GenerateTokensUseCase{
		userRepository: userRepository,
		tokenGenerator: tokenGenerator,
		config:         config,
	}
}

func (u *GenerateTokensUseCase) Execute(ctx context.Context, cmd GenerateTokensCommand) (*GenerateTokensOutput, error) {
	accessToken, err := u.tokenGenerator.Generate(cmd.UserID, cmd.Email, u.config.AccessTokenExpiration)
	if err != nil {
		return nil, err
	}

	refreshToken, err := u.tokenGenerator.Generate(cmd.UserID, cmd.Email, u.config.RefreshTokenExpiration)
	if err != nil {
		return nil, err
	}
	return &GenerateTokensOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
