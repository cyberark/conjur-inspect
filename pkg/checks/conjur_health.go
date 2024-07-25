// Package checks defines all of the possible Conjur Inspect checks that can
// be run.
package checks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/cyberark/conjur-inspect/pkg/check"
	"github.com/cyberark/conjur-inspect/pkg/container"
	"github.com/cyberark/conjur-inspect/pkg/log"
	"github.com/cyberark/conjur-inspect/pkg/shell"
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
func (ch *ConjurHealth) Run(runContext *check.RunContext) []check.Result {
	// If there is no container ID, return
	if strings.TrimSpace(runContext.ContainerID) == "" {
		return []check.Result{}
	}

	container := ch.Provider.Container(runContext.ContainerID)
	stdout, stderr, err := container.Exec(
		"curl", "-k", "https://localhost/health",
	)

	if err != nil {
		return check.ErrorResult(
			ch,
			fmt.Errorf(
				"failed to collect health data: %w (%s))",
				err,
				strings.TrimSpace(shell.ReadOrDefault(stderr, "N/A")),
			),
		)
	}

	// Read the stdout data to save and parse it
	healthJSONBytes, err := io.ReadAll(stdout)
	if err != nil {
		return check.ErrorResult(
			ch,
			fmt.Errorf("failed to read health data: %w)", err),
		)
	}

	// Save raw health output before parsing, in case there are parsing errors
	outputFileName := fmt.Sprintf(
		"conjur-health-%s.json",
		strings.ToLower(ch.Provider.Name()),
	)
	_, err = runContext.OutputStore.Save(
		outputFileName,
		bytes.NewReader(healthJSONBytes),
	)
	if err != nil {
		log.Warn(
			"Failed to save %s Conjur health output: %s",
			ch.Provider.Name(),
			err,
		)
	}

	// Parse the health JSON to return the report data.
	conjurHealthData := &ConjurHealthData{}
	err = json.Unmarshal(healthJSONBytes, conjurHealthData)
	if err != nil {
		return check.ErrorResult(
			ch,
			fmt.Errorf("failed to parse health data: %w)", err),
		)
	}

	return []check.Result{
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
}
