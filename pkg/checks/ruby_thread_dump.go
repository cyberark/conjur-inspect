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

// RubyThreadDump collects thread dumps from Ruby processes running in a container
type RubyThreadDump struct {
	Provider container.ContainerProvider
}

// Describe provides a textual description of what this check gathers info on
func (rtd *RubyThreadDump) Describe() string {
	return fmt.Sprintf("Ruby service thread dumps (%s)", rtd.Provider.Name())
}

// Run performs the Ruby thread dump collection
func (rtd *RubyThreadDump) Run(runContext *check.RunContext) []check.Result {
	// If there is no container ID, return
	if strings.TrimSpace(runContext.ContainerID) == "" {
		return []check.Result{}
	}

	// Check if the container runtime is available
	runtimeKey := strings.ToLower(rtd.Provider.Name())
	if !IsRuntimeAvailable(runContext, runtimeKey) {
		if runContext.VerboseErrors {
			return check.ErrorResult(
				rtd,
				fmt.Errorf("container runtime not available"),
			)
		}
		return []check.Result{}
	}

	containerInstance := rtd.Provider.Container(runContext.ContainerID)

	// Discover Ruby process PIDs
	stdout, stderr, err := containerInstance.Exec(
		"sh", "-c", "pgrep -f ruby || true",
	)
	if err != nil {
		if runContext.VerboseErrors {
			return check.ErrorResult(
				rtd,
				fmt.Errorf("failed to discover Ruby processes: %w", err),
			)
		}
		return []check.Result{}
	}

	// Read the PIDs from stdout
	pidsBytes, err := io.ReadAll(stdout)
	if err != nil {
		if runContext.VerboseErrors {
			return check.ErrorResult(
				rtd,
				fmt.Errorf("failed to read Ruby PIDs: %w", err),
			)
		}
		return []check.Result{}
	}

	// Read any stderr for logging
	stderrBytes, _ := io.ReadAll(stderr)
	if len(stderrBytes) > 0 {
		log.Debug("stderr while discovering Ruby processes: %s", string(stderrBytes))
	}

	pidsStr := strings.TrimSpace(string(pidsBytes))
	if pidsStr == "" {
		// No Ruby processes found
		log.Debug("no Ruby processes found in container")
		return []check.Result{}
	}

	// Split PIDs by newlines
	pids := strings.Split(pidsStr, "\n")
	log.Debug("found %d Ruby process(es): %s", len(pids), strings.Join(pids, ", "))

	// Collect thread dump for each PID
	results := []check.Result{}
	for _, pid := range pids {
		pid = strings.TrimSpace(pid)
		if pid == "" {
			continue
		}

		result := rtd.collectThreadDump(containerInstance, pid, runContext)
		if result != nil {
			results = append(results, *result)
		}
	}

	return results
}

// collectThreadDump sends SIGCONT to a Ruby process, waits for sigdump to write,
// reads the thread dump file, saves it to the output store, and cleans up
func (rtd *RubyThreadDump) collectThreadDump(
	containerInstance container.Container,
	pid string,
	runContext *check.RunContext,
) *check.Result {
	// Single atomic command: send signal, wait, read file, cleanup
	dumpPath := fmt.Sprintf("/tmp/sigdump-%s.log", pid)
	command := fmt.Sprintf(
		"kill -CONT %s && sleep 1 && cat %s && rm %s",
		pid, dumpPath, dumpPath,
	)

	stdout, stderr, err := containerInstance.Exec("sh", "-c", command)
	if err != nil {
		if runContext.VerboseErrors {
			return &check.Result{
				Title:   rtd.Describe(),
				Status:  check.StatusError,
				Value:   "N/A",
				Message: fmt.Sprintf("failed to collect thread dump for PID %s: %s", pid, err),
			}
		}
		return nil
	}

	// Read the thread dump from stdout
	dumpBytes, err := io.ReadAll(stdout)
	if err != nil {
		if runContext.VerboseErrors {
			return &check.Result{
				Title:   rtd.Describe(),
				Status:  check.StatusError,
				Value:   "N/A",
				Message: fmt.Sprintf("failed to read thread dump for PID %s: %s", pid, err),
			}
		}
		return nil
	}

	// Read any stderr for logging
	stderrBytes, _ := io.ReadAll(stderr)
	if len(stderrBytes) > 0 {
		log.Debug("stderr while collecting thread dump for PID %s: %s", pid, string(stderrBytes))
	}

	// Check if we got any output
	if len(dumpBytes) == 0 {
		if runContext.VerboseErrors {
			return &check.Result{
				Title:   rtd.Describe(),
				Status:  check.StatusError,
				Value:   "N/A",
				Message: fmt.Sprintf("no thread dump output for PID %s (sigdump may not be installed or enabled)", pid),
			}
		}
		return nil
	}

	// Save thread dump to output store
	outputFileName := fmt.Sprintf("ruby-dump-%s.txt", pid)
	_, err = runContext.OutputStore.Save(
		outputFileName,
		strings.NewReader(string(dumpBytes)),
	)
	if err != nil {
		if runContext.VerboseErrors {
			return &check.Result{
				Title:   rtd.Describe(),
				Status:  check.StatusError,
				Value:   "N/A",
				Message: fmt.Sprintf("failed to save thread dump for PID %s: %s", pid, err),
			}
		}
		return nil
	}

	log.Debug("successfully collected thread dump for PID %s", pid)
	return nil
}
