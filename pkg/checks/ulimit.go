package checks

import (
	"bufio"
	"io"
	"strings"

	"github.com/cyberark/conjur-inspect/pkg/check"
	"github.com/cyberark/conjur-inspect/pkg/shell"
)

var executeUlimitInfoFunc func() (stderr, stdout io.Reader, err error) = executeUlimitInfo

// Ulimit collects information on the systems avalible resources.
type Ulimit struct{}

// Describe provides a textual description of what this check gathers info on
func (*Ulimit) Describe() string {
	return "ulimit"
}

// Run performs the Ulimit collection
func (ulimit *Ulimit) Run(context *check.RunContext) []check.Result {
	ulimitOutput, stderr, err := executeUlimitInfoFunc()

	// In case of an error, return a check result with an error status.
	if err != nil {
		return []check.Result{
			{
				Title:   "Ulimit Error",
				Status:  check.StatusError,
				Value:   "N/A",
				Message: shell.ReadOrDefault(stderr, "N/A"),
			},
		}
	}

	// A slice of all ulimit results.
	results := []check.Result{}

	// Iterate over output lines
	scanner := bufio.NewScanner(ulimitOutput)
	for scanner.Scan() {
		// Extract the resource name and value
		fields := strings.Fields(scanner.Text())

		// Extract the resource name by joining all elements before the last element in fields.
		resourceName := strings.Join(fields[:len(fields)-1], " ")

		// Extract the resource value as the last element in fields.
		resourceValue := fields[len(fields)-1]

		result := check.Result{
			Title:  resourceName,
			Status: check.StatusInfo,
			Value:  resourceValue,
		}

		results = append(results, result)
	}

	return results
}

func executeUlimitInfo() (stdout, stderr io.Reader, err error) {
	return shell.NewCommandWrapper(
		"sh",
		"-c",
		"ulimit -a",
	).Run()
}
