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
func (ulimit *Ulimit) Run(context *check.RunContext) <-chan []check.Result {
	future := make(chan []check.Result)

	go func() {
		ulimitOutput, stderr, err := executeUlimitInfoFunc()

		// In case of an error, return a check result with an error status.
		if err != nil {
			future <- []check.Result{
				{
					Title:   "Ulimit Error",
					Status:  check.StatusError,
					Value:   "N/A",
					Message: shell.ReadOrDefault(stderr, "N/A"),
				},
			}
			return
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

		future <- results
	}() // async

	return future
}

func executeUlimitInfo() (stdout, stderr io.Reader, err error) {
	return shell.NewCommandWrapper(
		"sh",
		"-c",
		"ulimit -a",
	).Run()
}
