package entities

import (
	"time"

	uuid "github.com/google/uuid"
)

type HealthCheck struct {
	ID             uuid.UUID
	MonitorID      uuid.UUID
	StatusCode     *uint32
	ResponseTimeMS *uint32
	IsSuccess      *bool
	ErrorMessage   *string
	CheckedAt      *time.Time
	CreatedAt      time.Time
}

func NewHealthCheck(monitorID uuid.UUID) *HealthCheck {
	checkedAt := time.Now()
	return &HealthCheck{
		ID:             uuid.New(),
		MonitorID:      monitorID,
		StatusCode:     nil,
		ResponseTimeMS: nil,
		IsSuccess:      nil,
		ErrorMessage:   nil,
		CheckedAt:      &checkedAt,
		CreatedAt:      time.Now(),
	}
}

func (h *HealthCheck) SetStatusCode(statusCode uint32) {
	h.StatusCode = &statusCode
}

func (h *HealthCheck) SetResponseTimeMS(responseTimeMS uint32) {
	h.ResponseTimeMS = &responseTimeMS
}

func (h *HealthCheck) SetIsSuccess(isSuccess bool) {
	h.IsSuccess = &isSuccess
}

func (h *HealthCheck) SetErrorMessage(errorMessage string) {
	h.ErrorMessage = &errorMessage
}
