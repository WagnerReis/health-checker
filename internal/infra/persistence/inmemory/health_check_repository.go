package repository

import (
	"context"
	entities "health-checker/internal/domain/entity"
	"sync"

	"github.com/google/uuid"
)

type HealthCheckRepositoryInMemory struct {
	healthChecks map[uuid.UUID]*entities.HealthCheck
	mu           sync.Mutex
}

func NewHealthCheckRepositoryInMemory() *HealthCheckRepositoryInMemory {
	return &HealthCheckRepositoryInMemory{
		healthChecks: make(map[uuid.UUID]*entities.HealthCheck),
		mu:           sync.Mutex{},
	}
}

func (r *HealthCheckRepositoryInMemory) Create(ctx context.Context, healthCheck *entities.HealthCheck) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.healthChecks[healthCheck.ID] = healthCheck
	return nil
}
