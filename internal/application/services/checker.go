package services

import (
	"context"
	"errors"
	"fmt"
	application "health-checker/internal/application/logger"
	"health-checker/internal/domain/repository"
	"time"

	"github.com/google/uuid"
)

type CheckerService struct {
	monitorRepository repository.MonitorRepository
	logger            application.Logger
}

func NewCheckerService(monitorRepository repository.MonitorRepository, logger application.Logger) *CheckerService {
	return &CheckerService{monitorRepository: monitorRepository, logger: logger}
}

func (s *CheckerService) Check(ctx context.Context, monitorID uuid.UUID) error {
	// TODO: implementar
	ticker := time.NewTicker(time.Second * 2)
	for {
		select {
		case <-ctx.Done():
			return errors.New("context canceled")
		case <-ticker.C:
			fmt.Println("Checking monitor", monitorID)
		}
	}
}
