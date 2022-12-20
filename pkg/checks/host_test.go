package checks_test

import (
	"testing"

	"github.com/conjurinc/conjur-preflight/pkg/checks"
	"github.com/conjurinc/conjur-preflight/pkg/framework"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
)

func TestHostRun(t *testing.T) {
	testCheck := &checks.Host{}
	resultChan := testCheck.Run()
	results := <-resultChan

	hostname := GetResultByTitle(results, "Hostname")
	assert.NotNil(t, hostname, "Includes 'Hostname'")
	assert.Equal(t, framework.STATUS_INFO, hostname.Status)
	assert.NotEmpty(t, hostname.Value)

	uptime := GetResultByTitle(results, "Uptime")
	assert.NotNil(t, uptime, "Includes 'Uptime'")
	assert.Equal(t, framework.STATUS_INFO, uptime.Status)
	assert.NotEmpty(t, uptime.Value)

	os := GetResultByTitle(results, "OS")
	assert.NotNil(t, os, "Includes 'OS'")
	assert.Equal(t, framework.STATUS_INFO, os.Status)
	assert.NotEmpty(t, os.Value)

	virtualization := GetResultByTitle(results, "Virtualization")
	assert.NotNil(t, virtualization, "Includes 'Virtualization'")
	assert.Equal(t, framework.STATUS_INFO, virtualization.Status)
	assert.NotEmpty(t, virtualization.Value)
}

func GetResultByTitle(
	results []framework.CheckResult,
	title string,
) *framework.CheckResult {
	idx := slices.IndexFunc(
		results,
		func(c framework.CheckResult) bool { return c.Title == title },
	)

	if idx < 0 {
		return nil
	}

	return &results[idx]
}
