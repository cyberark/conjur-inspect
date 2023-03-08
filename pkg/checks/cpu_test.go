package checks

import (
	"regexp"
	"testing"

	"github.com/cyberark/conjur-inspect/pkg/check"
	"github.com/stretchr/testify/assert"
)

func TestCpuRun(t *testing.T) {
	testCheck := &Cpu{}
	resultChan := testCheck.Run()
	results := <-resultChan

	// Ensure the result includes a CPU Cores value
	cpuCores := GetResultByTitle(results, "CPU Cores")
	assert.NotNil(t, cpuCores, "CPU results includes 'CPU Cores'")
	assert.Equal(t, check.STATUS_INFO, cpuCores.Status)
	assert.Regexp(t, regexp.MustCompile(`\d+`), cpuCores.Value, "CPU cores are in the expected format")

	// Ensure the result includes a CPU archiecture value
	cpuArchitecture := GetResultByTitle(results, "CPU Architecture")
	assert.NotNil(t, cpuArchitecture, "CPU results includes 'CPU Cores'")
	assert.Equal(t, check.STATUS_INFO, cpuArchitecture.Status)
	assert.NotEmpty(t, cpuArchitecture.Value, "CPU architecture is not empty")
}
