// Package container defines the providers for concrete container engines
// (e.g. Docker, Podman)
package container

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cyberark/conjur-inspect/pkg/shell"
)

// Function variable for dependency injection
var executePodmanInfoFunc = executePodmanInfo

// PodmanProvider is a concrete implementation of the
// ContainerProvider interface for Podman
type PodmanProvider struct{}

// Name returns the name of the Podman provider
func (*PodmanProvider) Name() string {
	return "Podman"
}

// Info returns the Podman runtime info
func (*PodmanProvider) Info() (ContainerProviderInfo, error) {
	stdout, stderr, err := executePodmanInfoFunc()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to inspect Podman runtime: %w (%s)",
			err,
			strings.TrimSpace(string(stderr)),
		)
	}

	// Parse the JSON output
	podmanInfo := &PodmanInfo{}
	err = json.Unmarshal(stdout, podmanInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Podman info output: %w", err)
	}

	podmanProviderInfo := &PodmanProviderInfo{
		rawData: stdout,
		info:    podmanInfo,
	}

	return podmanProviderInfo, nil
}

// Container returns a Podman container instance for the given ID or name
func (*PodmanProvider) Container(containerID string) Container {
	return &PodmanContainer{ContainerID: containerID}
}

func executePodmanInfo() (stdout, stderr []byte, err error) {
	return shell.NewCommandWrapper(
		"podman",
		"info",
		"--debug",
		"--format",
		"{{json .}}",
	).Run()
}
