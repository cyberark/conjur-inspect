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

func TestConjurInfoRun(t *testing.T) {
	infoJSON := `{"version": "1.2.3", "release": "4.5.6"}`

	containerProvider := &test.ContainerProvider{
		ExecStdout: strings.NewReader(infoJSON),
	}

	// Create the ConjurInfo instance
	conjurInfo := &ConjurInfo{
		Provider: containerProvider,
	}

	runContext := test.NewRunContext("test-container-id")

	// Run the function
	results := conjurInfo.Run(&runContext)

	// Check the results
	expectedResults := []check.Result{
		{
			Title:  "Version (Test Container Provider)",
			Status: check.StatusInfo,
			Value:  "1.2.3",
		},
		{
			Title:  "Release (Test Container Provider)",
			Status: check.StatusInfo,
			Value:  "4.5.6",
		},
	}
	assert.Equal(t, expectedResults, results)

	// Check the output store
	outputStoreItems, err := runContext.OutputStore.Items()
	assert.NoError(t, err)
	assert.Len(t, outputStoreItems, 1)

	outputStoreItem := outputStoreItems[0]

	outputStoreItemInfo, err := outputStoreItem.Info()
	assert.NoError(t, err)
	assert.Equal(
		t,
		"conjur-info-test container provider.json",
		outputStoreItemInfo.Name(),
	)

	outputStoreItemReader, cleanup, err := outputStoreItem.Open()
	defer cleanup()
	assert.NoError(t, err)

	outputStoreItemData, err := io.ReadAll(outputStoreItemReader)
	assert.NoError(t, err)
	assert.Equal(t, infoJSON, string(outputStoreItemData))
}

func TestConjurInfoRun_NoContainerID(t *testing.T) {
	containerProvider := &test.ContainerProvider{}

	// Create the ConjurHealth instance
	conjurHealth := &ConjurInfo{
		Provider: containerProvider,
	}

	runContext := test.NewRunContext("")

	// Run the function
	results := conjurHealth.Run(&runContext)

	// Check the results
	expectedResults := []check.Result{}
	assert.Equal(t, expectedResults, results)

	// Check the output store
	outputStoreItems, err := runContext.OutputStore.Items()
	assert.NoError(t, err)
	assert.Empty(t, outputStoreItems)
}

func TestConjurInfoRun_ExecError(t *testing.T) {
	containerProvider := &test.ContainerProvider{
		ExecStderr: strings.NewReader("test stderr"),
		ExecError:  errors.New("test error"),
	}

	// Create the ConjurHealth instance
	conjurHealth := &ConjurInfo{
		Provider: containerProvider,
	}

	runContext := test.NewRunContext("test-container-id")

	// Run the function
	results := conjurHealth.Run(&runContext)

	// Check the results
	expectedResults := []check.Result{
		{
			Title:   "Conjur Info (Test Container Provider)",
			Status:  check.StatusError,
			Value:   "N/A",
			Message: "failed to collect info data: test error (test stderr))",
		},
	}
	assert.Equal(t, expectedResults, results)

	// Check the output store
	outputStoreItems, err := runContext.OutputStore.Items()
	assert.NoError(t, err)
	assert.Empty(t, outputStoreItems)
}

func TestConjurInfoRun_UnmarshalError(t *testing.T) {
	infoJSON := `{"version": 1, "release": false}`

	containerProvider := &test.ContainerProvider{
		ExecStdout: strings.NewReader(infoJSON),
	}

	// Create the ConjurHealth instance
	conjurInfo := &ConjurInfo{
		Provider: containerProvider,
	}

	runContext := test.NewRunContext("test-container-id")

	// Run the function
	results := conjurInfo.Run(&runContext)

	// Check the results
	expectedResults := []check.Result{
		{
			Title:   "Conjur Info (Test Container Provider)",
			Status:  check.StatusError,
			Value:   "N/A",
			Message: "failed to parse info data: json: cannot unmarshal number into Go struct field ConjurInfoData.version of type string)",
		},
	}
	assert.Equal(t, expectedResults, results)

	// Check the output store. The raw output should be saved even with a
	// parse error.
	outputStoreItems, err := runContext.OutputStore.Items()
	assert.NoError(t, err)
	assert.Len(t, outputStoreItems, 1)

	outputStoreItem := outputStoreItems[0]

	outputStoreItemInfo, err := outputStoreItem.Info()
	assert.NoError(t, err)
	assert.Equal(
		t,
		"conjur-info-test container provider.json",
		outputStoreItemInfo.Name(),
	)

	outputStoreItemReader, cleanup, err := outputStoreItem.Open()
	defer cleanup()
	assert.NoError(t, err)

	outputStoreItemData, err := io.ReadAll(outputStoreItemReader)
	assert.NoError(t, err)
	assert.Equal(t, infoJSON, string(outputStoreItemData))
}
