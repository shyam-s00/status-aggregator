package engine

import (
	"context"
	"fmt"
	"status-aggregator/internal/models"
	"status-aggregator/internal/providers"

	"golang.org/x/sync/errgroup"
)

type Engine struct {
	systems []models.SystemConfig
}

func NewEngine(systems []models.SystemConfig) *Engine {
	return &Engine{systems: systems}
}

type Result struct {
	SystemId          string
	SystemName        string
	Incidents         []models.Incident
	HasActiveIncident bool
	Error             error
}

func (e *Engine) Run(ctx context.Context) <-chan Result {
	results := make(chan Result, len(e.systems))

	g, ctx := errgroup.WithContext(ctx)

	// start goroutines for each system
	for _, sys := range e.systems {
		sys := sys
		g.Go(func() error {
			var provider providers.Provider
			// TODO: Explore to see if there is a better way to do this...
			switch sys.Type {
			case "rss":
				provider = providers.NewRSSProvider()
			case "html":
				provider = providers.NewHtmlProvider()
			default:
				results <- Result{
					SystemId:          sys.Id,
					SystemName:        sys.Name,
					Incidents:         nil,
					HasActiveIncident: false,
					Error:             fmt.Errorf("unknown provider type %s", sys.Type),
				}
				return nil
			}

			incidents, err := provider.Fetch(ctx, sys)

			hasActiveIncident := false
			if err == nil {
				for _, inc := range incidents {
					if inc.IsOngoing {
						hasActiveIncident = true
						break
					}
				}
			}
			results <- Result{
				SystemId:          sys.Id,
				SystemName:        sys.Name,
				Incidents:         incidents,
				HasActiveIncident: hasActiveIncident,
				Error:             err,
			}
			return nil
		})
	}

	//close the results channel when done
	go func() {
		_ = g.Wait()
		close(results)
	}()

	return results
}
