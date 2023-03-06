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

func TestLatencyCheck(t *testing.T) {
	testCheck := &LatencyCheck{
		fioNewJob: newSuccessfulLatencyFioJob,
	}
	resultChan := testCheck.Run()
	results := <-resultChan

	assert.Equal(t, 3, len(results), "There are disk latency results present")

	assertReadLatencyResult(t, results[0], framework.STATUS_INFO)
	assertWriteLatencyResult(t, results[1], framework.STATUS_INFO)
	assertSyncLatencyResult(t, results[2], framework.STATUS_INFO)
}

func TestLatencyCheckWithPoorPerformance(t *testing.T) {
	testCheck := &LatencyCheck{
		fioNewJob: newPoorLatencyPerformanceFioJob,
	}
	resultChan := testCheck.Run()
	results := <-resultChan

	assert.Equal(t, 3, len(results), "There are disk latency results present")

	assertReadLatencyResult(t, results[0], framework.STATUS_WARN)
	assertWriteLatencyResult(t, results[1], framework.STATUS_WARN)
	assertSyncLatencyResult(t, results[2], framework.STATUS_WARN)
}

func TestLatencyWithError(t *testing.T) {
	testCheck := &LatencyCheck{
		fioNewJob: newErrorFioJob,
	}
	resultChan := testCheck.Run()
	results := <-resultChan

	// Expect only the error result
	assert.Equal(t, 1, len(results))

	assert.Equal(t, "FIO Latency", results[0].Title)
	assert.Equal(t, framework.STATUS_ERROR, results[0].Status)
	assert.Equal(t, "N/A", results[0].Value)
	assert.Equal(t, "test error", results[0].Message)
}

func TestLatencyWithNoJobs(t *testing.T) {
	testCheck := &LatencyCheck{
		fioNewJob: newEmptyFioJob,
	}
	resultChan := testCheck.Run()
	results := <-resultChan

	// Expect only the error result
	assert.Equal(t, 1, len(results))

	assert.Equal(t, "FIO Latency", results[0].Title)
	assert.Equal(t, framework.STATUS_ERROR, results[0].Status)
	assert.Equal(t, "N/A", results[0].Value)
	assert.Equal(t, "No job results returned by 'fio'", results[0].Message)
}

func TestLatencyWithWorkingDirectoryError(t *testing.T) {
	// Double the working directory function to simulate it failing with an error
	originalWorkingDirectoryFunc := getWorkingDirectory
	getWorkingDirectory = failedWorkingDir
	defer func() {
		getWorkingDirectory = originalWorkingDirectoryFunc
	}()

	testCheck := &LatencyCheck{
		fioNewJob: newSuccessfulLatencyFioJob,
	}
	resultChan := testCheck.Run()
	results := <-resultChan

	assert.Equal(t, 3, len(results), "There are disk latency results present")

	assertReadLatencyResult(t, results[0], framework.STATUS_INFO)
	assertWriteLatencyResult(t, results[1], framework.STATUS_INFO)
	assertSyncLatencyResult(t, results[2], framework.STATUS_INFO)
}

func assertReadLatencyResult(
	t *testing.T,
	result framework.CheckResult,
	expectedStatus string,
) {
	assertLatencyResult(
		t,
		result,
		`FIO - Read Latency \(99%, .+\)`,
		expectedStatus,
	)
}

func assertWriteLatencyResult(
	t *testing.T,
	result framework.CheckResult,
	expectedStatus string,
) {
	assertLatencyResult(
		t,
		result,
		`FIO - Write Latency \(99%, .+\)`,
		expectedStatus,
	)
}

func assertSyncLatencyResult(
	t *testing.T,
	result framework.CheckResult,
	expectedStatus string,
) {
	assertLatencyResult(
		t,
		result,
		`FIO - Sync Latency \(99%, .+\)`,
		expectedStatus,
	)
}

func assertLatencyResult(
	t *testing.T,
	result framework.CheckResult,
	expectedTitleRegex string,
	expectedStatus string,
) {
	assert.Regexp(t, regexp.MustCompile(expectedTitleRegex), result.Title)
	assert.Equal(t, expectedStatus, result.Status)
	assert.Regexp(t, regexp.MustCompile(`.+ ms`), result.Value)
}

func newSuccessfulLatencyFioJob(jobName string, args []string) fio.Executable {
	return &mockFioJob{
		Result: fio.Result{
			Jobs: []fio.JobResult{
				{
					Read: fio.JobModeResult{
						LatNs: fio.ResultStats{
							Percentile: fio.Percentile{
								NinetyNinth: 10 * 1e6,
							},
						},
					},
					Write: fio.JobModeResult{
						LatNs: fio.ResultStats{
							Percentile: fio.Percentile{
								NinetyNinth: 10 * 1e6,
							},
						},
					},
					Sync: fio.JobModeResult{
						LatNs: fio.ResultStats{
							Percentile: fio.Percentile{
								NinetyNinth: 10 * 1e6,
							},
						},
					},
				},
			},
		},
	}
}

func newPoorLatencyPerformanceFioJob(
	jobName string,
	args []string,
) fio.Executable {
	return &mockFioJob{
		Result: fio.Result{
			Jobs: []fio.JobResult{
				{
					Read: fio.JobModeResult{
						LatNs: fio.ResultStats{
							Percentile: fio.Percentile{
								NinetyNinth: 11 * 1e6,
							},
						},
					},
					Write: fio.JobModeResult{
						LatNs: fio.ResultStats{
							Percentile: fio.Percentile{
								NinetyNinth: 11 * 1e6,
							},
						},
					},
					Sync: fio.JobModeResult{
						LatNs: fio.ResultStats{
							Percentile: fio.Percentile{
								NinetyNinth: 11 * 1e6,
							},
						},
					},
				},
			},
		},
	}
}

func failedWorkingDir() (string, error) {
	return "", errors.New("working directory error")
}
