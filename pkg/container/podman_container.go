// Package container defines the providers for concrete container engines
// (e.g. Docker, Podman)
package container

// PodmanContainer is a concrete implementation of the Container interface
// for Podman
type PodmanContainer struct {
	ContainerID string
}

// ID returns the container ID
func (container *PodmanContainer) ID() string {
	return container.ContainerID
}
