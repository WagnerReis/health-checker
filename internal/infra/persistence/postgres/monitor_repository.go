package postgres

import (
	"context"
	"database/sql"
	entities "health-checker/internal/domain/entity"
	"health-checker/internal/infra/persistence/database/sqlc"
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
