// Package container defines the providers for concrete container engines
// (e.g. Docker, Podman)
package container

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDockerContainerInspect(t *testing.T) {
	rawOutput := []byte(`{"Test Key":"Test Value"}`)

	// Mock dependencies
	oldFunc := executeDockerInspectFunc
	executeDockerInspectFunc = func(string) (stdout, stderr []byte, err error) {
		stdout = rawOutput
		return stdout, stderr, err
	}
	defer func() {
		executeDockerInspectFunc = oldFunc
	}()

	dockerContainer := &DockerContainer{
		ContainerID: "test-container",
	}

	inspectResult, err := dockerContainer.Inspect()
	assert.NoError(t, err)

	assert.Equal(t, rawOutput, inspectResult)
}

func TestDockerContainerInspectError(t *testing.T) {
	testError := errors.New("fake error")
	// Mock dependencies
	oldFunc := executeDockerInspectFunc
	executeDockerInspectFunc = func(string) (stdout, stderr []byte, err error) {
		err = testError
		return stdout, stderr, err
	}
	defer func() {
		executeDockerInspectFunc = oldFunc
	}()

	dockerContainer := &DockerContainer{
		ContainerID: "test-container",
	}

	inspectResult, err := dockerContainer.Inspect()
	assert.Error(t, testError, err)
	assert.Nil(t, inspectResult)
}
