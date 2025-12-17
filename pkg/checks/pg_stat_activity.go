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

// PgStatActivity collects the output of PostgreSQL's pg_stat_activity view
type PgStatActivity struct {
	Provider container.ContainerProvider
}

// Describe provides a textual description of what this check gathers info on
func (psa *PgStatActivity) Describe() string {
	return fmt.Sprintf("PostgreSQL pg_stat_activity (%s)", psa.Provider.Name())
}

// Run performs the PostgreSQL pg_stat_activity check
func (psa *PgStatActivity) Run(runContext *check.RunContext) []check.Result {
	// If there is no container ID, return
	if strings.TrimSpace(runContext.ContainerID) == "" {
		return []check.Result{}
	}

	containerInstance := psa.Provider.Container(runContext.ContainerID)

	// Execute psql command to get pg_stat_activity as the conjur user
	stdout, stderr, err := containerInstance.ExecAsUser(
		"conjur",
		"psql",
		"-c",
		"select * from pg_stat_activity",
	)

	if err != nil {
		// Capture stderr for error context if available
		stderrText := "N/A"
		if stderr != nil {
			stderrBytes, _ := io.ReadAll(stderr)
			if len(stderrBytes) > 0 {
				stderrText = string(stderrBytes)
			}
		}

		// Save the error output for reference, then return error result
		if len(stderrText) > 0 && stderrText != "N/A" {
			outputFileName := fmt.Sprintf(
				"%s-pg-stat-activity-error.log",
				strings.ToLower(psa.Provider.Name()),
			)
			_, _ = runContext.OutputStore.Save(
				outputFileName,
				strings.NewReader(fmt.Sprintf("Error: %s\n%s", err, stderrText)),
			)
		}

		return check.ErrorResult(
			psa,
			fmt.Errorf("failed to retrieve pg_stat_activity: %w", err),
		)
	}

	// Read stdout from the command execution
	activityBytes, err := io.ReadAll(stdout)
	if err != nil {
		return check.ErrorResult(
			psa,
			fmt.Errorf("failed to read pg_stat_activity output: %w", err),
		)
	}

	activityOutput := string(activityBytes)

	// Save pg_stat_activity output to file
	outputFileName := "pg_stat_activity.log"
	_, err = runContext.OutputStore.Save(
		outputFileName,
		strings.NewReader(activityOutput),
	)
	if err != nil {
		log.Warn("failed to save pg_stat_activity output: %s", err)
	}

	// Return empty results - this check only produces raw output
	return []check.Result{}
}
