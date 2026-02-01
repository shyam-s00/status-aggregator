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

func (p *HtmlProvider) FetchHistory(_ context.Context, _ models.SystemConfig) ([]models.Incident, error) {
	return []models.Incident{}, nil
}

func (p *HtmlProvider) FetchStatus(ctx context.Context, url string, config map[string]string) (string, bool, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", false, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("User-Agent", "StatusAggregator/1.0")

	resp, err := p.client.Do(req)
	if err != nil {
		return "", false, fmt.Errorf("error fetching HTML page: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("error closing response body: %v\n", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", false, fmt.Errorf("unexpected status code %d fetching HTML page", resp.StatusCode)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", false, fmt.Errorf("error parsing HTML page: %w", err)
	}

	//body, err := io.ReadAll(resp.Body)
	//if err != nil {
	//	return nil, fmt.Errorf("error reading response body for %s: %w", sys.Name, err)
	//}

	var text string
	node := findNode(doc, config)

	if node != nil {
		text = extractText(node)
	} else {
		text = extractText(doc) // no selector found, so read the whole page
	}

	content := strings.ToLower(text)
	badKeywords := []string{"maintenance", "outage", "down", "degraded", "unavailable", "failure"}
	isOngoing := false

	for _, kw := range badKeywords {
		if strings.Contains(content, kw) {
			isOngoing = true
			break
		}
	}
	statusText := strings.TrimSpace(text)
	if len(statusText) > 100 {
		statusText = statusText[:100] + "..."
	}
	if statusText == "" {
		statusText = "No status information available"
	}

	return statusText, isOngoing, nil
}

func findNode(n *html.Node, config map[string]string) *html.Node {
	id := config["id"]
	class := config["class"]
	tag := config["tag"]

	// if no config is provided, still go head and search everything. though may not get the desired output
	if id == "" && class == "" && tag == "" {
		return n
	}

	if n.Type == html.ElementNode {

		matches := true

		//check tag
		if tag != "" && n.Data != tag {
			matches = false
		}

		//check id
		if matches && id != "" {
			idVal := getAttrValue(n, "id")
			if idVal != id {
				matches = false
			}
		}

		// check class
		if matches && class != "" {
			classVal := getAttrValue(n, "class")
			if !strings.Contains(classVal, class) {
				matches = false
			}
		}

		if matches {
			return n
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		res := findNode(c, config)
		if res != nil {
			return res
		}
	}

	return nil
}

func getAttrValue(n *html.Node, attr string) string {
	for _, a := range n.Attr {
		if a.Key == attr {
			return a.Val
		}
	}
	return ""
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
