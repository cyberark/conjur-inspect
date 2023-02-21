package disk_test

import (
	"regexp"
	"testing"

	"github.com/conjurinc/conjur-preflight/pkg/checks/disk"
	"github.com/conjurinc/conjur-preflight/pkg/framework"
	"github.com/stretchr/testify/assert"
)

func TestSpaceCheck(t *testing.T) {
	testCheck := &disk.SpaceCheck{}
	resultChan := testCheck.Run()
	results := <-resultChan

	assert.Greater(t, len(results), 0, "There are disk space results present")

	for _, result := range results {
		assert.Regexp(
			t,
			regexp.MustCompile(`Disk Space (.+, .+)`),
			result.Title,
			"Disk space title matches the expected format",
		)
		assert.Equal(t, framework.STATUS_INFO, result.Status)
		assert.Regexp(
			t,
			regexp.MustCompile(`.+ Total, .+ Used \( ?\d+%\), .+ Free`),
			result.Value,
			"Disk space is in the expected format",
		)
	}
}
