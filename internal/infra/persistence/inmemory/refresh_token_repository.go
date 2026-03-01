package repository

import (
	"context"
	entities "health-checker/internal/domain/entity"
	domainerrors "health-checker/internal/domain/errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

type RefreshTokenRepositoryInMemory struct {
	refreshTokens map[uuid.UUID]*entities.RefreshToken
	mu            sync.Mutex
}

func NewRefreshTokenRepositoryInMemory() *RefreshTokenRepositoryInMemory {
	return &RefreshTokenRepositoryInMemory{
		refreshTokens: make(map[uuid.UUID]*entities.RefreshToken),
		mu:            sync.Mutex{},
	}
}

func (r *RefreshTokenRepositoryInMemory) Create(ctx context.Context, refreshToken *entities.RefreshToken) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.refreshTokens[refreshToken.ID] = refreshToken
	return nil
}

func (r *RefreshTokenRepositoryInMemory) FindByTokenHash(ctx context.Context, tokenHash string) (*entities.RefreshToken, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, refreshToken := range r.refreshTokens {
		if refreshToken.TokenHash == tokenHash {
			return refreshToken, nil
		}
	}
	return nil, domainerrors.ErrRefreshTokenNotFound
}

func (r *RefreshTokenRepositoryInMemory) Revoke(ctx context.Context, refreshTokenHash string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, refreshToken := range r.refreshTokens {
		if refreshToken.TokenHash == refreshTokenHash {
			refreshToken.Revoked = true
			return nil
		}
	}
	return domainerrors.ErrRefreshTokenNotFound
}

func (r *RefreshTokenRepositoryInMemory) RevokeAllByUser(ctx context.Context, userID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, refreshToken := range r.refreshTokens {
		if refreshToken.UserID == userID {
			refreshToken.Revoked = true
		}
	}
	return nil
}

func (r *RefreshTokenRepositoryInMemory) DeleteAllExpired(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, refreshToken := range r.refreshTokens {
		if refreshToken.ExpiresAt.Before(time.Now()) {
			delete(r.refreshTokens, refreshToken.ID)
		}
	}
	return nil
}
