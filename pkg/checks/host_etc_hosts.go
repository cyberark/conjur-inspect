// Package checks defines all of the possible Conjur Inspect checks that can
// be run.
package checks

import (
	"bytes"
	"os"

	"github.com/cyberark/conjur-inspect/pkg/check"
	"github.com/cyberark/conjur-inspect/pkg/log"
)

// HostEtcHosts collects the contents of /etc/hosts from the host machine
type HostEtcHosts struct{}

// Describe provides a textual description of what this check gathers info on
func (*HostEtcHosts) Describe() string {
	return "Host /etc/hosts"
}

// Run performs the host /etc/hosts collection
func (h *HostEtcHosts) Run(runContext *check.RunContext) []check.Result {
	fileBytes, err := os.ReadFile("/etc/hosts")
	if err != nil {
		if runContext.VerboseErrors {
			return check.ErrorResult(h, err)
		}
		log.Warn("failed to read /etc/hosts: %s", err)
		return []check.Result{}
	}

	// Save the file contents to output store
	_, err = runContext.OutputStore.Save(
		"host-etc-hosts.txt",
		bytes.NewReader(fileBytes),
	)
	if err != nil {
		log.Warn("failed to save /etc/hosts output: %s", err)
	}

	// Return empty results on success (output is saved)
	return []check.Result{}
}
