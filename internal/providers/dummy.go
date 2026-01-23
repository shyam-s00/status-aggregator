package providers

import (
	"fmt"
	"status-aggregator/internal/models"
	"time"
)

type DummyProvider struct{}

func NewDummyProvider() *DummyProvider {
	return &DummyProvider{}
}

func (d *DummyProvider) Fetch(sys models.SystemConfig) ([]models.Incident, error) {
	//TODO implement me
	time.Sleep(2 * time.Second)
	return []models.Incident{
		{
			SystemId:   sys.Id,
			Provider:   sys.Type,
			IncidentId: fmt.Sprintf("%s-incident-001", sys.Id),
			Title:      fmt.Sprintf("Incident 001 on %s", sys.Name),
			Status:     "investigating",
			IsOngoing:  true,
			UpdatedAt:  time.Now(),
			Url:        sys.Url,
		},
	}, nil
}
