// Package sanitize provides utilities for redacting sensitive information from strings.
package sanitize

import (
	"regexp"
	"strings"
)

// Pattern represents a regex pattern used to identify and redact sensitive values
type Pattern struct {
	Name    string
	Regex   *regexp.Regexp
	Replace string // Replacement format: $0 for full match, $1+ for groups
}

// SensitivePatterns contains regex patterns to match common sensitive data
var SensitivePatterns = []Pattern{
	{
		Name:    "API Key",
		Regex:   regexp.MustCompile(`(?i)(api[_-]?key\s*[:=\s]\s*)\S+`),
		Replace: "$1[REDACTED]",
	},
	{
		Name:    "Password",
		Regex:   regexp.MustCompile(`(?i)((?:^|[^a-zA-Z])p(?:assword)?\s*[:=\s]\s*)\S+`),
		Replace: "$1[REDACTED]",
	},
	{
		Name:    "Bearer Token",
		Regex:   regexp.MustCompile(`(?i)(bearer\s+)\S+`),
		Replace: "$1[REDACTED]",
	},
	{
		Name:    "Authorization Header",
		Regex:   regexp.MustCompile(`(?i)(Authorization[:=\s]\s*)\S+`),
		Replace: "$1[REDACTED]",
	},
	{
		Name:    "Token Assignment",
		Regex:   regexp.MustCompile(`(?i)(token\s*[:=\s]\s*)\S+`),
		Replace: "$1[REDACTED]",
	},
	{
		Name:    "Secret Assignment",
		Regex:   regexp.MustCompile(`(?i)(secret\s*[:=\s]\s*)\S+`),
		Replace: "$1[REDACTED]",
	},
	{
		Name:    "Key Assignment",
		Regex:   regexp.MustCompile(`(?i)(key\s*[:=\s]\s*)\S+`),
		Replace: "$1[REDACTED]",
	},
	{
		Name:    "Database URL",
		Regex:   regexp.MustCompile(`(?i)((?:postgres|mysql|mongodb|redis)://)[^\s@]+@`),
		Replace: "$1[REDACTED]@",
	},
	{
		Name:    "Basic Auth",
		Regex:   regexp.MustCompile(`(?i)(basic\s+)\S+`),
		Replace: "$1[REDACTED]",
	},
}

// Redactor redacts sensitive values from text
type Redactor struct {
	patterns []Pattern
}

// NewRedactor creates a new Redactor with default sensitive patterns
func NewRedactor() *Redactor {
	return &Redactor{
		patterns: SensitivePatterns,
	}
}

// RedactString redacts all sensitive patterns from the input string
func (r *Redactor) RedactString(input string) string {
	result := input
	for _, pattern := range r.patterns {
		result = pattern.Regex.ReplaceAllString(result, pattern.Replace)
	}
	return result
}

// RedactLines redacts sensitive values from each line in the input string
// Returns the redacted content
func (r *Redactor) RedactLines(input string) string {
	lines := strings.Split(input, "\n")
	redactedLines := make([]string, len(lines))
	for i, line := range lines {
		redactedLines[i] = r.RedactString(line)
	}
	return strings.Join(redactedLines, "\n")
}
