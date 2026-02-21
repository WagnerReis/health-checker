package repository

import (
	"context"
	entities "health-checker/internal/domain/entity"

	"github.com/google/uuid"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, refreshToken *entities.RefreshToken) error
	FindByTokenHash(ctx context.Context, tokenHash string) (*entities.RefreshToken, error)
	Revoke(ctx context.Context, refreshTokenHash string) error
	RevokeAllByUser(ctx context.Context, userID uuid.UUID) error
	DeleteAllExpired(ctx context.Context) error
}
