package test

import "github.com/cyberark/conjur-inspect/pkg/check"

// NewRunContext returns a test run context pre-configured with an in-memory
// output data store.
func NewRunContext() check.RunContext {
	return check.RunContext{
		OutputStore: NewOutputStore(),
	}
}
