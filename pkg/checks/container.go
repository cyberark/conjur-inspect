// Package checks defines all of the possible Conjur Inspect checks that can
// be run.
package checks

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/cyberark/conjur-inspect/pkg/check"
	"github.com/cyberark/conjur-inspect/pkg/container"
	"github.com/cyberark/conjur-inspect/pkg/log"
)

// Docker collects the information on the version of Docker on the system
type Container struct {
	Provider container.ContainerProvider
}

// Describe provides a textual description of what this check gathers info on
func (container *Container) Describe() string {
	return fmt.Sprintf("%s runtime", container.Provider.Name())
}

// Run performs the Docker inspection checks
func (container *Container) Run(context *check.RunContext) <-chan []check.Result {
	future := make(chan []check.Result)

	go func() {
		containerInfo, err := container.Provider.Info()
		if err != nil {
			future <- []check.Result{
				{
					Title:   container.Provider.Name(),
					Status:  check.StatusError,
					Value:   "N/A",
					Message: err.Error(),
				},
			}

			return
		}

		// Save raw container info output
		outputReader := bytes.NewReader(containerInfo.RawData())
		outputFileName := fmt.Sprintf(
			"%s-info.json",
			strings.ToLower(container.Provider.Name()),
		)
		err = context.OutputStore.Save(outputFileName, outputReader)
		if err != nil {
			log.Warn(
				"Failed to save %s info output: %s",
				container.Provider.Name(),
				err,
			)
		}

		future <- containerInfo.Results()
	}() // async

	return future
}
