package postgres

import (
	"context"
	"database/sql"
	entities "health-checker/internal/domain/entity"
	domainerrors "health-checker/internal/domain/errors"
	"health-checker/internal/infra/persistence/database/sqlc"

	"github.com/google/uuid"
)

type RefreshTokenRepository struct {
	queries *sqlc.Queries
}

func NewRefreshTokenRepository(db *sql.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{queries: sqlc.New(db)}
}

func (r *RefreshTokenRepository) Create(ctx context.Context, refreshToken *entities.RefreshToken) error {
	err := r.queries.CreateRefreshToken(ctx, sqlc.CreateRefreshTokenParams{
		ID:        refreshToken.ID,
		UserID:    refreshToken.UserID,
		TokenHash: refreshToken.TokenHash,
		ExpiresAt: refreshToken.ExpiresAt,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *RefreshTokenRepository) FindByTokenHash(ctx context.Context, tokenHash string) (*entities.RefreshToken, error) {
	refreshToken, err := r.queries.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, err
	}
	return &entities.RefreshToken{
		ID:        refreshToken.ID,
		UserID:    refreshToken.UserID,
		TokenHash: refreshToken.TokenHash,
		ExpiresAt: refreshToken.ExpiresAt,
		Revoked:   refreshToken.Revoked.Bool,
		CreatedAt: refreshToken.CreatedAt,
		UpdatedAt: refreshToken.UpdatedAt.Time,
	}, nil
}

func (r *RefreshTokenRepository) Revoke(ctx context.Context, refreshTokenHash string) error {
	affectedRows, err := r.queries.RevokeRefreshToken(ctx, refreshTokenHash)
	if affectedRows == 0 {
		return domainerrors.ErrRefreshTokenNotFound
	}
	if err != nil {
		if IsNoRowsError(err) {
			return domainerrors.ErrRefreshTokenNotFound
		}
		return err
	}
	return nil
}

func (r *RefreshTokenRepository) RevokeAllByUser(ctx context.Context, userID uuid.UUID) error {
	err := r.queries.RevokeAllByUser(ctx, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *RefreshTokenRepository) DeleteAllExpired(ctx context.Context) error {
	err := r.queries.DeleteAllExpired(ctx)
	if err != nil {
		return err
	}
	return nil
}
