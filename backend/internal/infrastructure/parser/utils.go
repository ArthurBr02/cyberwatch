package parser

import "strings"

// cleanHTML removes script and style tags to make the structure clearer for the LLM
func cleanHTML(html string) string {
	s := html
	// Remove script tags
	for {
		start := strings.Index(strings.ToLower(s), "<script")
		if start == -1 { break }
		end := strings.Index(strings.ToLower(s[start:]), "</script>")
		if end == -1 { break }
		s = s[:start] + s[start+end+9:]
	}
	// Remove style tags
	for {
		start := strings.Index(strings.ToLower(s), "<style")
		if start == -1 { break }
		end := strings.Index(strings.ToLower(s[start:]), "</style>")
		if end == -1 { break }
		s = s[:start] + s[start+end+8:]
	}
	return s
}
