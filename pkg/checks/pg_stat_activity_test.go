// Package checks defines all of the possible Conjur Inspect checks that can
// be run.
package checks

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cyberark/conjur-inspect/pkg/check"
	"github.com/cyberark/conjur-inspect/pkg/test"
)

func TestPgStatActivityDescribe(t *testing.T) {
	pgStatActivity := &PgStatActivity{
		Provider: &test.ContainerProvider{},
	}

	description := pgStatActivity.Describe()
	assert.Equal(t, "PostgreSQL pg_stat_activity (Test Container Provider)", description)
}

func TestPgStatActivityRun(t *testing.T) {
	pgStatActivityOutput := `   pid   | usesysid | usename  |     application_name      | client_addr | state  |                  query
---------+----------+----------+---------------------------+-------------+--------+------------------------------------------
  123456 |    16384 | postgres | psql                      |             | active | select * from pg_stat_activity
  123457 |    16384 | conjur   | conjur api                | 127.0.0.1   | idle   | <idle>
(2 rows)`

	provider := &test.ContainerProvider{
		ExecAsUserResponses: map[string]test.ExecResponse{
			"conjur psql -c select * from pg_stat_activity": {
				Stdout: strings.NewReader(pgStatActivityOutput),
				Stderr: strings.NewReader(""),
				Error:  nil,
			},
		},
	}

	pgStatActivity := &PgStatActivity{
		Provider: provider,
	}

	runContext := test.NewRunContext("test-container-id")
	results := pgStatActivity.Run(&runContext)

	// Should return no results since this check only produces raw output
	assert.Len(t, results, 0)

	// Verify the output file was saved
	items, err := runContext.OutputStore.Items()
	require.NoError(t, err)
	require.Len(t, items, 1)
	info, err := items[0].Info()
	require.NoError(t, err)
	assert.Equal(t, "pg_stat_activity.log", info.Name())
}

func TestPgStatActivityRunEmptyContainerID(t *testing.T) {
	provider := &test.ContainerProvider{}

	pgStatActivity := &PgStatActivity{
		Provider: provider,
	}

	runContext := test.NewRunContext("")
	results := pgStatActivity.Run(&runContext)

	// Should return no results when container ID is empty
	assert.Len(t, results, 0)
}

func TestPgStatActivityRunExecError(t *testing.T) {
	execError := fmt.Errorf("connection refused")

	provider := &test.ContainerProvider{
		ExecAsUserResponses: map[string]test.ExecResponse{
			"conjur psql -c select * from pg_stat_activity": {
				Stdout: nil,
				Stderr: strings.NewReader("psql: error in connecting to database\n"),
				Error:  execError,
			},
		},
	}

	pgStatActivity := &PgStatActivity{
		Provider: provider,
	}

	runContext := test.NewRunContext("test-container-id")
	results := pgStatActivity.Run(&runContext)

	// Should return an error result
	require.Len(t, results, 1)
	assert.Equal(t, check.StatusError, results[0].Status)
	assert.Contains(t, results[0].Message, "failed to retrieve pg_stat_activity")

	// Should save error output file
	items, err := runContext.OutputStore.Items()
	require.NoError(t, err)
	require.Len(t, items, 1)
	info, err := items[0].Info()
	require.NoError(t, err)
	// Provider name is "Test Container Provider" and gets ToLower() -> "test container provider"
	assert.Equal(t, "test container provider-pg-stat-activity-error.log", info.Name())
}

func TestPgStatActivityRunStderr(t *testing.T) {
	pgStatActivityOutput := `   pid   | usesysid | usename  |     application_name      | client_addr | state  |                  query
---------+----------+----------+---------------------------+-------------+--------+------------------------------------------
  123456 |    16384 | postgres | psql                      |             | active | select * from pg_stat_activity
(1 row)`

	provider := &test.ContainerProvider{
		ExecAsUserResponses: map[string]test.ExecResponse{
			"conjur psql -c select * from pg_stat_activity": {
				Stdout: strings.NewReader(pgStatActivityOutput),
				Stderr: strings.NewReader("warning: some non-fatal warning"),
				Error:  nil,
			},
		},
	}

	pgStatActivity := &PgStatActivity{
		Provider: provider,
	}

	runContext := test.NewRunContext("test-container-id")
	results := pgStatActivity.Run(&runContext)

	// Should still succeed with empty results despite stderr
	assert.Empty(t, results)

	// Verify the output was saved
	items, err := runContext.OutputStore.Items()
	require.NoError(t, err)
	require.Len(t, items, 1)
}
