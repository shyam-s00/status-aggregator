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

type SystemConfig struct {
	Id         string            `json:"id"`
	Name       string            `json:"name"`
	Url        string            `json:"url"`
	Type       string            `json:"type"`
	AuthToken  string            `json:"auth_token,omitempty"`
	HtmlConfig map[string]string `json:"html_config,omitempty"`
}
