// Package container defines the providers for concrete container engines
// (e.g. Docker, Podman)
package container

import (
	"errors"
	"io"
	"strings"
	"testing"
	"time"

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

func TestPodmanContainerLogs(t *testing.T) {
	logOut := "test logs"

	// Mock dependencies
	oldFunc := podmanCombinedOutputFunc
	podmanCombinedOutputFunc = func(...string) (output io.Reader, err error) {
		output = strings.NewReader(logOut)
		return output, err
	}
	defer func() {
		podmanCombinedOutputFunc = oldFunc
	}()

	podmanContainer := &PodmanContainer{
		ContainerID: "test-container",
	}

	output, err := podmanContainer.Logs(time.Duration(0))
	assert.NoError(t, err)

	outputBytes, err := io.ReadAll(output)
	assert.NoError(t, err)
	assert.Equal(t, logOut, string(outputBytes))
}

func TestPodmanContainerLogsError(t *testing.T) {
	testError := errors.New("fake error")

	// Mock dependencies
	oldFunc := podmanCombinedOutputFunc
	podmanCombinedOutputFunc = func(...string) (output io.Reader, err error) {
		err = testError
		return output, err
	}
	defer func() {
		podmanCombinedOutputFunc = oldFunc
	}()

	podmanContainer := &PodmanContainer{
		ContainerID: "test-container",
	}

	output, err := podmanContainer.Logs(time.Duration(0))
	assert.Error(t, testError, err)
	assert.Nil(t, output)
}
func TestPodmanContainerExecAsUser(t *testing.T) {
	standardOut := strings.NewReader("test standard output")
	standardErr := strings.NewReader("test standard error")
	capturedArgs := []string{}

	// Mock dependencies
	oldFunc := podmanFunc
	podmanFunc = func(args ...string) (stdout, stderr io.Reader, err error) {
		capturedArgs = args
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

	execStdout, execStderr, err := podmanContainer.ExecAsUser("conjur", "psql", "-c", "SELECT 1")
	assert.NoError(t, err)
	assert.Equal(t, standardOut, execStdout)
	assert.Equal(t, standardErr, execStderr)

	// Verify that --user flag was passed
	assert.Equal(t, []string{"exec", "--user", "conjur", "test-container", "psql", "-c", "SELECT 1"}, capturedArgs)
}

func TestPodmanContainerExecAsUserError(t *testing.T) {
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

	stdout, stderr, err := podmanContainer.ExecAsUser("conjur", "psql", "-c", "SELECT 1")
	assert.Error(t, testError, err)
	assert.Nil(t, stdout)
	assert.Nil(t, stderr)
}
