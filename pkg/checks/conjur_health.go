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

// ConjurHealth collects the output of Conjur's health API (/health)
type ConjurHealth struct {
	Provider container.ContainerProvider
}

// ConjurHealthData represents the fields from the Conjur health API's JSON
// response that we need to parse for the report.
type ConjurHealthData struct {
	OK       bool `json:"ok"`
	Degraded bool `json:"degraded"`
}

// Describe provides a textual description of what this check gathers info on
func (ch *ConjurHealth) Describe() string {
	return fmt.Sprintf("Conjur Health (%s)", ch.Provider.Name())
}

// Run performs the Conjur health check
func (ch *ConjurHealth) Run(context *check.RunContext) <-chan []check.Result {
	future := make(chan []check.Result)

	go func() {

		// If there is no container ID, return
		if strings.TrimSpace(context.ContainerID) == "" {
			future <- []check.Result{}

			return
		}

		container := ch.Provider.Container(context.ContainerID)

		stdout, stderr, err := container.Exec(
			"curl", "-k", "https://localhost/health",
		)

		if err != nil {
			future <- []check.Result{
				{
					Title:  fmt.Sprintf("Conjur Health (%s)", ch.Provider.Name()),
					Status: check.StatusError,
					Value:  "N/A",
					Message: fmt.Sprintf(
						"failed to collect health data: %s (%s))",
						err,
						strings.TrimSpace(string(stderr)),
					),
				},
			}

			return
		}

		// Save raw health output before parsing, in case there are parsing errors
		outputReader := bytes.NewReader(stdout)
		outputFileName := fmt.Sprintf(
			"conjur-health-%s.json",
			strings.ToLower(ch.Provider.Name()),
		)
		err = context.OutputStore.Save(outputFileName, outputReader)
		if err != nil {
			log.Warn(
				"Failed to save %s Conjur health output: %s",
				ch.Provider.Name(),
				err,
			)
		}

		// Parse the health JSON to return the report data
		conjurHealthData := &ConjurHealthData{}
		err = json.Unmarshal(stdout, conjurHealthData)
		if err != nil {
			future <- []check.Result{
				{
					Title:   fmt.Sprintf("Conjur Health (%s)", ch.Provider.Name()),
					Status:  check.StatusError,
					Value:   "N/A",
					Message: fmt.Sprintf("failed to parse health data: %s)", err.Error()),
				},
			}

			return
		}

		future <- []check.Result{
			{
				Title:  fmt.Sprintf("Healthy (%s)", ch.Provider.Name()),
				Status: check.StatusInfo,
				Value:  fmt.Sprintf("%t", conjurHealthData.OK),
			},
			{
				Title:  fmt.Sprintf("Degraded (%s)", ch.Provider.Name()),
				Status: check.StatusInfo,
				Value:  fmt.Sprintf("%t", conjurHealthData.Degraded),
			},
		}
	}() // async

	return future
}
