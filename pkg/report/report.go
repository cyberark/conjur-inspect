package report

import "time"

// Report contains an array of all sections and their reports
type Report interface {
	ID() string
	Run(config RunConfig) Result
}

// RunConfig contains the report run parameters
type RunConfig struct {
	ContainerID string
	Since       time.Duration
}
