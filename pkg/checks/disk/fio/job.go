package fio

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/cyberark/conjur-inspect/pkg/log"
	"github.com/cyberark/conjur-inspect/pkg/shell"
)

const fioExecutable = "fio"

var executeFioFunc func(args ...string) (stdout, stderr io.Reader, err error) = executeFio

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
}

// NewJob constructs a Job with the default dependencies
func NewJob(name string, args []string) Executable {
	return &Job{
		Name: name,
		Args: args,
	}
}

// Exec runs the given fio job in a temporary directory
func (job *Job) Exec() (*Result, error) {
	// Create the directory for running the fio test. We have this return the
	// cleanup method as well to simplify deferring this task when the function
	// finishes.
	cleanup, err := usingJobDirectory(job.Name)
	if err != nil {
		return nil, fmt.Errorf("unable to create test directory: %w", err)
	}
	defer cleanup()

	// Run 'fio' command
	stdout, stderr, err := executeFioFunc(job.Args...)
	if err != nil {
		log.Debug("Unable to execute 'fio' job:")

		// If we can read the standard error from the command, log it
		stderrStr, stderrErr := io.ReadAll(stderr)
		if stderrErr == nil {
			log.Debug(string(stderrStr))
		}

		return nil, fmt.Errorf("unable to execute 'fio' job: %w", err)
	}

	stdoutBytes, err := io.ReadAll(stdout)
	if err != nil {
		return nil, fmt.Errorf("unable to read 'fio' output: %w", err)
	}

	// If there is a configured result listener, notify it of the result output
	if job.rawOutputCallback != nil {
		job.rawOutputCallback(stdoutBytes)
	}

	// Parse the result JSON
	jsonResult := Result{}
	err = json.Unmarshal(stdoutBytes, &jsonResult)
	if err != nil {
		return nil, fmt.Errorf("unable to parse 'fio' output: %w", err)
	}

	return &jsonResult, nil
}

// OnRawOutput sets the callback to receive standard output from the fio
// command.
func (job *Job) OnRawOutput(callback func([]byte)) {
	job.rawOutputCallback = callback
}

func usingJobDirectory(jobName string) (func(), error) {
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

func executeFio(args ...string) (stdout, stderr io.Reader, err error) {
	return shell.NewCommandWrapper(fioExecutable, args...).Run()
}
