package checks

import (
	"runtime"
	"strconv"

	"github.com/cyberark/conjur-inspect/pkg/check"
)

// Cpu collects inspection information on the host machines CPU cores and
// architecture
type Cpu struct {
}

// Describe provides a textual description of the info this check gathers
func (*Cpu) Describe() string {
	return "CPU"
}

// Run executes the CPU inspection checks
func (cpu *Cpu) Run(_context *check.RunContext) <-chan []check.Result {
	future := make(chan []check.Result)

	// TODO: Can we return avg recent utilization?
	go func() {
		future <- []check.Result{
			{
				Title:   "CPU Cores",
				Status:  check.StatusInfo,
				Value:   strconv.Itoa(runtime.NumCPU()),
				Message: "",
			},
			{
				Title:  "CPU Architecture",
				Status: check.StatusInfo,
				Value:  runtime.GOARCH,
			},
		}
	}() // async

	return future
}
