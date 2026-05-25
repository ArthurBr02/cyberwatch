package usecase

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/cyberwatch/backend/internal/domain/entity"
	"github.com/cyberwatch/backend/internal/domain/repository"
	"github.com/cyberwatch/backend/internal/infrastructure/notifier"
	"github.com/cyberwatch/backend/internal/infrastructure/parser"
	"github.com/cyberwatch/backend/internal/infrastructure/scraper"
	"github.com/tidwall/gjson"
)

// ScrapeUseCase handles scraping-related use cases
type ScrapeUseCase struct {
	sourceRepo  repository.SourceRepository
	articleRepo repository.ArticleRepository
	scraper     *scraper.HTTPScraper
	llmParser   parser.LLMParser
	notifier    *notifier.DiscordNotifier
}

// NewScrapeUseCase creates a new ScrapeUseCase instance
func NewScrapeUseCase(sourceRepo repository.SourceRepository, articleRepo repository.ArticleRepository, scraper *scraper.HTTPScraper, llmParser parser.LLMParser, notifier *notifier.DiscordNotifier) *ScrapeUseCase {
	return &ScrapeUseCase{
		sourceRepo:  sourceRepo,
		articleRepo: articleRepo,
		scraper:     scraper,
		llmParser:   llmParser,
		notifier:    notifier,
	}
}

// ScrapeAll scrapes all active sources
func (uc *ScrapeUseCase) ScrapeAll() ([]*entity.Article, error) {
	sources, err := uc.sourceRepo.GetActive()
	if err != nil {
		return nil, err
	}

	allNewArticles := []*entity.Article{}

	for _, source := range sources {
		newArticles, err := uc.ScrapeSource(source.ID)
		if err != nil {
			return nil, err
		}
		allNewArticles = append(allNewArticles, newArticles...)
	}

	return allNewArticles, nil
}

// ScrapeSource scrapes a single source
func (uc *ScrapeUseCase) ScrapeSource(sourceID int64) ([]*entity.Article, error) {
	source, err := uc.sourceRepo.GetByID(sourceID)
	if err != nil {
		return nil, err
	}

	rawContent, err := uc.scraper.Fetch(source.URL, source.FetchType)
	if err != nil {
		return nil, err
	}

	articles, err := parseArticles(rawContent, source.LLMRules, source.FetchType, sourceID, source.URL)
	if err != nil {
		return nil, err
	}

	newArticles := []*entity.Article{}
	for _, article := range articles {
		exists, err := uc.articleRepo.ExistsByURL(article.URL)
		if err != nil {
			return nil, err
		}

		if !exists {
			err = uc.articleRepo.Create(article)
			if err != nil {
				return nil, err
			}
			newArticles = append(newArticles, article)
			
			// Notify about new article
			if uc.notifier != nil {
				err = uc.notifier.Notify(source.Name, article.Title, article.URL)
				if err == nil {
					// Mark as sent in database if notification was successful
					_ = uc.articleRepo.MarkAsSent(article.ID)
					article.IsSent = true // Update local object
				}
			}
		}
	}

	return newArticles, nil
}

// GetArticles returns articles, optionally filtered by source_id
func (uc *ScrapeUseCase) GetArticles(sourceID int64) ([]*entity.Article, error) {
	if sourceID == 0 {
		return uc.articleRepo.GetAll()
	}
	return uc.articleRepo.GetBySourceID(sourceID)
}

