package providers

import "status-aggregator/internal/models"

type StatusProvider interface {
	Fetch(sys models.SystemConfig) ([]models.Incident, error)
}
