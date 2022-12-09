package checks_test

import (
	"testing"

	"github.com/conjurinc/conjur-preflight/pkg/checks"
	"github.com/conjurinc/conjur-preflight/pkg/framework"
	"github.com/stretchr/testify/assert"
)

func TestMemoryRun(t *testing.T) {
	testCheck := &checks.Memory{}
	resultChan := testCheck.Run()
	results := <-resultChan

	memoryTotal := GetResultByTitle(results, "Memory Total")
	assert.NotNil(t, memoryTotal, "Includes 'Memory Total'")
	assert.Equal(t, framework.STATUS_INFO, memoryTotal.Status)
	assert.NotEmpty(t, memoryTotal.Value)

	memoryFree := GetResultByTitle(results, "Memory Free")
	assert.NotNil(t, memoryFree, "Includes 'Memory Free'")
	assert.Equal(t, framework.STATUS_INFO, memoryFree.Status)
	assert.NotEmpty(t, memoryFree.Value)

	memoryUsed := GetResultByTitle(results, "Memory Used")
	assert.NotNil(t, memoryUsed, "Includes 'Memory Used'")
	assert.Equal(t, framework.STATUS_INFO, memoryUsed.Status)
	assert.NotEmpty(t, memoryUsed.Value)
}
