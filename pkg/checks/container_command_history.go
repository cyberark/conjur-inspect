// Package checks defines all of the possible Conjur Inspect checks that can
// be run.
package checks

import (
	"fmt"
	"io"
	"strings"

	"github.com/cyberark/conjur-inspect/pkg/check"
	"github.com/cyberark/conjur-inspect/pkg/checks/sanitize"
	"github.com/cyberark/conjur-inspect/pkg/container"
	"github.com/cyberark/conjur-inspect/pkg/log"
)

// ContainerCommandHistory collects recent command history from inside a container
type ContainerCommandHistory struct {
	Provider container.ContainerProvider
}

// Describe provides a textual description of what this check gathers info on
func (cch *ContainerCommandHistory) Describe() string {
	return fmt.Sprintf("%s command history", cch.Provider.Name())
}

// Run performs the container command history collection
func (cch *ContainerCommandHistory) Run(runContext *check.RunContext) []check.Result {
	// If there is no container ID, return
	if strings.TrimSpace(runContext.ContainerID) == "" {
		return []check.Result{}
	}

	// Check if the container runtime is available
	runtimeKey := strings.ToLower(cch.Provider.Name())
	if !IsRuntimeAvailable(runContext, runtimeKey) {
		if runContext.VerboseErrors {
			return check.ErrorResult(
				cch,
				fmt.Errorf("container runtime not available"),
			)
		}
		return []check.Result{}
	}

	containerInstance := cch.Provider.Container(runContext.ContainerID)

	// Execute tail command to get last 100 lines of bash history
	// Use a shell command that won't fail if the file doesn't exist
	stdout, stderr, err := containerInstance.Exec(
		"sh", "-c", "tail -n 100 /root/.bash_history 2>/dev/null || true",
	)
	if err != nil {
		return check.ErrorResult(
			cch,
			fmt.Errorf("failed to retrieve command history from container: %w", err),
		)
	}

	// Read stdout from the command execution
	historyBytes, err := io.ReadAll(stdout)
	if err != nil {
		return check.ErrorResult(
			cch,
			fmt.Errorf("failed to read command history output: %w", err),
		)
	}

	// Read any stderr for logging
	stderrBytes, _ := io.ReadAll(stderr)
	if len(stderrBytes) > 0 {
		log.Warn("stderr while reading container command history: %s", string(stderrBytes))
	}

	historyContent := string(historyBytes)

	// Sanitize the history to redact sensitive values
	redactor := sanitize.NewRedactor()
	sanitizedContent := redactor.RedactLines(historyContent)

	// Save history to output store
	outputFileName := fmt.Sprintf(
		"%s-command-history.txt",
		strings.ToLower(cch.Provider.Name()),
	)
	_, err = runContext.OutputStore.Save(
		outputFileName,
		strings.NewReader(sanitizedContent),
	)
	if err != nil {
		log.Warn("failed to save container command history output: %s", err)
	}

	// Return empty results on success (output is saved)
	return []check.Result{}
}
