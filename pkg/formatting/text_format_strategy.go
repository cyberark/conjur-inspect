package formatting

import "github.com/TwiN/go-color"

// TextFormatStrategy describes what operations a particular format must
// implement.
type TextFormatStrategy interface {
	Bold(str string) string
	Color(str, color string) string
}

// RichANSIFormatStrategy will format report text for an interactive terminal
// using ANSI escape codes for color and font style.
type RichANSIFormatStrategy struct{}

// Bold wraps the provided string with the ANSI
// escape codes for bold text.
func (*RichANSIFormatStrategy) Bold(str string) string {
	return color.InBold(str)
}

// Color wraps the provided string with the ANSI
// escape codes for the given color.
func (*RichANSIFormatStrategy) Color(str, colorStr string) string {
	return color.With(colorStr, str)
}

// PlainFormatStrategy will format report text for file or plain output
// and will apply no styling to the text.
type PlainFormatStrategy struct{}

// Bold returns the original string with no modification.
func (*PlainFormatStrategy) Bold(str string) string {
	// Plaintext formatting is a no-op
	return str
}

// Color returns the original string with no modification.
func (*PlainFormatStrategy) Color(str, colorStr string) string {
	// Plaintext formatting is a no-op
	return str
}
