package config

import (
	"os"
)

// Config holds the application configuration
type Config struct {
	DiscordBotToken   string
	DiscordChannelID  string
	LLMProvider       string
	OpenRouterAPIKey  string
	OpenRouterModel   string
	MistralAPIKey     string
	MistralModel      string
	DatabasePath      string
	ServerPort        string
	ScrapeInterval    string
}

// Load reads environment variables and returns a Config struct
func Load() *Config {
	return &Config{
		DiscordBotToken:   os.Getenv("DISCORD_BOT_TOKEN"),
		DiscordChannelID:  os.Getenv("DISCORD_CHANNEL_ID"),
		LLMProvider:       getEnvOrDefault("LLM_PROVIDER", "openrouter"),
		OpenRouterAPIKey:  os.Getenv("OPENROUTER_API_KEY"),
		OpenRouterModel:   getEnvOrDefault("OPENROUTER_MODEL", "mistralai/mistral-7b-instruct"),
		MistralAPIKey:     os.Getenv("MISTRAL_API_KEY"),
		MistralModel:      getEnvOrDefault("MISTRAL_MODEL", "mistral-tiny"),
		DatabasePath:      getEnvOrDefault("DATABASE_PATH", "./data/cyberwatch.db"),
		ServerPort:        getEnvOrDefault("SERVER_PORT", "8080"),
		ScrapeInterval:    getEnvOrDefault("SCRAPE_INTERVAL", "1h"),
	}
}

// getEnvOrDefault returns the value of an environment variable or a default value if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}