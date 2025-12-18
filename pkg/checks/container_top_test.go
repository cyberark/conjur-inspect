package checks

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/cyberark/conjur-inspect/pkg/check"
	"github.com/cyberark/conjur-inspect/pkg/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContainerTopDescribe(t *testing.T) {
	provider := &test.ContainerProvider{}
	ct := &ContainerTop{Provider: provider}
	assert.Equal(t, "Container top (Test Container Provider)", ct.Describe())
}

func TestContainerTopRunSuccessful(t *testing.T) {
	// Create a mock container provider with top output
	topOutput := strings.Join([]string{
		"top - 10:00:00 up 1 day,  2:34,  0 users,  load average: 0.52, 0.58, 0.59",
		"Tasks:   4 total,   1 running,   3 sleeping,   0 stopped,   0 zombie",
		"%Cpu(s):  2.3 us,  1.0 sy,  0.0 ni, 96.5 id,  0.2 wa,  0.0 hi,  0.0 si,  0.0 st",
		"MiB Mem :  15953.8 total,   1234.5 free,   8765.4 used,   5953.9 buff/cache",
		"MiB Swap:   2048.0 total,   2048.0 free,      0.0 used.   6543.2 avail Mem",
		"",
		"    PID USER      PR  NI    VIRT    RES    SHR S  %CPU  %MEM     TIME+ COMMAND",
		"      1 root      20   0  715584  23456  15360 S   0.0   0.1   0:01.23 conjur",
		"     10 conjur    20   0 1234567  98765  43210 S   1.5   0.6   1:23.45 nginx",
		"     11 postgres  20   0 2345678 123456  54321 S   2.3   0.8   2:34.56 postgres",
	}, "\n")

	provider := &test.ContainerProvider{
		ExecResponses: map[string]test.ExecResponse{
			"top -b -c -H -w 512 -n 1": {
				Stdout: strings.NewReader(topOutput),
				Stderr: strings.NewReader(""),
				Error:  nil,
			},
		},
	}

	ct := &ContainerTop{Provider: provider}
	runContext := test.NewRunContext("container123")
	results := ct.Run(&runContext)

	// Should return empty results on success
	assert.Empty(t, results)

	// Verify the file was saved to the output store
	items, err := runContext.OutputStore.Items()
	require.NoError(t, err)
	require.Len(t, items, 1)
	info, err := items[0].Info()
	require.NoError(t, err)
	// Provider name is "Test Container Provider" and gets ToLower() -> "test container provider"
	assert.Equal(t, "test container provider-container-top.log", info.Name())

	// Verify the content was saved correctly
	outputStoreItemReader, cleanup, err := items[0].Open()
	defer cleanup()
	require.NoError(t, err)

	savedContent, err := io.ReadAll(outputStoreItemReader)
	require.NoError(t, err)
	assert.Equal(t, topOutput, string(savedContent))
}

func TestContainerTopEmptyContainerID(t *testing.T) {
	provider := &test.ContainerProvider{}
	ct := &ContainerTop{Provider: provider}
	runContext := test.NewRunContext("")
	results := ct.Run(&runContext)

	// Should return empty results when container ID is empty
	assert.Empty(t, results)

	// Verify nothing was saved
	items, err := runContext.OutputStore.Items()
	require.NoError(t, err)
	assert.Len(t, items, 0)
}

func TestContainerTopWhitespaceContainerID(t *testing.T) {
	provider := &test.ContainerProvider{}
	ct := &ContainerTop{Provider: provider}
	runContext := test.NewRunContext("   \t\n  ")
	results := ct.Run(&runContext)

	// Should return empty results when container ID is only whitespace
	assert.Empty(t, results)

	// Verify nothing was saved
	items, err := runContext.OutputStore.Items()
	require.NoError(t, err)
	assert.Len(t, items, 0)
}

func TestContainerTopRuntimeUnavailable(t *testing.T) {
	provider := &test.ContainerProvider{}
	ct := &ContainerTop{Provider: provider}
	runContext := test.NewRunContext("container123")

	// Simulate runtime unavailability
	runContext.ContainerRuntimeAvailability = map[string]check.RuntimeAvailability{
		"test container provider": {
			Available: false,
			Error:     fmt.Errorf("docker not found"),
		},
	}

	// Without VerboseErrors, should return empty results
	runContext.VerboseErrors = false
	results := ct.Run(&runContext)
	assert.Empty(t, results)

	// Verify nothing was saved
	items, err := runContext.OutputStore.Items()
	require.NoError(t, err)
	assert.Len(t, items, 0)
}

func TestContainerTopRuntimeUnavailableVerboseErrors(t *testing.T) {
	provider := &test.ContainerProvider{}
	ct := &ContainerTop{Provider: provider}
	runContext := test.NewRunContext("container123")

	// Simulate runtime unavailability
	runContext.ContainerRuntimeAvailability = map[string]check.RuntimeAvailability{
		"test container provider": {
			Available: false,
			Error:     fmt.Errorf("docker not found"),
		},
	}

	// With VerboseErrors, should return error result
	runContext.VerboseErrors = true
	results := ct.Run(&runContext)

	require.Len(t, results, 1)
	assert.Equal(t, "Container top (Test Container Provider)", results[0].Title)
	assert.Equal(t, check.StatusError, results[0].Status)
	assert.Contains(t, results[0].Message, "container runtime not available")
}

func TestContainerTopExecError(t *testing.T) {
	provider := &test.ContainerProvider{
		ExecResponses: map[string]test.ExecResponse{
			"top -b -c -H -w 512 -n 1": {
				Stdout: nil,
				Stderr: strings.NewReader("error: command not found"),
				Error:  fmt.Errorf("exec failed"),
			},
		},
	}

	ct := &ContainerTop{Provider: provider}
	runContext := test.NewRunContext("container123")
	results := ct.Run(&runContext)

	// Should return error result
	require.Len(t, results, 1)
	assert.Equal(t, "Container top (Test Container Provider)", results[0].Title)
	assert.Equal(t, check.StatusError, results[0].Status)
	assert.Contains(t, results[0].Message, "failed to retrieve top output from container")
}

func TestContainerTopStderr(t *testing.T) {
	topOutput := "top - 10:00:00 up 1 day,  2:34,  0 users,  load average: 0.52, 0.58, 0.59\nTasks:   4 total,   1 running,   3 sleeping,   0 stopped,   0 zombie"
	stderrOutput := "warning: some non-fatal warning"

	provider := &test.ContainerProvider{
		ExecResponses: map[string]test.ExecResponse{
			"top -b -c -H -w 512 -n 1": {
				Stdout: strings.NewReader(topOutput),
				Stderr: strings.NewReader(stderrOutput),
				Error:  nil,
			},
		},
	}

	ct := &ContainerTop{Provider: provider}
	runContext := test.NewRunContext("container123")
	results := ct.Run(&runContext)

	// Should still succeed with empty results
	assert.Empty(t, results)

	// Verify the output was saved despite stderr
	items, err := runContext.OutputStore.Items()
	require.NoError(t, err)
	require.Len(t, items, 1)
}
