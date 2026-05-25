package router

import (
    "github.com/cyberwatch/backend/internal/interface/handler"
    "github.com/cyberwatch/backend/internal/interface/middleware"
    "github.com/go-chi/chi/v5"
)

func Setup(sourceHandler *handler.SourceHandler, articleHandler *handler.ArticleHandler) *chi.Mux {
    r := chi.NewRouter()
    
    // Global middleware
    r.Use(middleware.CORS)
    
    // Health check
    r.Get("/api/v1/health", articleHandler.Health)
    
    // Sources
    r.Route("/api/v1/sources", func(r chi.Router) {
        r.Get("/", sourceHandler.List)          // GET /api/v1/sources
        r.Post("/", sourceHandler.Create)        // POST /api/v1/sources
        r.Get("/{id}", sourceHandler.Get)        // GET /api/v1/sources/{id}
        r.Put("/{id}", sourceHandler.Update)     // PUT /api/v1/sources/{id}
        r.Delete("/{id}", sourceHandler.Delete)  // DELETE /api/v1/sources/{id}
        r.Post("/{id}/generate-rules", sourceHandler.GenerateRules) // POST /api/v1/sources/{id}/generate-rules
        r.Post("/{id}/test", sourceHandler.TestScrape) // POST /api/v1/sources/{id}/test
    })
    
    // Scrape
    r.Post("/api/v1/scrape/all", articleHandler.ScrapeAll) // POST /api/v1/scrape/all
    
    // Articles
    r.Get("/api/v1/articles", articleHandler.List) // GET /api/v1/articles
    
    return r
}