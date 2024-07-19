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

// ContainerInspect collects the output of the container runtime's
// inspect API and saves it to the output store.
type ContainerInspect struct {
	Provider container.ContainerProvider
}

// Describe provides a textual description of what this check gathers info on
func (inspect *ContainerInspect) Describe() string {
	return fmt.Sprintf("%s inspect", inspect.Provider.Name())
}

// Run performs the Docker inspection checks
func (inspect *ContainerInspect) Run(context *check.RunContext) []check.Result {
	// If there is no container ID, return
	if strings.TrimSpace(context.ContainerID) == "" {
		return []check.Result{}
	}

	container := inspect.Provider.Container(context.ContainerID)

	inspectResult, err := container.Inspect()
	if err != nil {
		return []check.Result{
			{
				Title:   fmt.Sprintf("%s inspect", inspect.Provider.Name()),
				Status:  check.StatusError,
				Value:   "N/A",
				Message: err.Error(),
			},
		}
	}

	// Save raw container info output
	outputFileName := fmt.Sprintf(
		"%s-inspect.json",
		strings.ToLower(inspect.Provider.Name()),
	)
	_, err = context.OutputStore.Save(outputFileName, inspectResult)
	if err != nil {
		log.Warn(
			"Failed to save %s inspect output: %s",
			inspect.Provider.Name(),
			err,
		)
	}

	return []check.Result{}
}
