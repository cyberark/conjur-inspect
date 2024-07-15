// Package container defines the providers for concrete container engines
// (e.g. Docker, Podman)
package container

import (
	"fmt"
	"io"
	"strings"

	"github.com/cyberark/conjur-inspect/pkg/shell"
)

// Function variable for dependency injection
var podmanFunc = podman

// PodmanContainer is a concrete implementation of the Container interface
// for Podman
type PodmanContainer struct {
	ContainerID string
}

// ID returns the container ID
func (container *PodmanContainer) ID() string {
	return container.ContainerID
}

// Inspect returns the JSON output of the `podman inspect` command
func (container *PodmanContainer) Inspect() (io.Reader, error) {
	stdout, stderr, err := podmanFunc(
		"container",
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
			strings.TrimSpace(shell.ReadOrDefault(stderr, "N/A")),
		)
	}

	return stdout, nil
}

// Exec runs a command inside the container
func (container *PodmanContainer) Exec(
	command ...string,
) (stdout, stderr io.Reader, err error) {
	args := append([]string{"exec", container.ContainerID}, command...)
	return podmanFunc(args...)
}

func podman(
	command ...string,
) (stdout, stderr io.Reader, err error) {
	return shell.NewCommandWrapper("podman", command...).Run()
}
