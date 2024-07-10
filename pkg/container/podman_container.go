// Package container defines the providers for concrete container engines
// (e.g. Docker, Podman)
package container

import (
	"fmt"
	"strings"

	"github.com/cyberark/conjur-inspect/pkg/shell"
)

// Function variable for dependency injection
var executePodmanInspectFunc = executePodmanInspect

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
func (container *PodmanContainer) Inspect() ([]byte, error) {
	stdout, stderr, err := executePodmanInspectFunc(container.ContainerID)
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

func executePodmanInspect(containerID string) (stdout, stderr []byte, err error) {
	return shell.NewCommandWrapper(
		"podman",
		"container",
		"inspect",
		"--format",
		"json",
		"--size",
		containerID,
	).Run()
}
