package engine

import (
	"context"
	"fmt"
	"status-aggregator/internal/models"
	"time"

	"github.com/mmcdole/gofeed"
)

type Result struct {
	SystemId  string
	Incidents []models.Incident
	Error     error
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

				incidents, err := fetchRSS(ctx, s)
				results <- Result{SystemId: s.Id, Incidents: incidents, Error: err}
			}(sys)
		}
		for range systems {
			<-done
		}
	}()

	return results
}

func fetchRSS(ctx context.Context, sys models.SystemConfig) ([]models.Incident, error) {
	fp := gofeed.NewParser()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	feed, err := fp.ParseURLWithContext(sys.Url, ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetching RSS feed %s: %w", sys.Name, err)
	}

	var incidents []models.Incident
	for _, item := range feed.Items {
		inc := models.Incident{
			IncidentId: item.GUID,
			Title:      item.Title,
			Status:     "Unknown",
			Url:        item.Link,
			UpdatedAt:  *item.PublishedParsed,
			IsOngoing:  false,
		}

		// if published within last 24 hours, may be relevant?
		// or just return all items from the feed?
		incidents = append(incidents, inc)
	}
	return incidents, nil
}
