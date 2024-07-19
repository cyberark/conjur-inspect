// Package checks defines all of the possible Conjur Inspect checks that can
// be run.
package checks

import (
	"fmt"
	"io"
	"strings"

	"github.com/cyberark/conjur-inspect/pkg/check"
	"github.com/cyberark/conjur-inspect/pkg/container"
	"github.com/cyberark/conjur-inspect/pkg/output"
)

// ContainerInspect collects the output of the container runtime's
// inspect API and saves it to the output store.
type ContainerInspect struct {
	Provider container.ContainerProvider
}

// Describe provides a textual description of what this check gathers info on
func (ci *ContainerInspect) Describe() string {
	return fmt.Sprintf("%s inspect", ci.Provider.Name())
}

// Run performs the Docker inspection checks
func (ci *ContainerInspect) Run(context *check.RunContext) []check.Result {
	// If there is no container ID, return
	if strings.TrimSpace(context.ContainerID) == "" {
		return []check.Result{}
	}

	container := ci.Provider.Container(context.ContainerID)

	inspectResult, err := container.Inspect()
	if err != nil {
		return check.ErrorResult(
			ci,
			fmt.Errorf("failed to inspect container: %w", err),
		)
	}

	err = ci.saveOutput(context.OutputStore, inspectResult)
	if err != nil {
		return check.ErrorResult(
			ci,
			fmt.Errorf("failed to save inspect output: %w", err),
		)
	}

	return []check.Result{}
}

func (ci *ContainerInspect) saveOutput(
	outputStore output.Store,
	output io.Reader,
) error {
	outputFileName := fmt.Sprintf(
		"%s-inspect.json",
		strings.ToLower(ci.Provider.Name()),
	)
	_, err := outputStore.Save(outputFileName, output)

	return err
}
