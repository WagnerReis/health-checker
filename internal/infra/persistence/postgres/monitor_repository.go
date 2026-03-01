package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	entities "health-checker/internal/domain/entity"
	domainerrors "health-checker/internal/domain/errors"
	"health-checker/internal/infra/persistence/database/sqlc"
	"strings"
	"time"

	"github.com/google/uuid"
)

type MonitorRepository struct {
	queries *sqlc.Queries
}

func NewMonitorRepository(db *sql.DB) *MonitorRepository {
	return &MonitorRepository{queries: sqlc.New(db)}
}

func (r *MonitorRepository) Create(ctx context.Context, monitor *entities.Monitor) error {
	err := r.queries.CreateMonitor(ctx, sqlc.CreateMonitorParams{
		ID:                 monitor.ID,
		UserID:             monitor.UserID,
		Name:               monitor.Name,
		Url:                monitor.URL,
		Method:             string(monitor.Method),
		Headers:            NullRawMessage(monitor.Headers),
		Body:               NullString(monitor.Body),
		Interval:           int32(monitor.Interval.Seconds()),
		ExpectedStatusCode: NullInt32(monitor.ExpectedStatusCode),
		Timeout:            int32(monitor.Timeout.Seconds()),
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *MonitorRepository) FindByUserID(
	ctx context.Context,
	userID uuid.UUID,
	status *entities.MonitorStatus,
	limit,
	offset int32,
) ([]*entities.Monitor, error) {
	statusString := strings.ToUpper(status.String())
	monitors, err := r.queries.FindMonitorsByUserID(ctx, sqlc.FindMonitorsByUserIDParams{
		UserID:     userID,
		Status:     NullString(&statusString),
		PageLimit:  limit,
		PageOffset: offset,
	})
	if err != nil {
		return nil, err
	}
	monitorsEntities := make([]*entities.Monitor, len(monitors))
	for i, monitor := range monitors {
		timeout := time.Duration(monitor.Timeout) * time.Second

		var headers *map[string]string
		if monitor.Headers.Valid && len(monitor.Headers.RawMessage) > 0 {
			h := make(map[string]string)
			if err := json.Unmarshal(monitor.Headers.RawMessage, &h); err == nil {
				headers = &h
			}
			headers = &h

		}

		monitorsEntities[i] = &entities.Monitor{
			ID:                 monitor.ID,
			UserID:             monitor.UserID,
			Name:               monitor.Name,
			URL:                monitor.Url,
			Method:             entities.MonitorMethod(monitor.Method),
			Status:             entities.MonitorStatus(monitor.Status),
			Headers:            headers,
			Body:               &monitor.Body.String,
			Interval:           time.Duration(monitor.Interval) * time.Second,
			ExpectedStatusCode: &monitor.ExpectedStatusCode.Int32,
			Timeout:            &timeout,
			CreatedAt:          monitor.CreatedAt,
			UpdatedAt:          monitor.UpdatedAt,
			DeletedAt:          &monitor.DeletedAt.Time,
		}
	}
	return monitorsEntities, nil
}

func (r *MonitorRepository) CountByUserID(ctx context.Context, userID uuid.UUID, status *entities.MonitorStatus) (int64, error) {
	statusString := strings.ToUpper(status.String())
	count, err := r.queries.CountMonitorsByUserID(ctx, sqlc.CountMonitorsByUserIDParams{
		UserID: userID,
		Status: NullString(&statusString),
	})
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *MonitorRepository) GetAll(ctx context.Context) ([]*entities.Monitor, error) {
	monitors, err := r.queries.GetAllMonitors(ctx)
	if err != nil {
		return nil, err
	}
	monitorsEntities := make([]*entities.Monitor, len(monitors))
	for i, monitor := range monitors {
		timeout := time.Duration(monitor.Timeout) * time.Second

		monitorsEntities[i] = &entities.Monitor{
			ID:                 monitor.ID,
			UserID:             monitor.UserID,
			Name:               monitor.Name,
			URL:                monitor.Url,
			Method:             entities.MonitorMethod(monitor.Method),
			Status:             entities.MonitorStatus(monitor.Status),
			Body:               &monitor.Body.String,
			Interval:           time.Duration(monitor.Interval) * time.Second,
			ExpectedStatusCode: &monitor.ExpectedStatusCode.Int32,
			Timeout:            &timeout,
			CreatedAt:          monitor.CreatedAt,
			UpdatedAt:          monitor.UpdatedAt,
			DeletedAt:          &monitor.DeletedAt.Time,
		}
	}
	return monitorsEntities, nil
}

func (r *MonitorRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Monitor, error) {
	monitor, err := r.queries.FindMonitorByID(ctx, id)
	if err != nil {
		if IsNoRowsError(err) {
			return nil, domainerrors.ErrMonitorNotFound
		}
		return nil, err
	}
	timeout := time.Duration(monitor.Timeout) * time.Second
	return &entities.Monitor{
		ID:                 monitor.ID,
		UserID:             monitor.UserID,
		Name:               monitor.Name,
		URL:                monitor.Url,
		Method:             entities.MonitorMethod(monitor.Method),
		Status:             entities.MonitorStatus(monitor.Status),
		Body:               &monitor.Body.String,
		Interval:           time.Duration(monitor.Interval) * time.Second,
		ExpectedStatusCode: &monitor.ExpectedStatusCode.Int32,
		Timeout:            &timeout,
		CreatedAt:          monitor.CreatedAt,
		UpdatedAt:          monitor.UpdatedAt,
		DeletedAt:          &monitor.DeletedAt.Time,
	}, nil
}

func (r *MonitorRepository) Update(ctx context.Context, monitor *entities.Monitor) error {
	err := r.queries.UpdateMonitor(ctx, sqlc.UpdateMonitorParams{
		ID:                 monitor.ID,
		Status:             string(monitor.Status),
		Name:               monitor.Name,
		Url:                monitor.URL,
		Method:             string(monitor.Method),
		Headers:            NullRawMessage(monitor.Headers),
		Body:               NullString(monitor.Body),
		Interval:           int32(monitor.Interval.Seconds()),
		ExpectedStatusCode: NullInt32(monitor.ExpectedStatusCode),
		Timeout:            int32(monitor.Timeout.Seconds()),
	})
	if err != nil {
		return err
	}
	return nil
}
