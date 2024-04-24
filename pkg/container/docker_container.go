// Package container defines the providers for concrete container engines
// (e.g. Docker, Podman)
package container

// DockerContainer is a concrete implementation of the Container interface
type DockerContainer struct {
	ContainerID string
}

// ID returns the container ID
func (container *DockerContainer) ID() string {
	return container.ContainerID
}
