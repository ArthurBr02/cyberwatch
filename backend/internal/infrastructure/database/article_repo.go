package database

import (
	"database/sql"
	"fmt"

	"github.com/cyberwatch/backend/internal/domain/entity"
)

// ArticleRepository implements repository.ArticleRepository interface
type ArticleRepository struct {
	db *sql.DB
}

// NewArticleRepository creates a new ArticleRepository instance
func NewArticleRepository(db *sql.DB) *ArticleRepository {
	return &ArticleRepository{
		db: db,
	}
}

// Create creates a new article in the database
func (r *ArticleRepository) Create(article *entity.Article) error {
	query := `
		INSERT INTO articles (source_id, title, url, summary, published_at, sent_to_discord)
		VALUES (?, ?, ?, ?, ?, ?)`

	res, err := r.db.Exec(query, article.SourceID, article.Title, article.URL, article.Summary, article.PublishedAt, article.IsSent)
	if err != nil {
		return fmt.Errorf("failed to create article: %w", err)
	}

	id, err := res.LastInsertId()
	if err == nil {
		article.ID = id
	}

	return nil
}

// ExistsByURL checks if an article with the given URL already exists
func (r *ArticleRepository) ExistsByURL(url string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM articles WHERE url = ?)`

	var exists bool
	err := r.db.QueryRow(query, url).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if article exists: %w", err)
	}

	return exists, nil
}

// GetBySourceID retrieves all articles for a given source ID
func (r *ArticleRepository) GetBySourceID(sourceID int64) ([]*entity.Article, error) {
	query := `SELECT id, source_id, title, url, summary, published_at, sent_to_discord, created_at, updated_at FROM articles WHERE source_id = ?`

	rows, err := r.db.Query(query, sourceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get articles by source ID: %w", err)
	}
	defer rows.Close()

	articles := []*entity.Article{}
	for rows.Next() {
		article := &entity.Article{}
		err := rows.Scan(&article.ID, &article.SourceID, &article.Title, &article.URL, &article.Summary, &article.PublishedAt, &article.IsSent, &article.CreatedAt, &article.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan article row: %w", err)
		}
		articles = append(articles, article)
	}

	return articles, nil
}

// GetAll retrieves all articles
func (r *ArticleRepository) GetAll() ([]*entity.Article, error) {
	query := `SELECT id, source_id, title, url, summary, published_at, sent_to_discord, created_at, updated_at FROM articles`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all articles: %w", err)
	}
	defer rows.Close()

	articles := []*entity.Article{}
	for rows.Next() {
		article := &entity.Article{}
		err := rows.Scan(&article.ID, &article.SourceID, &article.Title, &article.URL, &article.Summary, &article.PublishedAt, &article.IsSent, &article.CreatedAt, &article.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan article row: %w", err)
		}
		articles = append(articles, article)
	}

	return articles, nil
}

// MarkAsSent marks an article as sent
func (r *ArticleRepository) MarkAsSent(id int64) error {
	query := `UPDATE articles SET sent_to_discord = 1 WHERE id = ?`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to mark article as sent: %w", err)
	}

	return nil
}