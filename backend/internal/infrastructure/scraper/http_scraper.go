package scraper

import (
	"fmt"
	"io"
	"net/http"
)

// HTTPScraper implements the scraper interface for HTTP requests
type HTTPScraper struct {
	client *http.Client
}

// NewHTTPScraper creates a new HTTPScraper instance
func NewHTTPScraper() *HTTPScraper {
	return &HTTPScraper{
		client: &http.Client{},
	}
}

// Fetch makes an HTTP GET request to the specified URL and returns the body as string
func (s *HTTPScraper) Fetch(url string, fetchType string) (string, error) {
	resp, err := s.client.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch URL %s: %w", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body from %s: %w", url, err)
	}

	return string(body), nil
}