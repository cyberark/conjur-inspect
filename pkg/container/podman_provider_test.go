// Package container defines the providers for concrete container engines
// (e.g. Docker, Podman)
package container

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/cyberark/conjur-inspect/pkg/check"
	"github.com/stretchr/testify/assert"
)

func TestPodmanProviderInfo(t *testing.T) {
	rawOutput := []byte(
		`{"version": {"version": "2.2.1"}, "store": {"graphDriverName": "overlay", "graphRoot": "/var/lib/containers/storage", "runRoot": "/run/user/0", "volumePath": "/var/lib/containers/storage/volumes"}}`,
	)

	// Mock executePodmanInfoFunc to return expected output
	originalFunc := executePodmanInfoFunc
	executePodmanInfoFunc = func() (stdout, stderr io.Reader, err error) {
		stdout = bytes.NewReader(rawOutput)
		return stdout, stderr, err
	}
	defer func() {
		executePodmanInfoFunc = originalFunc
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

	// Check the raw data
	infoOutputBytes, err := io.ReadAll(podmanInfo.RawData())
	assert.NoError(t, err)
	assert.Equal(t, rawOutput, infoOutputBytes)
}

func TestPodmanProviderInfoParseError(t *testing.T) {
	// Mock dependencies
	oldFunc := executePodmanInfoFunc
	executePodmanInfoFunc = func() (stdout, stderr io.Reader, err error) {
		stdout = strings.NewReader(`invalid json`)
		return stdout, stderr, err
	}
	defer func() {
		executePodmanInfoFunc = oldFunc
	}()

	// Get the info
	podman := &PodmanProvider{}
	podmanInfo, err := podman.Info()

	assert.Nil(t, podmanInfo)
	assert.ErrorContains(t, err, "failed to parse Podman info output: ")
}

func TestPodmanProviderInfoError(t *testing.T) {
	// Mock executePodmanInfoFunc to return an error
	originalFunc := executePodmanInfoFunc
	executePodmanInfoFunc = func() (stdout, stderr io.Reader, err error) {
		return stdout, stderr, errors.New("fake error")
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

func TestPodmanProviderContainer(t *testing.T) {
	containerID := "test-container"

	// Get the container
	podman := &PodmanProvider{}
	container := podman.Container(containerID)

	// Check the container
	assert.Equal(t, containerID, container.ID())
}
