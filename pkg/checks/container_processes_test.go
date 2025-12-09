package checks

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/cyberark/conjur-inspect/pkg/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContainerProcessesDescribe(t *testing.T) {
	provider := &test.ContainerProvider{}
	cp := &ContainerProcesses{Provider: provider}
	assert.Equal(t, "Container processes (Test Container Provider)", cp.Describe())
}

func TestContainerProcessesRunSuccessful(t *testing.T) {
	// Create a mock container provider with process list
	processOutput := strings.Join([]string{
		"UID        PID  PPID  C STIME TTY          TIME CMD",
		"root         1     0  0 10:00 ?        00:00:01 /bin/sh",
		"root        10     1  0 10:00 ?        00:00:00  \\_ nginx: master process",
		"www-data    11    10  0 10:00 ?        00:00:00      \\_ nginx: worker process",
	}, "\n")

	provider := &test.ContainerProvider{
		ExecResponses: map[string]test.ExecResponse{
			"ps -ef --forest": {
				Stdout: strings.NewReader(processOutput),
				Stderr: strings.NewReader(""),
				Error:  nil,
			},
		},
	}

	cp := &ContainerProcesses{Provider: provider}
	runContext := test.NewRunContext("container123")
	results := cp.Run(&runContext)

	// Should return empty results on success
	assert.Empty(t, results)

	// Verify the file was saved to the output store
	items, err := runContext.OutputStore.Items()
	require.NoError(t, err)
	require.Len(t, items, 1)
	info, err := items[0].Info()
	require.NoError(t, err)
	// Provider name is "Test Container Provider" and gets ToLower() -> "test container provider"
	assert.Equal(t, "test container provider-container-processes.log", info.Name())

	// Verify the content was saved correctly
	outputStoreItemReader, cleanup, err := items[0].Open()
	defer cleanup()
	require.NoError(t, err)

	savedContent, err := io.ReadAll(outputStoreItemReader)
	require.NoError(t, err)
	assert.Equal(t, processOutput, string(savedContent))
}

func TestContainerProcessesEmptyContainerID(t *testing.T) {
	provider := &test.ContainerProvider{}
	cp := &ContainerProcesses{Provider: provider}
	runContext := test.NewRunContext("")
	results := cp.Run(&runContext)

	// Should return empty results when container ID is empty
	assert.Empty(t, results)

	// Verify nothing was saved
	items, err := runContext.OutputStore.Items()
	require.NoError(t, err)
	assert.Len(t, items, 0)
}

func TestContainerProcessesExecError(t *testing.T) {
	provider := &test.ContainerProvider{
		ExecResponses: map[string]test.ExecResponse{
			"ps -ef --forest": {
				Stdout: nil,
				Stderr: strings.NewReader("error: process not found"),
				Error:  fmt.Errorf("exec failed"),
			},
		},
	}

	cp := &ContainerProcesses{Provider: provider}
	runContext := test.NewRunContext("container123")
	results := cp.Run(&runContext)

	// Should return error result
	require.Len(t, results, 1)
	assert.Equal(t, "Container processes (Test Container Provider)", results[0].Title)
	assert.Contains(t, results[0].Message, "failed to retrieve processes from container")
}

func TestContainerProcessesStderr(t *testing.T) {
	processOutput := "UID        PID  PPID  C STIME TTY          TIME CMD\nroot         1     0  0 10:00 ?        00:00:01 /bin/sh"
	stderrOutput := "warning: some non-fatal warning"

	provider := &test.ContainerProvider{
		ExecResponses: map[string]test.ExecResponse{
			"ps -ef --forest": {
				Stdout: strings.NewReader(processOutput),
				Stderr: strings.NewReader(stderrOutput),
				Error:  nil,
			},
		},
	}

	cp := &ContainerProcesses{Provider: provider}
	runContext := test.NewRunContext("container123")
	results := cp.Run(&runContext)

	// Should still succeed with empty results
	assert.Empty(t, results)

	// Verify the output was saved despite stderr
	items, err := runContext.OutputStore.Items()
	require.NoError(t, err)
	require.Len(t, items, 1)
}
