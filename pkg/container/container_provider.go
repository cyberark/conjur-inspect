// Package container defines the providers for concrete container engines
// (e.g. Docker, Podman)
package container

import (
	"io"
	"time"

	"github.com/cyberark/conjur-inspect/pkg/check"
)

// ContainerProvider is an interface for a concrete container
// engine (e.g. Docker, Podman)
type ContainerProvider interface {
	Name() string
	Info() (ContainerProviderInfo, error)
	Container(containerID string) Container
	NetworkInspect() (io.Reader, error)
}

// Container is an interface for a container instance
type Container interface {
	ID() string
	Inspect() (io.Reader, error)
	Exec(command ...string) (stdout, stderr io.Reader, err error)
	ExecAsUser(user string, command ...string) (stdout, stderr io.Reader, err error)
	Logs(since time.Duration) (io.Reader, error)
}

// ContainerProviderInfo is an interface for the results of
// gathering the container runtime info, but as the raw
// data, and specific reporting results for that runtime
type ContainerProviderInfo interface {
	Results() []check.Result
	RawData() io.Reader
}
