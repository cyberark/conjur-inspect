package shell

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandWrapper_Run(t *testing.T) {
	cmd := NewCommandWrapper("echo", "hello world")

	stdout, stderr, err := cmd.Run()
	assert.NoError(t, err)
	assert.Empty(t, stderr)
	assert.Equal(t, "hello world\n", string(stdout))
}

func TestCommandWrapper_Run_Error(t *testing.T) {
	cmd := NewCommandWrapper("invalid_command", "hello world")

	stdout, stderr, err := cmd.Run()
	assert.Error(t, err)
	assert.Empty(t, stderr)
	assert.Empty(t, stdout)
}
