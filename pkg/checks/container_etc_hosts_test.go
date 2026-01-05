package checks

import (
	"fmt"
	"strings"
	"testing"

	"github.com/cyberark/conjur-inspect/pkg/check"
	"github.com/cyberark/conjur-inspect/pkg/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContainerEtcHostsDescribe(t *testing.T) {
	provider := &test.ContainerProvider{}
	ceh := &ContainerEtcHosts{Provider: provider}
	assert.Equal(t, "Test Container Provider /etc/hosts", ceh.Describe())
}

func TestContainerEtcHostsRunSuccessful(t *testing.T) {
	hostsContent := "127.0.0.1\tlocalhost\n::1\tlocalhost\n172.17.0.2\tconjur\n"

	provider := &test.ContainerProvider{
		ExecResponses: map[string]test.ExecResponse{
			"cat /etc/hosts": {
				Stdout: strings.NewReader(hostsContent),
				Stderr: strings.NewReader(""),
				Error:  nil,
			},
		},
	}

	ceh := &ContainerEtcHosts{Provider: provider}
	runContext := test.NewRunContext("container123")

	results := ceh.Run(&runContext)

	// Should return empty results on success
	assert.Empty(t, results)

	// Verify the file was saved to the output store
	items, err := runContext.OutputStore.Items()
	require.NoError(t, err)
	require.Len(t, items, 1)
	info, err := items[0].Info()
	require.NoError(t, err)
	// Test provider name is "Test Container Provider" -> lowercase -> "test container provider"
	assert.Equal(t, "test container provider-etc-hosts.txt", info.Name())
}

func TestContainerEtcHostsEmptyContainerID(t *testing.T) {
	provider := &test.ContainerProvider{}

	ceh := &ContainerEtcHosts{Provider: provider}
	runContext := test.NewRunContext("")
	results := ceh.Run(&runContext)

	// Should return empty results when no container ID
	assert.Empty(t, results)

	// Verify no output was saved
	items, err := runContext.OutputStore.Items()
	require.NoError(t, err)
	require.Len(t, items, 0)
}

func TestContainerEtcHostsRuntimeNotAvailable(t *testing.T) {
	provider := &test.ContainerProvider{}

	ceh := &ContainerEtcHosts{Provider: provider}
	runContext := test.NewRunContext("container123")
	
	// Set runtime as not available
	runContext.ContainerRuntimeAvailability = map[string]check.RuntimeAvailability{
		"test container provider": {
			Available: false,
			Error:     fmt.Errorf("runtime not found"),
		},
	}

	results := ceh.Run(&runContext)

	// Should return empty results when runtime not available
	assert.Empty(t, results)

	// Verify no output was saved
	items, err := runContext.OutputStore.Items()
	require.NoError(t, err)
	require.Len(t, items, 0)
}

func TestContainerEtcHostsRuntimeNotAvailableWithVerboseErrors(t *testing.T) {
	provider := &test.ContainerProvider{}

	ceh := &ContainerEtcHosts{Provider: provider}
	runContext := test.NewRunContext("container123")
	runContext.VerboseErrors = true
	
	// Set runtime as not available
	runContext.ContainerRuntimeAvailability = map[string]check.RuntimeAvailability{
		"test container provider": {
			Available: false,
			Error:     fmt.Errorf("runtime not found"),
		},
	}

	results := ceh.Run(&runContext)

	// Should return error result with verbose errors enabled
	require.Len(t, results, 1)
	assert.Equal(t, check.StatusError, results[0].Status)
	assert.Contains(t, results[0].Message, "container runtime not available")
}

func TestContainerEtcHostsExecError(t *testing.T) {
	provider := &test.ContainerProvider{
		ExecResponses: map[string]test.ExecResponse{
			"cat /etc/hosts": {
				Stdout: strings.NewReader(""),
				Stderr: strings.NewReader("cat: /etc/hosts: Permission denied"),
				Error:  fmt.Errorf("exit status 1"),
			},
		},
	}

	ceh := &ContainerEtcHosts{Provider: provider}
	runContext := test.NewRunContext("container123")

	results := ceh.Run(&runContext)

	// Should return empty results without verbose errors
	assert.Empty(t, results)

	// Verify no output was saved
	items, err := runContext.OutputStore.Items()
	require.NoError(t, err)
	require.Len(t, items, 0)
}

func TestContainerEtcHostsExecErrorWithVerboseErrors(t *testing.T) {
	provider := &test.ContainerProvider{
		ExecResponses: map[string]test.ExecResponse{
			"cat /etc/hosts": {
				Stdout: strings.NewReader(""),
				Stderr: strings.NewReader("cat: /etc/hosts: Permission denied"),
				Error:  fmt.Errorf("exit status 1"),
			},
		},
	}

	ceh := &ContainerEtcHosts{Provider: provider}
	runContext := test.NewRunContext("container123")
	runContext.VerboseErrors = true

	results := ceh.Run(&runContext)

	// Should return error result with verbose errors enabled
	require.Len(t, results, 1)
	assert.Equal(t, check.StatusError, results[0].Status)
	assert.Contains(t, results[0].Message, "failed to read /etc/hosts")
	assert.Contains(t, results[0].Message, "Permission denied")
}
