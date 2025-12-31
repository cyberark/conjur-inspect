// Package checks defines all of the possible Conjur Inspect checks that can
// be run.
package checks

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/cyberark/conjur-inspect/pkg/check"
	"github.com/cyberark/conjur-inspect/pkg/container"
	"github.com/cyberark/conjur-inspect/pkg/log"
	"github.com/cyberark/conjur-inspect/pkg/shell"
)

const testDir = "/var/lib/conjur/etcd_performance_test"
const etcdLogFile = testDir + "/server.log"
const logFilePrefix = "etcdctl-perf-check"

var _ check.Check = EtcdPerfCheck{}

// EtcdPerfCheck runs etcdctl check perf in a container and parses its output.
type EtcdPerfCheck struct {
	Provider        container.ContainerProvider
	RunContext      *check.RunContext
	stdoutFileName  string
	stderrFileName  string
	etcdLogFileName string
}

// Describe provides a textual description of what this check gathers info on
func (c EtcdPerfCheck) Describe() string {
	return fmt.Sprintf("Etcd Performance Check (60s) (%s)", c.Provider.Name())
}

// Run executes the etcdctl check perf command in the container and returns results.
func (c EtcdPerfCheck) Run(runContext *check.RunContext) []check.Result {
	// Check if the container runtime is available
	runtimeKey := strings.ToLower(c.Provider.Name())
	if !IsRuntimeAvailable(runContext, runtimeKey) {
		if runContext.VerboseErrors {
			return check.ErrorResult(
				c,
				fmt.Errorf("container runtime not available"),
			)
		}
		return []check.Result{}
	}

	// store RunContext so we are not passing it to methods
	c.RunContext = runContext
	providerSuffix := strings.ToLower(c.Provider.Name())
	c.stderrFileName = fmt.Sprintf("%s-stderr-%s.txt", logFilePrefix, providerSuffix)
	c.stdoutFileName = fmt.Sprintf("%s-stdout-%s.txt", logFilePrefix, providerSuffix)
	c.etcdLogFileName = fmt.Sprintf("%s-etcdLog-%s.log", logFilePrefix, providerSuffix)

	// Verify that action is valid
	validationErrors := c.validateAction()
	if len(validationErrors) != 0 {
		if runContext.VerboseErrors {
			return validationErrors
		}
		return []check.Result{}
	}

	// start etcd
	etcdPid, err := c.startEtcd()
	defer c.killEtcd(etcdPid)
	if err != nil {
		return err
	}

	// Run "etcdctl check perf" in the container
	// Note: when test fails it returns error code so at that point we cannot
	//       distinguish failed test from execution error. Therefore, even if we
	//       get an error here we still need to parse an output from stdout.
	container := c.Provider.Container(runContext.ContainerID)
	stdout, stderr, etcdErr := container.Exec("env", "ETCDCTL_API=3", "etcdctl", "check", "perf",
		"--prefix", "/etcdctl-check-perf/")

	rawPerCheckResults := shell.ReadOrDefault(stdout, "")

	// Save raw performance results to OutputStore
	_, saveErr := runContext.OutputStore.Save(c.stdoutFileName,
		strings.NewReader(rawPerCheckResults))
	if saveErr != nil {
		log.Warn("Failed to save etcdctl stdout: %w", saveErr)
	}

	// If error was returned then also save stderr.
	if etcdErr != nil {
		_, saveErr = runContext.OutputStore.Save(c.stderrFileName, stderr)
		if saveErr != nil {
			log.Warn("Failed to save etcdctl stderr: %w", saveErr)
		}
	}

	// Copy etcd log to OutputStore
	out, errorResult := c.containerCall("cat", etcdLogFile)
	if errorResult != nil {
		return errorResult
	}
	_, saveErr = runContext.OutputStore.Save(c.etcdLogFileName, out)
	if saveErr != nil {
		log.Warn("Failed to save etcd server log: %s", saveErr)
	}

	return c.parse(rawPerCheckResults)
}

func (c EtcdPerfCheck) validateAction() []check.Result {
	// If there is no container ID, return
	if strings.TrimSpace(c.RunContext.ContainerID) == "" {
		return []check.Result{
			{
				Title:   c.Describe(),
				Value:   "N/A",
				Status:  check.StatusError,
				Message: "container id missing",
			},
		}
	}

	// Check if we can execute anything in the container
	_, errorResult := c.containerCall("echo")
	if errorResult != nil {
		return errorResult
	}

	// Check if etcd is installed
	_, errorResult = c.containerCall("which", "etcd")
	if errorResult != nil {
		return errorResult
	}

	// Check if etcdctl is installed
	_, errorResult = c.containerCall("which", "etcdctl")
	if errorResult != nil {
		return errorResult
	}

	// Check if services are running, we cannot make performance checks on configured container
	services := []string{"conjur", "pg", "etcd"}
	for _, service := range services {
		out, errorResult := c.containerCall("sv", "status", service)
		if errorResult != nil {
			return errorResult
		}
		if strings.HasPrefix(shell.ReadOrDefault(out, ""), "run:") {
			return []check.Result{
				{
					Title:   c.Describe(),
					Value:   "N/A",
					Status:  check.StatusError,
					Message: fmt.Sprintf("service is running: %s", service),
				},
			}
		}
	}

	// Sanity check, check if etcd process is running
	_, errorResult = c.containerCall("pgrep", "etcd")
	if errorResult == nil {
		return errorResult
	}

	// Validation is positive
	return []check.Result{}
}

