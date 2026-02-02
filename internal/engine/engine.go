package engine

import (
	"context"
	"status-aggregator/internal/models"
	"status-aggregator/internal/providers"

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
			result, err := e.fetch(ctx, sys)
			if err != nil {
				result.Error = err
			}

			select {
			case results <- result:
			case <-ctx.Done():
				return ctx.Err()
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

func (e *Engine) fetch(ctx context.Context, sys models.SystemConfig) (models.Result, error) {
	result := models.Result{
		SystemId:   sys.Id,
		SystemName: sys.Name,
	}

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		status, isOngoing, err := e.htmlProvider.FetchStatus(ctx, sys.StatusUrl, sys.HtmlConfig)
		if err != nil {
			return err
		}
		result.CurrentStatus = status
		result.HasActiveIncident = isOngoing
		return nil
	})

	g.Go(func() error {
		if sys.Type == "rss" && sys.FeedUrl != "" {
			incidents, err := e.rssProvider.FetchHistory(ctx, sys)
			if err != nil {
				return err
			}
			result.Incidents = incidents
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return result, err
	}

	return result, nil
}
