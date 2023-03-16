package formatting

import (
	"encoding/json"
	"io"

	"github.com/cyberark/conjur-inspect/pkg/report"
)

// JSON renders a report result as JSON
type JSON struct{}

func (*JSON) Write(
	writer io.Writer,
	result *report.Result,
) error {

	encoder := json.NewEncoder(writer)

	encoder.SetIndent("", " ")

	return encoder.Encode(result)
}
