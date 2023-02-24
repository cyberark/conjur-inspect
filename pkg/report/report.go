package report

import (
	"github.com/conjurinc/conjur-preflight/pkg/framework"
	"github.com/conjurinc/conjur-preflight/pkg/version"
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

	for i, section := range report.Sections {

		sectionResults := []framework.CheckResult{}

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
