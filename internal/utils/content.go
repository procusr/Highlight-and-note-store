package utils

import (
	"regexp"
	"strings"
)

// DetectContentType automatically detects the content type based on the content
func DetectContentType(content string) string {
	content = strings.TrimSpace(content)

	// Check for URL patterns (simple detection)
	urlPattern := regexp.MustCompile(`^https?:\/\/[^\s]+$`)
	if urlPattern.MatchString(content) {
		return "link"
	}

	// Check for common markdown patterns
	markdownPatterns := []*regexp.Regexp{
		regexp.MustCompile(`^#{1,6}\s+`),                    // Headers
		regexp.MustCompile(`\*\*[^*]+\*\*`),                 // Bold text
		regexp.MustCompile(`\*[^*]+\*`),                     // Italic text
		regexp.MustCompile(`\[[^\]]+\]\([^)]+\)`),           // Links
		regexp.MustCompile("`[^`]+`"),                       // Inline code
		regexp.MustCompile("```[\\s\\S]*?```"),              // Code blocks
		regexp.MustCompile(`^\s*[-*+]\s+`),                  // Unordered lists
		regexp.MustCompile(`^\s*\d+\.\s+`),                  // Ordered lists
		regexp.MustCompile(`^\s*>\s+`),                      // Blockquotes
		regexp.MustCompile(`\|.*\|`),                        // Tables
	}

	for _, pattern := range markdownPatterns {
		if pattern.MatchString(content) {
			return "markdown"
		}
	}

	// Default to text
	return "text"
}