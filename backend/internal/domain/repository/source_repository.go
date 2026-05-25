package repository

import (
	"github.com/cyberwatch/backend/internal/domain/entity"
)

// SourceRepository defines the interface for source persistence operations
type SourceRepository interface {
	// Create saves a new source to the repository
	Create(source *entity.Source) error

	// GetByID retrieves a source by its ID
	GetByID(id int64) (*entity.Source, error)

	// GetAll retrieves all sources
	GetAll() ([]*entity.Source, error)

	// Update updates an existing source
	Update(source *entity.Source) error

	// Delete removes a source by its ID
	Delete(id int64) error

	// GetActive retrieves all active sources
	GetActive() ([]*entity.Source, error)
}