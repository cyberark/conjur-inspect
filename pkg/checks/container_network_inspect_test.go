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

func TestContainerNetworkInspectRun(t *testing.T) {
	rawOutput := `[{"Name":"bridge","Id":"abc123"}]`

	testCheck := &ContainerNetworkInspect{
		Provider: &test.ContainerProvider{
			NetworkInspectResult: strings.NewReader(rawOutput),
		},
	}

	testOutputStore := test.NewOutputStore()

	results := testCheck.Run(
		&check.RunContext{
			OutputStore: testOutputStore,
			ContainerRuntimeAvailability: map[string]check.RuntimeAvailability{
				"test container provider": {Available: true},
			},
		},
	)

	assert.Empty(t, results)

	outputStoreItems, err := testOutputStore.Items()
	assert.NoError(t, err)

	assert.Equal(t, 1, len(outputStoreItems))

	itemInfo, err := outputStoreItems[0].Info()
	assert.NoError(t, err)
	assert.Equal(t, "test container provider-network-inspect.json", itemInfo.Name())

	reader, cleanup, err := outputStoreItems[0].Open()
	assert.NoError(t, err)
	defer cleanup()

	output, err := io.ReadAll(reader)
	assert.NoError(t, err)

	assert.Equal(t, rawOutput, string(output))
}

func TestContainerNetworkInspectRunEmpty(t *testing.T) {
	rawOutput := `[]`

	testCheck := &ContainerNetworkInspect{
		Provider: &test.ContainerProvider{
			NetworkInspectResult: strings.NewReader(rawOutput),
		},
	}

	testOutputStore := test.NewOutputStore()

	results := testCheck.Run(
		&check.RunContext{
			OutputStore: testOutputStore,
			ContainerRuntimeAvailability: map[string]check.RuntimeAvailability{
				"test container provider": {Available: true},
			},
		},
	)

	assert.Empty(t, results)

	outputStoreItems, err := testOutputStore.Items()
	assert.NoError(t, err)

	assert.Equal(t, 1, len(outputStoreItems))

	_, err = outputStoreItems[0].Info()
	assert.NoError(t, err)

	reader, cleanup, err := outputStoreItems[0].Open()
	assert.NoError(t, err)
	defer cleanup()

	output, err := io.ReadAll(reader)
	assert.NoError(t, err)

	assert.Equal(t, rawOutput, string(output))
}

func TestContainerNetworkInspectRunError(t *testing.T) {
	testCheck := &ContainerNetworkInspect{
		Provider: &test.ContainerProvider{
			NetworkInspectError: errors.New("network inspect failed"),
		},
	}

	results := testCheck.Run(
		&check.RunContext{
			ContainerRuntimeAvailability: map[string]check.RuntimeAvailability{
				"test container provider": {Available: true},
			},
		},
	)

	assert.Empty(t, results)
}

func TestContainerNetworkInspectRunErrorVerboseErrors(t *testing.T) {
	testCheck := &ContainerNetworkInspect{
		Provider: &test.ContainerProvider{
			NetworkInspectError: errors.New("network inspect failed"),
		},
	}

	results := testCheck.Run(
		&check.RunContext{
			VerboseErrors: true,
			ContainerRuntimeAvailability: map[string]check.RuntimeAvailability{
				"test container provider": {Available: true},
			},
		},
	)

	assert.Len(t, results, 1)
	assert.Equal(t, check.StatusError, results[0].Status)
	assert.Contains(t, results[0].Message, "network inspect failed")
}

func TestContainerNetworkInspectRuntimeNotAvailable(t *testing.T) {
	testCheck := &ContainerNetworkInspect{
		Provider: &test.ContainerProvider{},
	}

	results := testCheck.Run(
		&check.RunContext{
			ContainerRuntimeAvailability: map[string]check.RuntimeAvailability{
				"test container provider": {Available: false},
			},
		},
	)

	assert.Empty(t, results)
}

func TestContainerNetworkInspectRuntimeNotAvailableVerboseErrors(t *testing.T) {
	testCheck := &ContainerNetworkInspect{
		Provider: &test.ContainerProvider{},
	}

	results := testCheck.Run(
		&check.RunContext{
			VerboseErrors: true,
			ContainerRuntimeAvailability: map[string]check.RuntimeAvailability{
				"test container provider": {Available: false},
			},
		},
	)

	assert.Len(t, results, 1)
	assert.Equal(t, check.StatusError, results[0].Status)
	assert.Contains(t, results[0].Message, "runtime not available")
}
