// Package container defines the providers for concrete container engines
// (e.g. Docker, Podman)
package container

import (
	"fmt"
	"strings"

	"github.com/cyberark/conjur-inspect/pkg/shell"
)

// Function variable for dependency injection
var executeDockerInspectFunc = executeDockerInspect

// DockerContainer is a concrete implementation of the Container interface
type DockerContainer struct {
	ContainerID string
}

// ID returns the container ID
func (container *DockerContainer) ID() string {
	return container.ContainerID
}

// Inspect returns the JSON output of the `docker inspect` command
func (container *DockerContainer) Inspect() ([]byte, error) {
	stdout, stderr, err := executeDockerInspectFunc(container.ContainerID)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to inspect Podman container %s: %w (%s)",
			container.ContainerID,
			err,
			strings.TrimSpace(string(stderr)),
		)
	}

	return stdout, nil
}

func executeDockerInspect(containerID string) (stdout, stderr []byte, err error) {
	return shell.NewCommandWrapper(
		"docker",
		"inspect",
		"--format",
		"json",
		"--size",
		containerID,
	).Run()
}
