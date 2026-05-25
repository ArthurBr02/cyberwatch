package database

import (
	"database/sql"
	"fmt"

	"github.com/cyberwatch/backend/internal/domain/entity"
)

// SourceRepository implements repository.SourceRepository interface
type SourceRepository struct {
	db *sql.DB
}

// NewSourceRepository creates a new SourceRepository instance
func NewSourceRepository(db *sql.DB) *SourceRepository {
	return &SourceRepository{
		db: db,
	}
}

// Create creates a new source in the database
func (r *SourceRepository) Create(source *entity.Source) error {
	query := `
		INSERT INTO sources (name, url, fetch_type, is_active, llm_rules)
		VALUES (?, ?, ?, ?, ?)`

	result, err := r.db.Exec(query, source.Name, source.URL, source.FetchType, source.IsActive, source.LLMRules)
	if err != nil {
		return fmt.Errorf("failed to create source: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %w", err)
	}
	source.ID = id

	return nil
}

// GetByID retrieves a source by its ID
func (r *SourceRepository) GetByID(id int64) (*entity.Source, error) {
	query := `SELECT id, name, url, fetch_type, is_active, llm_rules, created_at, updated_at FROM sources WHERE id = ?`

	row := r.db.QueryRow(query, id)
	source := &entity.Source{}
	err := row.Scan(&source.ID, &source.Name, &source.URL, &source.FetchType, &source.IsActive, &source.LLMRules, &source.CreatedAt, &source.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("source with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get source by ID: %w", err)
	}

	return source, nil
}

// GetAll retrieves all sources
func (r *SourceRepository) GetAll() ([]*entity.Source, error) {
	query := `SELECT id, name, url, fetch_type, is_active, llm_rules, created_at, updated_at FROM sources`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all sources: %w", err)
	}
	defer rows.Close()

	sources := []*entity.Source{}
	for rows.Next() {
		source := &entity.Source{}
		err := rows.Scan(&source.ID, &source.Name, &source.URL, &source.FetchType, &source.IsActive, &source.LLMRules, &source.CreatedAt, &source.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan source row: %w", err)
		}
		sources = append(sources, source)
	}

	return sources, nil
}

// Update updates an existing source
func (r *SourceRepository) Update(source *entity.Source) error {
	query := `
		UPDATE sources
		SET name = ?, url = ?, fetch_type = ?, is_active = ?, llm_rules = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`

	_, err := r.db.Exec(query, source.Name, source.URL, source.FetchType, source.IsActive, source.LLMRules, source.ID)
	if err != nil {
		return fmt.Errorf("failed to update source: %w", err)
	}

	return nil
}

// Delete deletes a source by its ID
func (r *SourceRepository) Delete(id int64) error {
	query := `DELETE FROM sources WHERE id = ?`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete source: %w", err)
	}

	return nil
}

// GetActive retrieves all active sources
func (r *SourceRepository) GetActive() ([]*entity.Source, error) {
	query := `SELECT id, name, url, fetch_type, is_active, llm_rules, created_at, updated_at FROM sources WHERE is_active = 1`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get active sources: %w", err)
	}
	defer rows.Close()

	sources := []*entity.Source{}
	for rows.Next() {
		source := &entity.Source{}
		err := rows.Scan(&source.ID, &source.Name, &source.URL, &source.FetchType, &source.IsActive, &source.LLMRules, &source.CreatedAt, &source.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan source row: %w", err)
		}
		sources = append(sources, source)
	}

	return sources, nil
}