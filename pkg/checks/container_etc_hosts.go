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
)

// ContainerEtcHosts collects the contents of /etc/hosts from inside a container
type ContainerEtcHosts struct {
	Provider container.ContainerProvider
}

// Describe provides a textual description of what this check gathers info on
func (ceh *ContainerEtcHosts) Describe() string {
	return fmt.Sprintf("%s /etc/hosts", ceh.Provider.Name())
}

// Run performs the container /etc/hosts collection
func (ceh *ContainerEtcHosts) Run(runContext *check.RunContext) []check.Result {
	// If there is no container ID, return
	if strings.TrimSpace(runContext.ContainerID) == "" {
		return []check.Result{}
	}

	// Check if the container runtime is available
	runtimeKey := strings.ToLower(ceh.Provider.Name())
	if !IsRuntimeAvailable(runContext, runtimeKey) {
		if runContext.VerboseErrors {
			return check.ErrorResult(
				ceh,
				fmt.Errorf("container runtime not available"),
			)
		}
		return []check.Result{}
	}

	container := ceh.Provider.Container(runContext.ContainerID)

	// Execute cat /etc/hosts inside the container
	stdout, stderr, err := container.Exec("cat", "/etc/hosts")
	if err != nil {
		if runContext.VerboseErrors {
			stderrBytes, _ := io.ReadAll(stderr)
			return check.ErrorResult(
				ceh,
				fmt.Errorf("failed to read /etc/hosts: %w (stderr: %s)", err, string(stderrBytes)),
			)
		}
		log.Warn("failed to read /etc/hosts from container: %s", err)
		return []check.Result{}
	}

	// Read the stdout content
	fileBytes, err := io.ReadAll(stdout)
	if err != nil {
		if runContext.VerboseErrors {
			return check.ErrorResult(
				ceh,
				fmt.Errorf("failed to read command output: %w", err),
			)
		}
		log.Warn("failed to read /etc/hosts output: %s", err)
		return []check.Result{}
	}

	// Save the file contents to output store with runtime-specific filename
	outputFilename := fmt.Sprintf(
		"%s-etc-hosts.txt",
		strings.ToLower(ceh.Provider.Name()),
	)
	_, err = runContext.OutputStore.Save(outputFilename, bytes.NewReader(fileBytes))
	if err != nil {
		log.Warn("failed to save /etc/hosts output: %s", err)
	}

	// Return empty results on success (output is saved)
	return []check.Result{}
}
