package framework_test

import (
	"testing"

	"github.com/conjurinc/conjur-preflight/pkg/framework"
	"github.com/stretchr/testify/assert"
)

type TestCheck struct{}

func (*TestCheck) Run() <-chan []framework.CheckResult {
	channel := make(chan []framework.CheckResult)

	go func() {
		channel <- []framework.CheckResult{
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
	testReport := framework.Report{
		Sections: []framework.ReportSection{
			{
				Title: "Test section",
				Checks: []framework.Check{
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
		framework.CheckResult{
			Title:   "Test Check",
			Status:  "Test Status",
			Value:   "Test Value",
			Message: "Test Message",
		},
		testCheckResult,
	)

	textOutput := testReportResult.ToText()

	assert.Equal(
		t,
		"\033[1m========================================\n"+
			"Conjur Enterprise Preflight Qualification\n"+
			"Version: 0.1.0\n"+
			"========================================\n\033[0m"+
			"\033[1m\nTest section\n"+
			"------------\n\033[0m"+
			"Test Status - Test Check: Test Value (Test Message)\n\033[0m",
		textOutput,
	)
}
