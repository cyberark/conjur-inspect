package shell

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandWrapper_Run(t *testing.T) {
	cmd := NewCommandWrapper("echo", "hello world")

	stdoutReader, stderrReader, err := cmd.Run()
	assert.NoError(t, err)

	stderr, err := io.ReadAll(stderrReader)
	assert.NoError(t, err)
	assert.Empty(t, stderr)

	stdout, err := io.ReadAll(stdoutReader)
	assert.NoError(t, err)
	assert.Equal(t, "hello world\n", string(stdout))
}

func TestCommandWrapper_Run_Error(t *testing.T) {
	cmd := NewCommandWrapper("invalid_command", "hello world")

	stdoutReader, stderrReader, err := cmd.Run()
	assert.Error(t, err)

	stderr, err := io.ReadAll(stderrReader)
	assert.NoError(t, err)
	assert.Empty(t, stderr)

	stdout, err := io.ReadAll(stdoutReader)
	assert.NoError(t, err)
	assert.Empty(t, stdout)
}

func TestCommandWrapper_RunCombinedOutput(t *testing.T) {
	cmd := NewCommandWrapper("echo", "hello world")

	stdoutReader, err := cmd.RunCombinedOutput()
	assert.NoError(t, err)

	stdout, err := io.ReadAll(stdoutReader)
	assert.NoError(t, err)
	assert.Equal(t, "hello world\n", string(stdout))
}

func TestCommandWrapper_RunCombinedOutput_Error(t *testing.T) {
	cmd := NewCommandWrapper("invalid_command", "hello world")

	stdoutReader, err := cmd.RunCombinedOutput()
	assert.Error(t, err)

	stdout, err := io.ReadAll(stdoutReader)
	assert.NoError(t, err)
	assert.Empty(t, stdout)
}
