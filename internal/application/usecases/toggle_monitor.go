package usecases

import (
	"context"
	application "health-checker/internal/application/logger"
	entities "health-checker/internal/domain/entity"
	"health-checker/internal/domain/repository"
	register "health-checker/internal/infra/regiter"

	"github.com/google/uuid"
)

type ToggleMonitorUseCase struct {
	monitorRepository repository.MonitorRepository
	monitorRegister   *register.MonitorRegister
	logger            application.Logger
}

func NewToggleMonitorUseCase(
	monitorRepository repository.MonitorRepository,
	monitorRegister *register.MonitorRegister,
	logger application.Logger,
) *ToggleMonitorUseCase {
	return &ToggleMonitorUseCase{
		monitorRepository: monitorRepository,
		monitorRegister:   monitorRegister,
		logger:            logger,
	}
}

func (u *ToggleMonitorUseCase) Execute(ctx context.Context, id uuid.UUID) error {
	monitor, err := u.monitorRepository.FindByID(ctx, id)
	if err != nil {
		return err
	}
	var newStatus entities.MonitorStatus
	if monitor.Status == entities.MonitorStatusUP {
		newStatus = entities.MonitorStatusDOWN
	} else {
		newStatus = entities.MonitorStatusUP
	}
	monitor.Status = newStatus
	err = u.monitorRepository.Update(ctx, monitor)
	if err != nil {
		return err
	}
	err = u.monitorRegister.Toggle(id)
	if err != nil {
		return err
	}
	return nil
}
