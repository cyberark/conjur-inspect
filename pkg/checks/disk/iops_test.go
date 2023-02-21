// This can't be in the disk_test package because this requires access to the
// internal fioExec field on LatencyCheck.
package disk

import (
	"errors"
	"regexp"
	"testing"

	"github.com/conjurinc/conjur-preflight/pkg/checks/disk/fio"
	"github.com/conjurinc/conjur-preflight/pkg/framework"
	"github.com/stretchr/testify/assert"
)

type mockFioJob struct {
	Result fio.Result
	Error  error
}

func (job *mockFioJob) Exec() (*fio.Result, error) {
	return &job.Result, job.Error
}

func (job *mockFioJob) OnRawOutput(func([]byte)) {
	// no-op
}

func TestIopsCheck(t *testing.T) {
	testCheck := &IopsCheck{
		debug:     true,
		fioNewJob: newSuccessfulIopsFioJob,
	}
	resultChan := testCheck.Run()
	results := <-resultChan

	assert.Equal(
		t,
		2,
		len(results),
		"There are read and write IOPs results present",
	)

	assertReadIopsResult(t, results[0], framework.STATUS_INFO)
	assertWriteIopsResult(t, results[1], framework.STATUS_INFO)
}

func TestIopsCheckWithPoorPerformance(t *testing.T) {
	testCheck := &IopsCheck{
		fioNewJob: newPoorIopsPerformanceFioJob,
	}
	resultChan := testCheck.Run()
	results := <-resultChan

	assert.Equal(
		t,
		2,
		len(results),
		"There are read and write IOPs results present",
	)

	assertReadIopsResult(t, results[0], framework.STATUS_WARN)
	assertWriteIopsResult(t, results[1], framework.STATUS_WARN)
}

func TestIopsWithError(t *testing.T) {
	testCheck := &IopsCheck{
		fioNewJob: newErrorFioJob,
	}
	resultChan := testCheck.Run()
	results := <-resultChan

	// Expect only the error result
	assert.Equal(t, 1, len(results))

	assert.Equal(t, "FIO IOPs", results[0].Title)
	assert.Equal(t, framework.STATUS_ERROR, results[0].Status)
	assert.Equal(t, "N/A", results[0].Value)
	assert.Equal(t, "test error", results[0].Message)
}

func TestIopsWithNoJobs(t *testing.T) {
	testCheck := &IopsCheck{
		fioNewJob: newEmptyFioJob,
	}
	resultChan := testCheck.Run()
	results := <-resultChan

	// Expect only the error result
	assert.Equal(t, 1, len(results))

	assert.Equal(t, "FIO IOPs", results[0].Title)
	assert.Equal(t, framework.STATUS_ERROR, results[0].Status)
	assert.Equal(t, "N/A", results[0].Value)
	assert.Equal(t, "No job results returned by 'fio'", results[0].Message)
}

func assertReadIopsResult(
	t *testing.T,
	result framework.CheckResult,
	expectedStatus string,
) {
	assert.Regexp(
		t,
		regexp.MustCompile(`FIO - Read IOPs \(.+\)`),
		result.Title,
	)
	assert.Equal(t, expectedStatus, result.Status)
	assert.Regexp(
		t,
		regexp.MustCompile(`.+ \(Min: .+, Max: .+, StdDev: .+\)`),
		result.Value,
	)
}

func assertWriteIopsResult(
	t *testing.T,
	result framework.CheckResult,
	expectedStatus string,
) {
	assert.Regexp(
		t,
		regexp.MustCompile(`FIO - Write IOPs \(.+\)`),
		result.Title,
	)
	assert.Equal(t, expectedStatus, result.Status)
	assert.Regexp(
		t,
		regexp.MustCompile(`.+ \(Min: .+, Max: .+, StdDev: .+\)`),
		result.Value,
	)
}

func newSuccessfulIopsFioJob(jobName string, args []string) fio.Executable {
	return &mockFioJob{
		Result: fio.Result{
			Jobs: []fio.JobResult{
				{
					Read: fio.JobModeResult{
						Iops: 50,
					},
					Write: fio.JobModeResult{
						Iops: 50,
					},
				},
			},
		},
	}
}

func newPoorIopsPerformanceFioJob(
	jobName string,
	args []string,
) fio.Executable {
	return &mockFioJob{
		Result: fio.Result{
			Jobs: []fio.JobResult{
				{
					Read: fio.JobModeResult{
						Iops: 48,
					},
					Write: fio.JobModeResult{
						Iops: 48,
					},
				},
			},
		},
	}
}

func newErrorFioJob(jobName string, args []string) fio.Executable {
	return &mockFioJob{
		Error: errors.New("test error"),
	}
}

func newEmptyFioJob(jobName string, args []string) fio.Executable {
	return &mockFioJob{
		Result: fio.Result{
			Jobs: []fio.JobResult{},
		},
	}
}
