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

// ContainerProcesses collects the process list from inside a container
type ContainerProcesses struct {
	Provider container.ContainerProvider
}

// Describe provides a textual description of what this check gathers info on
func (cp *ContainerProcesses) Describe() string {
	return fmt.Sprintf("Container processes (%s)", cp.Provider.Name())
}

// Run performs the container process list collection
func (cp *ContainerProcesses) Run(runContext *check.RunContext) []check.Result {
	// If there is no container ID, return
	if strings.TrimSpace(runContext.ContainerID) == "" {
		return []check.Result{}
	}

	containerInstance := cp.Provider.Container(runContext.ContainerID)

	// Execute ps command to get process list with tree view
	stdout, stderr, err := containerInstance.Exec(
		"ps", "-ef", "--forest",
	)
	if err != nil {
		return check.ErrorResult(
			cp,
			fmt.Errorf("failed to retrieve processes from container: %w", err),
		)
	}

	// Read stdout from the command execution
	processBytes, err := io.ReadAll(stdout)
	if err != nil {
		return check.ErrorResult(
			cp,
			fmt.Errorf("failed to read processes output: %w", err),
		)
	}

	// Read any stderr for logging
	stderrBytes, _ := io.ReadAll(stderr)
	if len(stderrBytes) > 0 {
		log.Warn("stderr while reading container processes: %s", string(stderrBytes))
	}

	processOutput := string(processBytes)

	// Save process list to output store
	outputFileName := fmt.Sprintf(
		"%s-container-processes.log",
		strings.ToLower(cp.Provider.Name()),
	)
	_, err = runContext.OutputStore.Save(
		outputFileName,
		strings.NewReader(processOutput),
	)
	if err != nil {
		log.Warn("failed to save container processes output: %s", err)
	}

	// Return empty results - this check only produces raw output
	return []check.Result{}
}
