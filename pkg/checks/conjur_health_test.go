package checks

import (
	"errors"
	"io"
	"testing"

	"github.com/cyberark/conjur-inspect/pkg/check"
	"github.com/cyberark/conjur-inspect/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestConjurHealthRun(t *testing.T) {
	healthJSON := []byte(`{"ok": true, "degraded": false}`)

	containerProvider := &test.ContainerProvider{
		ExecStdout: healthJSON,
	}

	// Create the ConjurHealth instance
	conjurHealth := &ConjurHealth{
		Provider: containerProvider,
	}

	context := test.NewRunContext("test-container-id")

	// Run the function
	results := <-conjurHealth.Run(&context)

	// Check the results
	expectedResults := []check.Result{
		{
			Title:  "Healthy (Test Container Provider)",
			Status: check.StatusInfo,
			Value:  "true",
		},
		{
			Title:  "Degraded (Test Container Provider)",
			Status: check.StatusInfo,
			Value:  "false",
		},
	}
	assert.Equal(t, expectedResults, results)

	// Check the output store
	outputStoreItems, err := context.OutputStore.Items()
	assert.NoError(t, err)
	assert.Len(t, outputStoreItems, 1)

	outputStoreItem := outputStoreItems[0]

	outputStoreItemInfo, err := outputStoreItem.Info()
	assert.NoError(t, err)
	assert.Equal(
		t,
		"conjur-health-test container provider.json",
		outputStoreItemInfo.Name(),
	)

	outputStoreItemReader, cleanup, err := outputStoreItem.Open()
	defer cleanup()
	assert.NoError(t, err)

	outputStoreItemData, err := io.ReadAll(outputStoreItemReader)
	assert.NoError(t, err)
	assert.Equal(t, healthJSON, outputStoreItemData)
}

func TestConjurHealthRun_NoContainerID(t *testing.T) {
	containerProvider := &test.ContainerProvider{}

	// Create the ConjurHealth instance
	conjurHealth := &ConjurHealth{
		Provider: containerProvider,
	}

	context := test.NewRunContext("")

	// Run the function
	results := <-conjurHealth.Run(&context)

	// Check the results
	expectedResults := []check.Result{}
	assert.Equal(t, expectedResults, results)

	// Check the output store
	outputStoreItems, err := context.OutputStore.Items()
	assert.NoError(t, err)
	assert.Empty(t, outputStoreItems)
}

func TestConjurHealthRun_ExecError(t *testing.T) {
	containerProvider := &test.ContainerProvider{
		ExecStderr: []byte("test stderr"),
		ExecError:  errors.New("test error"),
	}

	// Create the ConjurHealth instance
	conjurHealth := &ConjurHealth{
		Provider: containerProvider,
	}

	context := test.NewRunContext("test-container-id")

	// Run the function
	results := <-conjurHealth.Run(&context)

	// Check the results
	expectedResults := []check.Result{
		{
			Title:   "Conjur Health (Test Container Provider)",
			Status:  check.StatusError,
			Value:   "N/A",
			Message: "failed to collect health data: test error (test stderr))",
		},
	}
	assert.Equal(t, expectedResults, results)

	// Check the output store
	outputStoreItems, err := context.OutputStore.Items()
	assert.NoError(t, err)
	assert.Empty(t, outputStoreItems)
}

func TestConjurHealthRun_UnmarshalError(t *testing.T) {
	healthJSON := []byte(`{"ok": "invalid", "degraded": false}`)

	containerProvider := &test.ContainerProvider{
		ExecStdout: healthJSON,
	}

	// Create the ConjurHealth instance
	conjurHealth := &ConjurHealth{
		Provider: containerProvider,
	}

	context := test.NewRunContext("test-container-id")

	// Run the function
	results := <-conjurHealth.Run(&context)

	// Check the results
	expectedResults := []check.Result{
		{
			Title:   "Conjur Health (Test Container Provider)",
			Status:  check.StatusError,
			Value:   "N/A",
			Message: "failed to parse health data: json: cannot unmarshal string into Go struct field ConjurHealthData.ok of type bool)",
		},
	}
	assert.Equal(t, expectedResults, results)

	// Check the output store. The raw output should be saved even with a
	// parse error.
	outputStoreItems, err := context.OutputStore.Items()
	assert.NoError(t, err)
	assert.Len(t, outputStoreItems, 1)

	outputStoreItem := outputStoreItems[0]

	outputStoreItemInfo, err := outputStoreItem.Info()
	assert.NoError(t, err)
	assert.Equal(
		t,
		"conjur-health-test container provider.json",
		outputStoreItemInfo.Name(),
	)

	outputStoreItemReader, cleanup, err := outputStoreItem.Open()
	defer cleanup()
	assert.NoError(t, err)

	outputStoreItemData, err := io.ReadAll(outputStoreItemReader)
	assert.NoError(t, err)
	assert.Equal(t, healthJSON, outputStoreItemData)
}
