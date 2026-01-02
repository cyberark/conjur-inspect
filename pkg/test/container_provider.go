// Package test defines utilities and mock implementations for testing
package test

import (
	"fmt"
	"io"
	"strings"
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

	NetworkInspectError  error
	NetworkInspectResult io.Reader

	ExecResponses       map[string]ExecResponse
	ExecAsUserResponses map[string]ExecResponse

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

	ExecResponses       map[string]ExecResponse
	ExecAsUserResponses map[string]ExecResponse

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

		ExecResponses:       cp.ExecResponses,
		ExecAsUserResponses: cp.ExecAsUserResponses,

		LogsOutput: cp.LogsOutput,
		LogsError:  cp.LogsError,
	}
}

// NetworkInspect returns the mock network inspect output
func (cp *ContainerProvider) NetworkInspect() (io.Reader, error) {
	if cp.NetworkInspectError != nil {
		return nil, cp.NetworkInspectError
	}

	return cp.NetworkInspectResult, nil
}

// ExecResponse allows for mocking responses to multiple exec calls against a
// container.
type ExecResponse struct {
	Error  error
	Stdout io.Reader
	Stderr io.Reader
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

	commandString := strings.Join(command, " ")

	response, exists := c.ExecResponses[commandString]

	// Return an error if there is no configured response for the given command
	if !exists {
		return nil, nil, fmt.Errorf("no exec response for: %s", commandString)
	}

	return response.Stdout, response.Stderr, response.Error
}

// ExecAsUser returns the JSON output of the mock `exec` command as a specific user
func (c *Container) ExecAsUser(
	user string,
	command ...string,
) (stdout, stderr io.Reader, err error) {

	// Build the command string with user prefix for lookup
	commandParts := append([]string{user}, command...)
	commandString := strings.Join(commandParts, " ")

	response, exists := c.ExecAsUserResponses[commandString]

	// Return an error if there is no configured response for the given command
	if !exists {
		return nil, nil, fmt.Errorf("no exec as user response for: %s", commandString)
	}

	return response.Stdout, response.Stderr, response.Error
}

// Logs returns the output of the mock `logs` command
func (c *Container) Logs(since time.Duration) (io.Reader, error) {
	return c.LogsOutput, c.LogsError
}
