package checks

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/cyberark/conjur-inspect/pkg/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRubyThreadDumpDescribe(t *testing.T) {
	provider := &test.ContainerProvider{}
	rtd := &RubyThreadDump{Provider: provider}
	assert.Equal(t, "Ruby service thread dumps (Test Container Provider)", rtd.Describe())
}

func TestRubyThreadDumpRunNoContainerID(t *testing.T) {
	provider := &test.ContainerProvider{}
	rtd := &RubyThreadDump{Provider: provider}

	// Create run context with empty container ID
	runContext := test.NewRunContext("")
	results := rtd.Run(&runContext)

	// Should return empty results
	assert.Empty(t, results)

	// Should not have called any container operations
	items, err := runContext.OutputStore.Items()
	require.NoError(t, err)
	assert.Empty(t, items)
}

func TestRubyThreadDumpRunNoRubyProcesses(t *testing.T) {
	// Mock pgrep returning no PIDs (empty output)
	provider := &test.ContainerProvider{
		ExecResponses: map[string]test.ExecResponse{
			"sh -c pgrep -f ruby || true": {
				Stdout: strings.NewReader(""),
				Stderr: strings.NewReader(""),
				Error:  nil,
			},
		},
	}

	rtd := &RubyThreadDump{Provider: provider}
	runContext := test.NewRunContext("container123")
	results := rtd.Run(&runContext)

	// Should return empty results when no Ruby processes found
	assert.Empty(t, results)

	// Should not have saved any files
	items, err := runContext.OutputStore.Items()
	require.NoError(t, err)
	assert.Empty(t, items)
}

func TestRubyThreadDumpRunSuccessfulSingleProcess(t *testing.T) {
	threadDumpContent := `Sigdump at 2025-12-09 10:00:00 +0000 process 1234
  Thread #<Thread:0x00000001424518> status=run priority=0
      /app/conjur/lib/conjur.rb:32:in 'method'
      /app/conjur/lib/conjur.rb:18:in 'block'
  GC stat:
      count: 34
  Built-in objects:
   367,492: TOTAL
`

	provider := &test.ContainerProvider{
		ExecResponses: map[string]test.ExecResponse{
			"sh -c pgrep -f ruby || true": {
				Stdout: strings.NewReader("1234"),
				Stderr: strings.NewReader(""),
				Error:  nil,
			},
			"sh -c kill -CONT 1234 && sleep 1 && cat /tmp/sigdump-1234.log && rm /tmp/sigdump-1234.log": {
				Stdout: strings.NewReader(threadDumpContent),
				Stderr: strings.NewReader(""),
				Error:  nil,
			},
		},
	}

	rtd := &RubyThreadDump{Provider: provider}
	runContext := test.NewRunContext("container123")
	results := rtd.Run(&runContext)

	// Should return empty results on success
	assert.Empty(t, results)

	// Verify the thread dump file was saved
	items, err := runContext.OutputStore.Items()
	require.NoError(t, err)
	require.Len(t, items, 1)

	info, err := items[0].Info()
	require.NoError(t, err)
	assert.Equal(t, "ruby-dump-1234.txt", info.Name())

	// Verify file content
	outputStoreItemReader, cleanup, err := items[0].Open()
	defer cleanup()
	require.NoError(t, err)

	savedContent, err := io.ReadAll(outputStoreItemReader)
	require.NoError(t, err)
	assert.Equal(t, threadDumpContent, string(savedContent))
}

func TestRubyThreadDumpRunSuccessfulMultipleProcesses(t *testing.T) {
	threadDump1 := "Thread dump for process 1234\n"
	threadDump2 := "Thread dump for process 5678\n"
	threadDump3 := "Thread dump for process 9012\n"

	provider := &test.ContainerProvider{
		ExecResponses: map[string]test.ExecResponse{
			"sh -c pgrep -f ruby || true": {
				Stdout: strings.NewReader("1234\n5678\n9012"),
				Stderr: strings.NewReader(""),
				Error:  nil,
			},
			"sh -c kill -CONT 1234 && sleep 1 && cat /tmp/sigdump-1234.log && rm /tmp/sigdump-1234.log": {
				Stdout: strings.NewReader(threadDump1),
				Stderr: strings.NewReader(""),
				Error:  nil,
			},
			"sh -c kill -CONT 5678 && sleep 1 && cat /tmp/sigdump-5678.log && rm /tmp/sigdump-5678.log": {
				Stdout: strings.NewReader(threadDump2),
				Stderr: strings.NewReader(""),
				Error:  nil,
			},
			"sh -c kill -CONT 9012 && sleep 1 && cat /tmp/sigdump-9012.log && rm /tmp/sigdump-9012.log": {
				Stdout: strings.NewReader(threadDump3),
				Stderr: strings.NewReader(""),
				Error:  nil,
			},
		},
	}

	rtd := &RubyThreadDump{Provider: provider}
	runContext := test.NewRunContext("container123")
	results := rtd.Run(&runContext)

	// Should return empty results on success
	assert.Empty(t, results)

	// Verify all three thread dump files were saved
	items, err := runContext.OutputStore.Items()
	require.NoError(t, err)
	require.Len(t, items, 3)

	// Check filenames
	fileNames := make([]string, 3)
	for i, item := range items {
		info, err := item.Info()
		require.NoError(t, err)
		fileNames[i] = info.Name()
	}
	assert.Contains(t, fileNames, "ruby-dump-1234.txt")
	assert.Contains(t, fileNames, "ruby-dump-5678.txt")
	assert.Contains(t, fileNames, "ruby-dump-9012.txt")
}

