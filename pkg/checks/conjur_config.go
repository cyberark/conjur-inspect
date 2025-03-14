// Package checks defines all of the possible Conjur Inspect checks that can
// be run.
package checks

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/cyberark/conjur-inspect/pkg/check"
	"github.com/cyberark/conjur-inspect/pkg/container"
	"github.com/cyberark/conjur-inspect/pkg/log"
	"github.com/cyberark/conjur-inspect/pkg/shell"
)

// Alias io.ReadAll as a variable so that we can stub it out for unit tests
var readAllFunc = io.ReadAll

// ConjurConfig collects the contents of Conjur's config files
type ConjurConfig struct {
	Provider container.ContainerProvider
}

// Describe provides a textual description of what this check gathers info on
func (cc *ConjurConfig) Describe() string {
	return fmt.Sprintf("Conjur Config (%s)", cc.Provider.Name())
}

// Run performs the Conjur configuration check
func (cc *ConjurConfig) Run(runContext *check.RunContext) []check.Result {
	// If there is no container ID, return
	if strings.TrimSpace(runContext.ContainerID) == "" {
		return []check.Result{}
	}

	container := cc.Provider.Container(runContext.ContainerID)

	configPaths := []string{
		"/etc/conjur/config/conjur.yml",
		"/opt/conjur/etc/conjur.conf",
		"/opt/conjur/etc/possum.conf",
		"/opt/conjur/etc/ui.conf",
		"/opt/conjur/etc/cluster.conf",

		// There are two possible locations for the Chef solo configuration file
		"/etc/cinc/solo.json",
		"/etc/chef/solo.json",

		"/etc/postgresql/15/main/postgresql.conf",
	}

	results := []check.Result{}

	// For each path in config Paths
	for _, path := range configPaths {
		result := cc.collectConfigFile(path, container, runContext)

		if result != nil {
			results = append(results, *result)
		}
	}

	return results
}

func (cc *ConjurConfig) collectConfigFile(
	path string,
	container container.Container,
	runContext *check.RunContext,
) *check.Result {
	stdout, stderr, err := container.Exec(
		"cat", path,
	)

	if err != nil {
		return &check.Result{
			Title:   cc.Describe(),
			Status:  check.StatusError,
			Value:   "N/A",
			Message: fmt.Sprintf(
				"failed to collect '%s' : %s (%s))",
				path,
				err,
				strings.TrimSpace(shell.ReadOrDefault(stderr, "N/A")),
			),
		}
	}

	// Read the stdout data to save it
	fileBytes, err := readAllFunc(stdout)
	if err != nil {
		return &check.Result{
			Title:   cc.Describe(),
			Status:  check.StatusError,
			Value:   "N/A",
			Message: fmt.Sprintf("failed to read '%s': %s", path, err),
		}
	}

	// replace path separates in path to underscores
	outputFilename := strings.TrimPrefix(strings.ReplaceAll(path, "/", "_"), "_")

	// Save raw health output before parsing, in case there are parsing errors
	_, err = runContext.OutputStore.Save(
		outputFilename,
		bytes.NewReader(fileBytes),
	)
	if err != nil {
		log.Warn(
			"Failed to save %s '%s': %s",
			cc.Provider.Name(),
			path,
			err,
		)
	}

	return nil
}
