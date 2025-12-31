// Package container defines the providers for concrete container engines
// (e.g. Docker, Podman)
package container

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/cyberark/conjur-inspect/pkg/shell"
)

// Function variable for dependency injection
var podmanFunc = podman
var podmanCombinedOutputFunc = podmanCombinedOutput

// PodmanContainer is a concrete implementation of the Container interface
// for Podman
type PodmanContainer struct {
	ContainerID string
}

// ID returns the container ID
func (pc *PodmanContainer) ID() string {
	return pc.ContainerID
}

// Inspect returns the JSON output of the `podman inspect` command
func (pc *PodmanContainer) Inspect() (io.Reader, error) {
	stdout, stderr, err := podmanFunc(
		"container",
		"inspect",
		"--format",
		"json",
		"--size",
		pc.ContainerID,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to inspect Podman container %s: %w (%s)",
			pc.ContainerID,
			err,
			strings.TrimSpace(shell.ReadOrDefault(stderr, "N/A")),
		)
	}

	return stdout, nil
}

// Exec runs a command inside the container
func (pc *PodmanContainer) Exec(
	command ...string,
) (stdout, stderr io.Reader, err error) {
	args := append([]string{"exec", pc.ContainerID}, command...)
	return podmanFunc(args...)
}

// ExecAsUser runs a command inside the container as a specific user
func (pc *PodmanContainer) ExecAsUser(
	user string,
	command ...string,
) (stdout, stderr io.Reader, err error) {
	args := append([]string{"exec", "--user", user, pc.ContainerID}, command...)
	return podmanFunc(args...)
}

// Logs returns the logs of the container
func (pc *PodmanContainer) Logs(since time.Duration) (io.Reader, error) {
	args := []string{"logs", fmt.Sprintf("--since=%s", since), pc.ContainerID}
	return podmanCombinedOutputFunc(args...)
}

func podman(
	command ...string,
) (stdout, stderr io.Reader, err error) {
	return shell.NewCommandWrapper("podman", command...).Run()
}

func podmanCombinedOutput(
	command ...string,
) (io.Reader, error) {
	return shell.NewCommandWrapper("podman", command...).RunCombinedOutput()
}
