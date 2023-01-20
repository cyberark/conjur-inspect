package framework

import "github.com/TwiN/go-color"

// FormatStrategy describes what operations a particular format must
// implement.
type FormatStrategy interface {
	FormatBold(str string) string
	FormatColor(str, color string) string
}

// RichTextFormatStrategy will format report text for an interactive terminal
// using ANSI escape codes for color and font style.
type RichTextFormatStrategy struct{}

// FormatBold wraps the provided string with the ANSI
// escape codes for bold text.
func (*RichTextFormatStrategy) FormatBold(str string) string {
	return color.InBold(str)
}

// FormatColor wraps the provided string with the ANSI
// escape codes for the given color.
func (*RichTextFormatStrategy) FormatColor(str, colorStr string) string {
	return color.With(colorStr, str)
}

// PlainTextFormatStrategy will format report text for file or plain output
// and will apply no styling to the text.
type PlainTextFormatStrategy struct{}

// FormatBold returns the original string with no
// modification.
func (*PlainTextFormatStrategy) FormatBold(str string) string {
	// Plaintext formatting is a no-op
	return str
}

// FormatColor returns the original string with no
// modification.
func (*PlainTextFormatStrategy) FormatColor(str, colorStr string) string {
	// Plaintext formatting is a no-op
	return str
}
