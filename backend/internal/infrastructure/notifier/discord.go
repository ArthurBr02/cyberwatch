package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// DiscordNotifier implements the notifier interface for Discord
type DiscordNotifier struct {
	token     string
	channelID string
	client    *http.Client
}

// NewDiscordNotifier creates a new DiscordNotifier instance
func NewDiscordNotifier(token, channelID string) *DiscordNotifier {
	return &DiscordNotifier{
		token:     token,
		channelID: channelID,
		client:    &http.Client{Timeout: 10 * time.Second},
	}
}

// Notify is a convenience method to send a formatted article notification
func (n *DiscordNotifier) Notify(sourceName, title, url string) error {
	if n.channelID == "" {
		return fmt.Errorf("discord channel ID not configured")
	}
	message := fmt.Sprintf("**[%s]** New article detected:\n**%s**\n%s", sourceName, title, url)
	return n.Send(n.channelID, message)
}

// Send sends a message to a Discord channel
func (n *DiscordNotifier) Send(channelID string, message string) error {
	// Prepare the request payload
	payload := map[string]interface{}{
		"content": message,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request payload: %w", err)
	}

	// Create the HTTP request
	url := fmt.Sprintf("https://discord.com/api/v10/channels/%s/messages", channelID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bot %s", n.token))

	// Send the request
	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request to Discord: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("Discord API returned error status: %d", resp.StatusCode)
	}

	return nil
}