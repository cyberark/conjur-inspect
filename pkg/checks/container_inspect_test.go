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

func TestContainerInspectRun(t *testing.T) {
	testCheck := &ContainerInspect{
		Provider: &test.ContainerProvider{
			InspectResult: strings.NewReader("test"),
		},
	}

	testOutputStore := test.NewOutputStore()

	results := testCheck.Run(
		&check.RunContext{
			ContainerID: "test",
			OutputStore: testOutputStore,
		},
	)

	assert.Empty(t, results)

	outputStoreItems, err := testOutputStore.Items()
	assert.NoError(t, err)

	assert.Equal(t, 1, len(outputStoreItems))

	itemInfo, err := outputStoreItems[0].Info()
	assert.NoError(t, err)
	assert.Equal(t, "test container provider-inspect.json", itemInfo.Name())

	reader, cleanup, err := outputStoreItems[0].Open()
	assert.NoError(t, err)
	defer cleanup()

	output, err := io.ReadAll(reader)
	assert.NoError(t, err)

	assert.Equal(t, "test", string(output))
}

func TestContainerInspectRunError(t *testing.T) {
	testCheck := &ContainerInspect{
		Provider: &test.ContainerProvider{
			InspectError: errors.New("Test error"),
		},
	}

	results := testCheck.Run(
		&check.RunContext{
			ContainerID: "test",
		},
	)

	assert.Equal(t, 1, len(results))

	assert.Equal(t, "Test Container Provider inspect", results[0].Title)
	assert.Equal(t, check.StatusError, results[0].Status)
	assert.Equal(t, "N/A", results[0].Value)
	assert.Equal(t, "failed to inspect container: Test error", results[0].Message)
}
