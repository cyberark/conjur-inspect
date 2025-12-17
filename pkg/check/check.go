package check

import (
	"time"

	"github.com/cyberark/conjur-inspect/pkg/output"
)

// StatusInfo means the result is informational only
const StatusInfo = "INFO"

// StatusPass means the result falls within the production operational requirements
const StatusPass = "PASS"

// StatusWarn means that the system is at risk for production operation
const StatusWarn = "WARN"

// StatusFail means the system is unacceptable for production operation
const StatusFail = "FAIL"

// StatusError means the result could not be obtained
const StatusError = "ERROR"

// Check represent a single operation (API call, external program execution,
// etc.) that returns one or more result.
type Check interface {
	Describe() string
	Run(*RunContext) []Result
}

// RunContext is container of other services available to checks within the
// context of a particular report run.
type RunContext struct {
	OutputStore output.Store

	ContainerID string
	Since       time.Duration

	// ContainerRuntimeAvailability caches the availability status of container runtimes
	// Maps provider names (e.g., "docker", "podman") to availability status and error
	ContainerRuntimeAvailability map[string]RuntimeAvailability

	// VerboseErrors controls whether to report errors for unavailable container runtimes
	VerboseErrors bool
}

// RuntimeAvailability represents the availability status of a container runtime
type RuntimeAvailability struct {
	Available bool
	Error     error // Error encountered when checking availability (e.g., executable not found)
}

// Result is the outcome of a particular check. A check may produce multiple
// results.
type Result struct {
	Title   string `json:"title"`
	Value   string `json:"value"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// ErrorResult returns a single result with an error message.
func ErrorResult(c Check, err error) []Result {
	return []Result{
		{
			Title:   c.Describe(),
			Status:  StatusError,
			Value:   "N/A",
			Message: err.Error(),
		},
	}
}
