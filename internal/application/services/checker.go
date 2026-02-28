package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	application "health-checker/internal/application/logger"
	entities "health-checker/internal/domain/entity"
	"health-checker/internal/domain/repository"
	"io"
	"net/http"
	"time"
)

type CheckerService struct {
	healthCheckRepository repository.HealthCheckRepository
	logger                application.Logger
}

func NewCheckerService(healthCheckRepository repository.HealthCheckRepository, logger application.Logger) *CheckerService {
	return &CheckerService{healthCheckRepository: healthCheckRepository, logger: logger}
}

func (s *CheckerService) Check(ctx context.Context, monitor *entities.Monitor) error {
	ticker := time.NewTicker(monitor.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return errors.New("context canceled")
		case <-ticker.C:
			healthCheck := entities.NewHealthCheck(monitor.ID)

			var timeout time.Duration
			if monitor.Timeout == nil {
				timeout = 15 * time.Second
			} else {
				timeout = time.Duration(*monitor.Timeout*2) * time.Second
			}

			ctxTimeout, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			request, err := s.buildRequest(ctxTimeout, monitor)
			if err != nil {
				s.logger.Error("failed to build request", application.Field{Key: "error", Value: err.Error()})
				continue
			}

			start := time.Now()

			response, err := http.DefaultClient.Do(request)
			if err != nil {
				s.logger.Error("failed to perform health check", application.Field{Key: "error", Value: err.Error()})

				healthCheck.SetStatusCode(500)
				healthCheck.SetIsSuccess(false)
				healthCheck.SetErrorMessage(err.Error())

				duration := time.Since(start)
				durationInMS := uint32(duration.Milliseconds())
				healthCheck.ResponseTimeMS = &durationInMS
				err := s.healthCheckRepository.Create(context.Background(), healthCheck)
				if err != nil {
					s.logger.Error("failed to save health check", application.Field{Key: "error", Value: err.Error()})
				}
				continue
			}
			defer response.Body.Close()

			healthCheck.SetStatusCode(uint32(response.StatusCode))
			duration := time.Since(start)
			durationInMS := uint32(duration.Milliseconds())
			healthCheck.ResponseTimeMS = &durationInMS
			if response.StatusCode >= 200 && response.StatusCode < 300 {
				healthCheck.SetIsSuccess(true)
			}

			err = s.healthCheckRepository.Create(context.Background(), healthCheck)
			if err != nil {
				s.logger.Error("failed to save health check", application.Field{Key: "error", Value: err.Error()})
			}

			s.logger.Info("health check completed", application.Field{Key: "status", Value: response.StatusCode})
		}
	}
}

func (s *CheckerService) buildRequest(ctx context.Context, monitor *entities.Monitor) (*http.Request, error) {
	var bodyReader io.Reader

	if monitor.Method.String() != "GET" && monitor.Body != nil {
		data, err := json.Marshal(monitor.Body)
		if err != nil {
			return nil, err
		}

		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		monitor.Method.String(),
		monitor.URL,
		bodyReader,
	)
	if err != nil {
		return nil, err
	}

	if monitor.Body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}
