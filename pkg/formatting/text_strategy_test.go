package formatting_test

import (
	"testing"

	"github.com/TwiN/go-color"
	"github.com/stretchr/testify/assert"

	"github.com/conjurinc/conjur-preflight/pkg/formatting"
)

func TestRichTextFormatStrategy(t *testing.T) {

	richText := formatting.RichANSIFormatStrategy{}

	inputText := "test"

	boldText := richText.Bold(inputText)
	assert.Equal(t, "\x1b[1mtest\x1b[0m", boldText)

	colorText := richText.Color(inputText, color.Red)
	assert.Equal(t, "\x1b[31mtest\x1b[0m", colorText)
}

func TestPlainTextFormatStrategy(t *testing.T) {

	plainText := formatting.PlainFormatStrategy{}

	inputText := "test"

	boldText := plainText.Bold(inputText)
	assert.Equal(t, inputText, boldText)

	colorText := plainText.Color(inputText, color.Red)
	assert.Equal(t, inputText, colorText)
}
