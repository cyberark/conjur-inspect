// Package container defines the providers for concrete container engines
// (e.g. Docker, Podman)
package container

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/cyberark/conjur-inspect/pkg/shell"
)

// Function variable for dependency injection
var executeDockerInfoFunc = executeDockerInfo
var executeDockerNetworkInspectFunc = executeDockerNetworkInspect

// DockerProvider is a concrete implementation of the
// ContainerProvider interface for Docker
type DockerProvider struct {
}

// Name returns the name of the Docker provider
func (*DockerProvider) Name() string {
	return "Docker"
}

// Info returns the Docker runtime info
func (*DockerProvider) Info() (ContainerProviderInfo, error) {
	stdout, stderr, err := executeDockerInfoFunc()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to inspect Docker runtime: %w (%s)",
			err,
			strings.TrimSpace(shell.ReadOrDefault(stderr, "N/A")),
		)
	}

	// Read the stdout into a byte slice
	stdoutBytes, err := io.ReadAll(stdout)
	if err != nil {
		return nil, fmt.Errorf("failed to read Docker info output: %w", err)
	}

	// Parse the JSON output
	dockerInfo := &DockerInfo{}
	err = json.Unmarshal(stdoutBytes, dockerInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Docker info output: %w", err)
	}

	// Check for server errors
	if len(dockerInfo.ServerErrors) > 0 {
		return nil, fmt.Errorf(
			"Docker runtime has server errors: %s",
			strings.Join(dockerInfo.ServerErrors, ", "),
		)
	}

	dockerProviderInfo := &DockerProviderInfo{
		rawData: stdoutBytes,
		info:    dockerInfo,
	}

	return dockerProviderInfo, nil
}

// Container returns a Docker container instance for the given ID or name
func (*DockerProvider) Container(containerID string) Container {
	return &DockerContainer{ContainerID: containerID}
}

// NetworkInspect returns the JSON output of all Docker networks
func (*DockerProvider) NetworkInspect() (io.Reader, error) {
	return executeDockerNetworkInspectFunc()
}

func executeDockerInfo() (stdout, stderr io.Reader, err error) {
	return shell.NewCommandWrapper(
		"docker",
		"--debug",
		"info",
		"--format",
		"{{json .}}",
	).Run()
}

func executeDockerNetworkInspect() (io.Reader, error) {
	// First, get the list of network IDs
	stdout, stderr, err := shell.NewCommandWrapper(
		"docker",
		"network",
		"ls",
		"-q",
	).Run()

	if err != nil {
		return nil, fmt.Errorf(
			"failed to list Docker networks: %w (%s)",
			err,
			strings.TrimSpace(shell.ReadOrDefault(stderr, "N/A")),
		)
	}

	// Read the network IDs
	networkIDsBytes, err := io.ReadAll(stdout)
	if err != nil {
		return nil, fmt.Errorf("failed to read Docker network IDs: %w", err)
	}

	networkIDs := strings.TrimSpace(string(networkIDsBytes))

	// If there are no networks, return empty JSON array
	if networkIDs == "" {
		return strings.NewReader("[]"), nil
	}

	// Split network IDs and inspect them all at once
	ids := strings.Fields(networkIDs)
	args := append([]string{"network", "inspect"}, ids...)

	stdout, stderr, err = shell.NewCommandWrapper("docker", args...).Run()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to inspect Docker networks: %w (%s)",
			err,
			strings.TrimSpace(shell.ReadOrDefault(stderr, "N/A")),
		)
	}

	return stdout, nil
}
