package postgres

import (
	"context"
	"database/sql"
	entities "health-checker/internal/domain/entity"
	domainerrors "health-checker/internal/domain/errors"
	"health-checker/internal/infra/persistence/database/sqlc"

	"github.com/google/uuid"
)

type UserRepository struct {
	queries *sqlc.Queries
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{queries: sqlc.New(db)}
}

func (r *UserRepository) Create(ctx context.Context, user *entities.User) error {
	err := r.queries.CreateUser(ctx, sqlc.CreateUserParams{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) Update(ctx context.Context, user *entities.User) error {
	err := r.queries.Update(ctx, sqlc.UpdateParams{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	user, err := r.queries.FindByID(ctx, id)
	if err != nil {
		if IsNoRowsError(err) {
			return nil, domainerrors.ErrUserNotFound
		}
		return nil, err
	}
	return &entities.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	user, err := r.queries.FindByEmail(ctx, email)
	if err != nil {
		if IsNoRowsError(err) {
			return nil, domainerrors.ErrUserNotFound
		}
		return nil, err
	}
	return &entities.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
