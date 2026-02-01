package models

import "time"

type Incident struct {
	SystemId   string    `json:"system_id"`
	Provider   string    `json:"provider"`
	IncidentId string    `json:"incident_id"`
	Title      string    `json:"title"`
	Status     string    `json:"status"`
	Message    string    `json:"message"`
	IsOngoing  bool      `json:"is_ongoing"`
	Url        string    `json:"url"`
	UpdatedAt  time.Time `json:"updated_at"`
}
