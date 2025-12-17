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
func (cr *ContainerRuntime) Describe() string {
	return fmt.Sprintf("%s runtime", cr.Provider.Name())
}

// Run performs the Docker inspection checks
func (cr *ContainerRuntime) Run(runContext *check.RunContext) []check.Result {
	// Check if the container runtime is available
	runtimeKey := strings.ToLower(cr.Provider.Name())
	if !IsRuntimeAvailable(runContext, runtimeKey) {
		if runContext.VerboseErrors {
			return check.ErrorResult(
				cr,
				fmt.Errorf("container runtime not available"),
			)
		}
		return []check.Result{}
	}

	containerInfo, err := cr.Provider.Info()
	if err != nil {
		return check.ErrorResult(
			cr,
			fmt.Errorf("failed to collect container runtime info: %w", err),
		)
	}

	// Save raw container info output
	outputFileName := fmt.Sprintf(
		"%s-info.json",
		strings.ToLower(cr.Provider.Name()),
	)
	_, err = runContext.OutputStore.Save(outputFileName, containerInfo.RawData())
	if err != nil {
		log.Warn(
			"Failed to save %s info output: %s",
			cr.Provider.Name(),
			err,
		)
	}

	return containerInfo.Results()
}
