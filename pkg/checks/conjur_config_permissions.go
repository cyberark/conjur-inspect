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
	"github.com/cyberark/conjur-inspect/pkg/shell"
)

// ConjurConfigPermissions collects the permissions of the Conjur configuration
// file (conjur.yml) and containing directory
type ConjurConfigPermissions struct {
	Provider container.ContainerProvider
}

// Describe provides a textual description of what this check gathers info on
func (ccp *ConjurConfigPermissions) Describe() string {
	return fmt.Sprintf("Conjur Config Permissions (%s)", ccp.Provider.Name())
}

// Run performs the Conjur configuration check
func (ccp *ConjurConfigPermissions) Run(runContext *check.RunContext) []check.Result {
	// If there is no container ID, return
	if strings.TrimSpace(runContext.ContainerID) == "" {
		return []check.Result{}
	}

	container := ccp.Provider.Container(runContext.ContainerID)

	results := []check.Result{}

	result := ccp.collectConjurConfigPermissions(container, runContext)
	if result != nil {
		results = append(results, *result)
	}

	return results
}

func (ccp *ConjurConfigPermissions) collectConjurConfigPermissions(
	container container.Container,
	runContext *check.RunContext,
) *check.Result {
	stdout, stderr, err := container.Exec(
		"ls", "-la", "/etc/conjur/config",
	)

	if err != nil {
		return &check.Result{
			Title:   ccp.Describe(),
			Status:  check.StatusError,
			Value:   "N/A",
			Message: fmt.Sprintf(
				"failed to collect Conjur config permissions: %s (%s))",
				err,
				strings.TrimSpace(shell.ReadOrDefault(stderr, "N/A")),
			),
		}
	}

	// Read the stdout data to save it
	fileBytes, err := readAllFunc(stdout)
	if err != nil {
		return &check.Result{
			Title:   ccp.Describe(),
			Status:  check.StatusError,
			Value:   "N/A",
			Message: fmt.Sprintf("failed to read Conjur config permissions: %s", err),
		}
	}

	// Save raw health output before parsing, in case there are parsing errors
	_, err = runContext.OutputStore.Save(
		"conjur_config_permissions.txt",
		bytes.NewReader(fileBytes),
	)
	if err != nil {
		log.Warn(
			"Failed to save %s Conjur config permissions: %s",
			ccp.Provider.Name(),
			err,
		)
	}

	return nil
}
