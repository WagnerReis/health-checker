package repository

import (
	"context"
	entities "health-checker/internal/domain/entity"

	"github.com/gofrs/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	Update(ctx context.Context, user *entities.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
}
