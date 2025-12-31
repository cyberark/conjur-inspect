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

func TestDockerContainerInspect(t *testing.T) {
	rawOutput := `{"Test Key":"Test Value"}`

	// Mock dependencies
	oldFunc := dockerFunc
	dockerFunc = func(...string) (stdout, stderr io.Reader, err error) {
		stdout = strings.NewReader(rawOutput)
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

	inspectBytes, err := io.ReadAll(inspectResult)
	assert.NoError(t, err)

	assert.Equal(t, rawOutput, string(inspectBytes))
}

func TestDockerContainerInspectError(t *testing.T) {
	testError := errors.New("fake error")
	// Mock dependencies
	oldFunc := dockerFunc
	dockerFunc = func(...string) (stdout, stderr io.Reader, err error) {
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
	standardOut := "test standard output"
	standardErr := "test standard error"

	// Mock dependencies
	oldFunc := dockerFunc
	dockerFunc = func(...string) (stdout, stderr io.Reader, err error) {
		stdout = strings.NewReader(standardOut)
		stderr = strings.NewReader(standardErr)
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

	stdoutBytes, err := io.ReadAll(execStdout)
	assert.NoError(t, err)
	assert.Equal(t, standardOut, string(stdoutBytes))

	stderrBytes, err := io.ReadAll(execStderr)
	assert.NoError(t, err)
	assert.Equal(t, standardErr, string(stderrBytes))
}

func TestDockerContainerExecError(t *testing.T) {
	testError := errors.New("fake error")

	// Mock dependencies
	oldFunc := dockerFunc
	dockerFunc = func(...string) (stdout, stderr io.Reader, err error) {
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

func TestDockerContainerLogs(t *testing.T) {
	logOut := "test logs"

	// Mock dependencies
	oldFunc := dockerCombinedOutputFunc
	dockerCombinedOutputFunc = func(...string) (output io.Reader, err error) {
		output = strings.NewReader(logOut)
		return output, err
	}
	defer func() {
		dockerCombinedOutputFunc = oldFunc
	}()

	dockerContainer := &DockerContainer{
		ContainerID: "test-container",
	}

	output, err := dockerContainer.Logs(time.Duration(0))
	assert.NoError(t, err)

	outputBytes, err := io.ReadAll(output)
	assert.NoError(t, err)
	assert.Equal(t, logOut, string(outputBytes))
}

func TestDockerContainerLogsError(t *testing.T) {
	testError := errors.New("fake error")

	// Mock dependencies
	oldFunc := dockerCombinedOutputFunc
	dockerCombinedOutputFunc = func(...string) (output io.Reader, err error) {
		err = testError
		return output, err
	}
	defer func() {
		dockerCombinedOutputFunc = oldFunc
	}()

	dockerContainer := &DockerContainer{
		ContainerID: "test-container",
	}

	output, err := dockerContainer.Logs(time.Duration(0))
	assert.Error(t, testError, err)
	assert.Nil(t, output)
}
func TestDockerContainerExecAsUser(t *testing.T) {
	standardOut := "test standard output"
	standardErr := "test standard error"
	capturedArgs := []string{}

	// Mock dependencies
	oldFunc := dockerFunc
	dockerFunc = func(args ...string) (stdout, stderr io.Reader, err error) {
		capturedArgs = args
		stdout = strings.NewReader(standardOut)
		stderr = strings.NewReader(standardErr)
		return stdout, stderr, err
	}
	defer func() {
		dockerFunc = oldFunc
	}()

	dockerContainer := &DockerContainer{
		ContainerID: "test-container",
	}

	execStdout, execStderr, err := dockerContainer.ExecAsUser("conjur", "psql", "-c", "SELECT 1")
	assert.NoError(t, err)

	stdoutBytes, err := io.ReadAll(execStdout)
	assert.NoError(t, err)
	assert.Equal(t, standardOut, string(stdoutBytes))

	stderrBytes, err := io.ReadAll(execStderr)
	assert.NoError(t, err)
	assert.Equal(t, standardErr, string(stderrBytes))

	// Verify that --user flag was passed
	assert.Equal(t, []string{"exec", "--user", "conjur", "test-container", "psql", "-c", "SELECT 1"}, capturedArgs)
}

func TestDockerContainerExecAsUserError(t *testing.T) {
	testError := errors.New("fake error")

	// Mock dependencies
	oldFunc := dockerFunc
	dockerFunc = func(...string) (stdout, stderr io.Reader, err error) {
		err = testError
		return stdout, stderr, err
	}
	defer func() {
		dockerFunc = oldFunc
	}()

	dockerContainer := &DockerContainer{
		ContainerID: "test-container",
	}

	stdout, stderr, err := dockerContainer.ExecAsUser("conjur", "psql", "-c", "SELECT 1")
	assert.Error(t, testError, err)
	assert.Nil(t, stdout)
	assert.Nil(t, stderr)
}
