// Package container defines the providers for concrete container engines
// (e.g. Docker, Podman)
package container

import (
	"fmt"
	"testing"

	"github.com/cyberark/conjur-inspect/pkg/check"
	"github.com/stretchr/testify/assert"
)

// Define a mock version of executePodmanInfo that returns some expected output
func mockExecutePodmanInfo() (stdout, stderr []byte, err error) {
	stdout = []byte(`{"version": {"version": "2.2.1"}, "store": {"graphDriverName": "overlay", "graphRoot": "/var/lib/containers/storage", "runRoot": "/run/user/0", "volumePath": "/var/lib/containers/storage/volumes"}}`)
	return stdout, stderr, nil
}

func TestPodmanProviderInfo(t *testing.T) {
	// Mock executePodmanInfoFunc to return expected output
	originalFunc := executePodmanInfoFunc
	executePodmanInfoFunc = mockExecutePodmanInfo
	defer func() {
		executeDockerInfoFunc = originalFunc
	}()

	// Get the info
	podman := &PodmanProvider{}
	podmanInfo, err := podman.Info()

	assert.NoError(t, err)

	// Check the results
	expected := []check.Result{
		{
			Title:  "Podman Version",
			Status: check.StatusInfo,
			Value:  "2.2.1",
		},
		{
			Title:  "Podman Driver",
			Status: check.StatusInfo,
			Value:  "overlay",
		},
		{
			Title:  "Podman Graph Root",
			Status: check.StatusInfo,
			Value:  "/var/lib/containers/storage",
		},
		{
			Title:  "Podman Run Root",
			Status: check.StatusInfo,
			Value:  "/run/user/0",
		},
		{
			Title:  "Podman Volume Path",
			Status: check.StatusInfo,
			Value:  "/var/lib/containers/storage/volumes",
		},
	}
	assert.Equal(t, expected, podmanInfo.Results())
}

func TestPodmanProviderInfoError(t *testing.T) {
	// Mock executePodmanInfoFunc to return an error
	originalFunc := executePodmanInfoFunc
	executePodmanInfoFunc = func() (stdout, stderr []byte, err error) {
		return stdout, stderr, fmt.Errorf("fake error")
	}
	defer func() {
		executePodmanInfoFunc = originalFunc
	}()

	// Get the info
	podman := &PodmanProvider{}
	podmanInfo, err := podman.Info()

	assert.Error(t, err)
	assert.Nil(t, podmanInfo)
}
