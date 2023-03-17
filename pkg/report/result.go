package report

import "github.com/cyberark/conjur-inspect/pkg/check"

// Result contains each sections check result
type Result struct {
	Version  string          `json:"version"`
	Sections []ResultSection `json:"sections"`
}

// ResultSection is the individual check and its result
type ResultSection struct {
	Title   string         `json:"title"`
	Results []check.Result `json:"results"`
}
