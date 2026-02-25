package usecases

import (
	"context"
	application "health-checker/internal/application/logger"
	entities "health-checker/internal/domain/entity"
	"health-checker/internal/domain/repository"

	"github.com/google/uuid"
)

type GetMonitorsCommand struct {
	UserID uuid.UUID
	Status entities.MonitorStatus
	Limit  int32
	Offset int32
}

type GetMonitorsUseCase struct {
	monitorRepository repository.MonitorRepository
	logger            application.Logger
}

func NewGetMonitorsUseCase(monitorRepository repository.MonitorRepository, logger application.Logger) *GetMonitorsUseCase {
	return &GetMonitorsUseCase{monitorRepository: monitorRepository, logger: logger}
}

func (u *GetMonitorsUseCase) Execute(ctx context.Context, cmd GetMonitorsCommand) ([]*entities.Monitor, error) {
	monitors, err := u.monitorRepository.FindByUserID(ctx, cmd.UserID, &cmd.Status, cmd.Limit, cmd.Offset)
	if err != nil {
		u.logger.Error("Failed to get monitors", application.Field{Key: "error", Value: err.Error()})
		return nil, err
	}
	return monitors, nil
}
