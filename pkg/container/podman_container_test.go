// Package container defines the providers for concrete container engines
// (e.g. Docker, Podman)
package container

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPodmanContainerInspect(t *testing.T) {
	rawOutput := strings.NewReader(`{"Test Key":"Test Value"}`)

	// Mock dependencies
	oldFunc := podmanFunc
	podmanFunc = func(...string) (stdout, stderr io.Reader, err error) {
		stdout = rawOutput
		return stdout, stderr, err
	}
	defer func() {
		podmanFunc = oldFunc
	}()

	podmanContainer := &PodmanContainer{
		ContainerID: "test-container",
	}

	inspectResult, err := podmanContainer.Inspect()
	assert.NoError(t, err)

	assert.Equal(t, rawOutput, inspectResult)
}

func TestPodmanContainerInspectError(t *testing.T) {
	testError := errors.New("fake error")
	// Mock dependencies
	oldFunc := podmanFunc
	podmanFunc = func(...string) (stdout, stderr io.Reader, err error) {
		err = testError
		return stdout, stderr, err
	}
	defer func() {
		podmanFunc = oldFunc
	}()

	podmanContainer := &PodmanContainer{
		ContainerID: "test-container",
	}

	inspectResult, err := podmanContainer.Inspect()
	assert.Error(t, testError, err)
	assert.Nil(t, inspectResult)
}

func TestPodmanContainerExec(t *testing.T) {
	standardOut := strings.NewReader("test standard output")
	standardErr := strings.NewReader("test standard error")

	// Mock dependencies
	oldFunc := podmanFunc
	podmanFunc = func(...string) (stdout, stderr io.Reader, err error) {
		stdout = standardOut
		stderr = standardErr
		return stdout, stderr, err
	}
	defer func() {
		podmanFunc = oldFunc
	}()

	podmanContainer := &PodmanContainer{
		ContainerID: "test-container",
	}

	execStdout, execStderr, err := podmanContainer.Exec("test")
	assert.NoError(t, err)
	assert.Equal(t, standardOut, execStdout)
	assert.Equal(t, standardErr, execStderr)
}

func TestPodmanContainerExecError(t *testing.T) {
	testError := errors.New("fake error")

	// Mock dependencies
	oldFunc := podmanFunc
	podmanFunc = func(...string) (stdout, stderr io.Reader, err error) {
		err = testError
		return stdout, stderr, err
	}
	defer func() {
		podmanFunc = oldFunc
	}()

	podmanContainer := &PodmanContainer{
		ContainerID: "test-container",
	}

	stdout, stderr, err := podmanContainer.Exec("test")
	assert.Error(t, testError, err)
	assert.Nil(t, stdout)
	assert.Nil(t, stderr)
}
