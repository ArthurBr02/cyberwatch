package repository

import (
	"github.com/cyberwatch/backend/internal/domain/entity"
)

// ArticleRepository defines the interface for article persistence operations
type ArticleRepository interface {
	// Create saves a new article to the repository
	Create(article *entity.Article) error

	// ExistsByURL checks if an article with the given URL already exists
	ExistsByURL(url string) (bool, error)

	// GetBySourceID retrieves all articles for a given source ID
	GetBySourceID(sourceID int64) ([]*entity.Article, error)

	// GetAll retrieves all articles
	GetAll() ([]*entity.Article, error)

	// MarkAsSent marks an article as sent to Discord
	MarkAsSent(id int64) error
}