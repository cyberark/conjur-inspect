package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRootCommand(t *testing.T) {
	rootCmd := newRootCommand()

	// Sanity test that the root command is not nil
	assert.NotNil(t, rootCmd)
}

// The only scenario we can't adequately test is when this is actually
// a terminal
func TestIsTerminal(t *testing.T) {
	// Test with a file writer
	file, err := os.CreateTemp("", "test")
	assert.Nil(t, err)
	defer os.Remove(file.Name())
	assert.False(t, isTerminal(file))

	// Test with a bytes buffer
	var buffer bytes.Buffer
	assert.False(t, isTerminal(&buffer))
}

func TestExecute(t *testing.T) {
	// Redirect stdout to a buffer so we can capture the output
	var stdout bytes.Buffer
	stderr := io.Discard

	// Execute the root command
	Execute(&stdout, stderr)

	// Test that the output is not empty
	assert.NotEmpty(t, stdout.String())
}

func TestExecuteWithJSONOutput(t *testing.T) {
	// Redirect stdout to a buffer so we can capture the output
	var stdout bytes.Buffer
	stderr := io.Discard

	// Execute the root command with the json flag set
	os.Args = []string{"conjur-inspect", "--json"}
	Execute(&stdout, stderr)

	// Test that the output is not empty
	assert.NotEmpty(t, stdout.String())

	// Test that the output is in JSON format
	assert.True(t, strings.HasPrefix(stdout.String(), "{"))
	assert.True(t, strings.HasSuffix(stdout.String(), "}\n"))
}

func TestExecuteWithDebugOutput(t *testing.T) {
	// Redirect stdout to a buffer so we can capture the output
	var stdout bytes.Buffer
	stderr := io.Discard

	// Execute the root command with the debug flag set
	os.Args = []string{"conjur-inspect", "--debug"}
	Execute(&stdout, stderr)

	// Test that the output is not empty
	assert.NotEmpty(t, stdout.String())
}
