package fio

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testJobArgs = []string{"test_arg_1", "test_arg_2"}

type mockCommand struct {
	OutputBytes []byte
	Error       error
}

func (mockCommand *mockCommand) Run() ([]byte, error) {
	return mockCommand.OutputBytes, mockCommand.Error
}

func TestNewJob(t *testing.T) {
	exec := NewJob("test_job", testJobArgs)
	if exec == nil {
		t.Error("Expected Job, got nil")
	}
}

func TestJob_Exec(t *testing.T) {
	var outputDestination []byte

	expectedOutput := []byte(`{"test_key": "test_value"}`)

	// Create Job with mocked dependencies
	exec := &Job{
		Name: "test_job",
		Args: testJobArgs,

		rawOutputCallback: func(data []byte) {
			outputDestination = data
		},

		execLookPath: func(string) (string, error) {
			return "fio", nil
		},

		newCommandWrapper: func(string, ...string) command {
			return &mockCommand{
				OutputBytes: expectedOutput,
			}
		},

		jsonUnmarshal: func(data []byte, v interface{}) error {
			return json.Unmarshal(data, v)
		},
	}

	result, err := exec.Exec()
	if err != nil {
		t.Errorf("Expected nil error, got %s", err)
	}

	if result == nil {
		t.Error("Expected Result, got nil")
	}

	assert.Equal(t, expectedOutput, outputDestination)
}

func TestJob_Exec_LookupError(t *testing.T) {
	exec := &Job{
		Name: "test_job",
		Args: testJobArgs,

		execLookPath: func(string) (string, error) {
			return "", fmt.Errorf("Test error")
		},
	}

	_, err := exec.Exec()
	assert.ErrorContains(t, err, "unable to find 'fio' path:")
}

func TestJob_Exec_CommandError(t *testing.T) {
	exec := &Job{
		Name: "test_job",
		Args: testJobArgs,

		execLookPath: func(string) (string, error) {
			return "", nil
		},

		newCommandWrapper: func(string, ...string) command {
			return &mockCommand{
				Error: fmt.Errorf("test error"),
			}
		},
	}

	_, err := exec.Exec()
	assert.ErrorContains(t, err, "unable to execute 'fio' job:")
}

func TestJob_Exec_ParseError(t *testing.T) {
	exec := &Job{
		Name: "test_job",
		Args: testJobArgs,

		execLookPath: func(string) (string, error) {
			return "", nil
		},

		newCommandWrapper: func(string, ...string) command {
			return &mockCommand{
				OutputBytes: []byte(`{"test_key": "test_value"}`),
			}
		},

		jsonUnmarshal: func(data []byte, v interface{}) error {
			return fmt.Errorf("test error")
		},
	}

	_, err := exec.Exec()
	assert.ErrorContains(t, err, "unable to parse 'fio' output:")
}

func TestUsingTestDirectory(t *testing.T) {
	cleanup, err := usingTestDirectory("test_dir")
	if err != nil {
		t.Errorf("Expected nil error, got %s", err)
	}
	cleanup()
}
