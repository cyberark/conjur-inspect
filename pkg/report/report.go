package report

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/cyberark/conjur-inspect/pkg/check"
	"github.com/cyberark/conjur-inspect/pkg/log"
	"github.com/cyberark/conjur-inspect/pkg/output"
	"github.com/cyberark/conjur-inspect/pkg/version"
	"github.com/schollz/progressbar/v3"
)

// Report contains an array of all sections and their reports
type Report interface {
	ID() string
	Run() Result
}

type report struct {
	id       string
	sections []Section

	outputStore   output.Store
	outputArchive output.Archive
}

// NewReport instantiates a new Report struct with the expected fields (e.g. ID)
func NewReport(
	id string,
	rawDataDir string,
	sections []Section,
) (Report, error) {

	storeDirectory := path.Join(rawDataDir, id)

	err := os.MkdirAll(storeDirectory, 0755)
	if err != nil {
		return nil, err
	}

	outputStore := output.NewDirectoryStore(storeDirectory)
	outputArchive := &output.TarGzipArchive{
		OutputDir: rawDataDir,
	}

	newReport := report{
		id:       id,
		sections: sections,

		outputStore:   outputStore,
		outputArchive: outputArchive,
	}

	return &newReport, nil
}

func (report *report) ID() string {
	return report.id
}

// Run starts each check and returns a report of the results
func (report *report) Run() Result {
	defer report.outputStore.Cleanup()

	result := Result{
		Version:  version.FullVersionName,
		Sections: make([]ResultSection, len(report.sections)),
	}

	// Initialize the progress indicator
	progress := newProgress(report.checkCount(), os.Stderr)

	for i, section := range report.sections {

		sectionResults := []check.Result{}

		for _, currentCheck := range section.Checks {
			// Update text in progress display
			progress.Describe(fmt.Sprintf("Checking %s...", currentCheck.Describe()))

			// Start check, this happens asynchronously
			checkResults := <-currentCheck.Run(
				&check.RunContext{
					OutputStore: report.outputStore,
				},
			)

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

	// Archive the raw outputs
	err := report.outputArchive.Archive(
		report.ID(),
		report.outputStore,
	)
	if err != nil {
		log.Error("Failed to save raw output: %s", err)
	}

	return result
}

func (report *report) checkCount() int {
	count := 0
	for _, section := range report.sections {
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
