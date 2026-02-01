package models

type Result struct {
	SystemId          string
	SystemName        string
	CurrentStatus     string
	Incidents         []Incident
	HasActiveIncident bool
	Error             error
}
