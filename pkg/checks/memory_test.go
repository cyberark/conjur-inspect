package checks

import (
	"errors"
	"testing"

	"github.com/cyberark/conjur-inspect/pkg/check"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/stretchr/testify/assert"
)

func TestMemoryRun(t *testing.T) {
	testCheck := &Memory{}
	results := testCheck.Run(&check.RunContext{})

	memoryTotal := GetResultByTitle(results, "Memory Total")
	assert.NotNil(t, memoryTotal, "Includes 'Memory Total'")
	assert.Equal(t, check.StatusInfo, memoryTotal.Status)
	assert.NotEmpty(t, memoryTotal.Value)

	memoryFree := GetResultByTitle(results, "Memory Free")
	assert.NotNil(t, memoryFree, "Includes 'Memory Free'")
	assert.Equal(t, check.StatusInfo, memoryFree.Status)
	assert.NotEmpty(t, memoryFree.Value)

	memoryUsed := GetResultByTitle(results, "Memory Used")
	assert.NotNil(t, memoryUsed, "Includes 'Memory Used'")
	assert.Equal(t, check.StatusInfo, memoryUsed.Status)
	assert.NotEmpty(t, memoryUsed.Value)
}

func TestMemoryRunError(t *testing.T) {
	// Double the virtual memory function to simulate an error
	originalVirtualMemory := getVirtualMemory
	getVirtualMemory = failedVirtualMemoryFunc
	defer func() {
		getVirtualMemory = originalVirtualMemory
	}()

	testCheck := &Memory{}
	results := testCheck.Run(&check.RunContext{})

	assert.Len(t, results, 1)

	errResult := results[0]
	assert.Equal(t, "memory", errResult.Title)
	assert.Equal(t, "failed to inspect memory: test virtual memory failure", errResult.Message)
}

func failedVirtualMemoryFunc() (*mem.VirtualMemoryStat, error) {
	return nil, errors.New("test virtual memory failure")
}
