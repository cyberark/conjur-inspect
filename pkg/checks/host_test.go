package checks

import (
	"errors"
	"testing"

	"github.com/cyberark/conjur-inspect/pkg/check"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
)

func TestHostRun(t *testing.T) {
	testCheck := &Host{}
	resultChan := testCheck.Run(&check.RunContext{})
	results := <-resultChan

	hostname := GetResultByTitle(results, "Hostname")
	assert.NotNil(t, hostname, "Includes 'Hostname'")
	assert.Equal(t, check.STATUS_INFO, hostname.Status)
	assert.NotEmpty(t, hostname.Value)

	uptime := GetResultByTitle(results, "Uptime")
	assert.NotNil(t, uptime, "Includes 'Uptime'")
	assert.Equal(t, check.STATUS_INFO, uptime.Status)
	assert.NotEmpty(t, uptime.Value)

	os := GetResultByTitle(results, "OS")
	assert.NotNil(t, os, "Includes 'OS'")
	assert.Equal(t, check.STATUS_INFO, os.Status)
	assert.NotEmpty(t, os.Value)

	virtualization := GetResultByTitle(results, "Virtualization")
	assert.NotNil(t, virtualization, "Includes 'Virtualization'")
	assert.Equal(t, check.STATUS_INFO, virtualization.Status)
	assert.NotEmpty(t, virtualization.Value)
}

func TestHostRunError(t *testing.T) {
	// Double the host info function to simulate an error
	originalHostInfo := getHostInfo
	getHostInfo = failedHostInfoFunc
	defer func() {
		getHostInfo = originalHostInfo
	}()

	testCheck := &Host{}
	resultChan := testCheck.Run(&check.RunContext{})
	results := <-resultChan

	errResult := results[0]
	assert.Equal(t, "Error", errResult.Title)
	assert.Equal(t, "test host failure", errResult.Value)
}

func TestHostRunNoVirtualization(t *testing.T) {
	// Double the host info function to simulate no virtualization
	originalHostInfo := getHostInfo
	getHostInfo = noVirtualizationHostInfoFunc
	defer func() {
		getHostInfo = originalHostInfo
	}()

	testCheck := &Host{}
	resultChan := testCheck.Run(&check.RunContext{})
	results := <-resultChan

	virtualization := GetResultByTitle(results, "Virtualization")
	assert.NotNil(t, virtualization)
	assert.Equal(t, check.STATUS_INFO, virtualization.Status)
	assert.NotEmpty(t, "None")
}

func failedHostInfoFunc() (*host.InfoStat, error) {
	return nil, errors.New("test host failure")
}

func noVirtualizationHostInfoFunc() (*host.InfoStat, error) {
	original, err := host.Info()
	// We assume this will work on a test host. If not, abort the test run.
	if err != nil {
		panic(err)
	}

	original.VirtualizationSystem = ""

	return original, nil
}

func GetResultByTitle(
	results []check.Result,
	title string,
) *check.Result {
	idx := slices.IndexFunc(
		results,
		func(c check.Result) bool { return c.Title == title },
	)

	if idx < 0 {
		return nil
	}

	return &results[idx]
}