// parseArticles uses LLM rules to extract articles from raw content
func parseArticles(rawContent, llmRules, fetchType string, sourceID int64, baseURL string) ([]*entity.Article, error) {
	if llmRules == "" {
		return []*entity.Article{}, nil
	}

	// Clean the LLM output in case it includes markdown blocks
	cleanRules := llmRules
	if strings.Contains(cleanRules, "```") {
		parts := strings.Split(cleanRules, "```")
		if len(parts) >= 3 {
			content := parts[1]
			lines := strings.Split(content, "\n")
			if len(lines) > 0 {
				firstLine := strings.ToLower(strings.TrimSpace(lines[0]))
				if firstLine == "json" || firstLine == "javascript" || firstLine == "js" {
					content = strings.Join(lines[1:], "\n")
				}
			}
			cleanRules = content
		}
	}
	cleanRules = strings.TrimSpace(cleanRules)

	var rules map[string]string
	if err := json.Unmarshal([]byte(cleanRules), &rules); err != nil {
		var ruleList []interface{}
		if errArray := json.Unmarshal([]byte(cleanRules), &ruleList); errArray == nil {
			return nil, fmt.Errorf("LLM returned a JSON array instead of an object. Please regenerate rules. (rules: %s)", cleanRules)
		}
		return nil, fmt.Errorf("failed to parse LLM rules as JSON object: %w (clean_rules: %s)", err, cleanRules)
	}

	container := rules["container"]
	titleSel := rules["title"]
	urlSel := rules["url"]
	summarySel := rules["summary"]

	articles := []*entity.Article{}

	// Parse base URL for normalization
	base, err := url.Parse(baseURL)
	if err != nil {
		log.Printf("[ERROR] Failed to parse base URL %s: %v", baseURL, err)
	}

	if fetchType == "html" {
		log.Printf("[DEBUG] Parsing HTML with container: %s, title: %s, url: %s", container, titleSel, urlSel)
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(rawContent))
		if err != nil {
			return nil, fmt.Errorf("failed to parse HTML: %w", err)
		}

		selection := doc.Find(container)
		log.Printf("[DEBUG] Found %d container elements", selection.Length())

		selection.Each(func(i int, s *goquery.Selection) {
			title := strings.TrimSpace(s.Find(titleSel).Text())
			
			var urlStr string
			var exists bool
			if urlSel == "" || urlSel == "self" || urlSel == "." {
				urlStr, exists = s.Attr("href")
				if !exists || urlStr == "" {
					firstLink := s.Find("a").First()
					urlStr, exists = firstLink.Attr("href")
				}
			} else {
				urlStr, exists = s.Find(urlSel).Attr("href")
				if (!exists || urlStr == "") && s.Find(urlSel).Is("a") {
					urlStr, exists = s.Find(urlSel).Attr("href")
				}
			}
			
			summary := strings.TrimSpace(s.Find(summarySel).Text())

			if title != "" && exists && urlStr != "" {
				// Normalize URL (handle relative links correctly)
				absoluteURL := urlStr
				if base != nil {
					rel, err := url.Parse(urlStr)
					if err == nil {
						absoluteURL = base.ResolveReference(rel).String()
						if absoluteURL != urlStr {
							log.Printf("[DEBUG] Normalized URL: %s -> %s", urlStr, absoluteURL)
						}
					}
				}
				
				now := time.Now()
				articles = append(articles, &entity.Article{
					SourceID:    sourceID,
					Title:       title,
					URL:         absoluteURL,
					Summary:     summary,
					PublishedAt: &now, // Default to fetch time
				})
			}
		})
	} else if fetchType == "json" {
		res := gjson.Parse(rawContent)
		var items gjson.Result
		if container == "" || container == "." || container == "root" {
			items = res
		} else {
			items = res.Get(container)
		}

		if items.IsArray() {
			items.ForEach(func(key, value gjson.Result) bool {
				title := value.Get(titleSel).String()
				urlStr := value.Get(urlSel).String()
				summary := value.Get(summarySel).String()

				if title != "" && urlStr != "" {
					// Normalize URL for JSON as well
					absoluteURL := urlStr
					if base != nil {
						rel, err := url.Parse(urlStr)
						if err == nil {
							absoluteURL = base.ResolveReference(rel).String()
						}
					}

					now := time.Now()
					articles = append(articles, &entity.Article{
						SourceID:    sourceID,
						Title:       title,
						URL:         absoluteURL,
						Summary:     summary,
						PublishedAt: &now, // Default to fetch time
					})
				}
				return true
			})
		} else if items.IsObject() {
			title := items.Get(titleSel).String()
			urlStr := items.Get(urlSel).String()
			summary := items.Get(summarySel).String()

			if title != "" && urlStr != "" {
				// Normalize URL for JSON as well
				absoluteURL := urlStr
				if base != nil {
					rel, err := url.Parse(urlStr)
					if err == nil {
						absoluteURL = base.ResolveReference(rel).String()
					}
				}

				now := time.Now()
				articles = append(articles, &entity.Article{
					SourceID:    sourceID,
					Title:       title,
					URL:         absoluteURL,
					Summary:     summary,
					PublishedAt: &now, // Default to fetch time
				})
			}
		}
	}

	log.Printf("[DEBUG] Total articles found: %d", len(articles))
	return articles, nil
}
