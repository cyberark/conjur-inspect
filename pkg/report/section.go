package report

import "github.com/cyberark/conjur-inspect/pkg/check"

// Section is the category of check
type Section struct {
	Title  string        `json:"title"`
	Checks []check.Check `json:"checks"`
}
