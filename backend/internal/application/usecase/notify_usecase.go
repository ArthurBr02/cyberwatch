package usecase

import (
	"fmt"

	"github.com/cyberwatch/backend/internal/domain/entity"
	"github.com/cyberwatch/backend/internal/domain/repository"
	"github.com/cyberwatch/backend/internal/infrastructure/notifier"
)

// NotifyUseCase handles notification-related use cases
type NotifyUseCase struct {
	articleRepo repository.ArticleRepository
	notifier    *notifier.DiscordNotifier
	channelID   string
}

// NewNotifyUseCase creates a new NotifyUseCase instance
func NewNotifyUseCase(articleRepo repository.ArticleRepository, notifier *notifier.DiscordNotifier, channelID string) *NotifyUseCase {
	return &NotifyUseCase{
		articleRepo: articleRepo,
		notifier:    notifier,
		channelID:   channelID,
	}
}

// SendNewArticles sends notifications for new articles
func (uc *NotifyUseCase) SendNewArticles(articles []*entity.Article) error {
	var firstErr error

	for _, article := range articles {
		message := fmt.Sprintf("**%s**\n%s\n%s", article.Title, article.URL, article.Summary)
		
		err := uc.notifier.Send(uc.channelID, message)
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue // Continue with other articles even if one fails
		}

		// Mark article as sent
		err = uc.articleRepo.MarkAsSent(article.ID)
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
		}
	}

	return firstErr
}