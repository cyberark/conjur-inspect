package reports

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/cyberark/conjur-inspect/pkg/check"
	"github.com/cyberark/conjur-inspect/pkg/formatting"
	"github.com/cyberark/conjur-inspect/pkg/log"
	"github.com/cyberark/conjur-inspect/pkg/output"
	"github.com/cyberark/conjur-inspect/pkg/report"
	"github.com/cyberark/conjur-inspect/pkg/version"
	"github.com/schollz/progressbar/v3"
)

type standardReport struct {
	id       string
	sections []report.Section

	outputStore   output.Store
	outputArchive output.Archive
}

// NewReport instantiates a new Report struct with the expected fields (e.g. ID)
func NewStandardReport(
	id string,
	rawDataDir string,
	sections []report.Section,
) (report.Report, error) {

	storeDirectory := path.Join(rawDataDir, id)

	err := os.MkdirAll(storeDirectory, 0755)
	if err != nil {
		return nil, err
	}

	outputStore := output.NewDirectoryStore(storeDirectory)
	outputArchive := &output.TarGzipArchive{
		OutputDir: rawDataDir,
	}

	newReport := standardReport{
		id:       id,
		sections: sections,

		outputStore:   outputStore,
		outputArchive: outputArchive,
	}

	return &newReport, nil
}

func (sr *standardReport) ID() string {
	return sr.id
}

// Run starts each check and returns a report of the results
func (sr *standardReport) Run(containerID string) report.Result {
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
					ContainerID: containerID,
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

func (sr *standardReport) checkCount() int {
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
