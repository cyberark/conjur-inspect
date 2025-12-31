package formatting

import (
	"fmt"
	"io"
	"strings"

	"github.com/cyberark/conjur-inspect/pkg/check"
	"github.com/cyberark/conjur-inspect/pkg/maybe"
	"github.com/cyberark/conjur-inspect/pkg/report"

	"github.com/TwiN/go-color"
)

// Text renders a report result as text, using a given format strategy
type Text struct {
	FormatStrategy TextFormatStrategy
}

func (text *Text) Write(
	writer io.Writer,
	result *report.Result,
) error {
	// Write the string parts to a buffer with maybe monads for streamlined
	// error handling.
	maybeWriter := maybe.NewWriter(writer)

	// Write report header
	formattedHeader := text.FormatStrategy.Bold(reportHeader(result.Version))
	maybeWriter.WriteString(formattedHeader)
	maybeWriter.WriteString("\n\n")

	// Filter out sections with no results
	nonEmptySections := []report.ResultSection{}
	for _, section := range result.Sections {
		if len(section.Results) > 0 {
			nonEmptySections = append(nonEmptySections, section)
		}
	}

	// Write each report section
	for sectionIndex, section := range nonEmptySections {
		formattedTitle := text.FormatStrategy.Bold(titleHeader(section.Title))
		maybeWriter.WriteString(formattedTitle)
		maybeWriter.WriteString("\n")

		for _, result := range section.Results {
			formattedResultLine := text.FormatStrategy.Color(
				resultLine(result),
				statusColor(result.Status),
			)
			maybeWriter.WriteString(formattedResultLine)
			maybeWriter.WriteString("\n")
		}

		// Extra space between sections (but not extra space at the end)
		if sectionIndex < len(nonEmptySections)-1 {
			maybeWriter.WriteString("\n")
		}
	}

	return maybeWriter.Error()
}

func reportHeader(version string) string {
	return strings.Join(
		[]string{
			"========================================",
			"Conjur Enterprise Inspection Report",
			fmt.Sprintf("Version: %s", version),
			"========================================",
		},
		"\n",
	)
}

func titleHeader(title string) string {
	return strings.Join(
		[]string{
			title,
			strings.Repeat("-", len(title)),
		},
		"\n",
	)
}

func resultLine(result check.Result) string {
	switch {
	case result.Message == "":
		return fmt.Sprintf(
			"%s - %s: %s",
			result.Status,
			result.Title,
			result.Value,
		)
	default:
		return fmt.Sprintf(
			"%s - %s: %s (%s)",
			result.Status,
			result.Title,
			result.Value,
			result.Message,
		)
	}
}

func statusColor(status string) string {
	switch status {
	case check.StatusError:
		return color.Red
	case check.StatusFail:
		return color.Red
	case check.StatusWarn:
		return color.Yellow
	case check.StatusPass:
		return color.Green
	}

	return ""
}
