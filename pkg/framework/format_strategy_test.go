package framework_test

import (
	"testing"

	"github.com/TwiN/go-color"
	"github.com/conjurinc/conjur-preflight/pkg/framework"
	"github.com/stretchr/testify/assert"
)

func TestRichTextFormatStrategy(t *testing.T) {

	richText := framework.RichTextFormatStrategy{}

	inputText := "test"

	boldText := richText.FormatBold(inputText)
	assert.Equal(t, "\x1b[1mtest\x1b[0m", boldText)

	colorText := richText.FormatColor(inputText, color.Red)
	assert.Equal(t, "\x1b[31mtest\x1b[0m", colorText)
}

func TestPlainTextFormatStrategy(t *testing.T) {

	plainText := framework.PlainTextFormatStrategy{}

	inputText := "test"

	boldText := plainText.FormatBold(inputText)
	assert.Equal(t, inputText, boldText)

	colorText := plainText.FormatColor(inputText, color.Red)
	assert.Equal(t, inputText, colorText)
}
