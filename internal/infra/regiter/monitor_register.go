package register

import (
	entities "health-checker/internal/domain/entity"
	"sync"

	domainerrors "health-checker/internal/domain/errors"

	"github.com/google/uuid"
)

type MonitorRegister struct {
	Monitors sync.Map
	mu       sync.Mutex
}

func NewMonitorRegister() *MonitorRegister {
	return &MonitorRegister{
		Monitors: sync.Map{},
		mu:       sync.Mutex{},
	}
}

func (r *MonitorRegister) Register(monitor *entities.Monitor) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.Monitors.Load(monitor.ID)
	if ok {
		return domainerrors.ErrMonitorAlreadyRegistered
	}
	r.Monitors.Store(monitor.ID, monitor)
	return nil
}

func (r *MonitorRegister) Toggle(monitorID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	value, ok := r.Monitors.Load(monitorID)
	if !ok {
		return domainerrors.ErrMonitorNotFound
	}

	monitor := value.(*entities.Monitor)

	if monitor.Status == entities.MonitorStatusUP {
		monitor.Status = entities.MonitorStatusDOWN
	} else {
		monitor.Status = entities.MonitorStatusUP
	}
	return nil
}
