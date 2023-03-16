package report_test

import (
	"testing"

	"github.com/cyberark/conjur-inspect/pkg/report"
	"github.com/stretchr/testify/assert"
)

func TestNewDefaultReport(t *testing.T) {
	// We don't want to re-define the default report structure here, so we just
	// call the constructor to ensure we don't introduce any runtime errors and
	// that there are some report sections.

	report := report.NewDefaultReport(false)

	assert.Greater(t, len(report.Sections), 0)
}
