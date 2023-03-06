package formatting

import (
	"io"

	"github.com/conjurinc/conjur-preflight/pkg/report"
)

// Writer represents an object than can render a
// ReportResult to an io.Writer (e.g. file, stdout)
type Writer interface {
	Write(io.Writer, *report.Result) error
}