func TestRubyThreadDumpRunIndividualPIDFailure(t *testing.T) {
	threadDump1 := "Thread dump for process 1234\n"
	threadDump3 := "Thread dump for process 9012\n"

	provider := &test.ContainerProvider{
		ExecResponses: map[string]test.ExecResponse{
			"sh -c pgrep -f ruby || true": {
				Stdout: strings.NewReader("1234\n5678\n9012"),
				Stderr: strings.NewReader(""),
				Error:  nil,
			},
			"sh -c kill -CONT 1234 && sleep 1 && cat /tmp/sigdump-1234.log && rm /tmp/sigdump-1234.log": {
				Stdout: strings.NewReader(threadDump1),
				Stderr: strings.NewReader(""),
				Error:  nil,
			},
			"sh -c kill -CONT 5678 && sleep 1 && cat /tmp/sigdump-5678.log && rm /tmp/sigdump-5678.log": {
				Stdout: strings.NewReader(""),
				Stderr: strings.NewReader("cat: /tmp/sigdump-5678.log: No such file or directory"),
				Error:  fmt.Errorf("command failed"),
			},
			"sh -c kill -CONT 9012 && sleep 1 && cat /tmp/sigdump-9012.log && rm /tmp/sigdump-9012.log": {
				Stdout: strings.NewReader(threadDump3),
				Stderr: strings.NewReader(""),
				Error:  nil,
			},
		},
	}

	rtd := &RubyThreadDump{Provider: provider}
	runContext := test.NewRunContext("container123")
	results := rtd.Run(&runContext)

	// Should return empty results (failures are logged, not returned as errors)
	assert.Empty(t, results)

	// Verify only two thread dump files were saved (PID 5678 failed)
	items, err := runContext.OutputStore.Items()
	require.NoError(t, err)
	require.Len(t, items, 2)

	// Check filenames - should have 1234 and 9012, but not 5678
	fileNames := make([]string, 2)
	for i, item := range items {
		info, err := item.Info()
		require.NoError(t, err)
		fileNames[i] = info.Name()
	}
	assert.Contains(t, fileNames, "ruby-dump-1234.txt")
	assert.Contains(t, fileNames, "ruby-dump-9012.txt")
	assert.NotContains(t, fileNames, "ruby-dump-5678.txt")
}

func TestRubyThreadDumpRunEmptyThreadDump(t *testing.T) {
	// Test case where sigdump produces no output (e.g., sigdump not installed)
	provider := &test.ContainerProvider{
		ExecResponses: map[string]test.ExecResponse{
			"sh -c pgrep -f ruby || true": {
				Stdout: strings.NewReader("1234"),
				Stderr: strings.NewReader(""),
				Error:  nil,
			},
			"sh -c kill -CONT 1234 && sleep 1 && cat /tmp/sigdump-1234.log && rm /tmp/sigdump-1234.log": {
				Stdout: strings.NewReader(""),
				Stderr: strings.NewReader("cat: /tmp/sigdump-1234.log: No such file or directory"),
				Error:  nil,
			},
		},
	}

	rtd := &RubyThreadDump{Provider: provider}
	runContext := test.NewRunContext("container123")
	results := rtd.Run(&runContext)

	// Should return empty results
	assert.Empty(t, results)

	// Should not have saved any files (empty output is logged as warning)
	items, err := runContext.OutputStore.Items()
	require.NoError(t, err)
	assert.Empty(t, items)
}

func TestRubyThreadDumpRunPgrepFailure(t *testing.T) {
	provider := &test.ContainerProvider{
		ExecResponses: map[string]test.ExecResponse{
			"sh -c pgrep -f ruby || true": {
				Stdout: strings.NewReader(""),
				Stderr: strings.NewReader("pgrep: command not found"),
				Error:  fmt.Errorf("pgrep command failed"),
			},
		},
	}

	rtd := &RubyThreadDump{Provider: provider}
	runContext := test.NewRunContext("container123")
	results := rtd.Run(&runContext)

	// Should return empty results (failure is logged as warning)
	assert.Empty(t, results)

	// Should not have saved any files
	items, err := runContext.OutputStore.Items()
	require.NoError(t, err)
	assert.Empty(t, items)
}

func TestRubyThreadDumpRunWhitespaceInPIDs(t *testing.T) {
	threadDump := "Thread dump content\n"

	provider := &test.ContainerProvider{
		ExecResponses: map[string]test.ExecResponse{
			"sh -c pgrep -f ruby || true": {
				Stdout: strings.NewReader("1234\n\n5678\n  \n9012  "),
				Stderr: strings.NewReader(""),
				Error:  nil,
			},
			"sh -c kill -CONT 1234 && sleep 1 && cat /tmp/sigdump-1234.log && rm /tmp/sigdump-1234.log": {
				Stdout: strings.NewReader(threadDump),
				Stderr: strings.NewReader(""),
				Error:  nil,
			},
			"sh -c kill -CONT 5678 && sleep 1 && cat /tmp/sigdump-5678.log && rm /tmp/sigdump-5678.log": {
				Stdout: strings.NewReader(threadDump),
				Stderr: strings.NewReader(""),
				Error:  nil,
			},
			"sh -c kill -CONT 9012 && sleep 1 && cat /tmp/sigdump-9012.log && rm /tmp/sigdump-9012.log": {
				Stdout: strings.NewReader(threadDump),
				Stderr: strings.NewReader(""),
				Error:  nil,
			},
		},
	}

	rtd := &RubyThreadDump{Provider: provider}
	runContext := test.NewRunContext("container123")
	results := rtd.Run(&runContext)

	// Should return empty results on success
	assert.Empty(t, results)

	// Should handle whitespace correctly and save 3 files
	items, err := runContext.OutputStore.Items()
	require.NoError(t, err)
	require.Len(t, items, 3)
}
