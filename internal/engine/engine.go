package engine

import (
	"context"
	"fmt"
	"status-aggregator/internal/models"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
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

				incidents, active, err := fetchRSS(ctx, s)
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

func fetchRSS(ctx context.Context, sys models.SystemConfig) ([]models.Incident, bool, error) {
	fp := gofeed.NewParser()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	feed, err := fp.ParseURLWithContext(sys.Url, ctx)
	if err != nil {
		return nil, false, fmt.Errorf("error fetching RSS feed %s: %w", sys.Name, err)
	}

	// for now let's make a naive check to identify if there is an active incident
	hasActiveIncident := false
	if len(feed.Items) > 0 {
		latestTitle := strings.ToLower(feed.Items[0].Title)
		if strings.Contains(latestTitle, "investigating") ||
			strings.Contains(latestTitle, "monitoring") ||
			strings.Contains(latestTitle, "identified") ||
			strings.Contains(latestTitle, "acknowledged") ||
			strings.Contains(latestTitle, "degraded") ||
			strings.Contains(latestTitle, "outage") ||
			strings.Contains(latestTitle, "incident") ||
			!strings.Contains(latestTitle, "resolved") ||
			!strings.Contains(latestTitle, "completed") {
			hasActiveIncident = true
		}
	}

	limit := 5 //TODO: make configurable per system through config file
	if len(feed.Items) < limit {
		limit = len(feed.Items)
	}

	var incidents []models.Incident

	for i := 0; i < limit; i++ {
		item := feed.Items[i]
		t := time.Now()

		if item.PublishedParsed != nil {
			t = *item.PublishedParsed
		} else if item.UpdatedParsed != nil {
			t = *item.UpdatedParsed
		}

		inc := models.Incident{
			IncidentId: item.GUID,
			Title:      item.Title,
			Status:     "Unknown",
			Url:        item.Link,
			UpdatedAt:  t,
			IsOngoing:  hasActiveIncident && i == 0,
		}
		incidents = append(incidents, inc)
	}

	//for _, item := range feed.Items {
	//	inc := models.Incident{
	//		IncidentId: item.GUID,
	//		Title:      item.Title,
	//		Status:     "Unknown",
	//		Url:        item.Link,
	//		UpdatedAt:  *item.PublishedParsed,
	//		IsOngoing:  false,
	//	}
	//
	//	// if published within last 24 hours, may be relevant?
	//	// or just return all items from the feed?
	//	incidents = append(incidents, inc)
	//}

	return incidents, hasActiveIncident, nil
}
