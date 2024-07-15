package checks

import (
	"io"
	"strings"
	"testing"

	"github.com/cyberark/conjur-inspect/pkg/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUlimitRun(t *testing.T) {
	// Mock dependencies
	oldFunc := executeUlimitInfoFunc
	executeUlimitInfoFunc = func() (stderr, stdout io.Reader, err error) {
		stdout = strings.NewReader(
			"core file size      (blocks, -c) 0\npipe size      (512 bytes, -p) 1\nopen files      (-n) 6140\n",
		)
		return stdout, stderr, err
	}
	defer func() {
		executeUlimitInfoFunc = oldFunc
	}()

	// Run the check
	ulimit := &Ulimit{}
	context := test.NewRunContext("")
	results := <-ulimit.Run(&context)

	coreFileSize := GetResultByTitle(results, "core file size (blocks, -c)")
	require.NotNil(t, coreFileSize, "Includes 'core file size (blocks, -c)'")
	assert.Equal(t, "INFO", coreFileSize.Status)
	assert.Equal(t, "0", coreFileSize.Value)

	pipeSize := GetResultByTitle(results, "pipe size (512 bytes, -p)")
	require.NotNil(t, pipeSize, "Includes 'pipe size (512 bytes, -p)'")
	assert.Equal(t, "INFO", pipeSize.Status)
	assert.Equal(t, "1", pipeSize.Value)

	openFiles := GetResultByTitle(results, "open files (-n)")
	require.NotNil(t, openFiles, "Includes 'open files (-n)'")
	assert.Equal(t, "INFO", openFiles.Status)
	assert.Equal(t, "6140", openFiles.Value)
}
