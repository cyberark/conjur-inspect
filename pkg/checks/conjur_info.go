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
func (ci *ConjurInfo) Describe() string {
	return fmt.Sprintf("Conjur Info (%s)", ci.Provider.Name())
}

// Run retrieves and parses the Conjur /info API endpoint
func (ci *ConjurInfo) Run(context *check.RunContext) []check.Result {
	// If there is no container ID, return
	if strings.TrimSpace(context.ContainerID) == "" {
		return []check.Result{}
	}

	container := ci.Provider.Container(context.ContainerID)

	stdout, stderr, err := container.Exec(
		"curl", "-k", "https://localhost/info",
	)

	if err != nil {
		return check.ErrorResult(
			ci,
			fmt.Errorf(
				"failed to collect info data: %w (%s))",
				err,
				strings.TrimSpace(shell.ReadOrDefault(stderr, "N/A")),
			),
		)
	}

	// Read the stdout data to save and parse it
	infoJSONBytes, err := io.ReadAll(stdout)
	if err != nil {
		return check.ErrorResult(
			ci,
			fmt.Errorf("failed to read info data: %w)", err),
		)
	}

	// Save raw info output
	outputFileName := fmt.Sprintf(
		"conjur-info-%s.json",
		strings.ToLower(ci.Provider.Name()),
	)
	_, err = context.OutputStore.Save(
		outputFileName,
		bytes.NewReader(infoJSONBytes),
	)
	if err != nil {
		log.Warn(
			"Failed to save %s Conjur info output: %s",
			ci.Provider.Name(),
			err,
		)
	}

	conjurInfoData := &ConjurInfoData{}
	err = json.Unmarshal(infoJSONBytes, conjurInfoData)
	if err != nil {
		return check.ErrorResult(
			ci,
			fmt.Errorf("failed to parse info data: %w)", err),
		)
	}

	return []check.Result{
		{
			Title:  fmt.Sprintf("Version (%s)", ci.Provider.Name()),
			Status: check.StatusInfo,
			Value:  conjurInfoData.Version,
		},
		{
			Title:  fmt.Sprintf("Release (%s)", ci.Provider.Name()),
			Status: check.StatusInfo,
			Value:  conjurInfoData.Release,
		},
	}
}
