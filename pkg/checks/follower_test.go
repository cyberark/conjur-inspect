package checks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFollowerRun(t *testing.T) {
	t.Setenv("MASTER_HOSTNAME", "http://example.com")
	testCheck := &Follower{}
	resultChan := testCheck.Run()
	results := <-resultChan

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
