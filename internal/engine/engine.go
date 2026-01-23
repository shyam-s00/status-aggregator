package engine

import (
	"context"
	"status-aggregator/internal/models"
	"status-aggregator/internal/providers"
)

type Result struct {
	SystemId          string
	SystemName        string
	Incidents         []models.Incident
	HasActiveIncident bool
	Error             error
}

func Scrape(ctx context.Context, systems []models.SystemConfig) <-chan Result {
	results := make(chan Result)

	go func() {
		defer close(results)

		done := make(chan struct{})
		for _, sys := range systems {
			go func(s models.SystemConfig) {
				defer func() { done <- struct{}{} }()

				select {
				case <-ctx.Done():
					return
				default:
				}

				// Determine a provider based on system config or default to RSS
				// For now, defaulting to RSS as per previous logic
				provider := providers.NewRSSProvider()
				incidents, err := provider.Fetch(ctx, s)

				active := false
				if len(incidents) > 0 {
					//TODO: this should be changed, as some feeds doesn't show active incident
					active = incidents[0].IsOngoing
				}

				results <- Result{
					SystemId:          s.Id,
					SystemName:        s.Name,
					Incidents:         incidents,
					HasActiveIncident: active,
					Error:             err,
				}
			}(sys)
		}
		for range systems {
			<-done
		}
	}()

	return results
}
