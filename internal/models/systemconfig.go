package models

type SystemConfig struct {
	Id         string            `json:"id"`
	Name       string            `json:"name"`
	StatusUrl  string            `json:"status_url"`
	FeedUrl    string            `json:"feed_url,omitempty"`
	Type       string            `json:"type"`
	AuthToken  string            `json:"auth_token,omitempty"`
	HtmlConfig map[string]string `json:"html_config,omitempty"`
}
