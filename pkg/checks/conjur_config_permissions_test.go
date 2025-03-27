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

func TestConjurConfigPermissions_Run_EmptyContainerID(t *testing.T) {
	provider := &test.ContainerProvider{}
	cc := &ConjurConfigPermissions{Provider: provider}

	results := cc.Run(&check.RunContext{
		ContainerID: "",
		OutputStore: test.NewOutputStore(),
	})

	assert.Empty(t, results)
}

func TestConjurConfigPermissions_Run_Success(t *testing.T) {
	testConfigFilePermissions := "drwxr-xr-x"

	provider := &test.ContainerProvider{
		ExecResponses: map[string]test.ExecResponse{
			"ls -la /etc/conjur/config": {
				Stdout: strings.NewReader(testConfigFilePermissions),
			},
		},
	}
	ccp := &ConjurConfigPermissions{Provider: provider}

	runContext := test.NewRunContext("test-container")
	results := ccp.Run(&runContext)

	assert.Empty(t, results) // Success case returns no results

	// Verify saved outputs
	outputStoreItems, err := runContext.OutputStore.Items()
	assert.NoError(t, err)
	assert.Len(t, outputStoreItems, 1)

	outputStoreItem := outputStoreItems[0]

	outputStoreItemInfo, err := outputStoreItem.Info()
	assert.NoError(t, err)
	assert.Equal(
		t,
		"conjur_config_permissions.txt",
		outputStoreItemInfo.Name(),
	)
}

func TestConjurConfigPermissions_Run_ExecError(t *testing.T) {
	provider := &test.ContainerProvider{
		ExecResponses: map[string]test.ExecResponse{
			"ls -la /etc/conjur/config": {
				Error: errors.New("file not found"),
				Stderr: strings.NewReader("permission denied"),
			},
		},
	}

	ccp := &ConjurConfigPermissions{Provider: provider}

	runContext := test.NewRunContext("test-container")
	results := ccp.Run(&runContext)

	assert.NotEmpty(t, results)
	assert.Equal(t, check.StatusError, results[0].Status)
	assert.Contains(t, results[0].Message, "failed to collect")
}


func TestConjurConfigPermissions_Run_ReadAllError(t *testing.T) {
	// Save original readAllFunc to restore after test
	originalReadAllFunc := readAllFunc
	defer func() { readAllFunc = originalReadAllFunc }()

	readAllError := errors.New("read error")
	readAllFunc = func(r io.Reader) ([]byte, error) {
		return nil, readAllError
	}

	provider := &test.ContainerProvider{
		ExecResponses: map[string]test.ExecResponse{
			"ls -la /etc/conjur/config": {
				Stdout: strings.NewReader("drwxr-xr-x"),
			},
		},
	}

	ccp := &ConjurConfigPermissions{Provider: provider}

	runContext := test.NewRunContext("test-container")
	results := ccp.Run(&runContext)

	assert.NotEmpty(t, results)
	assert.Equal(t, check.StatusError, results[0].Status)
	assert.Contains(t, results[0].Message, "failed to read Conjur config permissions")
	assert.Contains(t, results[0].Message, readAllError.Error())
}
