package parser

// LLMParser defines the interface for parsing raw content into structured rules
type LLMParser interface {
	// GenerateRules takes raw content and returns JSON string of parsing rules
	GenerateRules(rawContent string, fetchType string) (string, error)
}