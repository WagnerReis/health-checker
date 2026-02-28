package repository

import (
	"context"
	entities "health-checker/internal/domain/entity"
)

type HealthCheckRepository interface {
	Create(ctx context.Context, healthCheck *entities.HealthCheck) error
}
