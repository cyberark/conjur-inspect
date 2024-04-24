// Package container defines the providers for concrete container engines
// (e.g. Docker, Podman)
package container

import (
	"errors"
	"testing"

	"github.com/cyberark/conjur-inspect/pkg/check"
	"github.com/stretchr/testify/assert"
)

func TestDockerProviderInfo(t *testing.T) {
	rawOutput := []byte(`{"ServerVersion":"20.10.7","Driver":"overlay2","DockerRootDir":"/var/lib/docker"}`)

	// Mock dependencies
	oldFunc := executeDockerInfoFunc
	executeDockerInfoFunc = func() (stdout, stderr []byte, err error) {
		stdout = rawOutput
		return stdout, stderr, err
	}
	defer func() {
		executeDockerInfoFunc = oldFunc
	}()

	// Get the info
	docker := &DockerProvider{}
	dockerInfo, err := docker.Info()

	assert.NoError(t, err)

	// Check the results
	expected := []check.Result{
		{
			Title:  "Docker Version",
			Status: check.StatusInfo,
			Value:  "20.10.7",
		},
		{
			Title:  "Docker Driver",
			Status: check.StatusInfo,
			Value:  "overlay2",
		},
		{
			Title:  "Docker Root Directory",
			Status: check.StatusInfo,
			Value:  "/var/lib/docker",
		},
	}
	assert.Equal(t, expected, dockerInfo.Results())

	assert.Equal(t, rawOutput, dockerInfo.RawData())
}

func TestDockerProviderInfoParseError(t *testing.T) {
	// Mock dependencies
	oldFunc := executeDockerInfoFunc
	executeDockerInfoFunc = func() (stdout, stderr []byte, err error) {
		stdout = []byte(`invalid json`)
		return stdout, stderr, err
	}
	defer func() {
		executeDockerInfoFunc = oldFunc
	}()

	// Get the info
	docker := &DockerProvider{}
	dockerInfo, err := docker.Info()

	assert.Nil(t, dockerInfo)
	assert.ErrorContains(t, err, "failed to parse Docker info output: ")
}

func TestDockerProviderInfoFailure(t *testing.T) {
	// Mock dependencies
	oldFunc := executeDockerInfoFunc
	executeDockerInfoFunc = func() (stdout, stderr []byte, err error) {
		err = errors.New("fake error")
		return stdout, stderr, err
	}
	defer func() {
		executeDockerInfoFunc = oldFunc
	}()

	// Get the info
	docker := &DockerProvider{}
	dockerInfo, err := docker.Info()

	assert.Nil(t, dockerInfo)

	assert.Error(t, err)
}

func TestDockerProviderInfoServerError(t *testing.T) {
	// Mock dependencies
	oldFunc := executeDockerInfoFunc
	executeDockerInfoFunc = func() (stdout, stderr []byte, err error) {
		stdout = []byte(`{"ServerErrors": ["Test error"]}`)
		return stdout, stderr, err
	}
	defer func() {
		executeDockerInfoFunc = oldFunc
	}()

	// Get the info
	docker := &DockerProvider{}
	dockerInfo, err := docker.Info()

	assert.Nil(t, dockerInfo)

	assert.Error(t, err)
}

func TestDockerProviderContainer(t *testing.T) {
	containerID := "test-container"

	// Get the container
	docker := &DockerProvider{}
	container := docker.Container(containerID)

	// Check the container
	assert.Equal(t, containerID, container.ID())
}
