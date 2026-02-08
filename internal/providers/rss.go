package providers

import (
	"context"
	"fmt"
	"status-aggregator/internal/models"
	"time"

	"github.com/mmcdole/gofeed"
)

type RSSProvider struct{}

func NewRSSProvider() *RSSProvider {
	return &RSSProvider{}
}

func (p *RSSProvider) FetchHistory(ctx context.Context, sys models.SystemConfig) ([]models.Incident, error) {
	fp := gofeed.NewParser()

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	feed, err := fp.ParseURLWithContext(sys.FeedUrl, ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetching RSS feed %s: %w", sys.Name, err)
	}

	limit := sys.HistoryLimit
	if limit == 0 {
		limit = 5 // default to 5 if not specified
	}

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
			IsOngoing:  false, // we rely on HTML scrapers to get this.
		}
		incidents = append(incidents, inc)
	}

	return incidents, nil
}
