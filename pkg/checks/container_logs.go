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

// ContainerLogs collects the logs of a given container and saves them to the
// output store.
type ContainerLogs struct {
	Provider container.ContainerProvider
}

// Describe provides a textual description of what this check gathers info on
func (cl *ContainerLogs) Describe() string {
	return fmt.Sprintf("%s logs", cl.Provider.Name())
}

// Run performs the Docker inspection checks
func (cl *ContainerLogs) Run(runContext *check.RunContext) []check.Result {
	// If there is no container ID, return
	if strings.TrimSpace(runContext.ContainerID) == "" {
		return []check.Result{}
	}

	container := cl.Provider.Container(runContext.ContainerID)

	inspectResult, err := container.Logs(runContext.Since)
	if err != nil {
		return check.ErrorResult(
			cl,
			fmt.Errorf("failed to collect container logs: %w", err),
		)
	}

	// Save raw container info output
	outputFileName := fmt.Sprintf(
		"%s-container.log",
		strings.ToLower(cl.Provider.Name()),
	)
	_, err = runContext.OutputStore.Save(outputFileName, inspectResult)
	if err != nil {
		log.Warn(
			"Failed to save %s container logs: %s",
			cl.Provider.Name(),
			err,
		)
	}

	return []check.Result{}
}
