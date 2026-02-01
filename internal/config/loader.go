package config

import (
	"encoding/json"
	"fmt"
	"os"
	"status-aggregator/internal/models"
)

func Load(path string) ([]models.SystemConfig, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var systems []models.SystemConfig
	if err := json.Unmarshal(file, &systems); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	return systems, nil
}
