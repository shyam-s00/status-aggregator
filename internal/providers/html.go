package providers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"status-aggregator/internal/models"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type HtmlProvider struct {
	client *http.Client
}

func NewHtmlProvider() *HtmlProvider {
	return &HtmlProvider{client: &http.Client{
		Timeout: 15 * time.Second,
	}}
}

func (p *HtmlProvider) Fetch(ctx context.Context, sys models.SystemConfig) ([]models.Incident, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", sys.Url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request for %s: %w", sys.Name, err)
	}

	req.Header.Set("User-Agent", "StatusAggregator/1.0")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching HTML page %s: %w", sys.Name, err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("error closing response body: %v\n", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d fetching HTML page %s", resp.StatusCode, sys.Name)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error parsing HTML page %s: %w", sys.Name, err)
	}

	text := extractText(doc)

	//body, err := io.ReadAll(resp.Body)
	//if err != nil {
	//	return nil, fmt.Errorf("error reading response body for %s: %w", sys.Name, err)
	//}

	content := strings.ToLower(text)

	badKeywords := []string{"maintenance", "outage", "down", "degraded", "unavailable", "failure"}

	isOngoing := false
	statusText := "Operational"
	// TODO: This logic doesn't work for all status pages. should be revisited
	for _, kw := range badKeywords {
		if strings.Contains(content, kw) {
			isOngoing = true
			statusText = fmt.Sprintf("Potential Issue detected: %s", kw)
			break
		}
	}

	inc := models.Incident{
		Title:     statusText,
		Status:    "***",
		Url:       sys.Url,
		UpdatedAt: time.Now(),
		IsOngoing: isOngoing,
	}

	return []models.Incident{inc}, nil
}

func extractText(n *html.Node) string {
	if n.Type == html.ElementNode {
		if n.Data == "script" || n.Data == "style" || n.Data == "head" || n.Data == "noscript" {
			return ""
		}
	}

	if n.Type == html.TextNode {
		return strings.TrimSpace(n.Data)
	}

	var buf bytes.Buffer
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text := extractText(c)
		if text != "" {
			buf.WriteString(text)
		}
	}
	return buf.String()
}
