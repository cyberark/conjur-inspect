package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	// Save the original arguments
	originalArgs := os.Args

	// Sanity test for the program entrypoint
	os.Args = []string{"conjur-preflight"}

	// Use local buffers rather than actually standard out and error
	var localStdout, localStderr bytes.Buffer

	cmdStdout = &localStdout
	cmdStderr = &localStderr

	main()

	assert.Contains(
		t,
		localStdout.String(),
		"Conjur Enterprise Preflight Qualification",
	)

	assert.Empty(t, localStderr.String())

	// Restore the original arguments
	os.Args = originalArgs
}

func TestMainError(t *testing.T) {
	// Save the original arguments
	originalArgs := os.Args

	// Sanity test for the program entrypoint
	os.Args = []string{"conjur-preflight", "bogus"}

	// Use local buffers rather than actually standard out and error
	var localStdout, localStderr bytes.Buffer

	cmdStdout = &localStdout
	cmdStderr = &localStderr

	main()

	assert.Contains(
		t,
		localStdout.String(),
		"Conjur Enterprise Preflight Qualification",
	)

	assert.Empty(t, localStderr.String())

	// Restore the original arguments
	os.Args = originalArgs
}
