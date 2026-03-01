package repository

import (
	"context"
	entities "health-checker/internal/domain/entity"

	"github.com/google/uuid"
)

type MonitorRepository interface {
	Create(ctx context.Context, monitor *entities.Monitor) error
	FindByUserID(ctx context.Context, userID uuid.UUID, status *entities.MonitorStatus, limit, offset int32) ([]*entities.Monitor, error)
	CountByUserID(ctx context.Context, userID uuid.UUID, status *entities.MonitorStatus) (int64, error)
	GetAll(ctx context.Context) ([]*entities.Monitor, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entities.Monitor, error)
	Update(ctx context.Context, monitor *entities.Monitor) error
}
