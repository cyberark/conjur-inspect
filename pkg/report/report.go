package report

// Report contains an array of all sections and their reports
type Report interface {
	ID() string
	Run(ContainerID string) Result
}
