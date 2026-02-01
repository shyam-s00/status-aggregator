package providers

import (
	"context"
	"fmt"
	"status-aggregator/internal/models"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

type RSSProvider struct{}

func NewRSSProvider() *RSSProvider {
	return &RSSProvider{}
}

func (p *RSSProvider) FetchStatus(_ context.Context, _ string, _ map[string]string) (string, bool, error) {
	return "Operational", false, nil
}

func (p *RSSProvider) FetchHistory(ctx context.Context, sys models.SystemConfig) ([]models.Incident, error) {
	fp := gofeed.NewParser()

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	feed, err := fp.ParseURLWithContext(sys.FeedUrl, ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetching RSS feed %s: %w", sys.Name, err)
	}

	// TODO: for now let's make a naive check to identify if there is an active incident
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

	return incidents, nil
}
