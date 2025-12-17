package checks

import (
	"testing"

	"github.com/cyberark/conjur-inspect/pkg/check"
	"github.com/stretchr/testify/assert"
)

func TestFollowerRun(t *testing.T) {
	t.Setenv("MASTER_HOSTNAME", "http://example.com")
	testCheck := &Follower{}
	results := testCheck.Run(&check.RunContext{})

	leaderReplicationPort := GetResultByTitle(results, "Leader Replication Port")
	assert.NotNil(t, leaderReplicationPort)
	assert.NotEmpty(t, leaderReplicationPort.Status)
	assert.NotEmpty(t, leaderReplicationPort.Value)

	leaderAPIPort := GetResultByTitle(results, "Leader API Port")
	assert.NotNil(t, leaderAPIPort, "Includes 'Leader API Port'")
	assert.NotEmpty(t, leaderAPIPort.Status)
	assert.NotEmpty(t, leaderAPIPort.Value)

	leaderAuditPort := GetResultByTitle(results, "Leader Audit Forwarding Port")
	assert.NotNil(t, leaderAuditPort, "Includes 'Leader Audit Forwarding Port'")
	assert.NotEmpty(t, leaderAuditPort.Status)
	assert.NotEmpty(t, leaderAuditPort.Value)
}

func TestFollowerRunWithoutMasterHostnameVerboseErrors(t *testing.T) {
	t.Setenv("MASTER_HOSTNAME", "")
	testCheck := &Follower{}
	results := testCheck.Run(&check.RunContext{VerboseErrors: true})

	assert.NotEmpty(t, results)
	assert.Equal(t, check.StatusError, results[0].Status)
	assert.Equal(t, "N/A", results[0].Value)
	assert.Contains(t, results[0].Message, "Leader hostname is not set")
	assert.Contains(t, results[0].Message, "MASTER_HOSTNAME")
}

func TestFollowerRunWithoutMasterHostnameNoVerboseErrors(t *testing.T) {
	t.Setenv("MASTER_HOSTNAME", "")
	testCheck := &Follower{}
	results := testCheck.Run(&check.RunContext{VerboseErrors: false})

	assert.Empty(t, results)
}
