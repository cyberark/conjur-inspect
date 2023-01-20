package fio

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/conjurinc/conjur-preflight/pkg/log"
)

const fioExecutable = "fio"

// Executable represents an operation that can produce an fio result and
// emit raw output data.
type Executable interface {
	Exec() (*Result, error)
	OnRawOutput(func([]byte))
}

// Job is a concrete Executable for executing fio jobs
type Job struct {
	// Required fields:
	// ----------------------

	Name string
	Args []string

	// Optional fields:
	// ----------------------

	// OnRawOutput may be configured to receive the full standard output
	// of the fio command. For example, to write the full output to a file.
	rawOutputCallback func([]byte)

	// Injected dependencies:
	// ----------------------

	// Lookup function to return the full path for a command name
	execLookPath      func(string) (string, error)
	newCommandWrapper func(string, ...string) command
	jsonUnmarshal     func([]byte, any) error
}

// NewJob constructs a Job with the default dependencies
func NewJob(name string, args []string) Executable {
	return &Job{
		// Set required fields
		Name: name,
		Args: args,

		// Construct default dependencies
		execLookPath:      exec.LookPath,
		newCommandWrapper: newCommandWrapper,
		jsonUnmarshal:     json.Unmarshal,
	}
}

// Exec runs the given fio job in a temporary directory
func (job *Job) Exec() (*Result, error) {
	// Create the directory for running the fio test. We have this return the
	// cleanup method as well to simplify deferring this task when the function
	// finishes.
	cleanup, err := usingTestDirectory(job.Name)
	if err != nil {
		return nil, fmt.Errorf("unable to create test directory: %s", err)
	}
	defer cleanup()

	// Lookup full path for 'fio'
	fioPath, err := job.execLookPath(fioExecutable)
	if err != nil {
		return nil, fmt.Errorf("unable to find 'fio' path: %s", err)
	}

	// Run 'fio' command
	commandWrapper := job.newCommandWrapper(fioPath, job.Args...)
	output, err := commandWrapper.Run()

	if err != nil {
		return nil, fmt.Errorf("unable to execute 'fio' job: %s", err)
	}

	// If there is a configured result listener, notify it of the result output
	if job.rawOutputCallback != nil {
		job.rawOutputCallback(output)
	}

	// Parse the result JSON
	jsonResult := Result{}
	err = job.jsonUnmarshal(output, &jsonResult)
	if err != nil {
		return nil, fmt.Errorf("unable to parse 'fio' output: %s", err)
	}

	return &jsonResult, nil
}

// OnRawOutput sets the callback to receive standard output from the fio
// command.
func (job *Job) OnRawOutput(callback func([]byte)) {
	job.rawOutputCallback = callback
}

func usingTestDirectory(jobName string) (func(), error) {
	err := os.MkdirAll(jobName, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return func() {
		err := os.RemoveAll(jobName)
		if err != nil {
			log.Warn("Unable to clean up test directory for job: %s", jobName)
		}
	}, nil
}
