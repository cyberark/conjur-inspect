// Package test defines utilities and mock implementations for testing
package test

import (
	"io"
	"time"

	"github.com/cyberark/conjur-inspect/pkg/check"
	"github.com/cyberark/conjur-inspect/pkg/container"
)

// ContainerProvider is a mock implementation of the ContainerProvider interface
// for testing
type ContainerProvider struct {
	InspectError  error
	InspectResult io.Reader

	InfoError   error
	InfoRawData io.Reader
	InfoResults []check.Result

	ExecError  error
	ExecStdout io.Reader
	ExecStderr io.Reader

	LogsOutput io.Reader
	LogsError  error
}

// ContainerProviderInfo is a mock implementation of the ContainerProviderInfo
// interface for testing
type ContainerProviderInfo struct {
	InfoRawData io.Reader
	InfoResults []check.Result
}

// Container is a mock implementation of the Container interface for testing
type Container struct {
	ContainerID string

	InspectError  error
	InspectResult io.Reader

	ExecError  error
	ExecStdout io.Reader
	ExecStderr io.Reader

	LogsOutput io.Reader
	LogsError  error
}

// Name returns the name of the container provider
func (*ContainerProvider) Name() string {
	return "Test Container Provider"
}

// Info returns the container provider info
func (cp *ContainerProvider) Info() (container.ContainerProviderInfo, error) {
	if cp.InfoError != nil {
		return nil, cp.InfoError
	}

	return &ContainerProviderInfo{
		InfoRawData: cp.InfoRawData,
		InfoResults: cp.InfoResults,
	}, nil
}

// Container returns a container instance for the given ID
func (cp *ContainerProvider) Container(
	containerID string,
) container.Container {
	return &Container{
		ContainerID: containerID,

		InspectError:  cp.InspectError,
		InspectResult: cp.InspectResult,

		ExecError:  cp.ExecError,
		ExecStdout: cp.ExecStdout,
		ExecStderr: cp.ExecStderr,

		LogsOutput: cp.LogsOutput,
		LogsError:  cp.LogsError,
	}
}

// Results returns the check results
func (cpi *ContainerProviderInfo) Results() []check.Result {
	return cpi.InfoResults
}

// RawData returns the raw data
func (cpi *ContainerProviderInfo) RawData() io.Reader {
	return cpi.InfoRawData
}

// ID returns the container ID
func (c *Container) ID() string {
	return c.ContainerID
}

// Inspect returns the JSON output of the mock `inspect` command
func (c *Container) Inspect() (io.Reader, error) {
	if c.InspectError != nil {
		return nil, c.InspectError
	}

	return c.InspectResult, nil
}

// Exec returns the JSON output of the mock `exec` command
func (c *Container) Exec(
	command ...string,
) (stdout, stderr io.Reader, err error) {
	return c.ExecStdout, c.ExecStderr, c.ExecError
}

// Logs returns the output of the mock `logs` command
func (c *Container) Logs(since time.Duration) (io.Reader, error) {
	return c.LogsOutput, c.LogsError
}
