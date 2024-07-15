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

func TestContainerRuntimeRun(t *testing.T) {
	testCheck := &ContainerRuntime{
		Provider: &test.ContainerProvider{
			InfoRawData: strings.NewReader("test info"),
			InfoResults: []check.Result{
				{
					Title:   "Test Container Runtime Check",
					Status:  check.StatusInfo,
					Value:   "Test value",
					Message: "Test message",
				},
			},
		},
	}

	testOutputStore := test.NewOutputStore()

	resultChan := testCheck.Run(
		&check.RunContext{
			OutputStore: testOutputStore,
		},
	)
	results := <-resultChan

	assert.Equal(t, 1, len(results))
	assert.Equal(t, "Test Container Runtime Check", results[0].Title)
	assert.Equal(t, check.StatusInfo, results[0].Status)
	assert.Equal(t, "Test value", results[0].Value)
	assert.Equal(t, "Test message", results[0].Message)

	outputStoreItems, err := testOutputStore.Items()
	assert.NoError(t, err)

	assert.Equal(t, 1, len(outputStoreItems))

	reader, cleanup, err := outputStoreItems[0].Open()
	assert.NoError(t, err)
	defer cleanup()

	output, err := io.ReadAll(reader)
	assert.NoError(t, err)

	assert.Equal(t, "test info", string(output))
}

func TestContainerRuntimeRunError(t *testing.T) {
	testCheck := &ContainerRuntime{
		Provider: &test.ContainerProvider{
			InfoError: errors.New("Test error"),
		},
	}

	resultChan := testCheck.Run(&check.RunContext{})
	results := <-resultChan

	assert.Equal(t, 1, len(results))

	assert.Equal(t, "Test Container Provider", results[0].Title)
	assert.Equal(t, check.StatusError, results[0].Status)
	assert.Equal(t, "N/A", results[0].Value)
	assert.Equal(t, "Test error", results[0].Message)
}
