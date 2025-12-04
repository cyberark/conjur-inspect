// Package checks defines all of the possible Conjur Inspect checks that can
// be run.
package checks

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/cyberark/conjur-inspect/pkg/check"
	"github.com/cyberark/conjur-inspect/pkg/container"
	"github.com/cyberark/conjur-inspect/pkg/log"
	"github.com/cyberark/conjur-inspect/pkg/shell"
)

// Alias io.ReadAll as a variable so that we can stub it out for unit tests
var readAllFuncClusterMembers = io.ReadAll

// EtcdClusterMembers collects the current etcd cluster members by running
// `evoke cluster member list` in an enrolled cluster node
type EtcdClusterMembers struct {
	Provider container.ContainerProvider
}

// Describe provides a textual description of what this check gathers info on
func (ecm *EtcdClusterMembers) Describe() string {
	return fmt.Sprintf("Etcd Cluster Members (%s)", ecm.Provider.Name())
}

// Run executes the cluster member list command and saves the output
func (ecm *EtcdClusterMembers) Run(runContext *check.RunContext) []check.Result {
	// If there is no container ID, return
	if strings.TrimSpace(runContext.ContainerID) == "" {
		return []check.Result{}
	}

	container := ecm.Provider.Container(runContext.ContainerID)

	// Check if node is enrolled in a cluster
	isEnrolled, err := ecm.isNodeEnrolled(container)
	if err != nil {
		return check.ErrorResult(ecm, err)
	}

	// If not enrolled, return empty results
	if !isEnrolled {
		return []check.Result{}
	}

	// Run evoke cluster member list command
	stdout, stderr, err := container.Exec("evoke", "cluster", "member", "list")
	if err != nil {
		stderrMsg := shell.ReadOrDefault(stderr, "N/A")
		return check.ErrorResult(
			ecm,
			fmt.Errorf("failed to get cluster members: %w (stderr: %s)", err, stderrMsg),
		)
	}

	// Read the stdout data
	memberListOutput, err := readAllFuncClusterMembers(stdout)
	if err != nil {
		return check.ErrorResult(
			ecm,
			fmt.Errorf("failed to read cluster member list output: %w", err),
		)
	}

	// Save raw output to OutputStore
	providerSuffix := strings.ToLower(ecm.Provider.Name())
	outputFileName := fmt.Sprintf("etcd-cluster-members-%s.txt", providerSuffix)
	_, saveErr := runContext.OutputStore.Save(outputFileName, strings.NewReader(string(memberListOutput)))
	if saveErr != nil {
		log.Warn("Failed to save cluster member list output: %w", saveErr)
	}

	return []check.Result{}
}

// isNodeEnrolled checks if the node is enrolled in a cluster by reading
// /etc/cinc/solo.json and verifying conjur.cluster_name exists and is non-empty
func (ecm *EtcdClusterMembers) isNodeEnrolled(container container.Container) (bool, error) {
	stdout, stderr, err := container.Exec("cat", "/etc/cinc/solo.json")
	if err != nil {
		stderrMsg := shell.ReadOrDefault(stderr, "N/A")
		return false, fmt.Errorf("failed to read solo.json: %w (stderr: %s)", err, stderrMsg)
	}

	soloJSONBytes, err := readAllFuncClusterMembers(stdout)
	if err != nil {
		return false, fmt.Errorf("failed to read solo.json content: %w", err)
	}

	// Parse solo.json to extract conjur.cluster_name
	var soloConfig map[string]interface{}
	err = json.Unmarshal(soloJSONBytes, &soloConfig)
	if err != nil {
		return false, fmt.Errorf("failed to parse solo.json: %w", err)
	}

	// Check if conjur section exists
	conjurConfig, exists := soloConfig["conjur"]
	if !exists {
		return false, fmt.Errorf("conjur section not found in solo.json")
	}

	conjurMap, ok := conjurConfig.(map[string]interface{})
	if !ok {
		return false, fmt.Errorf("conjur section is not a valid object in solo.json")
	}

	// Check if cluster_name exists and is non-empty
	clusterName, exists := conjurMap["cluster_name"]
	if !exists {
		return false, nil // Not enrolled
	}

	clusterNameStr, ok := clusterName.(string)
	if !ok {
		return false, fmt.Errorf("cluster_name is not a string in solo.json")
	}

	return strings.TrimSpace(clusterNameStr) != "", nil
}
