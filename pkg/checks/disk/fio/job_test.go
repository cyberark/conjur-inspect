package fio

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testJobArgs = []string{"test_arg_1", "test_arg_2"}

func TestJobExec(t *testing.T) {
	var outputDestination []byte
	expectedOutput := `{"test_key": "test_value"}`

	originalFunc := executeFioFunc
	executeFioFunc = mockExecuteFioFunc(expectedOutput, "", nil)
	defer func() {
		executeFioFunc = originalFunc
	}()

	// Create Job with mocked dependencies
	exec := &Job{
		Name: "test_job",
		Args: testJobArgs,

		rawOutputCallback: func(data []byte) {
			outputDestination = data
		},
	}

	result, err := exec.Exec()

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedOutput, string(outputDestination))
}

func TestJobExecCommandError(t *testing.T) {
	originalFunc := executeFioFunc
	executeFioFunc = mockExecuteFioFunc("", "", fmt.Errorf("test error"))
	defer func() {
		executeFioFunc = originalFunc
	}()

	exec := NewJob("test_job", testJobArgs)

	_, err := exec.Exec()
	assert.ErrorContains(t, err, "unable to execute 'fio' job:")
}

func TestJobExecParseError(t *testing.T) {
	originalFunc := executeFioFunc
	executeFioFunc = mockExecuteFioFunc("invalid JSON", "", nil)
	defer func() {
		executeFioFunc = originalFunc
	}()

	exec := NewJob("test_job", testJobArgs)

	_, err := exec.Exec()
	assert.ErrorContains(t, err, "unable to parse 'fio' output:")
}

func TestUsingJobDirectory(t *testing.T) {
	cleanup, err := usingJobDirectory("test_dir")
	if err != nil {
		t.Errorf("Expected nil error, got %s", err)
	}
	cleanup()
}

func mockExecuteFioFunc(
	stdout string,
	stderr string,
	err error,
) func(args ...string) (io.Reader, io.Reader, error) {
	return func(args ...string) (io.Reader, io.Reader, error) {
		return strings.NewReader(stdout), strings.NewReader(stderr), err
	}
}
