// Package checks defines all of the possible Conjur Inspect checks that can
// be run.
package checks

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/cyberark/conjur-inspect/pkg/check"
	"github.com/cyberark/conjur-inspect/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestEtcdClusterMembersRunEnrolled(t *testing.T) {
	soloJSONContent := `{
		"conjur": {
			"cluster_name": "test-cluster"
		}
	}`

	memberListOutput := `ID        | Name       | ClientURLs       | PeerURLs
1         | node-1     | http://1.1.1.1   | http://1.1.1.1
2         | node-2     | http://2.2.2.2   | http://2.2.2.2`

	testCheck := &EtcdClusterMembers{
		Provider: &test.ContainerProvider{
			ExecResponses: map[string]test.ExecResponse{
				"cat /etc/cinc/solo.json": {
					Stdout: strings.NewReader(soloJSONContent),
					Stderr: strings.NewReader(""),
					Error:  nil,
				},
				"evoke cluster member list": {
					Stdout: strings.NewReader(memberListOutput),
					Stderr: strings.NewReader(""),
					Error:  nil,
				},
			},
		},
	}

	testOutputStore := test.NewOutputStore()

	results := testCheck.Run(
		&check.RunContext{
			ContainerID: "test",
			OutputStore: testOutputStore,
		},
	)

	assert.Empty(t, results)

	outputStoreItems, err := testOutputStore.Items()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(outputStoreItems))

	itemInfo, err := outputStoreItems[0].Info()
	assert.NoError(t, err)
	assert.Equal(t, "etcd-cluster-members-test container provider.txt", itemInfo.Name())

	reader, cleanup, err := outputStoreItems[0].Open()
	assert.NoError(t, err)
	defer cleanup()

	output, err := io.ReadAll(reader)
	assert.NoError(t, err)
	assert.Equal(t, memberListOutput, string(output))
}

func TestEtcdClusterMembersRunNotEnrolled(t *testing.T) {
	soloJSONContent := `{
		"conjur": {
			"cluster_name": ""
		}
	}`

	testCheck := &EtcdClusterMembers{
		Provider: &test.ContainerProvider{
			ExecResponses: map[string]test.ExecResponse{
				"cat /etc/cinc/solo.json": {
					Stdout: strings.NewReader(soloJSONContent),
					Stderr: strings.NewReader(""),
					Error:  nil,
				},
			},
		},
	}

	testOutputStore := test.NewOutputStore()

	results := testCheck.Run(
		&check.RunContext{
			ContainerID: "test",
			OutputStore: testOutputStore,
		},
	)

	assert.Empty(t, results)

	outputStoreItems, err := testOutputStore.Items()
	assert.NoError(t, err)
	assert.Equal(t, 0, len(outputStoreItems))
}

func TestEtcdClusterMembersRunNotEnrolledMissingAttribute(t *testing.T) {
	soloJSONContent := `{
		"conjur": {
			"other_attribute": "value"
		}
	}`

	testCheck := &EtcdClusterMembers{
		Provider: &test.ContainerProvider{
			ExecResponses: map[string]test.ExecResponse{
				"cat /etc/cinc/solo.json": {
					Stdout: strings.NewReader(soloJSONContent),
					Stderr: strings.NewReader(""),
					Error:  nil,
				},
			},
		},
	}

	testOutputStore := test.NewOutputStore()

	results := testCheck.Run(
		&check.RunContext{
			ContainerID: "test",
			OutputStore: testOutputStore,
		},
	)

	assert.Empty(t, results)

	outputStoreItems, err := testOutputStore.Items()
	assert.NoError(t, err)
	assert.Equal(t, 0, len(outputStoreItems))
}

