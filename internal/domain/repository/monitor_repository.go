package repository

import (
	"context"
	entities "health-checker/internal/domain/entity"
)

type MonitorRepository interface {
	Create(ctx context.Context, monitor *entities.Monitor) error
}
