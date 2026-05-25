package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cyberwatch/backend/internal/application/usecase"
	"github.com/cyberwatch/backend/internal/domain/entity"
)

type ArticleHandler struct {
    scrapeUC *usecase.ScrapeUseCase
}

func NewArticleHandler(scrapeUC *usecase.ScrapeUseCase) *ArticleHandler {
    return &ArticleHandler{
        scrapeUC: scrapeUC,
    }
}

// List returns articles, optionally filtered by source_id
func (h *ArticleHandler) List(w http.ResponseWriter, r *http.Request) {
    sourceIDStr := r.URL.Query().Get("source_id")
    var articles []*entity.Article
    var err error

    if sourceIDStr != "" {
        sourceID, err := strconv.ParseInt(sourceIDStr, 10, 64)
        if err != nil {
            http.Error(w, `{"error": "Invalid source_id"}`, http.StatusBadRequest)
            return
        }
        articles, err = h.scrapeUC.GetArticles(sourceID)
    } else {
        articles, err = h.scrapeUC.GetArticles(0)
    }

    if err != nil {
        http.Error(w, `{"error": "Failed to fetch articles"}`, http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(articles)
}

// ScrapeAll scrapes all sources
func (h *ArticleHandler) ScrapeAll(w http.ResponseWriter, r *http.Request) {
    articles, err := h.scrapeUC.ScrapeAll()
    if err != nil {
        http.Error(w, `{"error": "Failed to scrape all sources"}`, http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(articles)
}

// Health checks the health of the service
func (h *ArticleHandler) Health(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}