func (c EtcdPerfCheck) parse(output string) []check.Result {
	results := []check.Result{}

	const (
		PassIndicator  = "PASS:"
		ErrorIndicator = "ERROR:"
	)

	// Multiple fail indicators
	failIndicators := []string{
		"FAIL:",
		"Slowest request took too long:",
		"Stddev too high:",
	}

	trimLine := func(line string, key string) string {
		idx := strings.Index(line, key)
		if idx != -1 {
			return strings.TrimSpace(line[idx+len(key):])
		}
		return strings.TrimSpace(line)
	}

	lines := strings.SplitSeq(output, "\n")
	for line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Check fail indicators
		matchedFail := ""
		for _, fi := range failIndicators {
			if strings.Contains(line, fi) {
				matchedFail = fi
				break
			}
		}
		if matchedFail != "" {
			results = append(results, check.Result{
				Title:   c.Describe(),
				Value:   "BAD",
				Status:  check.StatusFail,
				Message: trimLine(line, matchedFail),
			})
			continue
		}

		if strings.Contains(line, PassIndicator) {
			results = append(results, check.Result{
				Title:   c.Describe(),
				Value:   "GOOD",
				Status:  check.StatusPass,
				Message: trimLine(line, PassIndicator),
			})
		} else if strings.Contains(line, ErrorIndicator) {
			results = append(results, check.Result{
				Title:   c.Describe(),
				Value:   "ERROR",
				Status:  check.StatusError,
				Message: trimLine(line, ErrorIndicator),
			})
		}
	}

	// after parsing we are expecting a single ERROR result or
	if len(results) == 1 && results[0].Status == check.StatusError {
		return results
	}

	// multiple FAIL/PASS results
	if len(results) > 1 {
		allPassOrFail := true
		for _, r := range results {
			if r.Status != check.StatusPass && r.Status != check.StatusFail {
				allPassOrFail = false
				break
			}
		}
		if allPassOrFail {
			return results
		}
	}

	// otherwise we return parsing error
	return []check.Result{
		{
			Title:   c.Describe(),
			Status:  check.StatusError,
			Message: "unexpected output from etdcctl",
		},
	}
}

func (c EtcdPerfCheck) startEtcd() (string, []check.Result) {
	// Create test directory to store etcd logs
	_, errorResult := c.containerCall("rm", "-rf", testDir)
	if errorResult != nil {
		return "", errorResult
	}
	_, errorResult = c.containerCall("mkdir", "-p", testDir)
	if errorResult != nil {
		return "", errorResult
	}

	// Start etcd in the background
	stdout, errorResult := c.containerCall("sh", "-c", fmt.Sprintf(
		"ETCD_DATA_DIR=%s ETCD_DEBUG=true etcd >%s 2>&1 & echo $!", testDir, etcdLogFile))
	if errorResult != nil {
		return "", errorResult
	}
	etcdPid := strings.TrimSpace(shell.ReadOrDefault(stdout, ""))

	// Wait for etcd to be responsive
	if !c.waitForEtcd() {
		return etcdPid, []check.Result{
			{
				Title:   c.Describe(),
				Status:  check.StatusError,
				Message: "etcd is not ready",
			},
		}
	}
	return etcdPid, nil
}

func (c EtcdPerfCheck) waitForEtcd() bool {
	timeoutSeconds := 60
	timeout := time.After(time.Duration(timeoutSeconds) * time.Second)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-timeout:
			return false
		case <-ticker.C:
			stdout, errorResult := c.containerCall("curl", "-s", "-o", "/dev/null",
				"-w", "%{http_code}", "http://127.0.0.1:2379/health")
			if errorResult != nil {
				return false
			}
			if strings.TrimSpace(shell.ReadOrDefault(stdout, "")) == "200" {
				return true
			}
		}
	}
}

func (c EtcdPerfCheck) killEtcd(etcdPid string) {
	if etcdPid != "" {
		_, errorResult := c.containerCall("kill", "-HUP", etcdPid)
		if len(errorResult) > 0 {
			log.Warn("failed to kill etcd: %s", errorResult[0].Message)
		}
	}
}

func (c EtcdPerfCheck) containerCall(args ...string) (io.Reader, []check.Result) {
	container := c.Provider.Container(c.RunContext.ContainerID)
	stdout, stderr, err := container.Exec(args...)
	if err != nil {
		rawStderr := shell.ReadOrDefault(stderr, "")
		_, saveErr := c.RunContext.OutputStore.Save(c.stderrFileName, strings.NewReader(rawStderr))
		if saveErr != nil {
			log.Warn("Failed to save %s: %w", c.stderrFileName, saveErr)
		}
		return nil, check.ErrorResult(
			c, fmt.Errorf("%s: %s (error: %w) (stderr: %s)", "error during container call", args, err, rawStderr))
	}
	return stdout, nil
}
