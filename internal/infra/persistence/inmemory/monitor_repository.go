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

func (r *MonitorRepositoryInMemory) FindByUserID(ctx context.Context, userID uuid.UUID, status entities.MonitorStatus, limit, offset int32) ([]*entities.Monitor, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	monitors := []*entities.Monitor{}
	for _, monitor := range r.monitors {
		if monitor.UserID == userID && monitor.Status == status {
			monitors = append(monitors, monitor)
		}
	}
	return monitors, nil
}

func (r *MonitorRepositoryInMemory) CountByUserID(ctx context.Context, userID uuid.UUID, status *entities.MonitorStatus) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	count := 0
	for _, monitor := range r.monitors {
		if monitor.UserID == userID {
			if status != nil && monitor.Status == *status {
				count++
			}
			if status == nil {
				count++
			}
		}
	}
	return int64(count), nil
}

func (r *MonitorRepositoryInMemory) GetAll(ctx context.Context) ([]*entities.Monitor, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	monitors := make([]*entities.Monitor, len(r.monitors))
	for _, m := range r.monitors {
		monitors = append(monitors, m)
	}
	return monitors, nil
}
