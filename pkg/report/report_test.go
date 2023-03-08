package report_test

import (
	"io"
	"strings"
	"testing"

	"github.com/cyberark/conjur-inspect/pkg/check"
	"github.com/cyberark/conjur-inspect/pkg/formatting"
	"github.com/cyberark/conjur-inspect/pkg/report"
	"github.com/stretchr/testify/assert"
)

type TestCheck struct{}

func (*TestCheck) Describe() string {
	return "Test"
}

func (*TestCheck) Run() <-chan []check.Result {
	channel := make(chan []check.Result)

	go func() {
		channel <- []check.Result{
			{
				Title:   "Test Check",
				Status:  "Test Status",
				Value:   "Test Value",
				Message: "Test Message",
			},
		}
	}()

	return channel
}

func TestReport(t *testing.T) {
	testReport := report.Report{
		Sections: []report.Section{
			{
				Title: "Test section",
				Checks: []check.Check{
					&TestCheck{},
				},
			},
		},
	}

	testReportResult := testReport.Run()

	assert.NotEmpty(t, testReportResult.Sections)

	testSection := testReportResult.Sections[0]
	assert.Equal(t, "Test section", testSection.Title)
	assert.NotEmpty(t, testSection.Results)

	testCheckResult := testSection.Results[0]
	assert.Equal(
		t,
		check.Result{
			Title:   "Test Check",
			Status:  "Test Status",
			Value:   "Test Value",
			Message: "Test Message",
		},
		testCheckResult,
	)

	builder := strings.Builder{}

	textWriter := formatting.Text{
		FormatStrategy: &formatting.RichANSIFormatStrategy{},
	}

	err := textWriter.Write(
		io.Writer(&builder),
		&testReportResult,
	)
	assert.Nil(t, err)

	assert.Equal(
		t,
		"\033[1m========================================\n"+
			"Conjur Enterprise Inspection Report\n"+
			"Version: unset-unset (Build unset)\n"+
			"========================================\033[0m\n\n"+
			"\033[1mTest section\n"+
			"------------\033[0m\n"+
			"Test Status - Test Check: Test Value (Test Message)\033[0m\n",
		builder.String(),
	)
}

func TestJSONReport(t *testing.T) {
	testReport := report.Report{
		Sections: []report.Section{
			{
				Title: "Test section",
				Checks: []check.Check{
					&TestCheck{},
				},
			},
		},
	}

	testReportResult := testReport.Run()

	assert.NotEmpty(t, testReportResult.Sections)

	testSection := testReportResult.Sections[0]
	assert.Equal(t, "Test section", testSection.Title)
	assert.NotEmpty(t, testSection.Results)

	testCheckResult := testSection.Results[0]
	assert.Equal(
		t,
		check.Result{
			Title:   "Test Check",
			Status:  "Test Status",
			Value:   "Test Value",
			Message: "Test Message",
		},
		testCheckResult,
	)

	builder := strings.Builder{}

	jsonWriter := formatting.JSON{}

	err := jsonWriter.Write(
		io.Writer(&builder),
		&testReportResult,
	)

	assert.Nil(t, err)

	assert.JSONEq(t,
		`{
            "version": "unset-unset (Build unset)",
            "sections": [
            {
                "title": "Test section",
                "results": [
                {
                    "title": "Test Check",
                    "value": "Test Value",
                    "status": "Test Status",
                    "message": "Test Message"
                }
               ]
              }
             ]
            }`,
		builder.String(),
	)
}
