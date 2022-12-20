package framework

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/TwiN/go-color"
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

func (result *ReportResult) ToText() string {
	buf := new(bytes.Buffer)

	fmt.Fprintf(
		buf,
		color.InBold(
			"%s\n%s\n%s\n%s\n",
		),
		"========================================",
		"Conjur Enterprise Preflight Qualification",
		"Version: 0.1.0",
		"========================================",
	)

	for _, section := range result.Sections {
		fmt.Fprintf(
			buf,
			color.InBold(
				"\n%s\n%s\n",
			),
			section.Title,
			strings.Repeat("-", len(section.Title)),
		)

		for _, result := range section.Results {

			if result.Message == "" {
				fmt.Fprintf(
					buf,
					color.With(
						statusColor(result.Status),
						"%s - %s: %s\n",
					),
					result.Status,
					result.Title,
					result.Value,
				)
			} else {
				fmt.Fprintf(
					buf,
					color.With(
						statusColor(result.Status),
						"%s - %s: %s (%s)\n",
					),
					result.Status,
					result.Title,
					result.Value,
					result.Message,
				)
			}
		}
	}

	return buf.String()
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
