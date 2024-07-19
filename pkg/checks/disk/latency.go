package disk

import (
	"bytes"
	"fmt"

	"github.com/cyberark/conjur-inspect/pkg/check"
	"github.com/cyberark/conjur-inspect/pkg/checks/disk/fio"
	"github.com/cyberark/conjur-inspect/pkg/log"
	"github.com/cyberark/conjur-inspect/pkg/output"
)

// LatencyCheck is a inspection check to report the read, write, and sync
// latency for the directory in which `conjur-inspect` is run.
type LatencyCheck struct {
	// We inject the fio command execution as a dependency that we can swap for
	// unit testing
	fioNewJob func(string, []string) fio.Executable
}

// NewLatencyCheck instantiates a Latency check with the default dependencies
func NewLatencyCheck() *LatencyCheck {
	return &LatencyCheck{

		// Default dependencies
		fioNewJob: fio.NewJob,
	}
}

// Describe provides a textual description of what this check gathers info on
func (*LatencyCheck) Describe() string {
	return "disk latency"
}

// Run executes the LatencyCheck by running `fio` and processing its output
func (latencyCheck *LatencyCheck) Run(
	context *check.RunContext,
) []check.Result {
	fioResult, err := latencyCheck.runFioLatencyTest(
		context.OutputStore,
	)

	if err != nil {
		return []check.Result{
			{
				Title:   "FIO Latency",
				Status:  check.StatusError,
				Value:   "N/A",
				Message: err.Error(),
			},
		}
	}

	// Make sure a job exists in the fio results
	if len(fioResult.Jobs) < 1 {
		return []check.Result{
			{
				Title:   "FIO Latency",
				Status:  check.StatusError,
				Value:   "N/A",
				Message: "No job results returned by 'fio'",
			},
		}
	}

	return []check.Result{
		fioReadLatencyResult(&fioResult.Jobs[0]),
		fioWriteLatencyResult(&fioResult.Jobs[0]),
		fioSyncLatencyResult(&fioResult.Jobs[0]),
	}
}

func fioReadLatencyResult(jobResult *fio.JobResult) check.Result {
	// Convert the nanosecond result to milliseconds for readability
	latMs := float64(jobResult.Read.LatNs.Percentile.NinetyNinth) / 1e6

	latMsStr := fmt.Sprintf("%0.2f ms", latMs)

	status := check.StatusInfo
	if latMs > 10.0 {
		status = check.StatusWarn
	}

	path, err := getWorkingDirectory()
	if err != nil {
		log.Debug("Unable to get working directory: %s", err)
		path = "working directory"
	}

	return check.Result{
		Title:  fmt.Sprintf("FIO - Read Latency (99%%, %s)", path),
		Status: status,
		Value:  latMsStr,
	}
}

func fioWriteLatencyResult(jobResult *fio.JobResult) check.Result {
	// Convert the nanosecond result to milliseconds for readability
	latMs := float64(jobResult.Write.LatNs.Percentile.NinetyNinth) / 1e6

	latMsStr := fmt.Sprintf("%0.2f ms", latMs)

	status := check.StatusInfo
	if latMs > 10.0 {
		status = check.StatusWarn
	}

	path, err := getWorkingDirectory()
	if err != nil {
		log.Debug("Unable to get working directory: %s", err)
		path = "working directory"
	}

	return check.Result{
		Title:  fmt.Sprintf("FIO - Write Latency (99%%, %s)", path),
		Status: status,
		Value:  latMsStr,
	}
}

func fioSyncLatencyResult(jobResult *fio.JobResult) check.Result {
	// Convert the nanosecond result to milliseconds for readability
	latMs := float64(jobResult.Sync.LatNs.Percentile.NinetyNinth) / 1e6

	latMsStr := fmt.Sprintf("%0.2f ms", latMs)

	status := check.StatusInfo
	if latMs > 10.0 {
		status = check.StatusWarn
	}

	path, err := getWorkingDirectory()
	if err != nil {
		log.Debug("Unable to get working directory: %s", err)
		path = "working directory"
	}

	return check.Result{
		Title:  fmt.Sprintf("FIO - Sync Latency (99%%, %s)", path),
		Status: status,
		Value:  latMsStr,
	}
}

func (latencyCheck *LatencyCheck) runFioLatencyTest(
	store output.Store,
) (*fio.Result, error) {
	job := latencyCheck.fioNewJob(
		"conjur-fio-latency",
		[]string{
			"--rw=readwrite",
			"--ioengine=sync",
			"--fdatasync=1",
			"--directory=conjur-fio-latency",
			"--size=22m",
			"--bs=2300",
			"--name=conjur-fio-latency",
			"--output-format=json",
		},
	)

	// Save the full `fio` output to the results store
	job.OnRawOutput(func(data []byte) {
		store.Save("conjur-fio-latency", bytes.NewReader(data))
	})

	return job.Exec()
}
