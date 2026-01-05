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
var executePodmanInfoFunc = executePodmanInfo
var executePodmanNetworkInspectFunc = executePodmanNetworkInspect

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
			strings.TrimSpace(shell.ReadOrDefault(stderr, "N/A")),
		)
	}

	// Read the stdout into a byte slice
	stdoutBytes, err := io.ReadAll(stdout)
	if err != nil {
		return nil, fmt.Errorf("failed to read Podman info output: %w", err)
	}

	// Parse the JSON output
	podmanInfo := &PodmanInfo{}
	err = json.Unmarshal(stdoutBytes, podmanInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Podman info output: %w", err)
	}

	podmanProviderInfo := &PodmanProviderInfo{
		rawData: stdoutBytes,
		info:    podmanInfo,
	}

	return podmanProviderInfo, nil
}

// Container returns a Podman container instance for the given ID or name
func (*PodmanProvider) Container(containerID string) Container {
	return &PodmanContainer{ContainerID: containerID}
}

// NetworkInspect returns the JSON output of all Podman networks
func (*PodmanProvider) NetworkInspect() (io.Reader, error) {
	return executePodmanNetworkInspectFunc()
}

func executePodmanInfo() (stdout, stderr io.Reader, err error) {
	return shell.NewCommandWrapper(
		"podman",
		"info",
		"--debug",
		"--format",
		"{{json .}}",
	).Run()
}

func executePodmanNetworkInspect() (io.Reader, error) {
	// First, get the list of network IDs
	stdout, stderr, err := shell.NewCommandWrapper(
		"podman",
		"network",
		"ls",
		"-q",
	).Run()

	if err != nil {
		return nil, fmt.Errorf(
			"failed to list Podman networks: %w (%s)",
			err,
			strings.TrimSpace(shell.ReadOrDefault(stderr, "N/A")),
		)
	}

	// Read the network IDs
	networkIDsBytes, err := io.ReadAll(stdout)
	if err != nil {
		return nil, fmt.Errorf("failed to read Podman network IDs: %w", err)
	}

	networkIDs := strings.TrimSpace(string(networkIDsBytes))

	// If there are no networks, return empty JSON array
	if networkIDs == "" {
		return strings.NewReader("[]"), nil
	}

	// Split network IDs and inspect them all at once
	ids := strings.Fields(networkIDs)
	args := append([]string{"network", "inspect"}, ids...)

	stdout, stderr, err = shell.NewCommandWrapper("podman", args...).Run()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to inspect Podman networks: %w (%s)",
			err,
			strings.TrimSpace(shell.ReadOrDefault(stderr, "N/A")),
		)
	}

	return stdout, nil
}
