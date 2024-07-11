package cmd_test

import (
	"testing"

	"github.com/cyberark/conjur-inspect/pkg/cmd"
	"github.com/stretchr/testify/assert"
)

func TestNewDefaultReport(t *testing.T) {
	// We don't want to re-define the default report structure here, so we just
	// call the constructor to ensure we don't introduce any runtime errors and
	// that there are some report sections.

	id := "test-id"

	report, err := cmd.NewDefaultReport(id, ".")

	assert.Equal(t, id, report.ID())
	assert.NotNil(t, report)
	assert.Nil(t, err)
}
