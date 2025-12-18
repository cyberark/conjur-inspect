// Package checks defines all of the possible Conjur Inspect checks that can
// be run.
package checks

import (
	"fmt"
	"io"
	"strings"

	"github.com/cyberark/conjur-inspect/pkg/check"
	"github.com/cyberark/conjur-inspect/pkg/container"
	"github.com/cyberark/conjur-inspect/pkg/log"
)

// ContainerTop collects resource usage information from inside a container using top
type ContainerTop struct {
	Provider container.ContainerProvider
}

// Describe provides a textual description of what this check gathers info on
func (ct *ContainerTop) Describe() string {
	return fmt.Sprintf("Container top (%s)", ct.Provider.Name())
}

// Run performs the container top resource usage collection
func (ct *ContainerTop) Run(runContext *check.RunContext) []check.Result {
	// If there is no container ID, return
	if strings.TrimSpace(runContext.ContainerID) == "" {
		return []check.Result{}
	}

	// Check if the container runtime is available
	runtimeKey := strings.ToLower(ct.Provider.Name())
	if !IsRuntimeAvailable(runContext, runtimeKey) {
		if runContext.VerboseErrors {
			return check.ErrorResult(
				ct,
				fmt.Errorf("container runtime not available"),
			)
		}
		return []check.Result{}
	}

	containerInstance := ct.Provider.Container(runContext.ContainerID)

	// Execute top command to get resource usage snapshot
	// -b flag: batch mode (non-interactive)
	// -n 2: run for 2 iterations
	stdout, stderr, err := containerInstance.Exec(
		"top", "-b", "-c", "-H", "-w", "512", "-n", "1",
	)
	if err != nil {
		return check.ErrorResult(
			ct,
			fmt.Errorf("failed to retrieve top output from container: %w", err),
		)
	}

	// Read stdout from the command execution
	topBytes, err := io.ReadAll(stdout)
	if err != nil {
		return check.ErrorResult(
			ct,
			fmt.Errorf("failed to read top output: %w", err),
		)
	}

	// Read any stderr for logging
	stderrBytes, _ := io.ReadAll(stderr)
	if len(stderrBytes) > 0 {
		log.Warn("stderr while reading container top: %s", string(stderrBytes))
	}

	topOutput := string(topBytes)

	// Save top output to output store
	outputFileName := fmt.Sprintf(
		"%s-container-top.log",
		strings.ToLower(ct.Provider.Name()),
	)
	_, err = runContext.OutputStore.Save(
		outputFileName,
		strings.NewReader(topOutput),
	)
	if err != nil {
		log.Warn("failed to save container top output: %s", err)
	}

	// Return empty results - this check only produces raw output
	return []check.Result{}
}
