package worker

import (
	"context"
	application "health-checker/internal/application/logger"
	"health-checker/internal/application/services"
	entities "health-checker/internal/domain/entity"
	"sync"
)

type WorkerPool struct {
	monitors       chan entities.Monitor
	checkerService services.CheckerService
	maxWorkers     int32
	wg             *sync.WaitGroup
	logger         application.Logger
}

func NewWorkerPool(
	monitors chan entities.Monitor,
	checkerService services.CheckerService,
	maxWorkers uint32,
	logger application.Logger,
) *WorkerPool {
	return &WorkerPool{
		monitors:       monitors,
		checkerService: checkerService,
		maxWorkers:     int32(maxWorkers),
		wg:             &sync.WaitGroup{},
		logger:         logger,
	}
}

func (wp *WorkerPool) Start() {
	for range wp.maxWorkers {
		wp.wg.Add(1)
		go func() {
			defer wp.wg.Done()
			for monitor := range wp.monitors {
				err := wp.checkerService.Check(context.Background(), monitor.ID)
				if err != nil {
					wp.logger.Error("Error checking monitor", application.Field{Key: "error", Value: err.Error()})
				}
				wp.logger.Info("Monitor started", application.Field{Key: "monitor_id", Value: monitor.ID.String()})
			}
		}()
	}
}

func (wp *WorkerPool) Shutdown() {
	close(wp.monitors)
	wp.wg.Wait()
	wp.logger.Info("worker pool encerrado")
}
