package repository

import (
	"context"
	entities "health-checker/internal/domain/entity"
	"sync"

	"github.com/google/uuid"
)

type MonitorRepositoryInMemory struct {
	monitors map[uuid.UUID]*entities.Monitor
	mu       sync.Mutex
}

func NewMonitorRepositoryInMemory() *MonitorRepositoryInMemory {
	return &MonitorRepositoryInMemory{
		monitors: make(map[uuid.UUID]*entities.Monitor),
		mu:       sync.Mutex{},
	}
}

func (r *MonitorRepositoryInMemory) Create(ctx context.Context, monitor *entities.Monitor) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.monitors[monitor.ID] = monitor
	return nil
}
