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

func TestConjurConfig_Run_EmptyContainerID(t *testing.T) {
	provider := &test.ContainerProvider{}
	cc := &ConjurConfig{Provider: provider}

	results := cc.Run(&check.RunContext{
		ContainerID: "",
		OutputStore: test.NewOutputStore(),
	})

	assert.Empty(t, results)
}

func TestConjurConfig_Run_Success(t *testing.T) {
	testConfigFileContents := "test config content"

	provider := &test.ContainerProvider{
		ExecResponses: map[string]test.ExecResponse{
			"cat /etc/conjur/config/conjur.yml": {
				Stdout: strings.NewReader(testConfigFileContents),
			},
			"cat /etc/postgresql/15/main/postgresql.conf": {
				Stdout: strings.NewReader(testConfigFileContents),
			},
			"cat /opt/conjur/etc/conjur.conf": {
				Stdout: strings.NewReader(testConfigFileContents),
			},
			"cat /opt/conjur/etc/possum.conf": {
				Stdout: strings.NewReader(testConfigFileContents),
			},
			"cat /opt/conjur/etc/ui.conf": {
				Stdout: strings.NewReader(testConfigFileContents),
			},
			"cat /opt/conjur/etc/cluster.conf": {
				Stdout: strings.NewReader(testConfigFileContents),
			},
			"cat /etc/cinc/solo.json": {
				Stdout: strings.NewReader(testConfigFileContents),
			},
			"cat /etc/chef/solo.json": {
				Stdout: strings.NewReader(testConfigFileContents),
			},
		},
	}
	cc := &ConjurConfig{Provider: provider}

	runContext := test.NewRunContext("test-container")
	results := cc.Run(&runContext)

	assert.Empty(t, results) // Success case returns no results

	// Verify saved outputs
	items, err := runContext.OutputStore.Items()
	assert.NoError(t, err)
	assert.NotEmpty(t, items)

	// Verify config file saved
	assert.Condition(t, func() bool {
		for _, item := range items {
			info, err := item.Info()
			assert.NoError(t, err)
			if strings.Contains(info.Name(), "etc_conjur_config_conjur.yml") {
				return true
			}
		}
		return false
	}, "Expected config file not found in output items")
}

func TestConjurConfig_Run_FileReadError(t *testing.T) {
	provider := &test.ContainerProvider{
		ExecResponses: map[string]test.ExecResponse{
			"cat /etc/conjur/config/conjur.yml": {
				Error:  errors.New("file not found"),
				Stderr: strings.NewReader("permission denied"),
			},
		},
	}

	cc := &ConjurConfig{Provider: provider}

	runContext := test.NewRunContext("test-container")
	runContext.VerboseErrors = true
	results := cc.Run(&runContext)

	assert.NotEmpty(t, results)
	assert.Equal(t, check.StatusError, results[0].Status)
	assert.Contains(t, results[0].Message, "failed to collect")
}

func TestConjurConfig_Run_FileReadErrorNoVerboseErrors(t *testing.T) {
	provider := &test.ContainerProvider{
		ExecResponses: map[string]test.ExecResponse{
			"cat /etc/conjur/config/conjur.yml": {
				Error:  errors.New("file not found"),
				Stderr: strings.NewReader("permission denied"),
			},
		},
	}

	cc := &ConjurConfig{Provider: provider}

	runContext := test.NewRunContext("test-container")
	runContext.VerboseErrors = false
	results := cc.Run(&runContext)

	// Expect no results when VerboseErrors is false
	assert.Empty(t, results)
}

func TestConjurConfig_Run_ReadAllError(t *testing.T) {
	// Save original readAllFunc to restore after test
	originalReadAllFunc := readAllFunc
	defer func() { readAllFunc = originalReadAllFunc }()

	readAllError := errors.New("read error")
	readAllFunc = func(r io.Reader) ([]byte, error) {
		return nil, readAllError
	}

	provider := &test.ContainerProvider{
		ExecResponses: map[string]test.ExecResponse{
			"cat /etc/conjur/config/conjur.yml": {
				Stdout: strings.NewReader("test content"),
			},
		},
	}

	cc := &ConjurConfig{Provider: provider}

	runContext := test.NewRunContext("test-container")
	results := cc.Run(&runContext)

	assert.NotEmpty(t, results)
	assert.Equal(t, check.StatusError, results[0].Status)
	assert.Contains(t, results[0].Message, "failed to read")
	assert.Contains(t, results[0].Message, readAllError.Error())
}
