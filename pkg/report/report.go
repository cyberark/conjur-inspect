package report

import (
	"fmt"
	"io"
	"os"

	"github.com/conjurinc/conjur-preflight/pkg/framework"
	"github.com/conjurinc/conjur-preflight/pkg/version"
	"github.com/schollz/progressbar/v3"
)

// Report contains an array of all sections and their reports
type Report struct {
	Sections []Section `json:"sections"`
}

// Run starts each check and returns a report of the results
func (report *Report) Run() Result {
	result := Result{
		Version:  version.FullVersionName,
		Sections: make([]ResultSection, len(report.Sections)),
	}

	// Initialize the progress indicator
	progress := newProgress(report.checkCount(), os.Stderr)

	for i, section := range report.Sections {

		sectionResults := []framework.CheckResult{}

		for _, check := range section.Checks {
			// Update text in progress display
			progress.Describe(fmt.Sprintf("Checking %s...", check.Describe()))

			// Start check, this happens asynchronously
			checkResults := <-check.Run()

			// Add the results to the report section
			sectionResults = append(sectionResults, checkResults...)

			// Increment progress
			progress.Add(1)
		}

		result.Sections[i] = ResultSection{
			Title:   section.Title,
			Results: sectionResults,
		}
	}

	progress.Finish()

	return result
}

func (report *Report) checkCount() int {
	count := 0
	for _, section := range report.Sections {
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
