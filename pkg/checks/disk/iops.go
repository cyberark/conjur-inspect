package disk

import (
	"fmt"
	"os"

	"github.com/conjurinc/conjur-preflight/pkg/checks/disk/fio"
	"github.com/conjurinc/conjur-preflight/pkg/framework"
	"github.com/conjurinc/conjur-preflight/pkg/log"
)

const iopsJobName = "conjur-fio-iops"

// IopsCheck is a pre-flight check to report the read and write IOPs for the
// directory in which `conjur-preflight` is run.
type IopsCheck struct {
	// When debug mode is enabled, the IOPs check will write the full fio
	// results to a file.
	debug bool

	// We inject the fio command execution as a dependency that we can swap for
	// unit testing
	fioNewJob func(string, []string) fio.Executable
}

var getWorkingDirectory func() (string, error) = os.Getwd

// NewIopsCheck instantiates an Iops check with the default dependencies
func NewIopsCheck(debug bool) *IopsCheck {
	return &IopsCheck{
		debug: debug,

		// Default dependencies
		fioNewJob: fio.NewJob,
	}
}

// Describe provides a textual description of what this check gathers info on
func (*IopsCheck) Describe() string {
	return "disk IOPs"
}

// Run executes the IopsCheck by running `fio` and processing its output
func (iopsCheck *IopsCheck) Run() <-chan []framework.CheckResult {
	future := make(chan []framework.CheckResult)

	go func() {

		fioResult, err := iopsCheck.runFioIopsTest()

		if err != nil {
			future <- []framework.CheckResult{
				{
					Title:   "FIO IOPs",
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
					Title:   "FIO IOPs",
					Status:  framework.STATUS_ERROR,
					Value:   "N/A",
					Message: "No job results returned by 'fio'",
				},
			}

			return
		}

		future <- []framework.CheckResult{
			fioReadIopsResult(&fioResult.Jobs[0]),
			fioWriteIopsResult(&fioResult.Jobs[0]),
		}
	}() // async

	return future
}

func fioReadIopsResult(job *fio.JobResult) framework.CheckResult {

	// 50 iops min from https://etcd.io/docs/v3.3/op-guide/hardware/
	status := framework.STATUS_INFO
	if job.Read.Iops < 50 {
		status = framework.STATUS_WARN
	}

	// Format title
	path, err := getWorkingDirectory()
	if err != nil {
		log.Debug("Unable to get working directory: %s", err)
		path = "working directory"
	}
	titleStr := fmt.Sprintf("FIO - Read IOPs (%s)", path)

	// Format value
	valueStr := fmt.Sprintf(
		"%0.2f (Min: %d, Max: %d, StdDev: %0.2f)",
		job.Read.Iops,
		job.Read.IopsMin,
		job.Read.IopsMax,
		job.Read.IopsStddev,
	)

	return framework.CheckResult{
		Title:  titleStr,
		Status: status,
		Value:  valueStr,
	}
}

func fioWriteIopsResult(job *fio.JobResult) framework.CheckResult {

	// 50 iops min from https://etcd.io/docs/v3.3/op-guide/hardware/
	status := framework.STATUS_INFO
	if job.Write.Iops < 50 {
		status = framework.STATUS_WARN
	}

	// Format title
	path, err := getWorkingDirectory()
	if err != nil {
		log.Debug("Unable to get working directory: %s", err)
		path = "working directory"
	}
	titleStr := fmt.Sprintf("FIO - Write IOPs (%s)", path)

	// Format value
	valueStr := fmt.Sprintf(
		"%0.2f (Min: %d, Max: %d, StdDev: %0.2f)",
		job.Write.Iops,
		job.Write.IopsMin,
		job.Write.IopsMax,
		job.Write.IopsStddev,
	)

	return framework.CheckResult{
		Title:  titleStr,
		Status: status,
		Value:  valueStr,
	}
}

func (iopsCheck *IopsCheck) runFioIopsTest() (*fio.Result, error) {
	job := iopsCheck.fioNewJob(
		iopsJobName,
		[]string{
			"--filename=conjur-fio-iops/data",
			"--size=100MB",
			"--direct=1",
			"--rw=randrw",
			"--bs=4k",
			"--ioengine=libaio",
			"--iodepth=256",
			"--runtime=10",
			"--numjobs=4",
			"--time_based",
			"--group_reporting",
			"--output-format=json",
			"--name=conjur-fio-iops",
		},
	)

	// In debug mode, we'll write out the raw results from 'fio'
	if iopsCheck.debug {
		job.OnRawOutput(func(data []byte) { writeResultToFile(data, iopsJobName) })
	}

	return job.Exec()
}

func writeResultToFile(buffer []byte, jobName string) {
	outputFilename := fmt.Sprintf("%s.json", jobName)

	err := os.WriteFile(outputFilename, buffer, 0644)

	if err != nil {
		log.Warn("Failed to write result file for %s: %s", jobName, err)
	}
}
