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

// ContainerRuntime collects the information on the version of the
// container runtime on the system
type ContainerRuntime struct {
	Provider container.ContainerProvider
}

// Describe provides a textual description of what this check gathers info on
func (container *ContainerRuntime) Describe() string {
	return fmt.Sprintf("%s runtime", container.Provider.Name())
}

// Run performs the Docker inspection checks
func (container *ContainerRuntime) Run(context *check.RunContext) []check.Result {
	containerInfo, err := container.Provider.Info()
	if err != nil {
		return []check.Result{
			{
				Title:   container.Provider.Name(),
				Status:  check.StatusError,
				Value:   "N/A",
				Message: err.Error(),
			},
		}
	}

	// Save raw container info output
	outputFileName := fmt.Sprintf(
		"%s-info.json",
		strings.ToLower(container.Provider.Name()),
	)
	_, err = context.OutputStore.Save(outputFileName, containerInfo.RawData())
	if err != nil {
		log.Warn(
			"Failed to save %s info output: %s",
			container.Provider.Name(),
			err,
		)
	}

	return containerInfo.Results()
}
