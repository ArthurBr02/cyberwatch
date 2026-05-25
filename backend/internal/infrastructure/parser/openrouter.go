package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// OpenRouterParser implements LLMParser using OpenRouter API
type OpenRouterParser struct {
	apiKey string
	model  string
	client *http.Client
}

// NewOpenRouterParser creates a new OpenRouterParser instance
func NewOpenRouterParser(apiKey, model string) *OpenRouterParser {
	if model == "" {
		model = "mistralai/mistral-7b-instruct"
	}
	return &OpenRouterParser{
		apiKey: apiKey,
		model:  model,
		client: &http.Client{
			Timeout: 60 * time.Second, // Increased timeout for larger content
		},
	}
}

// GenerateRules sends raw content to OpenRouter LLM and returns JSON rules
func (p *OpenRouterParser) GenerateRules(rawContent string, fetchType string) (string, error) {
	// Clean HTML to save tokens and focus on structure
	analysisContent := rawContent
	if fetchType == "html" {
		analysisContent = cleanHTML(rawContent)
	}

	// Limit content to 50,000 characters (safe for LLM context, covers most articles)
	if len(analysisContent) > 50000 {
		analysisContent = analysisContent[:50000]
	}

	// Prepare the request payload
	payload := map[string]interface{}{
		"model": p.model,
		"messages": []map[string]string{
			{
				"role": "user",
				"content": fmt.Sprintf(`You are a specialized web scraping expert. Your task is to analyze the provided %s content and extract the EXACT structural rules needed to parse latest articles.

### MANDATORY OUTPUT FORMAT
Return ONLY a valid JSON object. No markdown blocks, no preamble, no explanation.
{
  "container": "CSS selector or JSON path for the repeating article element",
  "title": "CSS selector or JSON path for the title (relative to container)",
  "url": "CSS selector or JSON path for the link. Use 'self' if the container itself is the link",
  "summary": "CSS selector or JSON path for the description (relative to container)"
}

### EXTRACTION RULES
1. If the container element is an <a> tag, use "self" for the "url" field.
2. Be precise. Use classes or IDs that actually exist in the content.
3. For JSON, use standard GJSON path syntax.
4. Ensure the selectors are unique enough to only target articles.

### CONTENT TO ANALYZE
%s`, fetchType, analysisContent),
			},
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request payload: %w", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.apiKey))

	// Send the request
	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to OpenRouter: %w", err)
	}
	defer resp.Body.Close()

	// Read the response
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse the response to extract the JSON rules
	var result map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	// Extract the content from the response
	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("failed to extract choices from response")
	}

	choice, ok := choices[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("failed to extract choice[0] from response")
	}

	message, ok := choice["message"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("failed to extract message from choice")
	}

	content, ok := message["content"].(string)
	if !ok {
		return "", fmt.Errorf("failed to extract content from message")
	}

	return content, nil
}

