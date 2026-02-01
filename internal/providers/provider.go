package providers

import (
	"context"
	"status-aggregator/internal/models"
)

type HistoryProvider interface {
	FetchHistory(ctx context.Context, sys models.SystemConfig) ([]models.Incident, error)
}

type StatusProvider interface {
	FetchStatus(ctx context.Context, url string, config map[string]string) (string, bool, error)
}
