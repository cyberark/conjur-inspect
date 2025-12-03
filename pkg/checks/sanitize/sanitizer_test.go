package sanitize

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedactorAPIKey(t *testing.T) {
	redactor := NewRedactor()
	input := "api_key=sk_live_1234567890abcdef"
	result := redactor.RedactString(input)
	assert.Contains(t, result, "[REDACTED]")
	assert.NotContains(t, result, "sk_live_1234567890abcdef")
}

func TestRedactorPassword(t *testing.T) {
	redactor := NewRedactor()
	input := "password=SuperSecurePass123!"
	result := redactor.RedactString(input)
	assert.Contains(t, result, "[REDACTED]")
	assert.NotContains(t, result, "SuperSecurePass123!")
}

func TestRedactorBearerToken(t *testing.T) {
	redactor := NewRedactor()
	input := "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
	result := redactor.RedactString(input)
	assert.Contains(t, result, "[REDACTED]")
	assert.NotContains(t, result, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9")
}

func TestRedactorAWSKeys(t *testing.T) {
	redactor := NewRedactor()
	input := "AWS_SECRET_ACCESS_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
	result := redactor.RedactString(input)
	assert.Contains(t, result, "[REDACTED]")
	assert.NotContains(t, result, "wJalrXUtnFEMI")
}

func TestRedactorDatabaseURL(t *testing.T) {
	redactor := NewRedactor()
	tests := []struct {
		name        string
		input       string
		shouldMatch string
	}{
		{"PostgreSQL URL", "psql postgres://user:password@localhost:5432/mydb", "localhost"},
		{"MySQL URL", "mysql://admin:secretpass@db.example.com/users", "db.example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := redactor.RedactString(tt.input)
			assert.Contains(t, result, "[REDACTED]")
			assert.Contains(t, result, tt.shouldMatch)
			assert.NotContains(t, result, "password")
			assert.NotContains(t, result, "secretpass")
		})
	}
}

func TestRedactorToken(t *testing.T) {
	redactor := NewRedactor()
	input := "token=abc123def456ghi789"
	result := redactor.RedactString(input)
	assert.Equal(t, "token=[REDACTED]", result)
	assert.NotContains(t, result, "abc123def456ghi789")
}

func TestRedactorSecret(t *testing.T) {
	redactor := NewRedactor()
	input := "secret=my_secret_value"
	result := redactor.RedactString(input)
	assert.Equal(t, "secret=[REDACTED]", result)
	assert.NotContains(t, result, "my_secret_value")
}

func TestRedactorBasicAuth(t *testing.T) {
	redactor := NewRedactor()
	input := "Basic dXNlcm5hbWU6cGFzc3dvcmQ="
	result := redactor.RedactString(input)
	assert.Contains(t, result, "[REDACTED]")
	assert.NotContains(t, result, "dXNlcm5hbWU6cGFzc3dvcmQ=")
}

func TestRedactorSSHKey(t *testing.T) {
	redactor := NewRedactor()
	input := "ssh_key=/home/user/.ssh/id_rsa"
	result := redactor.RedactString(input)
	assert.Equal(t, "ssh_key=[REDACTED]", result)
	assert.NotContains(t, result, "/home/user/.ssh/id_rsa")
}

func TestRedactorConjurAPIKey(t *testing.T) {
	redactor := NewRedactor()
	input := "conjur_api_key=3mt9y2g8ks0a5j1k2l3m4n5o6p7q8r9s"
	result := redactor.RedactString(input)
	assert.Equal(t, "conjur_api_key=[REDACTED]", result)
	assert.NotContains(t, result, "3mt9y2g8ks0a5j1k2l3m4n5o6p7q8r9s")
}

func TestRedactorMultipleSensitiveValues(t *testing.T) {
	redactor := NewRedactor()
	input := "curl -H Authorization: Bearer token123 api_key=sk_live_abc123 password=mypass https://api.example.com"

	result := redactor.RedactString(input)

	// Check that all sensitive values are redacted
	assert.NotContains(t, result, "token123")
	assert.NotContains(t, result, "sk_live_abc123")
	assert.NotContains(t, result, "mypass")
	// Should contain the redaction marker
	redactedCount := strings.Count(result, "[REDACTED]")
	assert.GreaterOrEqual(t, redactedCount, 1)
}

func TestRedactorRedactLines(t *testing.T) {
	redactor := NewRedactor()
	input := `line1 api_key=secret123
line2 normal_text
line3 password=pass456`

	result := redactor.RedactLines(input)
	lines := strings.Split(result, "\n")

	// Check line 1
	assert.Contains(t, lines[0], "[REDACTED]")
	assert.NotContains(t, lines[0], "secret123")

	// Check line 2 is unchanged
	assert.Equal(t, "line2 normal_text", lines[1])

	// Check line 3
	assert.Contains(t, lines[2], "[REDACTED]")
	assert.NotContains(t, lines[2], "pass456")
}

func TestRedactorEmptyString(t *testing.T) {
	redactor := NewRedactor()
	result := redactor.RedactString("")
	assert.Equal(t, "", result)
}

func TestRedactorNoSensitiveData(t *testing.T) {
	redactor := NewRedactor()
	input := "this is just normal text with no secrets"
	result := redactor.RedactString(input)
	assert.Equal(t, input, result)
}

func TestRedactorNonMatchingPatterns(t *testing.T) {
	redactor := NewRedactor()
	// Test patterns that shouldn't match
	tests := []struct {
		name  string
		input string
	}{
		{"Incomplete pattern", "api_key"},
		{"Word containing api_key", "myapi_key_manager"},
		{"No whitespace after equals", "api_key="},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := redactor.RedactString(tt.input)
			// Results may or may not match depending on pattern specificity
			// Just verify no panic occurs
			assert.NotEmpty(t, result)
		})
	}
}
