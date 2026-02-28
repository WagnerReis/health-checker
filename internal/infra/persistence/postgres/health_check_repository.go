package postgres

import (
	"context"
	"database/sql"
	entities "health-checker/internal/domain/entity"
	"health-checker/internal/infra/persistence/database/sqlc"
)

type HealthCheckRepository struct {
	queries *sqlc.Queries
}

func NewHealthCheckRepository(db *sql.DB) *HealthCheckRepository {
	return &HealthCheckRepository{queries: sqlc.New(db)}
}

func (r *HealthCheckRepository) Create(ctx context.Context, healthCheck *entities.HealthCheck) error {
	statusCode := int32(*healthCheck.StatusCode)
	responseTimeMS := int32(*healthCheck.ResponseTimeMS)
	err := r.queries.CreateHealthCheck(ctx, sqlc.CreateHealthCheckParams{
		ID:             healthCheck.ID,
		MonitorID:      healthCheck.MonitorID,
		StatusCode:     NullInt32(&statusCode),
		ResponseTimeMs: NullInt32(&responseTimeMS),
		IsSuccess:      *healthCheck.IsSuccess,
		ErrorMessage:   NullString(healthCheck.ErrorMessage),
		CheckedAt:      *healthCheck.CheckedAt,
	})
	if err != nil {
		return err
	}
	return nil
}