func TestEtcdClusterMembersRunMissingSoloJson(t *testing.T) {
	testCheck := &EtcdClusterMembers{
		Provider: &test.ContainerProvider{
			ExecResponses: map[string]test.ExecResponse{
				"cat /etc/cinc/solo.json": {
					Stdout: strings.NewReader(""),
					Stderr: strings.NewReader("cat: /etc/cinc/solo.json: No such file or directory"),
					Error:  errors.New("exit code 1"),
				},
			},
		},
	}

	results := testCheck.Run(
		&check.RunContext{
			ContainerID: "test",
			OutputStore: test.NewOutputStore(),
		},
	)

	assert.Equal(t, 1, len(results))
	assert.Equal(t, "Etcd Cluster Members (Test Container Provider)", results[0].Title)
	assert.Equal(t, check.StatusError, results[0].Status)
	assert.Equal(t, "N/A", results[0].Value)
	assert.Contains(t, results[0].Message, "failed to read solo.json")
}

func TestEtcdClusterMembersRunInvalidSoloJson(t *testing.T) {
	soloJSONContent := `{ invalid json`

	testCheck := &EtcdClusterMembers{
		Provider: &test.ContainerProvider{
			ExecResponses: map[string]test.ExecResponse{
				"cat /etc/cinc/solo.json": {
					Stdout: strings.NewReader(soloJSONContent),
					Stderr: strings.NewReader(""),
					Error:  nil,
				},
			},
		},
	}

	results := testCheck.Run(
		&check.RunContext{
			ContainerID: "test",
			OutputStore: test.NewOutputStore(),
		},
	)

	assert.Equal(t, 1, len(results))
	assert.Equal(t, "Etcd Cluster Members (Test Container Provider)", results[0].Title)
	assert.Equal(t, check.StatusError, results[0].Status)
	assert.Equal(t, "N/A", results[0].Value)
	assert.Contains(t, results[0].Message, "failed to parse solo.json")
}

func TestEtcdClusterMembersRunEvokeCommandError(t *testing.T) {
	soloJSONContent := `{
		"conjur": {
			"cluster_name": "test-cluster"
		}
	}`

	testCheck := &EtcdClusterMembers{
		Provider: &test.ContainerProvider{
			ExecResponses: map[string]test.ExecResponse{
				"cat /etc/cinc/solo.json": {
					Stdout: strings.NewReader(soloJSONContent),
					Stderr: strings.NewReader(""),
					Error:  nil,
				},
				"evoke cluster member list": {
					Stdout: strings.NewReader(""),
					Stderr: strings.NewReader("error: failed to list members"),
					Error:  errors.New("exit code 1"),
				},
			},
		},
	}

	results := testCheck.Run(
		&check.RunContext{
			ContainerID: "test",
			OutputStore: test.NewOutputStore(),
		},
	)

	assert.Equal(t, 1, len(results))
	assert.Equal(t, "Etcd Cluster Members (Test Container Provider)", results[0].Title)
	assert.Equal(t, check.StatusError, results[0].Status)
	assert.Equal(t, "N/A", results[0].Value)
	assert.Contains(t, results[0].Message, "failed to get cluster members")
}

func TestEtcdClusterMembersRunNoContainerID(t *testing.T) {
	testCheck := &EtcdClusterMembers{
		Provider: &test.ContainerProvider{},
	}

	results := testCheck.Run(
		&check.RunContext{
			ContainerID: "",
			OutputStore: test.NewOutputStore(),
		},
	)

	assert.Empty(t, results)
}

func TestEtcdClusterMembersRunMissingConjurSection(t *testing.T) {
	soloJSONContent := `{
		"other_section": {
			"cluster_name": "test-cluster"
		}
	}`

	testCheck := &EtcdClusterMembers{
		Provider: &test.ContainerProvider{
			ExecResponses: map[string]test.ExecResponse{
				"cat /etc/cinc/solo.json": {
					Stdout: strings.NewReader(soloJSONContent),
					Stderr: strings.NewReader(""),
					Error:  nil,
				},
			},
		},
	}

	results := testCheck.Run(
		&check.RunContext{
			ContainerID: "test",
			OutputStore: test.NewOutputStore(),
		},
	)

	assert.Equal(t, 1, len(results))
	assert.Equal(t, "Etcd Cluster Members (Test Container Provider)", results[0].Title)
	assert.Equal(t, check.StatusError, results[0].Status)
	assert.Equal(t, "N/A", results[0].Value)
	assert.Contains(t, results[0].Message, "conjur section not found")
}
