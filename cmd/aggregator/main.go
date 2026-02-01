package main

import (
	"context"
	"fmt"
	"log"
	"status-aggregator/internal/config"
	"status-aggregator/internal/engine"
	"time"
)

func main() {

	systems, err := config.Load("config.json")
	if err != nil {
		log.Fatalf("❌ Could not load config: %v", err)
	}

	fmt.Printf("🚀 Starting Status Aggregator with %d systems...\n\n", len(systems))

	// Create a context that can be canceled (useful for grateful shutdown later)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	eng := engine.NewEngine(systems)
	results := eng.Run(ctx)

	// Main loop
	for result := range results {
		// start with a header
		fmt.Printf("\n🔹 System: %s (%s)\n", result.SystemName, result.SystemId)

		if result.Error != nil {
			fmt.Printf("❌ Error processing system %s: %v\n", result.SystemId, result.Error)
			fmt.Println("-----------------------------------------------------")
			continue
		}

		if result.HasActiveIncident {
			details := ""
			if len(result.Incidents) > 0 {
				details = fmt.Sprintf(" (most recent incident: %s)", result.Incidents[0].Title)
			}
			fmt.Println("   ⚠️  Status: ACTIVE INCIDENT DETECTED\n", details)
		} else {
			fmt.Println("   ✅  Status: Operational / No active incidents")
		}
		fmt.Println("-----------------------------------------------------")

		if len(result.Incidents) > 0 {
			fmt.Println("   ✅  Recent History:")
			for _, inc := range result.Incidents {
				fmt.Printf("  %s | %s\n", inc.UpdatedAt.Format("2006-01-02 15:04"), inc.Title)
			}
		} else {
			fmt.Println("    (No recent history available)")
		}

		fmt.Println("-----------------------------------------------------")
	}

	fmt.Println("\n🏁🏁 All Aggregation finished. 🏁🏁")

}
