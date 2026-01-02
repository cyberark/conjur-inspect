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

// ContainerNetworkInspect collects the network inspection data from the
// container runtime and saves it to the output store.
type ContainerNetworkInspect struct {
	Provider container.ContainerProvider
}

// Describe provides a textual description of what this check gathers info on
func (cni *ContainerNetworkInspect) Describe() string {
	return fmt.Sprintf("%s network inspect", cni.Provider.Name())
}

// Run performs the network inspection check
func (cni *ContainerNetworkInspect) Run(runContext *check.RunContext) []check.Result {
	// Check if the container runtime is available
	runtimeKey := strings.ToLower(cni.Provider.Name())
	if !IsRuntimeAvailable(runContext, runtimeKey) {
		if runContext.VerboseErrors {
			return check.ErrorResult(
				cni,
				fmt.Errorf("container runtime not available"),
			)
		}
		return []check.Result{}
	}

	networkInspectOutput, err := cni.Provider.NetworkInspect()
	if err != nil {
		if runContext.VerboseErrors {
			return check.ErrorResult(
				cni,
				fmt.Errorf("failed to inspect networks: %w", err),
			)
		}
		return []check.Result{}
	}

	// Read the output to save it
	outputBytes, err := io.ReadAll(networkInspectOutput)
	if err != nil {
		if runContext.VerboseErrors {
			return check.ErrorResult(
				cni,
				fmt.Errorf("failed to read network inspect output: %w", err),
			)
		}
		return []check.Result{}
	}

	// Save raw network inspect output
	outputFileName := fmt.Sprintf(
		"%s-network-inspect.json",
		strings.ToLower(cni.Provider.Name()),
	)
	_, err = runContext.OutputStore.Save(outputFileName, strings.NewReader(string(outputBytes)))
	if err != nil {
		log.Warn(
			"Failed to save %s network inspect output: %s",
			cni.Provider.Name(),
			err,
		)
	}

	return []check.Result{}
}
