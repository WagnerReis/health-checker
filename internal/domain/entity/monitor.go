package entities

import (
	valueobjects "health-checker/internal/shared/value-object"
	"time"

	domainerrors "health-checker/internal/domain/errors"

	uuid "github.com/google/uuid"
)

type MonitorStatus string

const (
	MonitorStatusUP   MonitorStatus = "ACTIVE"
	MonitorStatusDOWN MonitorStatus = "PAUSED"
)

func (s MonitorStatus) String() string {
	return string(s)
}

type MonitorMethod string

const (
	MonitorMethodGET  MonitorMethod = "GET"
	MonitorMethodPOST MonitorMethod = "POST"
	MonitorMethodHEAD MonitorMethod = "HEAD"
)

func (s MonitorMethod) String() string {
	return string(s)
}

type Monitor struct {
	ID                 uuid.UUID
	UserID             uuid.UUID
	Name               string
	URL                string
	Method             MonitorMethod
	Status             MonitorStatus
	Headers            *map[string]string
	Body               *string
	Interval           time.Duration
	ExpectedStatusCode *int32
	Timeout            *time.Duration
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          *time.Time
}

func NewMonitor(
	id, userID uuid.UUID,
	name, url string,
	method MonitorMethod,
	headers *map[string]string,
	body *string,
	interval time.Duration,
	expectedStatusCode *int32,
	timeout *time.Duration,
) (*Monitor, error) {
	if id == uuid.Nil {
		id = valueobjects.NewID(uuid.Nil).Value()
	}
	if userID == uuid.Nil {
		return nil, domainerrors.ErrUserIDRequired
	}
	if name == "" || url == "" || method == "" {
		return nil, domainerrors.ErrRequiredFields
	}
	if interval == 0 {
		return nil, domainerrors.ErrIntervalRequired
	}
	return &Monitor{
		ID:                 id,
		UserID:             userID,
		Name:               name,
		URL:                url,
		Method:             method,
		Headers:            headers,
		Body:               body,
		Interval:           interval,
		ExpectedStatusCode: expectedStatusCode,
		Status:             MonitorStatusUP,
		Timeout:            timeout,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}, nil
}
