package checks

import (
	"fmt"
	"testing"

	"github.com/cyberark/conjur-inspect/pkg/check"
	"github.com/cyberark/conjur-inspect/pkg/test"
	"github.com/stretchr/testify/assert"
)

// Define a mock version of executePodmanInfo that returns some expected output
func mockExecutePodmanInfo() (stdout, stderr []byte, err error) {
	stdout = []byte(`{"version": {"version": "2.2.1"}, "store": {"graphDriverName": "overlay", "graphRoot": "/var/lib/containers/storage", "runRoot": "/run/user/0", "volumePath": "/var/lib/containers/storage/volumes"}}`)
	return stdout, stderr, nil
}

func TestPodmanRun(t *testing.T) {
	// Mock executePodmanInfoFunc to return expected output
	originalFunc := executePodmanInfoFunc
	executePodmanInfoFunc = mockExecutePodmanInfo
	defer func() {
		executeDockerInfoFunc = originalFunc
	}()

	// Run the check
	podman := &Podman{}
	context := test.NewRunContext()
	results := podman.Run(&context)

	// Wait for the async results and check them using assert
	for _, res := range <-results {
		switch res.Title {
		case "Podman Version":
			assert.Equal(t, check.StatusInfo, res.Status)
			assert.Equal(t, "2.2.1", res.Value)
			assert.Empty(t, res.Message)
		case "Podman Driver":
			assert.Equal(t, check.StatusInfo, res.Status)
			assert.Equal(t, "overlay", res.Value)
			assert.Empty(t, res.Message)
		case "Podman Graph Root":
			assert.Equal(t, check.StatusInfo, res.Status)
			assert.Equal(t, "/var/lib/containers/storage", res.Value)
			assert.Empty(t, res.Message)
		case "Podman Run Root":
			assert.Equal(t, check.StatusInfo, res.Status)
			assert.Equal(t, "/run/user/0", res.Value)
			assert.Empty(t, res.Message)
		case "Podman Volume Path":
			assert.Equal(t, check.StatusInfo, res.Status)
			assert.Equal(t, "/var/lib/containers/storage/volumes", res.Value)
			assert.Empty(t, res.Message)
		default:
			assert.Failf(t, "unexpected result title", "got %q", res.Title)
		}
	}
}

func TestPodmanRunErrors(t *testing.T) {
	// Mock executePodmanInfoFunc to return an error
	originalFunc := executePodmanInfoFunc
	executePodmanInfoFunc = func() (stdout, stderr []byte, err error) {
		return stdout, stderr, fmt.Errorf("fake error")
	}
	defer func() {
		executeDockerInfoFunc = originalFunc
	}()

	// Run the check
	podman := &Podman{}
	context := test.NewRunContext()
	results := podman.Run(&context)

	// Wait for the async results and check the error result using assert
	for _, res := range <-results {
		assert.Equal(t, "Podman", res.Title)
		assert.Equal(t, check.StatusError, res.Status)
		assert.Equal(t, "N/A", res.Value)
		assert.Contains(t, res.Message, "fake error")
	}
}
