package framework

import (
	"fmt"
	"strings"

	"github.com/TwiN/go-color"
	"github.com/conjurinc/conjur-preflight/pkg/version"
)

type Report struct {
	Sections []ReportSection
}

type ReportSection struct {
	Title  string
	Checks []Check
}

type ReportResult struct {
	Version  string
	Sections []ResultSection
}

type ResultSection struct {
	Title   string
	Results []CheckResult
}

func (report *Report) Run() ReportResult {

	result := ReportResult{
		Sections: make([]ResultSection, len(report.Sections)),
	}

	for i, section := range report.Sections {

		sectionResults := []CheckResult{}

		for _, check := range section.Checks {
			// Start check, this happens asynchronously
			checkResults := <-check.Run()

			// Add the results to the report section
			sectionResults = append(sectionResults, checkResults...)
		}

		result.Sections[i] = ResultSection{
			Title:   section.Title,
			Results: sectionResults,
		}
	}

	return result
}

// ToText outputs the text for a given report result applying
// the designated format strategy.
func (result *ReportResult) ToText(format FormatStrategy) (string, error) {
	// Write the string parts to a buffer with maybe monads for streamlined
	// error handling.
	maybeBuffer := NewMaybeBuffer()

	// Write report header
	formattedHeader := format.FormatBold(reportHeader())
	maybeBuffer.WriteString(formattedHeader)
	maybeBuffer.WriteString("\n\n")

	// Write each report section
	for sectionIndex, section := range result.Sections {
		formattedTitle := format.FormatBold(titleHeader(section.Title))
		maybeBuffer.WriteString(formattedTitle)
		maybeBuffer.WriteString("\n")

		for _, result := range section.Results {
			formattedResultLine := format.FormatColor(
				resultLine(result),
				statusColor(result.Status),
			)
			maybeBuffer.WriteString(formattedResultLine)
			maybeBuffer.WriteString("\n")
		}

		// Extra space between sections (but not extra space at the end)
		if sectionIndex < len(result.Sections)-1 {
			maybeBuffer.WriteString("\n")
		}
	}

	return maybeBuffer.String()
}

func reportHeader() string {
	return strings.Join(
		[]string{
			"========================================",
			"Conjur Enterprise Preflight Qualification",
			fmt.Sprintf("Version: %s", version.FullVersionName),
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

func resultLine(result CheckResult) string {
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
	case STATUS_ERROR:
		return color.Red
	case STATUS_FAIL:
		return color.Red
	case STATUS_WARN:
		return color.Yellow
	case STATUS_PASS:
		return color.Green
	}

	return ""
}
