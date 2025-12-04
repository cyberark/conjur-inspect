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

func TestContainerCommandHistoryDescribe(t *testing.T) {
	provider := &test.ContainerProvider{}
	cch := &ContainerCommandHistory{Provider: provider}
	assert.Equal(t, "Test Container Provider command history", cch.Describe())
}

func TestContainerCommandHistoryRunSuccessful(t *testing.T) {
	// Create a mock container provider with bash history
	bashHistory := strings.Join([]string{
		"ls -la",
		"cd /tmp",
		"echo 'hello world'",
	}, "\n")

	provider := &test.ContainerProvider{
		ExecResponses: map[string]test.ExecResponse{
			"sh -c tail -n 100 /root/.bash_history 2>/dev/null || true": {
				Stdout: strings.NewReader(bashHistory),
				Stderr: strings.NewReader(""),
				Error:  nil,
			},
		},
	}

	cch := &ContainerCommandHistory{Provider: provider}
	runContext := test.NewRunContext("container123")
	results := cch.Run(&runContext)

	// Should return empty results on success
	assert.Empty(t, results)

	// Verify the file was saved to the output store
	items, err := runContext.OutputStore.Items()
	require.NoError(t, err)
	require.Len(t, items, 1)
	info, err := items[0].Info()
	require.NoError(t, err)
	// Provider name is "Test Container Provider" and gets ToLower() -> "test container provider"
	assert.Equal(t, "test container provider-command-history.txt", info.Name())
}

func TestContainerCommandHistorySanitization(t *testing.T) {
	// Create history with sensitive data that should be redacted
	bashHistory := strings.Join([]string{
		"mysql -u root -pMyPassword123",
		"api_key=sk_live_abc123xyz",
		"token=secret_token_value",
		"echo 'normal command'",
	}, "\n")

	provider := &test.ContainerProvider{
		ExecResponses: map[string]test.ExecResponse{
			"sh -c tail -n 100 /root/.bash_history 2>/dev/null || true": {
				Stdout: strings.NewReader(bashHistory),
				Stderr: strings.NewReader(""),
				Error:  nil,
			},
		},
	}

	cch := &ContainerCommandHistory{Provider: provider}
	runContext := test.NewRunContext("container123")
	results := cch.Run(&runContext)

	assert.Empty(t, results)

	// Verify sanitization occurred
	items, err := runContext.OutputStore.Items()
	require.NoError(t, err)
	require.Len(t, items, 1)

	outputStoreItemReader, cleanup, err := items[0].Open()
	defer cleanup()
	require.NoError(t, err)

	savedContent, err := io.ReadAll(outputStoreItemReader)
	require.NoError(t, err)

	// Verify sensitive data was redacted
	assert.NotContains(t, string(savedContent), "sk_live_abc123xyz")
	assert.NotContains(t, string(savedContent), "secret_token_value")
	// Normal command should still be present
	assert.Contains(t, string(savedContent), "echo 'normal command'")
	// The api_key and token patterns should have been redacted with [REDACTED]
	assert.Contains(t, string(savedContent), "[REDACTED]")
}

func TestContainerCommandHistoryEmptyContainerID(t *testing.T) {
	provider := &test.ContainerProvider{}
	cch := &ContainerCommandHistory{Provider: provider}
	runContext := test.NewRunContext("")
	results := cch.Run(&runContext)

	// Should return empty results when container ID is empty
	assert.Empty(t, results)

	// Verify nothing was saved
	items, err := runContext.OutputStore.Items()
	require.NoError(t, err)
	assert.Len(t, items, 0)
}

func TestContainerCommandHistoryExecError(t *testing.T) {
	provider := &test.ContainerProvider{
		ExecResponses: map[string]test.ExecResponse{
			"sh -c tail -n 100 /root/.bash_history 2>/dev/null || true": {
				Stdout: nil,
				Stderr: nil,
				Error:  fmt.Errorf("file not found"),
			},
		},
	}

	cch := &ContainerCommandHistory{Provider: provider}
	runContext := test.NewRunContext("container123")
	results := cch.Run(&runContext)

	// Should return error result
	require.Len(t, results, 1)
	assert.Equal(t, "Test Container Provider command history", results[0].Title)
	assert.Contains(t, results[0].Message, "failed to retrieve command history from container")
}

func TestContainerCommandHistoryNoExecResponse(t *testing.T) {
	// Provider with empty ExecResponses - should error on unknown command
	provider := &test.ContainerProvider{
		ExecResponses: map[string]test.ExecResponse{},
	}

	cch := &ContainerCommandHistory{Provider: provider}
	runContext := test.NewRunContext("container123")
	results := cch.Run(&runContext)

	// Should return error result
	require.Len(t, results, 1)
	assert.Equal(t, "Test Container Provider command history", results[0].Title)
	assert.Contains(t, results[0].Message, "failed to retrieve command history from container")
}

func TestContainerCommandHistoryEmptyHistory(t *testing.T) {
	// Test with empty history file
	provider := &test.ContainerProvider{
		ExecResponses: map[string]test.ExecResponse{
			"sh -c tail -n 100 /root/.bash_history 2>/dev/null || true": {
				Stdout: strings.NewReader(""),
				Stderr: strings.NewReader(""),
				Error:  nil,
			},
		},
	}

	cch := &ContainerCommandHistory{Provider: provider}
	runContext := test.NewRunContext("container123")
	results := cch.Run(&runContext)

	// Should return empty results (empty history is valid)
	assert.Empty(t, results)

	// Verify empty file was saved
	items, err := runContext.OutputStore.Items()
	require.NoError(t, err)
	require.Len(t, items, 1)
}

func TestContainerCommandHistoryStderr(t *testing.T) {
	// Test that stderr is logged but doesn't cause failure
	bashHistory := strings.Join([]string{
		"ls -la",
		"pwd",
	}, "\n")

	provider := &test.ContainerProvider{
		ExecResponses: map[string]test.ExecResponse{
			"sh -c tail -n 100 /root/.bash_history 2>/dev/null || true": {
				Stdout: strings.NewReader(bashHistory),
				Stderr: strings.NewReader("some warning message"),
				Error:  nil,
			},
		},
	}

	cch := &ContainerCommandHistory{Provider: provider}
	runContext := test.NewRunContext("container123")
	results := cch.Run(&runContext)

	// Should still succeed despite stderr
	assert.Empty(t, results)

	// Verify the file was saved
	items, err := runContext.OutputStore.Items()
	require.NoError(t, err)
	require.Len(t, items, 1)
}

func TestContainerCommandHistoryWhitespaceContainerID(t *testing.T) {
	provider := &test.ContainerProvider{}
	cch := &ContainerCommandHistory{Provider: provider}
	runContext := test.NewRunContext("   ")
	results := cch.Run(&runContext)

	// Should return empty results when container ID is only whitespace
	assert.Empty(t, results)

	// Verify nothing was saved
	items, err := runContext.OutputStore.Items()
	require.NoError(t, err)
	assert.Len(t, items, 0)
}
