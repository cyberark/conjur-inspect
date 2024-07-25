package disk

import (
	"errors"
	"regexp"
	"testing"

	"github.com/cyberark/conjur-inspect/pkg/check"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/stretchr/testify/assert"
)

func TestSpaceCheck(t *testing.T) {
	testCheck := &SpaceCheck{}
	results := testCheck.Run(&check.RunContext{})

	assert.Greater(t, len(results), 0, "There are disk space results present")

	for _, result := range results {
		assert.Regexp(
			t,
			regexp.MustCompile(`Disk Space (.+, .+)`),
			result.Title,
			"Disk space title matches the expected format",
		)
		assert.Equal(t, check.StatusInfo, result.Status)
		assert.Regexp(
			t,
			regexp.MustCompile(`.+ Total, .+ Used \( ?\d+%\), .+ Free`),
			result.Value,
			"Disk space is in the expected format",
		)
	}
}

func TestPartitionListError(t *testing.T) {
	// Double the usage function to simulate an error
	originalPartitionsFunc := getPartitions
	getPartitions = failedPartitionsFunc
	defer func() {
		getPartitions = originalPartitionsFunc
	}()

	testCheck := &SpaceCheck{}
	results := testCheck.Run(&check.RunContext{})

	assert.Len(t, results, 1)

	errResult := results[0]
	assert.Equal(t, "disk capacity", errResult.Title)
	assert.Equal(t, "unable to list disk partitions: test partitions failure", errResult.Message)
}

func TestDiskUsageError(t *testing.T) {
	// Double the usage function to simulate an error
	originalUsageFunc := getUsage
	getUsage = failedUsageFunc
	defer func() {
		getUsage = originalUsageFunc
	}()

	testCheck := &SpaceCheck{}
	results := testCheck.Run(&check.RunContext{})

	assert.Empty(t, results)
}

func failedPartitionsFunc(all bool) ([]disk.PartitionStat, error) {
	return nil, errors.New("test partitions failure")
}

func failedUsageFunc(path string) (*disk.UsageStat, error) {
	return nil, errors.New("test usage failure")
}
