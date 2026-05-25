package entity

import "time"

// Source represents a data source for articles
type Source struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	URL         string    `json:"url"`
	FetchType   string    `json:"fetch_type"` // "html" or "json"
	LLMRules    string    `json:"llm_rules"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}