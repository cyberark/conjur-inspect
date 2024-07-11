package test

import "github.com/cyberark/conjur-inspect/pkg/check"

// NewRunContext returns a test run context pre-configured with an in-memory
// output data store.
func NewRunContext(containerID string) check.RunContext {
	return check.RunContext{
		ContainerID: containerID,
		OutputStore: NewOutputStore(),
	}
}
