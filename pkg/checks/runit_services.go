// Package checks defines all of the possible Conjur Inspect checks that can
// be run.
package checks

import (
	"fmt"
	"strings"

	"github.com/cyberark/conjur-inspect/pkg/check"
	"github.com/cyberark/conjur-inspect/pkg/container"
	"github.com/cyberark/conjur-inspect/pkg/log"
)

// RunItServices collects the status of runit services running in the container
type RunItServices struct {
	Provider container.ContainerProvider
}

// Describe provides a textual description of what this check gathers info on
func (rs *RunItServices) Describe() string {
	return fmt.Sprintf("Runit Services (%s)", rs.Provider.Name())
}

// Run performs the runit services status check
func (rs *RunItServices) Run(runContext *check.RunContext) []check.Result {
	// If there is no container ID, return
	if strings.TrimSpace(runContext.ContainerID) == "" {
		return []check.Result{}
	}

	// Check if the container runtime is available
	runtimeKey := strings.ToLower(rs.Provider.Name())
	if !IsRuntimeAvailable(runContext, runtimeKey) {
		if runContext.VerboseErrors {
			return check.ErrorResult(
				rs,
				fmt.Errorf("container runtime not available"),
			)
		}
		return []check.Result{}
	}

	container := rs.Provider.Container(runContext.ContainerID)

	// Execute sv status command with shell globbing
	stdout, stderr, err := container.Exec(
		"sh", "-c", "sv status /etc/service/*",
	)

	// Read the stderr first to capture any error output
	stderrData, stderrErr := readAllFunc(stderr)
	if stderrErr != nil {
		return check.ErrorResult(
			rs,
			fmt.Errorf("failed to read runit services status error output: %w", stderrErr),
		)
	}

	if err != nil {
		return check.ErrorResult(
			rs,
			fmt.Errorf("failed to get runit services status: %w (%s)", err, strings.TrimSpace(string(stderrData))),
		)
	}

	// Read the output
	outputData, err := readAllFunc(stdout)
	if err != nil {
		return check.ErrorResult(
			rs,
			fmt.Errorf("failed to read runit services status output: %w", err),
		)
	}

	if len(stderrData) > 0 {
		log.Warn("runit services status stderr: %s", string(stderrData))
	}

	// Save raw output to archive
	_, err = runContext.OutputStore.Save(
		"runit-services-status.txt",
		strings.NewReader(string(outputData)),
	)
	if err != nil {
		log.Warn("Failed to save runit services status output: %s", err)
	}

	// Return empty results on success (output is saved)
	return []check.Result{}
}
