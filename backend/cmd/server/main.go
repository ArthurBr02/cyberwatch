package main

import (
	"log"
	"net/http"
	"time"

	"github.com/cyberwatch/backend/internal/application/usecase"
	"github.com/cyberwatch/backend/internal/infrastructure/config"
	"github.com/cyberwatch/backend/internal/infrastructure/database"
	"github.com/cyberwatch/backend/internal/infrastructure/notifier"
	"github.com/cyberwatch/backend/internal/infrastructure/parser"
	"github.com/cyberwatch/backend/internal/infrastructure/scraper"
	"github.com/cyberwatch/backend/internal/interface/handler"
	"github.com/cyberwatch/backend/internal/interface/router"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Open database connection
	db, err := database.NewSQLite(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create repositories
	sourceRepo := database.NewSourceRepository(db)
	articleRepo := database.NewArticleRepository(db)

	// Create scraper
	scraper := scraper.NewHTTPScraper()

	// Create parser
	var llmParser parser.LLMParser
	switch cfg.LLMProvider {
	case "mistral":
		if cfg.MistralAPIKey == "" {
			log.Fatal("Mistral API key is required when provider is mistral")
		}
		llmParser = parser.NewMistralParser(cfg.MistralAPIKey, cfg.MistralModel)
		log.Printf("Using Mistral provider with model: %s", cfg.MistralModel)
	case "openrouter":
		fallthrough
	default:
		if cfg.OpenRouterAPIKey == "" {
			log.Fatal("OpenRouter API key is required when provider is openrouter")
		}
		llmParser = parser.NewOpenRouterParser(cfg.OpenRouterAPIKey, cfg.OpenRouterModel)
		log.Printf("Using OpenRouter provider with model: %s", cfg.OpenRouterModel)
	}

	// Create notifier
	if cfg.DiscordBotToken == "" {
		log.Fatal("Discord token is required")
	}
	notifier := notifier.NewDiscordNotifier(cfg.DiscordBotToken, cfg.DiscordChannelID)

	// Create use cases
	sourceUC := usecase.NewSourceUseCase(sourceRepo, scraper, llmParser)
	scrapeUC := usecase.NewScrapeUseCase(sourceRepo, articleRepo, scraper, llmParser, notifier)

	// Create handlers
	sourceHandler := handler.NewSourceHandler(sourceUC, scrapeUC)
	articleHandler := handler.NewArticleHandler(scrapeUC)

	// Setup router
	r := router.Setup(sourceHandler, articleHandler)

	// Setup background scheduler
	go func() {
		interval, err := time.ParseDuration(cfg.ScrapeInterval)
		if err != nil {
			log.Printf("Invalid scrape interval %s, defaulting to 1 hour: %v", cfg.ScrapeInterval, err)
			interval = time.Hour
		}
		
		log.Printf("Background scheduler started with interval: %s", interval)
		ticker := time.NewTicker(interval)
		for range ticker.C {
			log.Println("Auto-scrape triggered...")
			articles, err := scrapeUC.ScrapeAll()
			if err != nil {
				log.Printf("Auto-scrape failed: %v", err)
			} else {
				log.Printf("Auto-scrape completed: found %d new articles", len(articles))
			}
		}
	}()

	// Start server
	log.Printf("Server starting on port %s", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, r))
}