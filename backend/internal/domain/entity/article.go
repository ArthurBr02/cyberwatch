package entity

import "time"

// Article represents an article fetched from a source
type Article struct {
	ID            int64     `json:"id"`
	SourceID      int64     `json:"source_id"`
	Title         string    `json:"title"`
	URL           string    `json:"url"`
	Summary       string    `json:"summary"`
	PublishedAt   *time.Time `json:"published_at,omitempty"`
	IsSent        bool      `json:"is_sent"`
	UpdatedAt     time.Time `json:"updated_at"`
	CreatedAt     time.Time `json:"created_at"`
}