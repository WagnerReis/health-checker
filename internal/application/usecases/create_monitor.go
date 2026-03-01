package usecases

import (
	"context"
	application "health-checker/internal/application/logger"
	entities "health-checker/internal/domain/entity"
	"health-checker/internal/domain/repository"
	register "health-checker/internal/infra/regiter"
	"time"

	"github.com/google/uuid"
)

type CreateMonitorCommand struct {
	UserID             uuid.UUID
	Name               string
	URL                string
	Method             string
	Headers            map[string]string
	Body               string
	Interval           int
	ExpectedStatusCode int32
	Timeout            int
}

type CreateMonitorUseCase struct {
	monitorRepository repository.MonitorRepository
	logger            application.Logger
	monitorRegister   *register.MonitorRegister
}

func NewCreateMonitorUseCase(monitorRepository repository.MonitorRepository, logger application.Logger, monitorRegister *register.MonitorRegister) *CreateMonitorUseCase {
	return &CreateMonitorUseCase{
		monitorRepository: monitorRepository,
		logger:            logger,
		monitorRegister:   monitorRegister,
	}
}

func (u *CreateMonitorUseCase) Execute(ctx context.Context, cmd CreateMonitorCommand) error {
	method := entities.MonitorMethod(cmd.Method)
	timeoutPtr := time.Duration(cmd.Timeout) * time.Second
	monitor, err := entities.NewMonitor(
		uuid.Nil,
		cmd.UserID,
		cmd.Name,
		cmd.URL,
		method,
		&cmd.Headers,
		&cmd.Body,
		time.Duration(cmd.Interval)*time.Second,
		&cmd.ExpectedStatusCode,
		&timeoutPtr,
	)
	if err != nil {
		return err
	}
	err = u.monitorRepository.Create(ctx, monitor)
	if err != nil {
		u.logger.Error("Failed to create monitor", application.Field{Key: "error", Value: err.Error()})
		return err
	}
	u.logger.Info("Monitor created successfully", application.Field{Key: "monitor_id", Value: monitor.ID.String()})
	err = u.monitorRegister.Register(monitor)
	if err != nil {
		u.logger.Error("Failed to register monitor", application.Field{Key: "error", Value: err.Error()})
		return err
	}
	return nil
}
