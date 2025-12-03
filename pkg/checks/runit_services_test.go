// Package checks defines all of the possible Conjur Inspect checks that can
// be run.
package checks

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/cyberark/conjur-inspect/pkg/check"
	"github.com/cyberark/conjur-inspect/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestRunItServices_Run_EmptyContainerID(t *testing.T) {
	provider := &test.ContainerProvider{}
	rs := &RunItServices{Provider: provider}

	results := rs.Run(&check.RunContext{
		ContainerID: "",
		OutputStore: test.NewOutputStore(),
	})

	assert.Empty(t, results)
}

func TestRunItServices_Run_Success(t *testing.T) {
	testStatusOutput := "run: /etc/sv/conjur: (pid 12345) 1234s\nrun: /etc/sv/postgres: (pid 67890) 5678s\n"

	provider := &test.ContainerProvider{
		ExecResponses: map[string]test.ExecResponse{
			"sh -c sv status /etc/service/*": {
				Stdout: strings.NewReader(testStatusOutput),
				Stderr: strings.NewReader(""),
			},
		},
	}
	rs := &RunItServices{Provider: provider}

	runContext := test.NewRunContext("test-container")
	results := rs.Run(&runContext)

	// Success case returns no results
	assert.Empty(t, results)

	// Verify output was saved
	items, err := runContext.OutputStore.Items()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(items))

	info, err := items[0].Info()
	assert.NoError(t, err)
	assert.Equal(t, "runit-services-status.txt", info.Name())

	reader, cleanup, err := items[0].Open()
	assert.NoError(t, err)
	defer cleanup()

	output, err := io.ReadAll(reader)
	assert.NoError(t, err)
	assert.Equal(t, testStatusOutput, string(output))
}

func TestRunItServices_Run_ExecError(t *testing.T) {
	provider := &test.ContainerProvider{
		ExecResponses: map[string]test.ExecResponse{
			"sh -c sv status /etc/service/*": {
				Stderr: strings.NewReader("command not found"),
				Error:  errors.New("exec failed"),
			},
		},
	}
	rs := &RunItServices{Provider: provider}

	runContext := test.NewRunContext("test-container")
	results := rs.Run(&runContext)

	assert.Equal(t, 1, len(results))
	assert.Equal(t, check.StatusError, results[0].Status)
	assert.Contains(t, results[0].Message, "failed to get runit services status")
	assert.Contains(t, results[0].Message, "command not found")
}

func TestRunItServices_Run_EmptyOutput(t *testing.T) {
	provider := &test.ContainerProvider{
		ExecResponses: map[string]test.ExecResponse{
			"sh -c sv status /etc/service/*": {
				Stdout: strings.NewReader(""),
				Stderr: strings.NewReader(""),
			},
		},
	}
	rs := &RunItServices{Provider: provider}

	runContext := test.NewRunContext("test-container")
	results := rs.Run(&runContext)

	// Success case returns no results
	assert.Empty(t, results)

	// Verify output was saved (even if empty)
	items, err := runContext.OutputStore.Items()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(items))

	reader, cleanup, err := items[0].Open()
	assert.NoError(t, err)
	defer cleanup()

	output, err := io.ReadAll(reader)
	assert.NoError(t, err)
	assert.Equal(t, "", string(output))
}

func TestRunItServices_Describe(t *testing.T) {
	provider := &test.ContainerProvider{}
	rs := &RunItServices{Provider: provider}
	assert.Equal(t, "Runit Services (Test Container Provider)", rs.Describe())
}
