package usecase

import (
	"github.com/cyberwatch/backend/internal/domain/entity"
	"github.com/cyberwatch/backend/internal/domain/repository"
	"github.com/cyberwatch/backend/internal/infrastructure/parser"
	"github.com/cyberwatch/backend/internal/infrastructure/scraper"
)

// SourceUseCase handles source-related use cases
type SourceUseCase struct {
	sourceRepo repository.SourceRepository
	scraper    *scraper.HTTPScraper
	llmParser  parser.LLMParser
}

// NewSourceUseCase creates a new SourceUseCase instance
func NewSourceUseCase(sourceRepo repository.SourceRepository, scraper *scraper.HTTPScraper, llmParser parser.LLMParser) *SourceUseCase {
	return &SourceUseCase{
		sourceRepo: sourceRepo,
		scraper:    scraper,
		llmParser:  llmParser,
	}
}

// CreateSource creates a new source
func (uc *SourceUseCase) CreateSource(name, url, fetchType string) (*entity.Source, error) {
	source := &entity.Source{
		Name:      name,
		URL:       url,
		FetchType: fetchType,
		IsActive:  true, // Default to active
	}

	err := uc.sourceRepo.Create(source)
	if err != nil {
		return nil, err
	}

	return source, nil
}

// GetSource retrieves a source by ID
func (uc *SourceUseCase) GetSource(id int64) (*entity.Source, error) {
	return uc.sourceRepo.GetByID(id)
}

// GetAllSources retrieves all sources
func (uc *SourceUseCase) GetAllSources() ([]*entity.Source, error) {
	return uc.sourceRepo.GetAll()
}

// UpdateSource updates an existing source
func (uc *SourceUseCase) UpdateSource(source *entity.Source) error {
	return uc.sourceRepo.Update(source)
}

// DeleteSource deletes a source by ID
func (uc *SourceUseCase) DeleteSource(id int64) error {
	return uc.sourceRepo.Delete(id)
}

// GenerateRules generates LLM rules for a source
func (uc *SourceUseCase) GenerateRules(id int64) (*entity.Source, error) {
	source, err := uc.GetSource(id)
	if err != nil {
		return nil, err
	}

	rawContent, err := uc.scraper.Fetch(source.URL, source.FetchType)
	if err != nil {
		return nil, err
	}

	rules, err := uc.llmParser.GenerateRules(rawContent, source.FetchType)
	if err != nil {
		return nil, err
	}

	source.LLMRules = rules

	err = uc.sourceRepo.Update(source)
	if err != nil {
		return nil, err
	}

	return source, nil
}