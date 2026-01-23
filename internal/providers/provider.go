package providers

import (
	"context"
	"status-aggregator/internal/models"
)

type StatusProvider interface {
	Fetch(ctx context.Context, sys models.SystemConfig) ([]models.Incident, error)
}
