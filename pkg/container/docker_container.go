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
var dockerFunc = docker
var dockerCombinedOutputFunc = dockerCombinedOutput

// DockerContainer is a concrete implementation of the Container interface
type DockerContainer struct {
	ContainerID string
}

// ID returns the container ID
func (dc *DockerContainer) ID() string {
	return dc.ContainerID
}

// Inspect returns the JSON output of the `docker inspect` command
func (dc *DockerContainer) Inspect() (io.Reader, error) {
	stdout, stderr, err := dockerFunc(
		"inspect",
		"--format",
		"json",
		"--size",
		dc.ContainerID,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to inspect Podman container %s: %w (%s)",
			dc.ContainerID,
			err,
			strings.TrimSpace(shell.ReadOrDefault(stderr, "N/A")),
		)
	}

	return stdout, nil
}

// Exec runs a command inside the container
func (dc *DockerContainer) Exec(
	command ...string,
) (stdout, stderr io.Reader, err error) {
	args := append([]string{"exec", dc.ContainerID}, command...)
	return dockerFunc(args...)
}

// ExecAsUser runs a command inside the container as a specific user
func (dc *DockerContainer) ExecAsUser(
	user string,
	command ...string,
) (stdout, stderr io.Reader, err error) {
	args := append([]string{"exec", "--user", user, dc.ContainerID}, command...)
	return dockerFunc(args...)
}

// Logs returns the logs of the container
func (dc *DockerContainer) Logs(since time.Duration) (io.Reader, error) {
	args := []string{"logs", fmt.Sprintf("--since=%s", since), dc.ContainerID}
	return dockerCombinedOutputFunc(args...)
}

func docker(
	command ...string,
) (stdout, stderr io.Reader, err error) {
	return shell.NewCommandWrapper("docker", command...).Run()
}

func dockerCombinedOutput(
	command ...string,
) (io.Reader, error) {
	return shell.NewCommandWrapper("docker", command...).RunCombinedOutput()
}
