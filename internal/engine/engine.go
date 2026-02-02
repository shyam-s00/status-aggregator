package engine

import (
	"context"
	"status-aggregator/internal/models"
	"status-aggregator/internal/providers"
	"sync"

	"golang.org/x/sync/errgroup"
)

type Engine struct {
	systems      []models.SystemConfig
	htmlProvider *providers.HtmlProvider
	rssProvider  *providers.RSSProvider
}

func NewEngine(systems []models.SystemConfig) *Engine {
	return &Engine{
		systems:      systems,
		htmlProvider: providers.NewHtmlProvider(),
		rssProvider:  providers.NewRSSProvider(),
	}
}

func (e *Engine) Run(ctx context.Context) <-chan models.Result {
	results := make(chan models.Result, len(e.systems))

	g, ctx := errgroup.WithContext(ctx)

	// start goroutines for each system
	for _, sys := range e.systems {
		sys := sys
		g.Go(func() error {
			var wg sync.WaitGroup

			result := models.Result{
				SystemId:   sys.Id,
				SystemName: sys.Name,
			}

			// TODO: Explore to see if there is a better way to do this...
			// 1. Fetch current status
			wg.Add(1)
			go func() {
				defer wg.Done()
				status, isOngoing, err := e.htmlProvider.FetchStatus(ctx, sys.StatusUrl, sys.HtmlConfig)
				if err != nil {
					result.Error = err
				} else {
					result.CurrentStatus = status
					result.HasActiveIncident = isOngoing
				}
			}()

			// 2. Fetch history
			wg.Add(1)
			go func() {
				defer wg.Done()
				var incidents []models.Incident
				var err error

				if sys.Type == "rss" && sys.FeedUrl != "" {
					incidents, err = e.rssProvider.FetchHistory(ctx, sys)
				}
				if err == nil {
					result.Incidents = incidents
				} else {
					result.Error = err
				}
			}()

			wg.Wait()
			results <- result
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
