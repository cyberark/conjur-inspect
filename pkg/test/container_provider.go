// Package test defines utilities and mock implementations for testing
package test

import (
	"io"

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
}

// Name returns the name of the container provider
func (provider *ContainerProvider) Name() string {
	return "Test Container Provider"
}

// Info returns the container provider info
func (provider *ContainerProvider) Info() (container.ContainerProviderInfo, error) {
	if provider.InfoError != nil {
		return nil, provider.InfoError
	}

	return &ContainerProviderInfo{
		InfoRawData: provider.InfoRawData,
		InfoResults: provider.InfoResults,
	}, nil
}

// Container returns a container instance for the given ID
func (provider *ContainerProvider) Container(
	containerID string,
) container.Container {
	return &Container{
		ContainerID: containerID,

		InspectError:  provider.InspectError,
		InspectResult: provider.InspectResult,

		ExecError:  provider.ExecError,
		ExecStdout: provider.ExecStdout,
		ExecStderr: provider.ExecStderr,
	}
}

// Results returns the check results
func (providerInfo *ContainerProviderInfo) Results() []check.Result {
	return providerInfo.InfoResults
}

// RawData returns the raw data
func (providerInfo *ContainerProviderInfo) RawData() io.Reader {
	return providerInfo.InfoRawData
}

// ID returns the container ID
func (container *Container) ID() string {
	return container.ContainerID
}

// Inspect returns the JSON output of the mock `inspect` command
func (container *Container) Inspect() (io.Reader, error) {
	if container.InspectError != nil {
		return nil, container.InspectError
	}

	return container.InspectResult, nil
}

// Exec returns the JSON output of the mock `exec` command
func (container *Container) Exec(
	command ...string,
) (stdout, stderr io.Reader, err error) {
	return container.ExecStdout, container.ExecStderr, container.ExecError
}
