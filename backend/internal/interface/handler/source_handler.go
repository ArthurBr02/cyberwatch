package handler

import (
    "encoding/json"
    "log"
    "net/http"
    "strconv"

    "github.com/cyberwatch/backend/internal/application/usecase"
    "github.com/cyberwatch/backend/internal/domain/entity"
    "github.com/go-chi/chi/v5"
)

type SourceHandler struct {
    sourceUC *usecase.SourceUseCase
    scrapeUC *usecase.ScrapeUseCase
}

func NewSourceHandler(sourceUC *usecase.SourceUseCase, scrapeUC *usecase.ScrapeUseCase) *SourceHandler {
    return &SourceHandler{
        sourceUC: sourceUC,
        scrapeUC: scrapeUC,
    }
}

// List returns all sources
func (h *SourceHandler) List(w http.ResponseWriter, r *http.Request) {
    sources, err := h.sourceUC.GetAllSources()
    if err != nil {
        log.Printf("Error fetching sources: %v", err)
        http.Error(w, `{"error": "Failed to fetch sources"}`, http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(sources)
}

// Create creates a new source
func (h *SourceHandler) Create(w http.ResponseWriter, r *http.Request) {
    var source entity.Source
    if err := json.NewDecoder(r.Body).Decode(&source); err != nil {
        http.Error(w, `{"error": "Invalid JSON"}`, http.StatusBadRequest)
        return
    }

    createdSource, err := h.sourceUC.CreateSource(source.Name, source.URL, source.FetchType)
    if err != nil {
        http.Error(w, `{"error": "Failed to create source"}`, http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(createdSource)
}

// Get returns a source by ID
func (h *SourceHandler) Get(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        http.Error(w, `{"error": "Invalid ID"}`, http.StatusBadRequest)
        return
    }

    source, err := h.sourceUC.GetSource(id)
    if err != nil {
        http.Error(w, `{"error": "Source not found"}`, http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(source)
}

// Update updates a source
func (h *SourceHandler) Update(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        http.Error(w, `{"error": "Invalid ID"}`, http.StatusBadRequest)
        return
    }

    var source entity.Source
    if err := json.NewDecoder(r.Body).Decode(&source); err != nil {
        http.Error(w, `{"error": "Invalid JSON"}`, http.StatusBadRequest)
        return
    }
    source.ID = id

    err = h.sourceUC.UpdateSource(&source)
    if err != nil {
        http.Error(w, `{"error": "Failed to update source"}`, http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(source)
}

// Delete deletes a source
func (h *SourceHandler) Delete(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        http.Error(w, `{"error": "Invalid ID"}`, http.StatusBadRequest)
        return
    }

    err = h.sourceUC.DeleteSource(id)
    if err != nil {
        http.Error(w, `{"error": "Failed to delete source"}`, http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

// GenerateRules generates rules for a source
func (h *SourceHandler) GenerateRules(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        http.Error(w, `{"error": "Invalid ID"}`, http.StatusBadRequest)
        return
    }

    updatedSource, err := h.sourceUC.GenerateRules(id)
    if err != nil {
        http.Error(w, `{"error": "Failed to generate rules"}`, http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(updatedSource)
}

// TestScrape tests scraping for a source
func (h *SourceHandler) TestScrape(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        http.Error(w, `{"error": "Invalid ID"}`, http.StatusBadRequest)
        return
    }

    articles, err := h.scrapeUC.ScrapeSource(id)
    if err != nil {
        log.Printf("Error testing scrape for source %d: %v", id, err)
        http.Error(w, `{"error": "Failed to scrape source"}`, http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(articles)
}