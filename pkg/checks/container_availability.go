package checks

import (
	"fmt"
	"os/exec"

	"github.com/cyberark/conjur-inspect/pkg/check"
	"github.com/cyberark/conjur-inspect/pkg/log"
)

// ContainerAvailability checks for the availability of container runtimes
// (Docker and Podman) and caches the results in the RunContext to prevent
// duplicate error messages for unavailable runtimes.
type ContainerAvailability struct{}

// Describe provides a textual description of what this check gathers info on
func (ca *ContainerAvailability) Describe() string {
	return "Container runtime availability"
}

// Run checks the availability of Docker and Podman runtimes
func (ca *ContainerAvailability) Run(runContext *check.RunContext) []check.Result {
	// If not already initialized, this shouldn't happen but be safe
	if runContext.ContainerRuntimeAvailability == nil {
		runContext.ContainerRuntimeAvailability = make(map[string]check.RuntimeAvailability)
	}

	// Check Docker availability
	dockerAvailability := ca.checkRuntimeAvailability("docker")
	runContext.ContainerRuntimeAvailability["docker"] = dockerAvailability

	// Check Podman availability
	podmanAvailability := ca.checkRuntimeAvailability("podman")
	runContext.ContainerRuntimeAvailability["podman"] = podmanAvailability

	results := []check.Result{}

	// Log availability for debugging
	if dockerAvailability.Available {
		log.Debug("Docker runtime is available")
	} else {
		log.Debug("Docker runtime is not available: %v", dockerAvailability.Error)
	}

	if podmanAvailability.Available {
		log.Debug("Podman runtime is available")
	} else {
		log.Debug("Podman runtime is not available: %v", podmanAvailability.Error)
	}

	// Only return results in verbose mode or if no runtimes are available
	if runContext.VerboseErrors {
		if !dockerAvailability.Available {
			results = append(results, check.Result{
				Title:   "Docker availability",
				Status:  check.StatusWarn,
				Value:   "N/A",
				Message: fmt.Sprintf("Docker is not available: %v", dockerAvailability.Error),
			})
		}

		if !podmanAvailability.Available {
			results = append(results, check.Result{
				Title:   "Podman availability",
				Status:  check.StatusWarn,
				Value:   "N/A",
				Message: fmt.Sprintf("Podman is not available: %v", podmanAvailability.Error),
			})
		}
	} else if !dockerAvailability.Available && !podmanAvailability.Available {
		// Warn if no container runtimes are available
		results = append(results, check.Result{
			Title:   "Container runtimes",
			Status:  check.StatusWarn,
			Value:   "N/A",
			Message: "No container runtimes (Docker or Podman) are available. Container-related checks will be skipped.",
		})
	}

	return results
}

// checkRuntimeAvailability checks if a runtime executable is available
func (ca *ContainerAvailability) checkRuntimeAvailability(runtimeName string) check.RuntimeAvailability {
	_, err := exec.LookPath(runtimeName)
	if err != nil {
		return check.RuntimeAvailability{
			Available: false,
			Error:     err,
		}
	}
	return check.RuntimeAvailability{
		Available: true,
		Error:     nil,
	}
}

// IsRuntimeAvailable is a helper function to check if a runtime is available
// from any check
func IsRuntimeAvailable(runContext *check.RunContext, runtimeName string) bool {
	if runContext.ContainerRuntimeAvailability == nil {
		return true // Assume available if cache not initialized
	}
	availability, exists := runContext.ContainerRuntimeAvailability[runtimeName]
	if !exists {
		return true // Assume available if not cached
	}
	return availability.Available
}
