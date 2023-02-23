package report

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/TwiN/go-color"
	"github.com/conjurinc/conjur-preflight/pkg/framework"
	"github.com/conjurinc/conjur-preflight/pkg/maybe"
	"github.com/conjurinc/conjur-preflight/pkg/version"
)

// Result contains each sections check result
type Result struct {
	Version  string          `json:"version"`
	Sections []ResultSection `json:"sections"`
}

// ResultSection is the individual check and its result
type ResultSection struct {
	Title   string                  `json:"title"`
	Results []framework.CheckResult `json:"results"`
}

// ToJSON outputs a JSON formated report.
func (result *Result) ToJSON() (string, error) {
	//Generate the JSON representation of the report
	out, err := json.MarshalIndent(result, "", " ")
	if err != nil {
		return "", err
	}

	return string(out), err
}

// ToText outputs the text for a given report result applying
// the designated format strategy.
func (result *Result) ToText(format framework.FormatStrategy) (string, error) {
	// Write the string parts to a buffer with maybe monads for streamlined
	// error handling.
	maybeBuffer := maybe.NewBuffer()

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

func resultLine(result framework.CheckResult) string {
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
	case framework.STATUS_ERROR:
		return color.Red
	case framework.STATUS_FAIL:
		return color.Red
	case framework.STATUS_WARN:
		return color.Yellow
	case framework.STATUS_PASS:
		return color.Green
	}

	return ""
}
