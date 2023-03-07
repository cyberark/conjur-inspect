package disk

import (
	"fmt"

	"github.com/conjurinc/conjur-preflight/pkg/checks/disk/fio"
	"github.com/conjurinc/conjur-preflight/pkg/framework"
	"github.com/conjurinc/conjur-preflight/pkg/log"
)

// LatencyCheck is a pre-flight check to report the read, write, and sync
// latency for the directory in which `conjur-preflight` is run.
type LatencyCheck struct {
	// When debug mode is enabled, the latency check will write the full fio
	// results to a file.
	debug bool

	// We inject the fio command execution as a dependency that we can swap for
	// unit testing
	fioNewJob func(string, []string) fio.Executable
}

// NewLatencyCheck instantiates a Latency check with the default dependencies
func NewLatencyCheck(debug bool) *LatencyCheck {
	return &LatencyCheck{
		debug: debug,

		// Default dependencies
		fioNewJob: fio.NewJob,
	}
}

// Describe provides a textual description of what this check gathers info on
func (*LatencyCheck) Describe() string {
	return "disk latency"
}

// Run executes the LatencyCheck by running `fio` and processing its output
func (latencyCheck *LatencyCheck) Run() <-chan []framework.CheckResult {
	future := make(chan []framework.CheckResult)

	go func() {
		fioResult, err := latencyCheck.runFioLatencyTest()

		if err != nil {
			future <- []framework.CheckResult{
				{
					Title:   "FIO Latency",
					Status:  framework.STATUS_ERROR,
					Value:   "N/A",
					Message: err.Error(),
				},
			}

			return
		}

		// Make sure a job exists in the fio results
		if len(fioResult.Jobs) < 1 {
			future <- []framework.CheckResult{
				{
					Title:   "FIO Latency",
					Status:  framework.STATUS_ERROR,
					Value:   "N/A",
					Message: "No job results returned by 'fio'",
				},
			}

			return
		}

		future <- []framework.CheckResult{
			fioReadLatencyResult(&fioResult.Jobs[0]),
			fioWriteLatencyResult(&fioResult.Jobs[0]),
			fioSyncLatencyResult(&fioResult.Jobs[0]),
		}
	}() // async

	return future
}

func fioReadLatencyResult(jobResult *fio.JobResult) framework.CheckResult {
	// Convert the nanosecond result to milliseconds for readability
	latMs := float64(jobResult.Read.LatNs.Percentile.NinetyNinth) / 1e6

	latMsStr := fmt.Sprintf("%0.2f ms", latMs)

	status := framework.STATUS_INFO
	if latMs > 10.0 {
		status = framework.STATUS_WARN
	}

	path, err := getWorkingDirectory()
	if err != nil {
		log.Debug("Unable to get working directory: %s", err)
		path = "working directory"
	}

	return framework.CheckResult{
		Title:  fmt.Sprintf("FIO - Read Latency (99%%, %s)", path),
		Status: status,
		Value:  latMsStr,
	}
}

func fioWriteLatencyResult(jobResult *fio.JobResult) framework.CheckResult {
	// Convert the nanosecond result to milliseconds for readability
	latMs := float64(jobResult.Write.LatNs.Percentile.NinetyNinth) / 1e6

	latMsStr := fmt.Sprintf("%0.2f ms", latMs)

	status := framework.STATUS_INFO
	if latMs > 10.0 {
		status = framework.STATUS_WARN
	}

	path, err := getWorkingDirectory()
	if err != nil {
		log.Debug("Unable to get working directory: %s", err)
		path = "working directory"
	}

	return framework.CheckResult{
		Title:  fmt.Sprintf("FIO - Write Latency (99%%, %s)", path),
		Status: status,
		Value:  latMsStr,
	}
}

func fioSyncLatencyResult(jobResult *fio.JobResult) framework.CheckResult {
	// Convert the nanosecond result to milliseconds for readability
	latMs := float64(jobResult.Sync.LatNs.Percentile.NinetyNinth) / 1e6

	latMsStr := fmt.Sprintf("%0.2f ms", latMs)

	status := framework.STATUS_INFO
	if latMs > 10.0 {
		status = framework.STATUS_WARN
	}

	path, err := getWorkingDirectory()
	if err != nil {
		log.Debug("Unable to get working directory: %s", err)
		path = "working directory"
	}

	return framework.CheckResult{
		Title:  fmt.Sprintf("FIO - Sync Latency (99%%, %s)", path),
		Status: status,
		Value:  latMsStr,
	}
}

func (latencyCheck *LatencyCheck) runFioLatencyTest() (*fio.Result, error) {
	return latencyCheck.fioNewJob(
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
	).Exec()
}
