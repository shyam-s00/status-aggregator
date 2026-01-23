package main

import (
	"context"
	"fmt"
	"status-aggregator/internal/engine"
	"status-aggregator/internal/models"
)

func main() {
	systems := []models.SystemConfig{
		{Id: "ct01", Name: "CommerceTools", Url: "https://status.commercetools.com/pages/56e4295370fe4ece420002bb/rss", Type: "rss"},
		{Id: "og01", Name: "OrderGroove", Url: "https://status.ordergroove.com/history.rss", Type: "rss"},
	}

	fmt.Printf("🚀 Starting Status Aggregator with %d systems...\n\n", len(systems))

	// Create a context that can be canceled (useful for grateful shutdown later)
	ctx := context.Background()

	results := engine.Scrape(ctx, systems)

	// Main loop
	for result := range results {
		if result.Error != nil {
			fmt.Printf("❌ Error processing system %s: %v\n", result.SystemId, result.Error)
			continue
		}

		fmt.Printf("✅ Fetched %d items for %s\n", len(result.Incidents), result.SystemId)
		for _, inc := range result.Incidents {
			fmt.Printf("  %s | %s\n", inc.UpdatedAt.Format("2006-01-02 15:04"), inc.Title)
		}

		fmt.Println("-----------------------------------------------------")
	}

	fmt.Println("🏁🏁 All Aggregation finished. 🏁🏁")

}
