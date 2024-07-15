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

	// Parse the JSON output
	dockerInfo := &DockerInfo{}
	jsonDecoder := json.NewDecoder(stdout)
	err = jsonDecoder.Decode(dockerInfo)
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
		rawData: stdout,
		info:    dockerInfo,
	}

	return dockerProviderInfo, nil
}

// Container returns a Docker container instance for the given ID or name
func (*DockerProvider) Container(containerID string) Container {
	return &DockerContainer{ContainerID: containerID}
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
