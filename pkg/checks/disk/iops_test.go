// This can't be in the disk_test package because this requires access to the
// internal fioExec field on LatencyCheck.
package disk

import (
	"errors"
	"regexp"
	"testing"

	"github.com/cyberark/conjur-inspect/pkg/check"
	"github.com/cyberark/conjur-inspect/pkg/checks/disk/fio"
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
		fioNewJob: newSuccessfulIopsFioJob,
	}
	results := testCheck.Run(&check.RunContext{})

	assert.Equal(
		t,
		2,
		len(results),
		"There are read and write IOPs results present",
	)

	assertReadIopsResult(t, results[0], check.StatusInfo)
	assertWriteIopsResult(t, results[1], check.StatusInfo)
}

func TestIopsCheckWithPoorPerformance(t *testing.T) {
	testCheck := &IopsCheck{
		fioNewJob: newPoorIopsPerformanceFioJob,
	}
	results := testCheck.Run(&check.RunContext{})

	assert.Equal(
		t,
		2,
		len(results),
		"There are read and write IOPs results present",
	)

	assertReadIopsResult(t, results[0], check.StatusWarn)
	assertWriteIopsResult(t, results[1], check.StatusWarn)
}

func TestIopsWithError(t *testing.T) {
	testCheck := &IopsCheck{
		fioNewJob: newErrorFioJob,
	}
	results := testCheck.Run(&check.RunContext{VerboseErrors: true})

	// Expect only the error result
	assert.Equal(t, 1, len(results))

	assert.Equal(t, "FIO IOPs", results[0].Title)
	assert.Equal(t, check.StatusError, results[0].Status)
	assert.Equal(t, "N/A", results[0].Value)
	assert.Equal(t, "test error", results[0].Message)
}

func TestIopsWithErrorVerboseErrors(t *testing.T) {
	testCheck := &IopsCheck{
		fioNewJob: newErrorFioJob,
	}
	results := testCheck.Run(&check.RunContext{VerboseErrors: true})

	// Expect only the error result
	assert.Equal(t, 1, len(results))

	assert.Equal(t, "FIO IOPs", results[0].Title)
	assert.Equal(t, check.StatusError, results[0].Status)
	assert.Equal(t, "N/A", results[0].Value)
	assert.Equal(t, "test error", results[0].Message)
}

func TestIopsWithErrorNoVerboseErrors(t *testing.T) {
	testCheck := &IopsCheck{
		fioNewJob: newErrorFioJob,
	}
	results := testCheck.Run(&check.RunContext{VerboseErrors: false})

	// Expect no results when VerboseErrors is false
	assert.Equal(t, 0, len(results))
}

func TestIopsWithNoJobs(t *testing.T) {
	testCheck := &IopsCheck{
		fioNewJob: newEmptyFioJob,
	}
	results := testCheck.Run(&check.RunContext{})

	// Expect only the error result
	assert.Equal(t, 1, len(results))

	assert.Equal(t, "FIO IOPs", results[0].Title)
	assert.Equal(t, check.StatusError, results[0].Status)
	assert.Equal(t, "N/A", results[0].Value)
	assert.Equal(t, "No job results returned by 'fio'", results[0].Message)
}

func TestIopsWithWorkingDirectoryError(t *testing.T) {
	// Double the working directory function to simulate it failing with an error
	originalWorkingDirectoryFunc := getWorkingDirectory
	getWorkingDirectory = failedWorkingDir
	defer func() {
		getWorkingDirectory = originalWorkingDirectoryFunc
	}()

	testCheck := &IopsCheck{
		fioNewJob: newSuccessfulIopsFioJob,
	}
	results := testCheck.Run(&check.RunContext{})

	assert.Equal(
		t,
		2,
		len(results),
		"There are read and write IOPs results present",
	)

	assertReadIopsResult(t, results[0], check.StatusInfo)
	assertWriteIopsResult(t, results[1], check.StatusInfo)
}

func assertReadIopsResult(
	t *testing.T,
	result check.Result,
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
	result check.Result,
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
