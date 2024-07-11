// Package checks defines all of the possible Conjur Inspect checks that can
// be run.
package checks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cyberark/conjur-inspect/pkg/check"
	"github.com/cyberark/conjur-inspect/pkg/container"
	"github.com/cyberark/conjur-inspect/pkg/log"
)

// ConjurInfo collects the output of the Conjur Info API (/info)
type ConjurInfo struct {
	Provider container.ContainerProvider
}

// ConjurInfoData represents the fields from the Conjur info API's JSON
// response that we need to parse for the report.
type ConjurInfoData struct {
	Version string `json:"version"`
	Release string `json:"release"`
}

// Describe provides a textual description of what this check gathers info on
func (inspect *ConjurInfo) Describe() string {
	return fmt.Sprintf("Conjur Info (%s)", inspect.Provider.Name())
}

// Run retrieves and parses the Conjur /info API endpoint
func (inspect *ConjurInfo) Run(context *check.RunContext) <-chan []check.Result {
	future := make(chan []check.Result)

	go func() {

		// If there is no container ID, return
		if strings.TrimSpace(context.ContainerID) == "" {
			future <- []check.Result{}

			return
		}

		container := inspect.Provider.Container(context.ContainerID)

		stdout, stderr, err := container.Exec(
			"curl", "-k", "https://localhost/info",
		)

		if err != nil {
			future <- []check.Result{
				{
					Title:  fmt.Sprintf("Conjur Info (%s)", inspect.Provider.Name()),
					Status: check.StatusError,
					Value:  "N/A",
					Message: fmt.Sprintf(
						"failed to collect info data: %s (%s))",
						err,
						strings.TrimSpace(string(stderr)),
					),
				},
			}

			return
		}

		// Save raw info output
		outputReader := bytes.NewReader(stdout)
		outputFileName := fmt.Sprintf(
			"conjur-info-%s.json",
			strings.ToLower(inspect.Provider.Name()),
		)
		err = context.OutputStore.Save(outputFileName, outputReader)
		if err != nil {
			log.Warn(
				"Failed to save %s Conjur info output: %s",
				inspect.Provider.Name(),
				err,
			)
		}

		conjurInfoData := &ConjurInfoData{}
		err = json.Unmarshal(stdout, conjurInfoData)
		if err != nil {
			future <- []check.Result{
				{
					Title:   fmt.Sprintf("Conjur Info (%s)", inspect.Provider.Name()),
					Status:  check.StatusError,
					Value:   "N/A",
					Message: fmt.Sprintf("failed to parse info data: %s)", err.Error()),
				},
			}

			return
		}

		future <- []check.Result{
			{
				Title:  fmt.Sprintf("Version (%s)", inspect.Provider.Name()),
				Status: check.StatusInfo,
				Value:  conjurInfoData.Version,
			},
			{
				Title:  fmt.Sprintf("Release (%s)", inspect.Provider.Name()),
				Status: check.StatusInfo,
				Value:  conjurInfoData.Release,
			},
		}
	}() // async

	return future
}
