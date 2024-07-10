// Package container defines the providers for concrete container engines
// (e.g. Docker, Podman)
package container

import (
	"fmt"
	"strings"

	"github.com/cyberark/conjur-inspect/pkg/shell"
)

// Function variable for dependency injection
var dockerFunc = docker

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
	stdout, stderr, err := dockerFunc(
		"inspect",
		"--format",
		"json",
		"--size",
		container.ContainerID,
	)

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

// Exec runs a command inside the container
func (container *DockerContainer) Exec(
	command ...string,
) (stdout, stderr []byte, err error) {
	args := append([]string{"exec", container.ContainerID}, command...)
	return dockerFunc(args...)
}

func docker(
	command ...string,
) (stdout, stderr []byte, err error) {
	return shell.NewCommandWrapper("docker", command...).Run()
}
