// Package reports contains concrete report implementations
package reports

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/cyberark/conjur-inspect/pkg/check"
	"github.com/cyberark/conjur-inspect/pkg/formatting"
	"github.com/cyberark/conjur-inspect/pkg/log"
	"github.com/cyberark/conjur-inspect/pkg/output"
	"github.com/cyberark/conjur-inspect/pkg/report"
	"github.com/cyberark/conjur-inspect/pkg/version"
	"github.com/schollz/progressbar/v3"
)

// StandardReport is a report that runs a series of checks and reports the
// results.
type StandardReport struct {
	id       string
	sections []report.Section

	outputStore   output.Store
	outputArchive output.Archive
}

// NewStandardReport initializes and returns a new StandardReport.
func NewStandardReport(
	id string,
	sections []report.Section,
	outputStore output.Store,
	outputArchive output.Archive,
) report.Report {
	return &StandardReport{
		id:            id,
		sections:      sections,
		outputStore:   outputStore,
		outputArchive: outputArchive,
	}
}

// ID returns the given ID of the report
func (sr *StandardReport) ID() string {
	return sr.id
}

// Run starts each check and returns a report of the results
func (sr *StandardReport) Run(config report.RunConfig) report.Result {
	defer sr.outputStore.Cleanup()

	result := report.Result{
		Version:  version.FullVersionName,
		Sections: make([]report.ResultSection, len(sr.sections)),
	}

	// Initialize the progress indicator
	progress := newProgress(sr.checkCount(), os.Stderr)

	for i, section := range sr.sections {

		sectionResults := []check.Result{}

		for _, currentCheck := range section.Checks {
			// Update text in progress display
			progress.Describe(fmt.Sprintf("Checking %s...", currentCheck.Describe()))

			// Start check, this happens asynchronously
			checkResults := <-currentCheck.Run(
				&check.RunContext{
					ContainerID: config.ContainerID,
					Since:       config.Since,
					OutputStore: sr.outputStore,
				},
			)

			// Add the results to the report section
			sectionResults = append(sectionResults, checkResults...)

			// Increment progress
			progress.Add(1)
		}

		result.Sections[i] = report.ResultSection{
			Title:   section.Title,
			Results: sectionResults,
		}
	}

	progress.Finish()

	// Write the report result to the output archive
	err := sr.archiveReport(&result)
	if err != nil {
		log.Error("Failed to archive report: %s", err)
	}

	// Archive the raw outputs
	err = sr.outputArchive.Archive(
		sr.ID(),
		sr.outputStore,
	)
	if err != nil {
		log.Error("Failed to save raw output: %s", err)
	}

	return result
}

func (sr *StandardReport) archiveReport(result *report.Result) error {
	var buffer bytes.Buffer

	// Always use JSON for the archived report
	writer := &formatting.JSON{}

	// Write the report result
	err := writer.Write(&buffer, result)
	if err != nil {
		return err
	}

	// Save the report to the output store
	_, err = sr.outputStore.Save("conjur-inspect.json", &buffer)
	if err != nil {
		return err
	}

	return nil
}

func (sr *StandardReport) checkCount() int {
	count := 0
	for _, section := range sr.sections {
		count += len(section.Checks)
	}
	return count
}

func newProgress(checkCount int, writer io.Writer) *progressbar.ProgressBar {
	return progressbar.NewOptions(
		checkCount,
		progressbar.OptionSetWriter(writer),
		progressbar.OptionClearOnFinish(),
		progressbar.OptionShowCount(),
		progressbar.OptionSetElapsedTime(true),
		progressbar.OptionSetPredictTime(false),
	)
}
