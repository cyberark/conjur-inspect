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
	oldFunc := dockerFunc
	dockerFunc = func(...string) (stdout, stderr []byte, err error) {
		stdout = rawOutput
		return stdout, stderr, err
	}
	defer func() {
		dockerFunc = oldFunc
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
	oldFunc := dockerFunc
	dockerFunc = func(...string) (stdout, stderr []byte, err error) {
		err = testError
		return stdout, stderr, err
	}
	defer func() {
		dockerFunc = oldFunc
	}()

	dockerContainer := &DockerContainer{
		ContainerID: "test-container",
	}

	inspectResult, err := dockerContainer.Inspect()
	assert.Error(t, testError, err)
	assert.Nil(t, inspectResult)
}

func TestDockerContainerExec(t *testing.T) {
	standardOut := []byte("test standard output")
	standardErr := []byte("test standard error")

	// Mock dependencies
	oldFunc := dockerFunc
	dockerFunc = func(...string) (stdout, stderr []byte, err error) {
		stdout = standardOut
		stderr = standardErr
		return stdout, stderr, err
	}
	defer func() {
		dockerFunc = oldFunc
	}()

	dockerContainer := &DockerContainer{
		ContainerID: "test-container",
	}

	execStdout, execStderr, err := dockerContainer.Exec("test")
	assert.NoError(t, err)
	assert.Equal(t, standardOut, execStdout)
	assert.Equal(t, standardErr, execStderr)
}

func TestDockerContainerExecError(t *testing.T) {
	testError := errors.New("fake error")

	// Mock dependencies
	oldFunc := dockerFunc
	dockerFunc = func(...string) (stdout, stderr []byte, err error) {
		err = testError
		return stdout, stderr, err
	}
	defer func() {
		dockerFunc = oldFunc
	}()

	dockerContainer := &DockerContainer{
		ContainerID: "test-container",
	}

	stdout, stderr, err := dockerContainer.Exec("test")
	assert.Error(t, testError, err)
	assert.Nil(t, stdout)
	assert.Nil(t, stderr)
}
