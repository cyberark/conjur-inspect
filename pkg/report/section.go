package report

import "github.com/cyberark/conjur-inspect/pkg/framework"

// Section is the category of check
type Section struct {
	Title  string            `json:"title"`
	Checks []framework.Check `json:"checks"`
}
