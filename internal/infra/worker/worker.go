package worker

import (
	"context"
	application "health-checker/internal/application/logger"
	"health-checker/internal/application/services"
	entities "health-checker/internal/domain/entity"
	register "health-checker/internal/infra/regiter"
	"sync"
)

type WorkerPool struct {
	monitorRegister *register.MonitorRegister
	checkerService  services.CheckerService
	maxWorkers      int32
	wg              *sync.WaitGroup
	logger          application.Logger
}

func NewWorkerPool(
	monitorRegister *register.MonitorRegister,
	checkerService services.CheckerService,
	maxWorkers uint32,
	logger application.Logger,
) *WorkerPool {
	return &WorkerPool{
		monitorRegister: monitorRegister,
		checkerService:  checkerService,
		maxWorkers:      int32(maxWorkers),
		wg:              &sync.WaitGroup{},
		logger:          logger,
	}
}

func (wp *WorkerPool) Start() {
	for range wp.maxWorkers {
		wp.wg.Add(1)
		go func() {
			defer wp.wg.Done()
			wp.monitorRegister.Monitors.Range(func(key, value any) bool {
				monitor := value.(*entities.Monitor)

				err := wp.checkerService.Check(context.Background(), monitor)
				if err != nil {
					wp.logger.Error("Error checking monitor", application.Field{Key: "error", Value: err.Error()})
				}

				wp.logger.Info("Monitor started",
					application.Field{Key: "monitor_id", Value: monitor.ID.String()},
				)

				return true
			})
		}()
	}
}

func (wp *WorkerPool) Shutdown() {
	wp.wg.Wait()
	wp.logger.Info("worker pool encerrado")
}
